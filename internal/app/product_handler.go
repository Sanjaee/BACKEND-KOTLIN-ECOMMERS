package app

import (
	"fmt"
	"io"
	"net/http"
	"path/filepath"
	"strconv"
	"strings"

	"yourapp/internal/config"
	"yourapp/internal/service"
	"yourapp/internal/util"

	"github.com/gin-gonic/gin"
)

type ProductHandler struct {
	productService   service.ProductService
	cloudinaryUpload *util.CloudinaryUploader
}

func NewProductHandler(productService service.ProductService, cfg *config.Config) *ProductHandler {
	var uploader *util.CloudinaryUploader
	if cfg.CloudinaryCloudName != "" && cfg.CloudinaryAPIKey != "" && cfg.CloudinaryAPISecret != "" {
		uploader = util.NewCloudinaryUploader(cfg.CloudinaryCloudName, cfg.CloudinaryAPIKey, cfg.CloudinaryAPISecret)
	}

	return &ProductHandler{
		productService:   productService,
		cloudinaryUpload: uploader,
	}
}

// CreateProduct handles product creation
// POST /api/v1/products
func (h *ProductHandler) CreateProduct(c *gin.Context) {
	userID, exists := c.Get("userID")
	if !exists {
		util.Unauthorized(c, "User not authenticated")
		return
	}

	var req service.CreateProductRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		util.BadRequest(c, err.Error())
		return
	}

	product, err := h.productService.CreateProduct(userID.(string), req)
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

// UploadMultipleProductImages handles uploading multiple images to Cloudinary and saving to database
// POST /api/v1/products/:id/images/upload
func (h *ProductHandler) UploadMultipleProductImages(c *gin.Context) {
	productID := c.Param("id")
	if productID == "" {
		util.BadRequest(c, "Product ID is required")
		return
	}

	// Validate product exists
	_, err := h.productService.GetProductByID(productID)
	if err != nil {
		util.ErrorResponse(c, http.StatusBadRequest, "Product not found", nil)
		return
	}

	if h.cloudinaryUpload == nil {
		util.ErrorResponse(c, http.StatusInternalServerError, "Cloudinary is not configured", nil)
		return
	}

	// Parse multipart form (max 20MB)
	err = c.Request.ParseMultipartForm(20 << 20) // 20MB
	if err != nil {
		util.BadRequest(c, "Failed to parse multipart form: "+err.Error())
		return
	}

	// Get files from form
	files := c.Request.MultipartForm.File["images"]
	if len(files) == 0 {
		util.BadRequest(c, "No images provided")
		return
	}

	// Limit to 20 images
	if len(files) > 20 {
		util.BadRequest(c, "Maximum 20 images allowed")
		return
	}

	// Validate MIME types
	allowedMIMETypes := map[string]bool{
		"image/jpeg": true,
		"image/jpg":  true,
		"image/png":  true,
		"image/webp": true,
		"image/gif":  true,
	}

	var fileDataList []util.FileData
	for _, fileHeader := range files {
		// Validate MIME type
		contentType := fileHeader.Header.Get("Content-Type")
		if contentType == "" {
			// Try to detect from filename
			ext := strings.ToLower(filepath.Ext(fileHeader.Filename))
			mimeMap := map[string]string{
				".jpg":  "image/jpeg",
				".jpeg": "image/jpeg",
				".png":  "image/png",
				".webp": "image/webp",
				".gif":  "image/gif",
			}
			if m, ok := mimeMap[ext]; ok {
				contentType = m
			}
		}

		if !allowedMIMETypes[contentType] {
			util.BadRequest(c, fmt.Sprintf("File %s has invalid image format. Allowed: JPEG, PNG, WEBP, GIF", fileHeader.Filename))
			return
		}

		// Open file
		file, err := fileHeader.Open()
		if err != nil {
			util.ErrorResponse(c, http.StatusBadRequest, fmt.Sprintf("Failed to open file %s: %s", fileHeader.Filename, err.Error()), nil)
			return
		}

		// Read file data
		fileData, err := io.ReadAll(file)
		file.Close()
		if err != nil {
			util.ErrorResponse(c, http.StatusBadRequest, fmt.Sprintf("Failed to read file %s: %s", fileHeader.Filename, err.Error()), nil)
			return
		}

		// Validate file size (max 5MB per image)
		if len(fileData) > 5<<20 {
			util.BadRequest(c, fmt.Sprintf("File %s exceeds 5MB limit", fileHeader.Filename))
			return
		}

		fileDataList = append(fileDataList, util.FileData{
			Data: fileData,
			Name: fileHeader.Filename,
		})
	}

	// Upload to Cloudinary
	folder := fmt.Sprintf("products/%s", productID)
	urls, err := h.cloudinaryUpload.UploadMultipleImages(fileDataList, folder, 20)
	if err != nil {
		util.ErrorResponse(c, http.StatusInternalServerError, "Failed to upload images: "+err.Error(), nil)
		return
	}

	// Save to database
	for i, url := range urls {
		req := service.AddProductImageRequest{
			ImageURL:  url,
			SortOrder: func() *int { v := i; return &v }(),
		}
		_, err := h.productService.AddProductImage(productID, req)
		if err != nil {
			util.ErrorResponse(c, http.StatusInternalServerError, fmt.Sprintf("Failed to save image %d: %s", i+1, err.Error()), nil)
			return
		}
	}

	util.SuccessResponse(c, http.StatusCreated, fmt.Sprintf("%d images uploaded successfully", len(urls)), gin.H{
		"images": urls,
		"count":  len(urls),
	})
}
