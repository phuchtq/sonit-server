package businesslogic

import (
	"context"
	"database/sql"
	"errors"
	"log"
	"sonit_server/constant/currency"
	domain_status "sonit_server/constant/domain_status"
	"sonit_server/constant/noti"
	repo "sonit_server/data_access" // Order data access ~~ Order repository
	"sonit_server/data_access/db"
	db_server "sonit_server/data_access/db_server"
	business_logic "sonit_server/interface/business_logic"
	data_access "sonit_server/interface/data_access"
	"sonit_server/model/dto/request"
	"sonit_server/model/dto/response"
	"sonit_server/model/entity"
	"sonit_server/utils"
	"strings"
	"sync"
	"time"
)

type orderService struct {
	logger          *log.Logger
	userRepo        data_access.IUserRepo
	inventoryRepo   data_access.IProductInventoryRepo
	inventoryTxRepo data_access.IProductInventoryTransactionRepo
	shippingRepo    data_access.IShippingRepo
	cartRepo        data_access.ICartRepo
	orderRepo       data_access.IOrderRepo
}

func GenerateOrderService() (business_logic.IOrderService, error) {
	var logger = utils.GetLogConfig()

	cnn, err := db.ConnectDB(logger, db_server.InitializePostgreSQL())

	if err != nil {
		return nil, err
	}

	order_cnn = cnn

	return InitializeOrderService(cnn, logger), nil
}

func InitializeOrderService(db *sql.DB, logger *log.Logger) business_logic.IOrderService {
	return &orderService{
		logger:          logger,
		userRepo:        repo.InitializeUserRepo(db, logger),
		inventoryRepo:   repo.InitializeProductInventoryRepo(db, logger),
		inventoryTxRepo: repo.InitializeProductInventoryTransactionRepo(db, logger),
		cartRepo:        repo.InitializeCartRepo(db, logger),
		orderRepo:       repo.InitializeOrderRepo(db, logger),
	}
}

var order_cnn *sql.DB

// CreateOrder implements businesslogic.IOrderService.
func (o *orderService) CreateOrder(req request.CreateOrderRequest, ctx context.Context) error {
	var genericError error = errors.New(noti.GENERIC_ERROR_WARN_MSG)
	defer closeCnn(order_cnn)

	if !isEntityExist(o.userRepo, id_type, req.UserId, ctx) {
		return genericError
	}

	cart, _, err := o.cartRepo.GetCart(req.UserId, ctx)
	if err != nil {
		return err
	}

	var status string
	if utils.IsOrderStatusValid(req.Status) {
		status = req.Status
	} else {
		status = domain_status.ORDER_PENDING
	}

	var itemsInCart []response.CartItem = utils.JsonStringToObject[[]response.CartItem](cart.Items)

	var inventories []entity.ProductInventory = make([]entity.ProductInventory, len(req.Items))
	var inventoryTxs []entity.ProductInventoryTransaction = make([]entity.ProductInventoryTransaction, len(req.Items))
	var totalAmount float64
	var curTime time.Time = time.Now()

	for index, item := range req.Items {
		inventory, err := o.inventoryRepo.GetProductInventory(item.ProductId, ctx)
		if err != nil {
			return err
		}

		if inventory == nil {
			return genericError
		}

		if item.Quantity > int(inventory.CurrentQuantity) {
			return genericError
		}

		inventory.CurrentQuantity -= int64(item.Quantity)
		totalAmount += float64(item.Quantity) * item.Price

		if strings.Contains(cart.Items, item.ProductId) {
			for i, itemInCart := range itemsInCart {
				if itemInCart.ProductId == item.ProductId {
					itemsInCart[i].Quantity -= item.Quantity
					if itemsInCart[i].Quantity < 1 {
						itemsInCart = append(itemsInCart[:i], itemsInCart[i+1:]...)
					}
					break
				}
			}
		}

		inventories[index] = *inventory
		inventoryTxs[index] = entity.ProductInventoryTransaction{
			TransactionId: utils.GenerateId(),
			ProductId:     item.ProductId,
			Amount:        int64(item.Quantity),
			Action:        "Sale",
			Date:          curTime,
			CreatedAt:     curTime,
			UpdatedAt:     curTime,
		}
	}

	cart.Items = utils.ObjectToJsonString(itemsInCart)
	cart.UpdatedAt = curTime

	_, cancel := context.WithCancel(ctx)
	var wg sync.WaitGroup
	var mu sync.Mutex
	var capturedErr error

	wg.Add(4)

	// Create inventory transaction
	go func() {
		defer wg.Done()
		for _, tx := range inventoryTxs {
			if err := o.inventoryTxRepo.CreateProductInventoryTransaction(tx, ctx); err != nil {
				mu.Lock()
				if capturedErr == nil {
					capturedErr = err
					cancel()
				}
				mu.Unlock()
			}
		}
	}()

	// Update inventories
	go func() {
		defer wg.Done()
		for _, inventory := range inventories {
			if err := o.inventoryRepo.UpdateProductInventory(inventory, ctx); err != nil {
				mu.Lock()
				if capturedErr == nil {
					capturedErr = err
					cancel()
				}
				mu.Unlock()
			}
		}
	}()

	// Create order
	go func() {
		defer wg.Done()
		if err := o.orderRepo.CreateOrder(entity.Order{
			OrderId:     utils.GenerateId(),
			UserId:      req.UserId,
			Items:       utils.ObjectToJsonString(req.Items),
			TotalAmount: totalAmount,
			Currency:    currency.VIETNAM_DONG,
			Status:      status,
			Note:        req.Note,
			CreatedAt:   curTime,
			UpdatedAt:   curTime,
		}, ctx); err != nil {
			mu.Lock()
			if capturedErr == nil {
				capturedErr = err
				cancel()
			}
			mu.Unlock()
		}
	}()

	// Update cart if existed
	go func() {
		defer wg.Done()
		if cart != nil {
			if err := o.cartRepo.UpdateCart(*cart, ctx); err != nil {
				mu.Lock()
				if capturedErr == nil {
					capturedErr = err
					cancel()
				}
				mu.Unlock()
			}
		}
	}()

	wg.Wait()

	return capturedErr
}

// GetOrder implements businesslogic.IOrderService.
func (o *orderService) GetOrder(id string, ctx context.Context) (*response.ViewOrderResponse, error) {
	defer closeCnn(order_cnn)

	order, err := o.orderRepo.GetOrder(id, ctx)
	if err != nil {
		return nil, err
	}

	if order == nil {
		return nil, errors.New(noti.GENERIC_ERROR_WARN_MSG)
	}

	return &response.ViewOrderResponse{
		OrderId:     order.OrderId,
		UserId:      order.UserId,
		Items:       utils.JsonStringToObject[[]response.CartItem](order.Items),
		TotalAmount: order.TotalAmount,
		Currency:    order.Currency,
		Status:      order.Status,
		Note:        order.Note,
		CreatedAt:   order.CreatedAt,
		UpdatedAt:   order.UpdatedAt,
	}, nil
}

// GetOrders implements businesslogic.IOrderService.
func (o *orderService) GetOrders(req request.GetOrdersRequest, ctx context.Context) (response.PaginationDataResponse, error) {
	if req.Request.PageNumber < 1 {
		req.Request.PageNumber = 1
	}

	req.Request.FilterProp = utils.AssignFilterProperty(req.Request.FilterProp)
	req.Request.Order = utils.AssignOrder(req.Request.Order)
	defer closeCnn(order_cnn)

	orders, pages, err := o.orderRepo.GetOrders(req, ctx)

	var data []response.ViewOrderResponse
	if orders != nil {
		for _, order := range *orders {
			data = append(data, response.ViewOrderResponse{
				OrderId:     order.OrderId,
				UserId:      order.UserId,
				Items:       utils.JsonStringToObject[[]response.CartItem](order.Items),
				TotalAmount: order.TotalAmount,
				Currency:    order.Currency,
				Status:      order.Status,
				Note:        order.Note,
				CreatedAt:   order.CreatedAt,
				UpdatedAt:   order.UpdatedAt,
			})
		}
	}

	return response.PaginationDataResponse{
		Data:       data,
		TotalPages: pages,
		PageNumber: req.Request.PageNumber,
	}, err
}

// UpdateOrder implements businesslogic.IOrderService.
func (o *orderService) UpdateOrder(req request.UpdateOrderRequest, ctx context.Context) error {
	defer closeCnn(order_cnn)

	order, err := o.orderRepo.GetOrder(req.OrderId, ctx)
	if err != nil {
		return err
	}

	if order == nil {
		return errors.New(noti.GENERIC_ERROR_WARN_MSG)
	}

	if req.Currency != "" {
		order.Currency = utils.AssignCurrency(req.Currency)
	}

	if req.Status != "" {
		if !utils.IsOrderStatusValid(req.Status) {
			return errors.New(noti.GENERIC_ERROR_WARN_MSG)
		}

		order.Status = req.Status
	}

	if req.Note != "" {
		order.Note = req.Note
	}

	order.UpdatedAt = time.Now()

	return o.orderRepo.UpdateOrder(*order, ctx)
}
