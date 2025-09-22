package businesslogic

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"log"
	"os"
	"sonit_server/constant/currency"
	domain_status "sonit_server/constant/domain_status"
	payment_env "sonit_server/constant/env/payment"
	mail_const "sonit_server/constant/mail_const"
	"sonit_server/constant/noti"
	payment_method "sonit_server/constant/payment_method"
	repo "sonit_server/data_access"
	"sonit_server/data_access/db"
	db_server "sonit_server/data_access/db_server"
	business_logic "sonit_server/interface/business_logic"
	data_access "sonit_server/interface/data_access"
	"sonit_server/model/dto/request"
	"sonit_server/model/dto/response"
	"sonit_server/model/entity"
	"strings"
	"time"

	"sonit_server/utils"

	"github.com/payOSHQ/payos-lib-golang"
)

type paymentService struct {
	logger       *log.Logger
	userRepo     data_access.IUserRepo
	cartRepo     data_access.ICartRepo
	invetoryRepo data_access.IProductInventoryRepo
	productRepo  data_access.IProductRepo
	shippingRepo data_access.IShippingRepo
	orderRepo    data_access.IOrderRepo
	paymentRepo  data_access.IPaymentRepo
}

func InitializePaymentService(db *sql.DB, logger *log.Logger) business_logic.IPaymentService {
	return &paymentService{
		logger:       logger,
		userRepo:     repo.InitializeUserRepo(db, logger),
		cartRepo:     repo.InitializeCartRepo(db, logger),
		invetoryRepo: repo.InitializeProductInventoryRepo(db, logger),
		productRepo:  repo.InitializeProductRepo(db, logger),
		shippingRepo: repo.InitializeShippingRepo(db, logger),
		orderRepo:    repo.InitializeOrderRepo(db, logger),
		paymentRepo:  repo.InitializePaymentRepo(db, logger),
	}
}

func GeneratePaymentService() (business_logic.IPaymentService, error) {
	var logger = utils.GetLogConfig()

	cnn, err := db.ConnectDB(logger, db_server.InitializePostgreSQL())

	if err != nil {
		return nil, err
	}

	payment_cnn = cnn

	return InitializePaymentService(cnn, logger), nil
}

var payment_cnn *sql.DB

// GetPaymentById implements businesslogic.IPaymentService.
func (p *paymentService) GetPaymentById(id string, ctx context.Context) (*entity.Payment, error) {
	defer closeCnn(payment_cnn)
	return p.paymentRepo.GetPaymentById(id, ctx)
}

// GetPayments implements businesslogic.IPaymentService.
func (p *paymentService) GetPayments(req request.GetPaymentsRequest, ctx context.Context) (response.PaginationDataResponse, error) {
	if req.Request.PageNumber < 1 {
		req.Request.PageNumber = 1
	}

	req.Request.FilterProp = utils.AssignFilterProperty(req.Request.FilterProp)
	req.Request.Order = utils.AssignOrder(req.Request.Order)

	defer closeCnn(payment_cnn)

	data, pages, err := p.paymentRepo.GetPayments(req, ctx)

	return response.PaginationDataResponse{
		Data:       data,
		PageNumber: req.Request.PageNumber,
		TotalPages: pages,
	}, err
}

// UpdatePayment implements businesslogic.IPaymentService.
func (p *paymentService) UpdatePayment(req request.UpdatePaymentRequest, ctx context.Context) error {
	payment, err := p.paymentRepo.GetPaymentById(req.PaymentId, ctx)
	if err != nil {
		return err
	}

	if payment == nil {
		return errors.New(noti.GENERIC_ERROR_WARN_MSG)
	}

	if req.Currency != "" {
		payment.Currency = utils.AssignCurrency(req.Currency)
	}

	// Must validate(implement later)
	if req.Method != "" {
		payment.Method = req.Method
	}

	// Must validate(implement later)
	if req.Status != "" {
		payment.Status = req.Status
	}

	return p.paymentRepo.UpdatePayment(*payment, ctx)
}

// CreatePayment implements businesslogic.IPaymentService.
func (p *paymentService) CreatePaymentThroughCart(req request.CreatePaymentThroughCartRequest, ctx context.Context) (response.UrlAPIResponse, error) {
	var errRes error = errors.New(noti.GENERIC_ERROR_WARN_MSG)
	var res response.UrlAPIResponse

	defer closeCnn(payment_cnn)

	if !isEntityExist(p.userRepo, req.UserId, id_type, ctx) {
		return res, errRes
	}

	cart, _, err := p.cartRepo.GetCart(req.UserId, ctx)
	if err != nil {
		return res, err
	}

	// No items in cart but execute payment
	if cart == nil {
		return res, errRes
	}

	var itemsInCart []response.CartItem = utils.JsonStringToObject[[]response.CartItem](cart.Items)
	var totalAmount float64
	var items []payos.Item
	var invetories []entity.ProductInventory
	var formatItems []response.CartItem

	for _, prod := range req.Items {
		// Item not existed in cart
		if !strings.Contains(cart.Items, prod.ProductId) {
			return res, errRes
		}

		inventory, err := p.invetoryRepo.GetProductInventory(prod.ProductId, ctx)
		if err != nil {
			return res, err
		}

		if inventory == nil {
			return res, errRes
		}

		product, err := p.productRepo.GetProductById(prod.ProductId, ctx)
		if err != nil {
			return res, err
		}

		if product == nil {
			return res, errRes
		}

		if inventory.CurrentQuantity < int64(prod.Quantity) {
			return res, errRes
		}

		inventory.CurrentQuantity -= int64(prod.Quantity)
		totalAmount += float64(prod.Quantity) * product.Price
		invetories = append(invetories, *inventory)

		items = append(items, payos.Item{
			Name:     product.ProductName,
			Quantity: prod.Quantity,
			Price:    int(product.Price),
		})

		formatItems = append(formatItems, response.CartItem{
			ProductId: prod.ProductId,
			Name:      product.ProductName,
			ImageUrl:  product.Image,
			Quantity:  prod.Quantity,
			Price:     product.Price,
			Currency:  product.Currency,
		})

		for index, item := range itemsInCart {
			if item.ProductId == prod.ProductId {
				itemsInCart[index].Quantity -= prod.Quantity
				if itemsInCart[index].Quantity <= 0 { // Remove item from cart
					itemsInCart = append(itemsInCart[:index], itemsInCart[index+1:]...)
				}

				break
			}
		}
	}

	var paymentId string = utils.GenerateId()
	var orderCode int = utils.GenerateNumber()

	// Create transaction url
	data, err := payos.CreatePaymentLink(payos.CheckoutRequestType{
		OrderCode:   int64(orderCode),
		Amount:      int(totalAmount),
		Items:       items,
		Description: fmt.Sprint(orderCode),
		ReturnUrl:   os.Getenv(payment_env.PAYMENT_CALLBACK_SUCCESS_PROCESS) + paymentId,
		CancelUrl:   os.Getenv(payment_env.PAYMENT_CALLBACK_CANCEL_PROCESS) + paymentId,
	})

	if err != nil {
		p.logger.Println("Err: ", err.Error())
		return res, errors.New(noti.INTERNALL_ERR_MSG)
	}

	// Update product inventories
	for _, inventory := range invetories {
		if err := p.invetoryRepo.UpdateProductInventory(inventory, ctx); err != nil {
			return res, err
		}
	}

	var curTime time.Time = time.Now()

	// No more items in cart
	if len(itemsInCart) == 0 {
		p.cartRepo.RemoveCart(req.UserId, ctx)
	} else {
		cart.Items = utils.ObjectToJsonString(itemsInCart)
		cart.UpdatedAt = curTime
		cart.ExpiredAt = curTime.AddDate(0, 0, 7)
		p.cartRepo.UpdateCart(*cart, ctx)
	}

	var orderId string = utils.GenerateId()

	// Create order
	if err := p.orderRepo.CreateOrder(entity.Order{
		OrderId:     orderId,
		UserId:      req.UserId,
		Items:       utils.ObjectToJsonString(formatItems),
		TotalAmount: totalAmount,
		Currency:    currency.VIETNAM_DONG,
		Note:        req.Note,
		Status:      domain_status.ORDER_PENDING,
		CreatedAt:   curTime,
		UpdatedAt:   curTime,
	}, ctx); err != nil {
		return res, err
	}

	// Create payment
	if err := p.paymentRepo.CreatePayment(entity.Payment{
		PaymentId:     paymentId,
		OrderId:       orderId,
		UserId:        req.UserId,
		TransactionId: fmt.Sprint(orderCode),
		Amount:        totalAmount,
		Currency:      currency.VIETNAM_DONG,
		Status:        domain_status.PAYMENT_PENDING,
		Method:        payment_method.PAYOS,
		CreatedAt:     curTime,
		UpdatedAt:     curTime,
	}, ctx); err != nil {
		return res, err
	}

	// Get user data
	user, _ := p.userRepo.GetUser(req.UserId, ctx)
	var fullName string
	if user != nil {
		fullName = user.FullName
	}

	// Create ship
	p.shippingRepo.CreateShipping(entity.Shipping{
		OrderId: orderId,
		ShippingDetail: utils.ObjectToJsonString(response.ShippingDetail{
			RecipientName: fullName,
			Address:       req.Address,
			PhoneNumber:   req.PhoneNumber,
		}),
		CreatedAt: curTime,
		UpdatedAt: curTime,
	}, ctx)

	res.Url = data.CheckoutUrl

	return res, nil
}

// CreatePaymentDirect implements businesslogic.IPaymentService.
func (p *paymentService) CreatePaymentDirect(req request.CreatePaymentDirectRequest, ctx context.Context) (response.UrlAPIResponse, error) {
	var errRes error = errors.New(noti.GENERIC_ERROR_WARN_MSG)
	var res response.UrlAPIResponse
	defer closeCnn(payment_cnn)

	if !isEntityExist(p.userRepo, req.UserId, id_type, ctx) {
		return res, errRes
	}

	product, err := p.productRepo.GetProductById(req.Product.ProductId, ctx)
	if err != nil {
		return res, err
	}

	if product == nil {
		return res, errRes
	}

	inventory, err := p.invetoryRepo.GetProductInventory(product.ProductId, ctx)
	if err != nil {
		return res, err
	}

	if inventory == nil {
		return res, errRes
	}

	if inventory.CurrentQuantity < int64(req.Product.Quantity) {
		return res, errRes
	}

	var paymentId string = utils.GenerateId()
	var orderCode int = utils.GenerateNumber()
	var totalAmount float64 = product.Price * float64(req.Product.Quantity)

	// Create transaction url
	data, err := payos.CreatePaymentLink(payos.CheckoutRequestType{
		OrderCode: int64(orderCode),
		Amount:    int(totalAmount),
		Items: []payos.Item{
			{
				Name:     product.ProductName,
				Quantity: req.Product.Quantity,
				Price:    int(product.Price),
			},
		},
		Description: fmt.Sprint(orderCode),
		ReturnUrl:   os.Getenv(payment_env.PAYMENT_CALLBACK_SUCCESS_PROCESS) + paymentId,
		CancelUrl:   os.Getenv(payment_env.PAYMENT_CALLBACK_CANCEL_PROCESS) + paymentId,
	})

	if err != nil {
		p.logger.Println(fmt.Sprintf(noti.PAYMENT_GENERATE_TRANSACTION_URL_ERR_MSG, payment_method.PAYOS) + err.Error())
		return res, errors.New(noti.INTERNALL_ERR_MSG)
	}

	inventory.CurrentQuantity -= int64(req.Product.Quantity)

	// Update inventory
	if err := p.invetoryRepo.UpdateProductInventory(*inventory, ctx); err != nil {
		return res, err
	}

	var orderId string = utils.GenerateId()
	var curTime time.Time = time.Now()

	// Create order
	if err := p.orderRepo.CreateOrder(entity.Order{
		OrderId: orderId,
		UserId:  req.UserId,
		Items: utils.ObjectToJsonString([]response.CartItem{
			{
				ProductId: req.Product.ProductId,
				Name:      product.ProductName,
				ImageUrl:  product.Image,
				Quantity:  req.Product.Quantity,
				Price:     product.Price,
				Currency:  product.Currency,
			},
		}),
		TotalAmount: totalAmount,
		Currency:    product.Currency,
		Status:      domain_status.ORDER_PENDING,
		Note:        req.Note,
		CreatedAt:   curTime,
		UpdatedAt:   curTime,
	}, ctx); err != nil {
		return res, err
	}

	// Create payment
	if err := p.paymentRepo.CreatePayment(entity.Payment{
		PaymentId:     paymentId,
		UserId:        req.UserId,
		OrderId:       orderId,
		TransactionId: orderId,
		Amount:        totalAmount,
		Currency:      product.Currency,
		Status:        domain_status.PAYMENT_PENDING,
		Method:        payment_method.PAYOS,
		CreatedAt:     curTime,
		UpdatedAt:     curTime,
	}, ctx); err != nil {
		return res, err
	}

	// Get user data
	user, _ := p.userRepo.GetUser(req.UserId, ctx)
	var fullName string
	if user != nil {
		fullName = user.FullName
	}

	// Create ship
	p.shippingRepo.CreateShipping(entity.Shipping{
		OrderId: orderId,
		ShippingDetail: utils.ObjectToJsonString(response.ShippingDetail{
			RecipientName: fullName,
			Address:       req.Address,
			PhoneNumber:   req.PhoneNumber,
		}),
		CreatedAt: curTime,
		UpdatedAt: curTime,
	}, ctx)

	// Transaction URL
	res.Url = data.CheckoutUrl
	return res, nil
}

// CallbackPaymentSuccess implements businesslogic.IPaymentService.
func (p *paymentService) CallbackPaymentSuccess(id string, ctx context.Context) (string, error) {
	var payment *entity.Payment
	var order *entity.Order
	var capturedErr error
	defer closeCnn(payment_cnn)

	// Get payment
	for i := 1; i <= 3; i++ {
		payment, capturedErr = p.paymentRepo.GetPaymentById(id, ctx)
		if capturedErr == nil {
			break
		}
	}

	if payment == nil {
		return "", errors.New(noti.GENERIC_ERROR_WARN_MSG)
	}

	// Get order
	for i := 1; i <= 3; i++ {
		order, capturedErr = p.orderRepo.GetOrder(payment.OrderId, ctx)
		if capturedErr == nil {
			break
		}
	}

	if capturedErr != nil {
		return "", capturedErr
	}

	payment.Status = domain_status.PAYMENT_PAID
	order.Status = domain_status.ORDER_CONFIRMED

	var curTime time.Time = time.Now()

	// Update payment
	payment.UpdatedAt = curTime
	if err := p.paymentRepo.UpdatePayment(*payment, ctx); err != nil {
		return "", err
	}

	// Update order
	order.UpdatedAt = curTime
	if err := p.orderRepo.UpdateOrder(*order, ctx); err != nil {
		return "", err
	}

	// // Get user data
	// user, _ := p.userRepo.GetUser(payment.UserId, ctx)
	// var fullName string
	// if user != nil {
	// 	fullName = user.FullName
	// }

	// // Create ship
	// p.shippingRepo.CreateShipping(entity.Shipping{
	// 	OrderId: order.OrderId,
	// 	ShippingDetail: utils.ObjectToJsonString(response.ShippingDetail{
	// 		RecipientName: fullName,
	// 	}),
	// 	CreatedAt: curTime,
	// 	UpdatedAt: curTime,
	// }, ctx)

	if user, _ := p.userRepo.GetUser(payment.UserId, ctx); user != nil {
		utils.SendMail(request.SendMailRequest{
			Body: request.MailBody{
				Email:    user.Email,
				Subject:  noti.NOTI_PAYMENT_MAIL_SUBJECT,
				Username: user.FullName,
				OrderId:  order.OrderId,
			},
			TemplatePath: mail_const.PAYMENT_CALLBACK_SUCCESS_TEMPLATE,
			Logger:       p.logger,
		})
	}

	// return &response.ViewOrderCreatedSuccessReponse{
	// 	OrderId:     order.OrderId,
	// 	Status:      domain_status.ORDER_PROCESSING,
	// 	TotalAmount: order.TotalAmount,
	// 	Items:       utils.JsonStringToObject[[]response.CartItem](order.Items),
	// }, nil

	return os.Getenv(payment_env.PAYMENT_CALLBACK_SUCCESS) + id, nil
}

// CallbackPaymentCancel implements businesslogic.IPaymentService.
func (p *paymentService) CallbackPaymentCancel(id string, ctx context.Context) (string, error) {
	var payment *entity.Payment
	var order *entity.Order
	var capturedErr error
	defer closeCnn(payment_cnn)

	// Get payment
	for i := 1; i <= 3; i++ {
		payment, capturedErr = p.paymentRepo.GetPaymentById(id, ctx)
		if capturedErr == nil {
			break
		}
	}

	if payment == nil {
		return "", errors.New(noti.GENERIC_ERROR_WARN_MSG)
	}

	// Get order
	for i := 1; i <= 3; i++ {
		order, capturedErr = p.orderRepo.GetOrder(payment.OrderId, ctx)
		if capturedErr == nil {
			break
		}
	}

	if capturedErr != nil {
		return "", capturedErr
	}

	payment.Status = domain_status.PAYMENT_CANCELLED
	order.Status = domain_status.ORDER_CANCELLED

	var curTime time.Time = time.Now()

	// Update payment
	payment.UpdatedAt = curTime
	if err := p.paymentRepo.UpdatePayment(*payment, ctx); err != nil {
		return "", err
	}

	// Update order
	order.UpdatedAt = curTime
	if err := p.orderRepo.UpdateOrder(*order, ctx); err != nil {
		return "", err
	}

	// Refund product amount
	for _, item := range utils.JsonStringToObject[[]response.CartItem](order.Items) {
		inventory, _ := p.invetoryRepo.GetProductInventory(item.ProductId, ctx)
		if inventory != nil {
			inventory.CurrentQuantity += int64(item.Quantity)
			p.invetoryRepo.UpdateProductInventory(*inventory, ctx)
		}
	}

	p.shippingRepo.RemoveShipping(order.OrderId, ctx)

	if user, _ := p.userRepo.GetUser(payment.UserId, ctx); user != nil {
		utils.SendMail(request.SendMailRequest{
			Body: request.MailBody{
				Email:    user.Email,
				Subject:  noti.NOTI_PAYMENT_MAIL_SUBJECT,
				Username: user.FullName,
				OrderId:  order.OrderId,
			},
			TemplatePath: mail_const.PAYMENT_CALLBACK_CANCEL_TEMPLATE,
			Logger:       p.logger,
		})
	}

	return os.Getenv(payment_env.PAYMENT_CALLBACK_CANCEL) + id, nil
}
