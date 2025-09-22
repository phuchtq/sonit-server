package businesslogic

import (
	"context"
	"database/sql"
	"errors"
	"log"
	"sonit_server/constant/noti"
	repo "sonit_server/data_access" // Category data access ~~ Category repository
	"sonit_server/data_access/db"
	db_server "sonit_server/data_access/db_server"
	business_logic "sonit_server/interface/business_logic"
	data_access "sonit_server/interface/data_access" // Category repo interface
	"sonit_server/model/dto/request"
	entity "sonit_server/model/entity"
	"sonit_server/utils"
	"time"
)

type categoryservice struct {
	categoryRepo data_access.ICategoryRepo
	logger       *log.Logger
}

func GenerateCategoryService() (business_logic.ICategoryService, error) {
	var logger = utils.GetLogConfig()

	cnn, err := db.ConnectDB(logger, db_server.InitializePostgreSQL())

	if err != nil {
		return nil, err
	}

	category_cnn = cnn

	return InitializeCategoryService(cnn, logger), nil
}

func InitializeCategoryService(db *sql.DB, logger *log.Logger) business_logic.ICategoryService {
	return &categoryservice{
		categoryRepo: repo.InitializeCategoryRepo(db, logger),
		logger:       logger,
	}
}

var category_cnn *sql.DB

// ActivateCategory implements business_logic.ICategorieservice.
func (c *categoryservice) ActivateCategory(id string, ctx context.Context) error {
	defer closeCnn(category_cnn)
	return c.categoryRepo.ActivateCategory(id, ctx)
}

// CreateCategory implements business_logic.ICategorieservice.
func (c *categoryservice) CreateCategory(req request.CreateCategoryRequest, ctx context.Context) error {
	defer closeCnn(category_cnn)

	// Category name existed
	if isEntityExist(c.categoryRepo, req.CategoryName, name_type, ctx) {
		return errors.New(noti.ITEM_EXISTED_WARN_MSG)
	}
	//---------------------------------------
	var curTime time.Time = time.Now()
	return c.categoryRepo.CreateCategory(entity.Category{
		CategoryId:   utils.GenerateId(),
		CategoryName: req.CategoryName,
		Description:  req.Description,
		ActiveStatus: true,
		CreatedAt:    curTime,
		UpdatedAt:    curTime,
	}, ctx)
}

// GetAllCategories implements business_logic.ICategorieservice.
func (c *categoryservice) GetAllCategories(ctx context.Context) (*[]entity.Category, error) {
	defer closeCnn(category_cnn)
	return c.categoryRepo.GetAllCategories(ctx)
}

// GetCategoryById implements business_logic.ICategorieservice.
func (c *categoryservice) GetCategoryById(id string, ctx context.Context) (*entity.Category, error) {
	defer closeCnn(category_cnn)

	if id == "" {
		return nil, errors.New(noti.GENERIC_ERROR_WARN_MSG)
	}
	//---------------------------------------
	res, err := c.categoryRepo.GetCategoryById(id, ctx)
	if err != nil {
		return nil, err
	}
	if res == nil {
		return nil, errors.New(entity.GetCategoryTable() + noti.UNDEFINED_OBJECT_WARN_MSG)
	}
	//---------------------------------------
	return res, nil
}

// GetCategoriesByName implements business_logic.ICategorieservice.
func (c *categoryservice) GetCategoriesByName(name string, ctx context.Context) (*[]entity.Category, error) {
	defer closeCnn(category_cnn)

	if name == "" {
		return c.categoryRepo.GetAllCategories(ctx)
	}

	return c.categoryRepo.GetCategoriesByName(utils.ToNormalizedString(name), ctx)
}

// GetCategoriesByStatus implements business_logic.ICategorieservice.
func (c *categoryservice) GetCategoriesByStatus(status *bool, ctx context.Context) (*[]entity.Category, error) {
	defer closeCnn(category_cnn)

	if status == nil {
		return c.categoryRepo.GetAllCategories(ctx)
	}

	return c.categoryRepo.GetCategoriesByStatus(*status, ctx)
}

// RemoveCategory implements business_logic.ICategorieservice.
func (c *categoryservice) RemoveCategory(id string, ctx context.Context) error {
	defer closeCnn(category_cnn)

	if id == "" {
		return errors.New(noti.GENERIC_ERROR_WARN_MSG)
	}

	return c.categoryRepo.RemoveCategory(id, ctx)
}

// UpdateCategory implements business_logic.ICategorieservice.
func (c *categoryservice) UpdateCategory(req request.UpdateCategoryRequest, ctx context.Context) error {
	defer closeCnn(category_cnn)

	category, err := c.categoryRepo.GetCategoryById(req.CategoryId, ctx)
	if err != nil {
		return err
	}

	if req.CategoryName != "" && !isEntityExist(c.categoryRepo, req.CategoryName, name_type, ctx) {
		category.CategoryName = req.CategoryName
	}

	if req.Description != "" {
		category.Description = req.Description
	}
	category.UpdatedAt = time.Now()

	return c.categoryRepo.UpdateCategory(*category, ctx)
}
