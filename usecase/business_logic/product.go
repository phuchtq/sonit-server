package businesslogic

import (
	"context"
	"database/sql"
	"errors"
	"log"
	"math"
	"sonit_server/constant/noti"
	repo "sonit_server/data_access"
	"sonit_server/data_access/db"
	db_server "sonit_server/data_access/db_server"
	business_logic "sonit_server/interface/business_logic"
	data_access "sonit_server/interface/data_access"
	"sonit_server/model/dto/request"
	"sonit_server/model/dto/response"
	"sonit_server/model/entity"
	"sync"
	"time"

	"sonit_server/utils"
)

type productService struct {
	logger               *log.Logger
	categoryRepo         data_access.ICategoryRepo
	collectionRepo       data_access.ICollectionRepo
	productInventoryRepo data_access.IProductInventoryRepo
	productRepo          data_access.IProductRepo
}

func InitializeProductService(db *sql.DB, logger *log.Logger) business_logic.IProductService {
	return &productService{
		logger:               logger,
		categoryRepo:         repo.InitializeCategoryRepo(db, logger),
		collectionRepo:       repo.InitializeCollectionRepo(db, logger),
		productInventoryRepo: repo.InitializeProductInventoryRepo(db, logger),
		productRepo:          repo.InitializeProductRepo(db, logger),
	}
}

func GenerateProductService() (business_logic.IProductService, error) {
	var logger = utils.GetLogConfig()

	cnn, err := db.ConnectDB(logger, db_server.InitializePostgreSQL())

	if err != nil {
		return nil, err
	}

	product_cnn = cnn

	return InitializeProductService(cnn, logger), nil
}

var product_cnn *sql.DB

// GetProductsCustomerUI implements businesslogic.IProductService.
func (p *productService) GetProductsCustomerUI(req request.GetProductsCustomerUI, ctx context.Context) (response.PaginationDataResponse, error) {
	if req.Pagination.PageNumber < 1 {
		req.Pagination.PageNumber = 1
	}

	req.Pagination.FilterProp = utils.AssignFilterProperty(req.Pagination.FilterProp)
	req.Pagination.Order = utils.AssignOrder(req.Pagination.Order)

	// var null string = "NULL"
	// if req.CategoryId == "" {
	// 	req.CategoryId = null
	// }

	// if req.CollectionId == "" {
	// 	req.CollectionId = null
	// }
	defer closeCnn(product_cnn)
	data, pages, err := p.productRepo.GetProductsCustomerUI(req, ctx)

	return response.PaginationDataResponse{
		Data:       data,
		PageNumber: req.Pagination.PageNumber,
		TotalPages: pages,
	}, err
}

// ActivateProduct implements business_logic.IProductService.
func (p *productService) ActivateProduct(id string, ctx context.Context) error {
	defer closeCnn(product_cnn)
	return p.productRepo.ActivateProduct(id, ctx)
}

// CreateProduct implements business_logic.IProductService.
func (p *productService) CreateProduct(req request.CreateProductRequest, ctx context.Context) error {
	defer closeCnn(product_cnn)

	// Collection name existed
	if isEntityExist(p.collectionRepo, req.ProductName, name_type, ctx) {
		return errors.New(noti.ITEM_EXISTED_WARN_MSG)
	}

	// Validate category if existed
	if req.CategoryId != "" {
		if !isEntityExist(p.categoryRepo, req.CategoryId, id_type, ctx) {
			return errors.New(noti.ITEM_EXISTED_WARN_MSG)
		}
	}

	// Validate collection if existed
	if req.CategoryId != "" {
		if !isEntityExist(p.collectionRepo, req.CollectionId, id_type, ctx) {
			return errors.New(noti.ITEM_EXISTED_WARN_MSG)
		}
	}

	var productId string = utils.GenerateId()

	// Error for capturing in creating new product process
	var capturedErr error

	_, cancel := context.WithCancel(ctx)
	var wg sync.WaitGroup
	var mu sync.Mutex

	// Add 2 process: creating 2 entities product and product quantity
	wg.Add(2)

	var curTime time.Time = time.Now()

	// Create new product
	go func() {
		defer wg.Done()

		if err := p.productRepo.CreateProduct(entity.Product{
			ProductId:    productId,
			CategoryId:   req.CategoryId,
			CollectionId: req.CollectionId,
			ProductName:  req.ProductName,
			Description:  req.Description,
			Image:        req.Image,
			Size:         req.Size,
			Color:        req.Color,
			Price:        req.Price,
			Currency:     utils.AssignCurrency(req.Currency),
			ActiveStatus: true,
			CreatedAt:    curTime,
			UpdatedAt:    curTime,
		}, ctx); err != nil {
			mu.Lock()
			if capturedErr == nil {
				// Assigning error
				capturedErr = err
				// Make a flag to cancel other goroutines
				cancel()
			}

			mu.Unlock()
		}
	}()

	// Create product quantity
	go func() {
		defer wg.Done()

		if err := p.productInventoryRepo.CreateProductInventory(entity.ProductInventory{
			ProductId:       productId,
			CurrentQuantity: req.Amount,
		}, ctx); err != nil {
			mu.Lock()
			if capturedErr == nil {
				// Assigning error
				capturedErr = err
				// Make a flag to cancel other goroutines
				cancel()
			}

			mu.Unlock()
		}
	}()

	// Wait for 2 goroutines to finish
	wg.Wait()

	return capturedErr
}

// GetAllProducts implements business_logic.IProductService.
func (p *productService) GetAllProducts(pageNumber int, ctx context.Context) (response.PaginationDataResponse, error) {
	// Fetch data
	if pageNumber <= 0 {
		pageNumber = 1
	}

	defer closeCnn(product_cnn)
	data, pages, err := p.productRepo.GetAllProducts(pageNumber, ctx)

	return response.PaginationDataResponse{
		Data:       data,
		PageNumber: pageNumber,
		TotalPages: pages,
	}, err
}

// GetProductById implements business_logic.IProductService.
func (p *productService) GetProductById(id string, ctx context.Context) (*entity.Product, error) {
	defer closeCnn(product_cnn)
	res, err := p.productRepo.GetProductById(id, ctx)
	return res, err
}

// GetProductsByCategory implements business_logic.IProductService.
func (p *productService) GetProductsByCategory(pageNumber int, id string, ctx context.Context) (response.PaginationDataResponse, error) {
	var res response.PaginationDataResponse
	if pageNumber <= 0 {
		pageNumber = 1
	}

	var data interface{}
	var pages int
	var err error
	defer closeCnn(product_cnn)

	if id != "" {
		data, pages, err = p.productRepo.GetAllProducts(pageNumber, ctx)
	} else {
		if !isEntityExist(p.categoryRepo, id, id_type, ctx) {
			return res, errors.New(noti.GENERIC_ERROR_WARN_MSG)
		}

		data, pages, err = p.productRepo.GetProductsByCategory(pageNumber, id, ctx)
	}

	if data != nil {
		res.Data = data
		res.PageNumber = pageNumber
		res.TotalPages = pages
	}

	return res, err
}

// GetProductsByCollection implements business_logic.IProductService.
func (p *productService) GetProductsByCollection(pageNumber int, id string, ctx context.Context) (response.PaginationDataResponse, error) {
	var res response.PaginationDataResponse
	if pageNumber <= 0 {
		pageNumber = 1
	}

	var data interface{}
	var pages int
	var err error
	defer closeCnn(product_cnn)

	if id == "" {
		data, pages, err = p.productRepo.GetAllProducts(pageNumber, ctx)
	} else {
		if !isEntityExist(p.collectionRepo, id, id_type, ctx) {
			return res, errors.New(noti.GENERIC_ERROR_WARN_MSG)
		}

		data, pages, err = p.productRepo.GetProductsByCollection(pageNumber, id, ctx)
	}

	if data != nil {
		res.Data = data
		res.PageNumber = pageNumber
		res.TotalPages = pages
	}

	return res, err
}

// GetProductsByDescription implements business_logic.IProductService.
func (p *productService) GetProductsByDescription(pageNumber int, description string, ctx context.Context) (response.PaginationDataResponse, error) {
	var res response.PaginationDataResponse
	if pageNumber <= 0 {
		pageNumber = 1
	}

	var data interface{}
	var pages int
	var err error
	defer closeCnn(product_cnn)

	if description != "" {
		data, pages, err = p.productRepo.GetAllProducts(pageNumber, ctx)
	} else {
		data, pages, err = p.productRepo.GetProductsByDescription(pageNumber, description, ctx)
	}

	if data != nil {
		res.Data = data
		res.PageNumber = pageNumber
		res.TotalPages = pages
	}

	return res, err
}

// GetProductsByKeyword implements business_logic.IProductService.
func (p *productService) GetProductsByKeyword(pageNumber int, keyword string, ctx context.Context) (response.PaginationDataResponse, error) {
	var res response.PaginationDataResponse
	if pageNumber <= 0 {
		pageNumber = 1
	}

	var data interface{}
	var pages int
	var err error
	defer closeCnn(product_cnn)

	if keyword != "" {
		data, pages, err = p.productRepo.GetAllProducts(pageNumber, ctx)
	} else {
		data, pages, err = p.productRepo.GetProductsByKeyword(pageNumber, keyword, ctx)
	}

	if data != nil {
		res.Data = data
		res.PageNumber = pageNumber
		res.TotalPages = pages
	}

	return res, err
}

// GetProductsByName implements business_logic.IProductService.
func (p *productService) GetProductsByName(pageNumber int, name string, ctx context.Context) (response.PaginationDataResponse, error) {
	var res response.PaginationDataResponse
	if pageNumber <= 0 {
		pageNumber = 1
	}

	var data interface{}
	var pages int
	var err error
	defer closeCnn(product_cnn)

	if name != "" {
		data, pages, err = p.productRepo.GetAllProducts(pageNumber, ctx)
	} else {
		data, pages, err = p.productRepo.GetProductsByKeyword(pageNumber, name, ctx)
	}

	if data != nil {
		res.Data = data
		res.PageNumber = pageNumber
		res.TotalPages = pages
	}

	return res, err
}

// GetProductsByPriceInterval implements business_logic.IProductService.
func (p *productService) GetProductsByPriceInterval(pageNumber int, maxPrice int64, minPrice int64, ctx context.Context) (response.PaginationDataResponse, error) {
	if pageNumber <= 0 {
		pageNumber = 1
	}

	if minPrice < 0 {
		minPrice = 0
	}

	if maxPrice <= 0 {
		maxPrice = math.MaxInt64
	}

	var res response.PaginationDataResponse
	defer closeCnn(product_cnn)

	data, pages, err := p.productRepo.GetProductsByPriceInterval(pageNumber, maxPrice, minPrice, ctx)
	if data != nil {
		res.Data = data
		res.PageNumber = pageNumber
		res.TotalPages = pages
	}

	return res, err
}

// GetProductsByStatus implements business_logic.IProductService.
func (p *productService) GetProductsByStatus(pageNumber int, status *bool, ctx context.Context) (response.PaginationDataResponse, error) {
	var res response.PaginationDataResponse
	if pageNumber <= 0 {
		pageNumber = 1
	}

	var data interface{}
	var pages int
	var err error
	defer closeCnn(product_cnn)

	if status == nil {
		data, pages, err = p.productRepo.GetAllProducts(pageNumber, ctx)
	} else {
		data, pages, err = p.productRepo.GetProductsByStatus(pageNumber, *status, ctx)
	}

	if data != nil {
		res.Data = data
		res.PageNumber = pageNumber
		res.TotalPages = pages
	}

	return res, err
}

// RemoveProduct implements business_logic.IProductService.
func (p *productService) RemoveProduct(id string, ctx context.Context) error {
	defer closeCnn(product_cnn)
	return p.productRepo.RemoveProduct(id, ctx)
}

// UpdateProduct implements business_logic.IProductService.
func (p *productService) UpdateProduct(req request.UpdateProductRequest, ctx context.Context) error {
	defer closeCnn(product_cnn)
	product, err := p.productRepo.GetProductById(req.ProductId, ctx)

	if err != nil {
		return err
	}

	if req.Amount != nil {
		go func() {
			inventory, _ := p.productInventoryRepo.GetProductInventory(req.ProductId, ctx)
			inventory.CurrentQuantity = *req.Amount
			p.productInventoryRepo.UpdateProductInventory(*inventory, ctx)
		}()
	}

	// Verify category if not empty
	if req.CategoryId != "" {
		if !isEntityExist(p.categoryRepo, req.CategoryId, id_type, ctx) {
			return errors.New(noti.GENERIC_ERROR_WARN_MSG)
		}

		product.CategoryId = req.CategoryId
	}

	// Verify collection if not empty
	if req.CollectionId != "" {
		if !isEntityExist(p.collectionRepo, req.CollectionId, id_type, ctx) {
			return errors.New(noti.GENERIC_ERROR_WARN_MSG)
		}

		product.CollectionId = req.CollectionId
	}

	if req.ProductName != "" {
		product.ProductName = req.ProductName
	}

	if req.Description != "" {
		product.Description = req.Description
	}

	if req.Image != "" {
		product.Image = req.Image
	}

	if req.Color != "" {
		product.Color = req.Color
	}

	if req.Price > 0 {
		product.Price = req.Price
	}

	if req.Currency != "" {
		product.Currency = utils.AssignCurrency(req.Currency)
	}

	product.UpdatedAt = time.Now()

	return p.productRepo.UpdateProduct(*product, ctx)
}
