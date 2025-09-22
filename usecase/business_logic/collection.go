package businesslogic

import (
	"context"
	"database/sql"
	"errors"
	"log"
	"sonit_server/constant/noti"
	repo "sonit_server/data_access" // Collection data access ~~ Collection repository
	"sonit_server/data_access/db"
	db_server "sonit_server/data_access/db_server"
	business_logic "sonit_server/interface/business_logic"
	data_access "sonit_server/interface/data_access" // Collection repo interface
	"sonit_server/model/dto/request"
	entity "sonit_server/model/entity"
	"sonit_server/utils"
	"time"
)

type collectionService struct {
	collectionRepo data_access.ICollectionRepo
	logger         *log.Logger
}

func GenerateCollectionService() (business_logic.ICollectionService, error) {
	var logger = utils.GetLogConfig()

	cnn, err := db.ConnectDB(logger, db_server.InitializePostgreSQL())

	if err != nil {
		return nil, err
	}

	collection_cnn = cnn

	return InitializeCollectionService(cnn, logger), nil
}

func InitializeCollectionService(db *sql.DB, logger *log.Logger) business_logic.ICollectionService {
	return &collectionService{
		collectionRepo: repo.InitializeCollectionRepo(db, logger),
		logger:         logger,
	}
}

var collection_cnn *sql.DB

// ActivateCollection implements business_logic.ICollectionService.
func (c *collectionService) ActivateCollection(id string, ctx context.Context) error {
	defer closeCnn(collection_cnn)
	return c.collectionRepo.ActivateCollection(id, ctx)
}

// CreateCollection implements business_logic.ICollectionService.
func (c *collectionService) CreateCollection(req request.CreateCollectionRequest, ctx context.Context) error {
	defer closeCnn(collection_cnn)

	// Collection name existed
	if isEntityExist(c.collectionRepo, req.CollectionName, name_type, ctx) {
		return errors.New(noti.ITEM_EXISTED_WARN_MSG)
	}
	//---------------------------------------
	var curTime time.Time = time.Now()
	return c.collectionRepo.CreateCollection(entity.Collection{
		CollectionId:   utils.GenerateId(),
		CollectionName: req.CollectionName,
		Description:    req.Description,
		ActiveStatus:   true,
		CreatedAt:      curTime,
		UpdatedAt:      curTime,
	}, ctx)
}

// GetAllCollections implements business_logic.ICollectionService.
func (c *collectionService) GetAllCollections(ctx context.Context) (*[]entity.Collection, error) {
	defer closeCnn(collection_cnn)
	return c.collectionRepo.GetAllCollections(ctx)
}

// GetCollectionById implements business_logic.ICollectionService.
func (c *collectionService) GetCollectionById(id string, ctx context.Context) (*entity.Collection, error) {
	defer closeCnn(collection_cnn)

	if id == "" {
		return nil, errors.New(noti.GENERIC_ERROR_WARN_MSG)
	}
	//---------------------------------------
	res, err := c.collectionRepo.GetCollectionById(id, ctx)
	if err != nil {
		return nil, err
	}
	if res == nil {
		return nil, errors.New(entity.GetCollectionTable() + noti.UNDEFINED_OBJECT_WARN_MSG)
	}
	//---------------------------------------
	return res, nil
}

// GetCollectionsByName implements business_logic.ICollectionService.
func (c *collectionService) GetCollectionsByName(name string, ctx context.Context) (*[]entity.Collection, error) {
	defer closeCnn(collection_cnn)

	if name == "" {
		return c.collectionRepo.GetAllCollections(ctx)
	}

	return c.collectionRepo.GetCollectionsByName(utils.ToNormalizedString(name), ctx)
}

// GetCollectionsByStatus implements business_logic.ICollectionService.
func (c *collectionService) GetCollectionsByStatus(status *bool, ctx context.Context) (*[]entity.Collection, error) {
	defer closeCnn(collection_cnn)

	if status == nil {
		return c.collectionRepo.GetAllCollections(ctx)
	}

	return c.collectionRepo.GetCollectionsByStatus(*status, ctx)
}

// RemoveCollection implements business_logic.ICollectionService.
func (c *collectionService) RemoveCollection(id string, ctx context.Context) error {
	defer closeCnn(collection_cnn)

	if id == "" {
		return errors.New(noti.GENERIC_ERROR_WARN_MSG)
	}

	return c.collectionRepo.RemoveCollection(id, ctx)
}

// UpdateCollection implements business_logic.ICollectionService.
func (c *collectionService) UpdateCollection(req request.UpdateCollectionRequest, ctx context.Context) error {
	defer closeCnn(collection_cnn)

	collection, err := c.collectionRepo.GetCollectionById(req.CollectionId, ctx)
	if err != nil {
		return err
	}

	if req.CollectionName != "" && !isEntityExist(c.collectionRepo, req.CollectionName, name_type, ctx) {
		collection.CollectionName = req.CollectionName
	}

	if req.Description != "" {
		collection.Description = req.Description
	}
	collection.UpdatedAt = time.Now()

	return c.collectionRepo.UpdateCollection(*collection, ctx)
}
