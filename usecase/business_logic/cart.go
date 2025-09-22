package businesslogic

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"log"
	"math"
	"sonit_server/constant/noti"
	repo "sonit_server/data_access" // Category data access ~~ Category repository
	"sonit_server/data_access/db"
	db_server "sonit_server/data_access/db_server"
	business_logic "sonit_server/interface/business_logic"
	data_access "sonit_server/interface/data_access" // Category repo interface
	"sonit_server/model/dto/request"
	"sonit_server/model/dto/response"
	entity "sonit_server/model/entity"
	"sonit_server/utils"
	"strings"
	"time"
)

type cartService struct {
	productRepo   data_access.IProductRepo
	inventoryRepo data_access.IProductInventoryRepo
	cartRepo      data_access.ICartRepo
	logger        *log.Logger
}

func GenerateCartService() (business_logic.ICartService, error) {
	var logger = utils.GetLogConfig()

	cnn, err := db.ConnectDB(logger, db_server.InitializePostgreSQL())

	if err != nil {
		return nil, err
	}

	cart_cnn = cnn

	return InitializeCartService(cnn, logger), nil
}

func InitializeCartService(db *sql.DB, logger *log.Logger) business_logic.ICartService {
	return &cartService{
		productRepo:   repo.InitializeProductRepo(db, logger),
		inventoryRepo: repo.InitializeProductInventoryRepo(db, logger),
		cartRepo:      repo.InitializeCartRepo(db, logger),
		logger:        logger,
	}
}

var cart_cnn *sql.DB

// AddItemToCart implements businesslogic.ICartService.
func (c *cartService) AddItemToCart(req request.AddItemToCartRequest, ctx context.Context) error {
	// _, cancel := context.WithCancel(ctx)
	// var wg sync.WaitGroup
	// var mu sync.Mutex

	// wg.Add(3)
	// var capturedErr error
	// var cart *entity.Cart
	// var product *entity.Product
	// var inventory *entity.ProductInventory

	// go func() {
	// 	defer wg.Done()
	// 	product, capturedErr = c.productRepo.GetProductById(req.Request.ProductId, ctx)
	// 	if capturedErr != nil {
	// 		mu.Lock()
	// 		c.logger.Println("Product cooked")
	// 		cancel()
	// 		mu.Unlock()
	// 	}
	// }()

	// go func() {
	// 	defer wg.Done()
	// 	cart, _, capturedErr = c.cartRepo.GetCart(req.Request.UserId, ctx)
	// 	if capturedErr != nil {
	// 		mu.Lock()
	// 		c.logger.Println("Cart cooked")
	// 		cancel()
	// 		mu.Unlock()
	// 	}
	// }()

	// go func() {
	// 	defer wg.Done()
	// 	inventory, capturedErr = c.inventoryRepo.GetProductInventory(req.Request.ProductId, ctx)
	// 	if capturedErr != nil {
	// 		mu.Lock()
	// 		c.logger.Println("Product inventory cooked")
	// 		cancel()
	// 		mu.Unlock()
	// 	}
	// }()

	// // Wait for 2 goroutines to finish
	// wg.Wait()

	// // Error while creating account
	// if capturedErr != nil {
	// 	return capturedErr
	// }

	defer closeCnn(cart_cnn)

	product, err := c.productRepo.GetProductById(req.Request.ProductId, ctx)
	if err != nil {
		return err
	}

	if product == nil {
		return errors.New(noti.GENERIC_ERROR_WARN_MSG)
	}

	inventory, err := c.inventoryRepo.GetProductInventory(req.Request.ProductId, ctx)
	if err != nil {
		return err
	}

	if inventory == nil {
		return errors.New(noti.GENERIC_ERROR_WARN_MSG)
	}

	cart, _, err := c.cartRepo.GetCart(req.Request.UserId, ctx)
	if err != nil {
		return err
	}

	if product == nil {
		return errors.New(noti.GENERIC_ERROR_WARN_MSG)
	}

	var structItems []response.CartItem // Items in cart
	if cart != nil {
		structItems = utils.JsonStringToObject[[]response.CartItem](cart.Items)
	}

	var isFound bool = false
	var totalCartItemQuantity int = req.Quantity // Variable to store the total quantity of the product in user's cart
	for index, item := range structItems {
		// Item already in cart
		if item.ProductId == req.Request.ProductId {
			isFound = true
			totalCartItemQuantity += structItems[index].Quantity // Added the current quantity of cart to the total variable
			structItems[index].Quantity += req.Quantity
			break
		}
	}

	// Cart not existed or item not in cart yet?
	if cart == nil || !isFound {
		structItems = append(structItems, response.CartItem{
			ProductId: req.Request.ProductId,
			Name:      product.ProductName,
			ImageUrl:  product.Image,
			Quantity:  req.Quantity,
			Price:     float64(product.Price),
			Currency:  product.Currency,
		})
	}

	// If total items in cart bigger than the actual quantity in inventory
	if totalCartItemQuantity > int(inventory.CurrentQuantity) {
		return fmt.Errorf(noti.ITEM_OUT_OF_STOCK_WARN_MSG, totalCartItemQuantity)
	}

	var itemsStr string = utils.ObjectToJsonString(structItems)
	var curTime time.Time = time.Now()
	var expiredPeriod time.Time = curTime.AddDate(0, 0, 7)
	if cart == nil {
		return c.cartRepo.CreateCart(entity.Cart{
			UserId:    req.Request.UserId,
			Items:     itemsStr,
			CreatedAt: curTime,
			UpdatedAt: curTime,
			ExpiredAt: expiredPeriod,
		}, ctx)
	}

	cart.UpdatedAt = curTime
	cart.ExpiredAt = expiredPeriod
	cart.Items = itemsStr

	return c.cartRepo.UpdateCart(*cart, ctx)
}

// EditItemInCart implements businesslogic.ICartService.
func (c *cartService) EditItemInCart(req request.EditItemInCartRequest, ctx context.Context) error {
	defer closeCnn(cart_cnn)

	cart, _, err := c.cartRepo.GetCart(req.Request.UserId, ctx)
	if err != nil {
		return err
	}

	inventory, err := c.inventoryRepo.GetProductInventory(req.Request.ProductId, ctx)
	if err != nil {
		return err
	}

	// Cart not existed or items not available in cart
	if cart == nil || req.Quantity > int(inventory.CurrentQuantity) {
		return errors.New(noti.GENERIC_ERROR_WARN_MSG)
	}

	var structItems []response.CartItem = utils.JsonStringToObject[[]response.CartItem](cart.Items) // Items in cart
	var isFound bool = false
	for index, item := range structItems {
		// Item already in cart
		if item.ProductId == req.Request.ProductId {
			isFound = true
			structItems[index].Quantity = req.Quantity
			break
		}
	}

	// Product not found in cart
	if !isFound {
		return errors.New(noti.GENERIC_ERROR_WARN_MSG)
	}

	var itemsStr string = utils.ObjectToJsonString(structItems)
	var curTime time.Time = time.Now()

	cart.UpdatedAt = curTime
	cart.ExpiredAt = curTime.AddDate(0, 0, 7)
	cart.Items = itemsStr

	return c.cartRepo.UpdateCart(*cart, ctx)
}

// RemoveItem implements businesslogic.ICartService.
func (c *cartService) RemoveItem(req request.RemoveItemFromCartRequest, ctx context.Context) error {
	defer closeCnn(cart_cnn)

	cart, _, err := c.cartRepo.GetCart(req.UserId, ctx)
	if err != nil {
		return err
	}

	if cart == nil || !strings.Contains(cart.Items, req.ProductId) {
		return errors.New(noti.GENERIC_ERROR_WARN_MSG)
	}

	var items []response.CartItem = utils.JsonStringToObject[[]response.CartItem](cart.Items)
	for index, product := range items {
		if product.ProductId == req.ProductId {
			cart.Items = utils.ObjectToJsonString(append(items[:index], items[index+1:]...))
			break
		}
	}

	return c.cartRepo.UpdateCart(*cart, ctx)
}

// ViewCartDetail implements businesslogic.ICartService.
func (c *cartService) ViewCartDetail(id string, pageNumber int, ctx context.Context) (response.PaginationDataResponse, error) {
	if pageNumber <= 0 {
		pageNumber = 1
	}

	defer closeCnn(cart_cnn)

	cart, limitRecords, err := c.cartRepo.GetCart(id, ctx)
	if err != nil {
		return response.PaginationDataResponse{}, err
	}

	//var maxRecords int = limitRecords * pageNumber

	var items []response.CartItem
	if cart != nil && cart.Items != "" {
		items = utils.JsonStringToObject[[]response.CartItem](cart.Items)
	}

	// [maxRecords-limitRecords : maxRecords]
	return response.PaginationDataResponse{
		Data:       items,
		PageNumber: pageNumber,
		TotalPages: int(math.Ceil(float64(len(items)) / float64(limitRecords))),
	}, nil
}
