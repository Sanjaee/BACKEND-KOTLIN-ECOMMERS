package app

import (
	"net/http"
	"strconv"

	"yourapp/internal/service"
	"yourapp/internal/util"

	"github.com/gin-gonic/gin"
)

type ProductHandler struct {
	productService service.ProductService
}

func NewProductHandler(productService service.ProductService) *ProductHandler {
	return &ProductHandler{
		productService: productService,
	}
}

// CreateProduct handles product creation
// POST /api/v1/products
func (h *ProductHandler) CreateProduct(c *gin.Context) {
	var req service.CreateProductRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		util.BadRequest(c, err.Error())
		return
	}

	product, err := h.productService.CreateProduct(req)
	if err != nil {
		util.ErrorResponse(c, http.StatusBadRequest, err.Error(), nil)
		return
	}

	util.SuccessResponse(c, http.StatusCreated, "Product created successfully", product)
}

// GetProduct handles getting product by ID
// GET /api/v1/products/:id
func (h *ProductHandler) GetProduct(c *gin.Context) {
	id := c.Param("id")
	if id == "" {
		util.BadRequest(c, "Product ID is required")
		return
	}

	product, err := h.productService.GetProductByID(id)
	if err != nil {
		util.ErrorResponse(c, http.StatusBadRequest, err.Error(), nil)
		return
	}

	util.SuccessResponse(c, http.StatusOK, "Product retrieved successfully", product)
}

// GetProducts handles getting list of products
// GET /api/v1/products
func (h *ProductHandler) GetProducts(c *gin.Context) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "10"))

	categoryID := c.Query("category_id")
	featured := c.Query("featured")
	activeOnly := c.Query("active_only")

	var categoryIDPtr, featuredPtr, activeOnlyPtr *string
	if categoryID != "" {
		categoryIDPtr = &categoryID
	}
	if featured != "" {
		featuredPtr = &featured
	}
	if activeOnly != "" {
		activeOnlyPtr = &activeOnly
	}

	response, err := h.productService.GetProducts(page, limit, categoryIDPtr, featuredPtr, activeOnlyPtr)
	if err != nil {
		util.ErrorResponse(c, http.StatusBadRequest, err.Error(), nil)
		return
	}

	util.SuccessResponse(c, http.StatusOK, "Products retrieved successfully", response)
}

// UpdateProduct handles product update
// PUT /api/v1/products/:id
func (h *ProductHandler) UpdateProduct(c *gin.Context) {
	id := c.Param("id")
	if id == "" {
		util.BadRequest(c, "Product ID is required")
		return
	}

	var req service.UpdateProductRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		util.BadRequest(c, err.Error())
		return
	}

	product, err := h.productService.UpdateProduct(id, req)
	if err != nil {
		util.ErrorResponse(c, http.StatusBadRequest, err.Error(), nil)
		return
	}

	util.SuccessResponse(c, http.StatusOK, "Product updated successfully", product)
}

// DeleteProduct handles product deletion
// DELETE /api/v1/products/:id
func (h *ProductHandler) DeleteProduct(c *gin.Context) {
	id := c.Param("id")
	if id == "" {
		util.BadRequest(c, "Product ID is required")
		return
	}

	if err := h.productService.DeleteProduct(id); err != nil {
		util.ErrorResponse(c, http.StatusBadRequest, err.Error(), nil)
		return
	}

	util.SuccessResponse(c, http.StatusOK, "Product deleted successfully", nil)
}

// AddProductImage handles adding image to product
// POST /api/v1/products/:id/images
func (h *ProductHandler) AddProductImage(c *gin.Context) {
	productID := c.Param("id")
	if productID == "" {
		util.BadRequest(c, "Product ID is required")
		return
	}

	var req service.AddProductImageRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		util.BadRequest(c, err.Error())
		return
	}

	image, err := h.productService.AddProductImage(productID, req)
	if err != nil {
		util.ErrorResponse(c, http.StatusBadRequest, err.Error(), nil)
		return
	}

	util.SuccessResponse(c, http.StatusCreated, "Image added successfully", image)
}

// DeleteProductImage handles deleting product image
// DELETE /api/v1/products/images/:imageId
func (h *ProductHandler) DeleteProductImage(c *gin.Context) {
	imageID := c.Param("imageId")
	if imageID == "" {
		util.BadRequest(c, "Image ID is required")
		return
	}

	if err := h.productService.DeleteProductImage(imageID); err != nil {
		util.ErrorResponse(c, http.StatusBadRequest, err.Error(), nil)
		return
	}

	util.SuccessResponse(c, http.StatusOK, "Image deleted successfully", nil)
}
