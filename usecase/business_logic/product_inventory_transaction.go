package businesslogic

import (
	"context"
	"database/sql"
	"errors"
	"log"
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

type productInventoryTransactionService struct {
	logger                 *log.Logger
	productRepo            data_access.IProductRepo
	productInventoryRepo   data_access.IProductInventoryRepo
	productInventoryTxRepo data_access.IProductInventoryTransactionRepo
}

func InitializeProductInventoryTransactionService(db *sql.DB, logger *log.Logger) business_logic.IProductInventoryTransactionService {
	return &productInventoryTransactionService{
		logger:                 logger,
		productRepo:            repo.InitializeProductRepo(db, logger),
		productInventoryRepo:   repo.InitializeProductInventoryRepo(db, logger),
		productInventoryTxRepo: repo.InitializeProductInventoryTransactionRepo(db, logger),
	}
}

func GenerateProductInventoryTransactionService() (business_logic.IProductInventoryTransactionService, error) {
	var logger = utils.GetLogConfig()

	cnn, err := db.ConnectDB(logger, db_server.InitializePostgreSQL())

	if err != nil {
		return nil, err
	}

	inventoryTx_cnn = cnn

	return InitializeProductInventoryTransactionService(cnn, logger), nil
}

var inventoryTx_cnn *sql.DB

// CreateProductInventoryTransaction implements businesslogic.IProductInventoryTransactionService.
func (p *productInventoryTransactionService) CreateProductInventoryTransaction(req request.CreateProductInventoryTransactionRequest, ctx context.Context) error {
	var errRes error = errors.New(noti.GENERIC_ERROR_WARN_MSG)
	defer closeCnn(inventoryTx_cnn)

	if !isEntityExist(p.productRepo, req.ProductId, id_type, ctx) {
		return errRes
	}

	if !isInventoryActionValid(req.Action) {
		return errRes
	}

	var curTime time.Time = time.Now()
	return p.productInventoryTxRepo.CreateProductInventoryTransaction(entity.ProductInventoryTransaction{
		TransactionId: utils.GenerateId(),
		ProductId:     req.ProductId,
		Amount:        req.Amount,
		Action:        req.Action,
		Note:          req.Note,
		Date:          req.Date,
		CreatedAt:     curTime,
		UpdatedAt:     curTime,
	}, ctx)
}

// GetAllProductInventoryTransactions implements businesslogic.IProductInventoryTransactionService.
func (p *productInventoryTransactionService) GetAllProductInventoryTransactions(req request.GetProductInventoryTrasactionsRequest, ctx context.Context) (response.PaginationDataResponse, error) {
	panic("unimplemented")
}

// GetInventoryTransactionsByProduct implements businesslogic.IProductInventoryTransactionService.
func (p *productInventoryTransactionService) GetInventoryTransactionsByProduct(req request.GetProductInventoryTrasactionsByProductRequest, ctx context.Context) (response.PaginationDataResponse, error) {
	if req.InventoryTransaction.Pagination.PageNumber < 1 {
		req.InventoryTransaction.Pagination.PageNumber = 1
	}

	req.InventoryTransaction.Pagination.FilterProp = utils.AssignFilterProperty(req.InventoryTransaction.Pagination.FilterProp)
	req.InventoryTransaction.Pagination.Order = utils.AssignOrder(req.InventoryTransaction.Pagination.Order)

	defer closeCnn(inventoryTx_cnn)

	data, pages, err := p.productInventoryTxRepo.GetInventoryTransactionsByProduct(req, ctx)

	return response.PaginationDataResponse{
		Data:       data,
		PageNumber: req.InventoryTransaction.Pagination.PageNumber,
		TotalPages: pages,
	}, err
}

// GetProductInventoryTransaction implements businesslogic.IProductInventoryTransactionService.
func (p *productInventoryTransactionService) GetProductInventoryTransaction(id string, ctx context.Context) (*entity.ProductInventoryTransaction, error) {
	defer closeCnn(inventoryTx_cnn)
	return p.productInventoryTxRepo.GetProductInventoryTransaction(id, ctx)
}

// UpdateProductInventoryTransaction implements businesslogic.IProductInventoryTransactionService.
func (p *productInventoryTransactionService) UpdateProductInventoryTransaction(req request.UpdateProductInventoryTransactionRequest, ctx context.Context) error {
	var errRes error = errors.New(noti.GENERIC_ERROR_WARN_MSG)
	defer closeCnn(inventoryTx_cnn)

	// Verify transaction
	if !isEntityExist(p.productInventoryTxRepo, req.TransactionId, id_type, ctx) {
		return errRes
	}

	// Verify product
	if !isEntityExist(p.productRepo, req.UpdateData.ProductId, id_type, ctx) {
		return errRes
	}

	// Retrieve original transaction data
	originalTx, err := p.productInventoryTxRepo.GetProductInventoryTransaction(req.TransactionId, ctx)
	if err != nil {
		return err
	}

	// Retrieve product inventory
	inventory, err := p.productInventoryRepo.GetProductInventory(originalTx.ProductId, ctx)
	if err != nil {
		return err
	}

	// Reverse previous action which affected to the inventory
	inventory.CurrentQuantity += inverseInventoryTransaction(originalTx.Action, originalTx.Amount)

	originalTx.Action = req.UpdateData.Action
	originalTx.Amount = req.UpdateData.Amount

	if req.UpdateData.Note != "" {
		originalTx.Note = req.UpdateData.Note
	}

	var updatedAmount int64 = executeInventoryTransaction(req.UpdateData.Action, req.UpdateData.Amount)

	// Retrieve new affected product inventory if changed
	if req.UpdateData.ProductId != originalTx.ProductId {
		newAfftedInventory, err := p.productInventoryRepo.GetProductInventory(req.UpdateData.ProductId, ctx)
		if err != nil {
			return err
		}

		newAfftedInventory.CurrentQuantity += updatedAmount
		if err := p.productInventoryRepo.UpdateProductInventory(*newAfftedInventory, ctx); err != nil {
			return err
		}
	} else {
		inventory.CurrentQuantity += updatedAmount
	}

	var capturedErr error

	_, cancel := context.WithCancel(ctx)
	var wg sync.WaitGroup
	var mu sync.Mutex

	// Add 2 process: updating inventory and transaction
	wg.Add(2)

	originalTx.UpdatedAt = time.Now()

	// Update inventory
	go func() {
		defer wg.Done()

		if err := p.productInventoryRepo.UpdateProductInventory(*inventory, ctx); err != nil {
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

	// Update inventory transaction
	go func() {
		defer wg.Done()

		if err := p.productInventoryTxRepo.UpdateProductInventoryTransaction(*originalTx, ctx); err != nil {
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
