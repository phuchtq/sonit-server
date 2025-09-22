package handler

import (
	"log"
	action_type "sonit_server/constant/action_type"
	"sonit_server/model/dto/request"
	"sonit_server/model/dto/response"
	businesslogic "sonit_server/usecase/business_logic"
	"sonit_server/utils"
	"strconv"

	"github.com/gin-gonic/gin"
)

// GetAllProducts godoc
// @Summary      Get all products
// @Description  Retrieve paginated list of all products
// @Tags         products
// @Security     BearerAuth
// @Produce      json
// @Param        page query int false "Page number"
// @Success      200 {object} response.PaginationDataResponse
// @Failure 401 {object} response.MessageAPIResponse "You have no rights to access this action."
// @Failure 400 {object} response.MessageAPIResponse "Invalid data. Please try again."
// @Failure 500 {object} response.MessageAPIResponse "There is something wrong in the system during the process. Please try again later."
// @Router       /products [get]
func GetAllProducts(ctx *gin.Context) {
	service, err := businesslogic.GenerateProductService()
	if err != nil {
		utils.ProcessResponse(utils.GenerateInvalidRequestAndSystemProblemModel(ctx, err))
		return
	}

	pageNumber, _ := strconv.Atoi(ctx.Query("page"))

	res, err := service.GetAllProducts(pageNumber, ctx)

	utils.ProcessResponse(response.APIResponse{
		Data1:    res,
		Data2:    res,
		ErrMsg:   err,
		Context:  ctx,
		PostType: action_type.NON_POST,
	})
}

// GetProductsCustomerUI godoc
// @Summary      Get products for customer UI
// @Description  Retrieve filtered products for customer view
// @Tags         products
// @Accept       json
// @Produce      json
// @Param        page_number   query int false "Page number"
// @Param        keyword       query string false "Search keyword"
// @Param        filter_prop   query string false "Filter property (e.g. date, price)"
// @Param        order         query string false "Sort order (ASC or DESC)"
// @Param        category_id   query string false "Category ID"
// @Param        collection_id query string false "Collection ID"
// @Success      200 {object} response.PaginationDataResponse
// @Failure 400 {object} response.MessageAPIResponse "Invalid data. Please try again."
// @Failure 500 {object} response.MessageAPIResponse "There is something wrong in the system during the process. Please try again later."
// @Router       /products/customer-ui [get]
func GetProductsCustomerUI(ctx *gin.Context) {
	var request request.GetProductsCustomerUI
	if ctx.ShouldBindQuery(&request) != nil {
		utils.ProcessResponse(utils.GenerateInvalidRequestAndSystemProblemModel(ctx, nil))
		return
	}

	log.Println("Request: ", request)

	service, err := businesslogic.GenerateProductService()
	if err != nil {
		utils.ProcessResponse(utils.GenerateInvalidRequestAndSystemProblemModel(ctx, err))
		return
	}

	res, err := service.GetProductsCustomerUI(request, ctx)

	utils.ProcessResponse(response.APIResponse{
		Data1:    res,
		Data2:    res,
		ErrMsg:   err,
		Context:  ctx,
		PostType: action_type.NON_POST,
	})
}

// GetProductsByCategory godoc
// @Summary      Get products by category
// @Description  Retrieve products under a specific category
// @Tags         products
// @Produce      json
// @Security BearerAuth
// @Param        pageNumber query int false "Page number"
// @Param        id path string true "Category ID"
// @Success      200 {object} response.PaginationDataResponse
// @Failure 401 {object} response.MessageAPIResponse "You have no rights to access this action."
// @Failure 400 {object} response.MessageAPIResponse "Invalid data. Please try again."
// @Failure 500 {object} response.MessageAPIResponse "There is something wrong in the system during the process. Please try again later."
// @Router       /products/category/{id} [get]
func GetProductsByCategory(ctx *gin.Context) {
	service, err := businesslogic.GenerateProductService()
	if err != nil {
		utils.ProcessResponse(utils.GenerateInvalidRequestAndSystemProblemModel(ctx, err))
		return
	}

	pageNumber, _ := strconv.Atoi(ctx.Query("pageNumber"))

	res, err := service.GetProductsByCategory(pageNumber, ctx.Param("id"), ctx)

	utils.ProcessResponse(response.APIResponse{
		Data1:    res,
		Data2:    res,
		ErrMsg:   err,
		Context:  ctx,
		PostType: action_type.NON_POST,
	})
}

// GetProductsByCollection godoc
// @Summary      Get products by collection
// @Description  Retrieve products under a specific category
// @Tags         products
// @Produce      json
// @Security BearerAuth
// @Param        pageNumber query int false "Page number"
// @Param        id path string true "Collection ID"
// @Success      200 {object} response.PaginationDataResponse
// @Failure 401 {object} response.MessageAPIResponse "You have no rights to access this action."
// @Failure 400 {object} response.MessageAPIResponse "Invalid data. Please try again."
// @Failure 500 {object} response.MessageAPIResponse "There is something wrong in the system during the process. Please try again later."
// @Router       /products/collection/{id} [get]
func GetProductsByCollection(ctx *gin.Context) {
	service, err := businesslogic.GenerateProductService()
	if err != nil {
		utils.ProcessResponse(utils.GenerateInvalidRequestAndSystemProblemModel(ctx, err))
		return
	}

	pageNumber, _ := strconv.Atoi(ctx.Query("pageNumber"))

	res, err := service.GetProductsByCollection(pageNumber, ctx.Param("id"), ctx)

	utils.ProcessResponse(response.APIResponse{
		Data1:    res,
		Data2:    res,
		ErrMsg:   err,
		Context:  ctx,
		PostType: action_type.NON_POST,
	})
}

// GetProductsByPriceInterval godoc
// @Summary      Get products by price range
// @Description  Retrieve products within a specified price range
// @Tags         products
// @Produce      json
// @Security BearerAuth
// @Param        pageNumber query int false "Page number"
// @Param        minPrice query int false "Minimum price"
// @Param        maxPrice query int false "Maximum price"
// @Success      200 {object} response.PaginationDataResponse
// @Failure 401 {object} response.MessageAPIResponse "You have no rights to access this action."
// @Failure 400 {object} response.MessageAPIResponse "Invalid data. Please try again."
// @Failure 500 {object} response.MessageAPIResponse "There is something wrong in the system during the process. Please try again later."
// @Router       /products/price-interval [get]
func GetProductsByPriceInterval(ctx *gin.Context) {
	service, err := businesslogic.GenerateProductService()
	if err != nil {
		utils.ProcessResponse(utils.GenerateInvalidRequestAndSystemProblemModel(ctx, err))
		return
	}

	pageNumber, _ := strconv.Atoi(ctx.Query("pageNumber"))
	minPrice, _ := strconv.ParseInt(ctx.Query("minPrice"), 10, 64)
	maxPrice, _ := strconv.ParseInt(ctx.Query("maxPrice"), 10, 64)

	res, err := service.GetProductsByPriceInterval(pageNumber, maxPrice, minPrice, ctx)

	utils.ProcessResponse(response.APIResponse{
		Data1:    res,
		Data2:    res,
		ErrMsg:   err,
		Context:  ctx,
		PostType: action_type.NON_POST,
	})
}

// GetProductsByName godoc
// @Summary      Get products by name
// @Description  Retrieve products matching the name
// @Tags         products
// @Security     BearerAuth
// @Produce      json
// @Security BearerAuth
// @Param        pageNumber query int false "Page number"
// @Param        name path string true "Product name"
// @Success      200 {object} response.PaginationDataResponse
// @Failure 401 {object} response.MessageAPIResponse "You have no rights to access this action."
// @Failure 400 {object} response.MessageAPIResponse "Invalid data. Please try again."
// @Failure 500 {object} response.MessageAPIResponse "There is something wrong in the system during the process. Please try again later."
// @Router       /products/name/{name} [get]
func GetProductsByName(ctx *gin.Context) {
	service, err := businesslogic.GenerateProductService()
	if err != nil {
		utils.ProcessResponse(utils.GenerateInvalidRequestAndSystemProblemModel(ctx, err))
		return
	}

	pageNumber, _ := strconv.Atoi(ctx.Query("pageNumber"))

	res, err := service.GetProductsByCategory(pageNumber, ctx.Param("name"), ctx)

	utils.ProcessResponse(response.APIResponse{
		Data1:    res,
		Data2:    res,
		ErrMsg:   err,
		Context:  ctx,
		PostType: action_type.NON_POST,
	})
}

// GetProductsByDescription godoc
// @Summary      Get products by description
// @Description  Retrieve products matching the description
// @Tags         products
// @Security     BearerAuth
// @Produce      json
// @Security BearerAuth
// @Param        pageNumber query int false "Page number"
// @Param        description path string true "Product description"
// @Success      200 {object} response.PaginationDataResponse
// @Failure 401 {object} response.MessageAPIResponse "You have no rights to access this action."
// @Failure 400 {object} response.MessageAPIResponse "Invalid data. Please try again."
// @Failure 500 {object} response.MessageAPIResponse "There is something wrong in the system during the process. Please try again later."
// @Router       /products/description/{description} [get]
func GetProductsByDescription(ctx *gin.Context) {
	service, err := businesslogic.GenerateProductService()
	if err != nil {
		utils.ProcessResponse(utils.GenerateInvalidRequestAndSystemProblemModel(ctx, err))
		return
	}

	pageNumber, _ := strconv.Atoi(ctx.Query("pageNumber"))

	res, err := service.GetProductsByCategory(pageNumber, ctx.Param("description"), ctx)

	utils.ProcessResponse(response.APIResponse{
		Data1:    res,
		Data2:    res,
		ErrMsg:   err,
		Context:  ctx,
		PostType: action_type.NON_POST,
	})
}

// GetProductsByStatus godoc
// @Summary      Get products by status
// @Description  Retrieve products filtered by active status
// @Tags         products
// @Security     BearerAuth
// @Produce      json
// @Security BearerAuth
// @Param        pageNumber query int false "Page number"
// @Param        status path bool true "Product status"
// @Success      200 {object} response.PaginationDataResponse
// @Failure 401 {object} response.MessageAPIResponse "You have no rights to access this action."
// @Failure 400 {object} response.MessageAPIResponse "Invalid data. Please try again."
// @Failure 500 {object} response.MessageAPIResponse "There is something wrong in the system during the process. Please try again later."
// @Router       /products/status/{status} [get]
func GetProductsByStatus(ctx *gin.Context) {
	service, err := businesslogic.GenerateProductService()
	if err != nil {
		utils.ProcessResponse(utils.GenerateInvalidRequestAndSystemProblemModel(ctx, err))
		return
	}

	pageNumber, _ := strconv.Atoi(ctx.Query("pageNumber"))

	var pStatus *bool
	status, err := strconv.ParseBool(ctx.Param("status"))
	if err == nil {
		pStatus = &status
	}

	res, err := service.GetProductsByStatus(pageNumber, pStatus, ctx)

	utils.ProcessResponse(response.APIResponse{
		Data1:    res,
		Data2:    res,
		ErrMsg:   err,
		Context:  ctx,
		PostType: action_type.NON_POST,
	})
}

func GetProductsByKeyword(ctx *gin.Context) {
	service, err := businesslogic.GenerateProductService()
	if err != nil {
		utils.ProcessResponse(utils.GenerateInvalidRequestAndSystemProblemModel(ctx, err))
		return
	}

	pageNumber, _ := strconv.Atoi(ctx.Query("page"))

	res, err := service.GetProductsByCategory(pageNumber, ctx.Query("keyword"), ctx)

	utils.ProcessResponse(response.APIResponse{
		Data1:    res,
		Data2:    res,
		ErrMsg:   err,
		Context:  ctx,
		PostType: action_type.NON_POST,
	})
}

// GetProductById godoc
// @Summary      Get product by ID
// @Description  Retrieve a product by its ID
// @Tags         products
// @Produce      json
// @Param        id path string true "Product ID"
// @Success      200 {object} entity.Product
// @Failure 401 {object} response.MessageAPIResponse "You have no rights to access this action."
// @Failure 400 {object} response.MessageAPIResponse "Invalid data. Please try again."
// @Failure 500 {object} response.MessageAPIResponse "There is something wrong in the system during the process. Please try again later."
// @Router       /products/{id} [get]
func GetProductById(ctx *gin.Context) {
	service, err := businesslogic.GenerateProductService()
	if err != nil {
		utils.ProcessResponse(utils.GenerateInvalidRequestAndSystemProblemModel(ctx, err))
		return
	}

	res, err := service.GetProductById(ctx.Param("id"), ctx)

	utils.ProcessResponse(response.APIResponse{
		Data1:    res,
		Data2:    res,
		ErrMsg:   err,
		Context:  ctx,
		PostType: action_type.NON_POST,
	})
}

// CreateProduct godoc
// @Summary      Create a product
// @Description  Add a new product
// @Tags         products
// @Security     BearerAuth
// @Accept       json
// @Produce      json
// @Param        request body request.CreateProductRequest true "Product creation request"
// @Success      200 {object} response.MessageAPIResponse "success"
// @Failure 401 {object} response.MessageAPIResponse "You have no rights to access this action."
// @Failure 400 {object} response.MessageAPIResponse "Invalid data. Please try again."
// @Failure 500 {object} response.MessageAPIResponse "There is something wrong in the system during the process. Please try again later."
// @Router       /products/create [post]
func CreateProduct(ctx *gin.Context) {
	var request request.CreateProductRequest
	if ctx.ShouldBindJSON(&request) != nil {
		utils.ProcessResponse(utils.GenerateInvalidRequestAndSystemProblemModel(ctx, nil))
		return
	}

	service, err := businesslogic.GenerateProductService()
	if err != nil {
		utils.ProcessResponse(utils.GenerateInvalidRequestAndSystemProblemModel(ctx, err))
		return
	}

	utils.ProcessResponse(response.APIResponse{
		ErrMsg:   service.CreateProduct(request, ctx),
		Context:  ctx,
		PostType: action_type.INFORM,
	})
}

// UpdateProduct godoc
// @Summary      Update a product
// @Description  Modify an existing product
// @Tags         products
// @Security     BearerAuth
// @Accept       json
// @Produce      json
// @Param        request body request.UpdateProductRequest true "Product update request"
// @Success      200 {object} response.MessageAPIResponse "success"
// @Failure 401 {object} response.MessageAPIResponse "You have no rights to access this action."
// @Failure 400 {object} response.MessageAPIResponse "Invalid data. Please try again."
// @Failure 500 {object} response.MessageAPIResponse "There is something wrong in the system during the process. Please try again later."
// @Router       /products/update [put]
func UpdateProduct(ctx *gin.Context) {
	var request request.UpdateProductRequest
	if ctx.ShouldBindJSON(&request) != nil {
		utils.ProcessResponse(utils.GenerateInvalidRequestAndSystemProblemModel(ctx, nil))
		return
	}

	service, err := businesslogic.GenerateProductService()
	if err != nil {
		utils.ProcessResponse(utils.GenerateInvalidRequestAndSystemProblemModel(ctx, err))
		return
	}

	utils.ProcessResponse(response.APIResponse{
		ErrMsg:   service.UpdateProduct(request, ctx),
		Context:  ctx,
		PostType: action_type.INFORM,
	})
}

// RemoveProduct godoc
// @Summary      Remove a product
// @Description  Delete a product by ID
// @Tags         products
// @Security     BearerAuth
// @Produce      json
// @Param        id path string true "Product ID"
// @Success      200 {object} response.MessageAPIResponse "success"
// @Failure 401 {object} response.MessageAPIResponse "You have no rights to access this action."
// @Failure 400 {object} response.MessageAPIResponse "Invalid data. Please try again."
// @Failure 500 {object} response.MessageAPIResponse "There is something wrong in the system during the process. Please try again later."
// @Router       /products/remove/{id} [delete]
func RemoveProduct(ctx *gin.Context) {
	service, err := businesslogic.GenerateProductService()
	if err != nil {
		utils.ProcessResponse(utils.GenerateInvalidRequestAndSystemProblemModel(ctx, err))
		return
	}

	utils.ProcessResponse(response.APIResponse{
		ErrMsg:   service.RemoveProduct(ctx.Param("id"), ctx),
		Context:  ctx,
		PostType: action_type.INFORM,
	})
}

// ActivateProduct godoc
// @Summary      Activate a product
// @Description  Set a product's status to active
// @Tags         products
// @Security     BearerAuth
// @Produce      json
// @Param        id path string true "Product ID"
// @Success      200 {object} response.MessageAPIResponse "success"
// @Failure 401 {object} response.MessageAPIResponse "You have no rights to access this action."
// @Failure 400 {object} response.MessageAPIResponse "Invalid data. Please try again."
// @Failure 500 {object} response.MessageAPIResponse "There is something wrong in the system during the process. Please try again later."
// @Router       /products/activate/{id} [patch]
func ActivateProduct(ctx *gin.Context) {
	service, err := businesslogic.GenerateProductService()
	if err != nil {
		utils.ProcessResponse(utils.GenerateInvalidRequestAndSystemProblemModel(ctx, err))
		return
	}

	utils.ProcessResponse(response.APIResponse{
		ErrMsg:   service.ActivateProduct(ctx.Param("id"), ctx),
		Context:  ctx,
		PostType: action_type.INFORM,
	})
}
