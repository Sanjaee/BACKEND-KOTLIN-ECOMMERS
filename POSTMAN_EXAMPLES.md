# Postman API Testing Examples

## Base URL
```
http://localhost:5000
```

## 1. GET All Products

### Request
```
GET /api/v1/products?page=1&limit=10&active_only=true
```

### Query Parameters
- `page` (optional): Page number, default: 1
- `limit` (optional): Items per page, default: 10, max: 100
- `category_id` (optional): Filter by category ID
- `featured` (optional): Filter featured products (true/false)
- `active_only` (optional): Filter active products only (true/false)

### Example Response
```json
{
  "success": true,
  "message": "Products retrieved successfully",
  "data": {
    "products": [
      {
        "id": "uuid-here",
        "category_id": "category-uuid",
        "name": "Laptop Gaming ASUS ROG",
        "description": "Laptop gaming dengan processor Intel i7",
        "sku": "LAP-ASUS-ROG-001",
        "price": 15000000,
        "stock": 10,
        "weight": 2500,
        "thumbnail": "https://example.com/image.jpg",
        "is_active": true,
        "is_featured": true,
        "created_at": "2024-01-15T10:00:00Z",
        "updated_at": "2024-01-15T10:00:00Z",
        "category": {
          "id": "category-uuid",
          "name": "Laptop",
          "slug": "laptop"
        },
        "images": []
      }
    ],
    "total": 50,
    "page": 1,
    "limit": 10
  }
}
```

---

## 2. GET Product by ID

### Request
```
GET /api/v1/products/:id
```

### Example
```
GET /api/v1/products/123e4567-e89b-12d3-a456-426614174000
```

### Example Response
```json
{
  "success": true,
  "message": "Product retrieved successfully",
  "data": {
    "id": "123e4567-e89b-12d3-a456-426614174000",
    "category_id": "category-uuid",
    "name": "Laptop Gaming ASUS ROG",
    "description": "Laptop gaming dengan processor Intel i7",
    "sku": "LAP-ASUS-ROG-001",
    "price": 15000000,
    "stock": 10,
    "weight": 2500,
    "thumbnail": "https://example.com/image.jpg",
    "is_active": true,
    "is_featured": true,
    "created_at": "2024-01-15T10:00:00Z",
    "updated_at": "2024-01-15T10:00:00Z",
    "category": {
      "id": "category-uuid",
      "name": "Laptop",
      "slug": "laptop"
    },
    "images": [
      {
        "id": "image-uuid",
        "product_id": "product-uuid",
        "image_url": "https://example.com/detail1.jpg",
        "sort_order": 1,
        "created_at": "2024-01-15T10:00:00Z"
      }
    ]
  }
}
```

---

## 3. CREATE Product

### Request
```
POST /api/v1/products
Content-Type: application/json
```

### Body (Full Example)
```json
{
  "category_id": "category-uuid-here",
  "name": "Laptop Gaming ASUS ROG Strix G16",
  "description": "Laptop gaming dengan processor Intel Core i7-13650HX, GPU NVIDIA RTX 4060 8GB, RAM 16GB DDR5, SSD 512GB, Layar 16 inch FHD 165Hz",
  "sku": "LAP-ASUS-ROG-G16-001",
  "price": 18500000,
  "stock": 15,
  "weight": 2500,
  "thumbnail": "https://example.com/images/products/laptop-asus-rog-strix.jpg",
  "is_active": true,
  "is_featured": true
}
```

### Body (Minimal - Required Fields Only)
```json
{
  "category_id": "category-uuid-here",
  "name": "Laptop Gaming ASUS ROG",
  "sku": "LAP-ASUS-ROG-001",
  "price": 15000000
}
```

### Example Response
```json
{
  "success": true,
  "message": "Product created successfully",
  "data": {
    "id": "new-product-uuid",
    "category_id": "category-uuid",
    "name": "Laptop Gaming ASUS ROG Strix G16",
    "sku": "LAP-ASUS-ROG-G16-001",
    "price": 18500000,
    "stock": 15,
    "is_active": true,
    "is_featured": true,
    "created_at": "2024-01-15T10:00:00Z"
  }
}
```

---

## 4. UPDATE Product

### Request
```
PUT /api/v1/products/:id
Content-Type: application/json
```

### Body (Partial Update - Only fields to update)
```json
{
  "name": "Laptop Gaming ASUS ROG Updated Name",
  "price": 14500000,
  "stock": 8,
  "is_featured": false
}
```

### Body (Full Update Example)
```json
{
  "category_id": "new-category-uuid",
  "name": "Updated Product Name",
  "description": "Updated description",
  "sku": "NEW-SKU-001",
  "price": 12000000,
  "stock": 20,
  "weight": 2000,
  "thumbnail": "https://example.com/new-thumbnail.jpg",
  "is_active": true,
  "is_featured": true
}
```

### Example Response
```json
{
  "success": true,
  "message": "Product updated successfully",
  "data": {
    "id": "product-uuid",
    "name": "Laptop Gaming ASUS ROG Updated Name",
    "price": 14500000,
    "stock": 8,
    "is_featured": false,
    "updated_at": "2024-01-15T11:00:00Z"
  }
}
```

---

## 5. DELETE Product

### Request
```
DELETE /api/v1/products/:id
```

### Example
```
DELETE /api/v1/products/123e4567-e89b-12d3-a456-426614174000
```

### Example Response
```json
{
  "success": true,
  "message": "Product deleted successfully",
  "data": null
}
```

---

## 6. ADD Product Image

### Request
```
POST /api/v1/products/:id/images
Content-Type: application/json
```

### Body
```json
{
  "image_url": "https://example.com/images/product-detail-1.jpg",
  "sort_order": 1
}
```

### Body (Minimal - sort_order optional)
```json
{
  "image_url": "https://example.com/images/product-detail-2.jpg"
}
```

### Example Response
```json
{
  "success": true,
  "message": "Image added successfully",
  "data": {
    "id": "image-uuid",
    "product_id": "product-uuid",
    "image_url": "https://example.com/images/product-detail-1.jpg",
    "sort_order": 1,
    "created_at": "2024-01-15T10:00:00Z"
  }
}
```

---

## 7. DELETE Product Image

### Request
```
DELETE /api/v1/products/images/:imageId
```

### Example
```
DELETE /api/v1/products/images/123e4567-e89b-12d3-a456-426614174000
```

### Example Response
```json
{
  "success": true,
  "message": "Image deleted successfully",
  "data": null
}
```

---

## Error Responses

### Validation Error (400)
```json
{
  "success": false,
  "message": "validation error message",
  "error": null
}
```

### Product Not Found (400)
```json
{
  "success": false,
  "message": "product not found",
  "error": null
}
```

### Category Not Found (400)
```json
{
  "success": false,
  "message": "category not found",
  "error": null
}
```

### SKU Already Exists (400)
```json
{
  "success": false,
  "message": "SKU already exists",
  "error": null
}
```

---

## Quick Test Examples

### Test 1: Create a Product
```json
POST /api/v1/products
{
  "category_id": "your-category-id",
  "name": "Test Product",
  "sku": "TEST-001",
  "price": 100000
}
```

### Test 2: Get All Products (with filters)
```
GET /api/v1/products?page=1&limit=5&featured=true&active_only=true
```

### Test 3: Update Product Stock
```json
PUT /api/v1/products/:id
{
  "stock": 25
}
```

### Test 4: Add Multiple Images to Product
```json
POST /api/v1/products/:id/images
{
  "image_url": "https://example.com/image1.jpg",
  "sort_order": 1
}

POST /api/v1/products/:id/images
{
  "image_url": "https://example.com/image2.jpg",
  "sort_order": 2
}
```

---

## Notes

1. **UUID Format**: Semua ID menggunakan UUID format
2. **Price**: Harga dalam integer (dalam satuan terkecil, misal: 15000000 = Rp 15.000.000)
3. **Weight**: Berat dalam gram (integer)
4. **Stock**: Stock dalam integer, minimal 0
5. **Soft Delete**: Delete menggunakan soft delete, data masih ada di database
6. **Pagination**: Default page=1, limit=10, maksimal limit=100
7. **Optional Fields**: Fields yang optional bisa diabaikan saat create/update
