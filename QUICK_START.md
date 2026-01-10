# Quick Start Guide - Product API Testing

## Cara Menggunakan

### 1. Import Postman Collection

1. Buka Postman
2. Klik **Import** di kiri atas
3. Pilih file `POSTMAN_COLLECTION.json`
4. Collection akan muncul di sidebar kiri

### 2. Setup Environment Variable

1. Di Postman, buat Environment baru atau gunakan default
2. Tambahkan variable:
   - **Variable Name**: `base_url`
   - **Initial Value**: `http://localhost:5000`
   - **Current Value**: `http://localhost:5000`

3. Pilih environment tersebut di dropdown kanan atas Postman

### 3. Insert Dummy Data ke Database

**Opsi A: Menggunakan SQL Script**
```bash
# Masuk ke container PostgreSQL
docker exec -it yourapp_postgres psql -U yourapp_db -d yourapp

# Atau dari host
psql -h localhost -U yourapp_db -d yourapp -f INSERT_DUMMY_DATA.sql
```

**Opsi B: Menggunakan API (jika sudah ada endpoint create category)**

Pastikan server sudah running:
```bash
cd be
go run cmd/server/main.go
```

Atau dengan Docker:
```bash
docker-compose up -d
```

### 4. Testing API dengan Postman

#### Test 1: Get All Products
1. Buka request "Get All Products" dari collection
2. Klik **Send**
3. Response akan menampilkan list products

#### Test 2: Get Product by ID
1. Buka request "Get Product by ID"
2. Ganti `:id` dengan UUID product dari response sebelumnya
3. Atau gunakan UUID dari dummy data: `660e8400-e29b-41d4-a716-446655440001`
4. Klik **Send**

#### Test 3: Create Product
1. Buka request "Create Product"
2. Ganti `category_id` dengan UUID category yang valid (dari dummy data: `550e8400-e29b-41d4-a716-446655440004`)
3. Sesuaikan data product sesuai kebutuhan
4. Klik **Send**

#### Test 4: Update Product
1. Buka request "Update Product"
2. Ganti `:id` dengan UUID product yang ada
3. Update body dengan field yang ingin diubah
4. Klik **Send**

#### Test 5: Add Product Image
1. Buka request "Add Product Image"
2. Ganti `:id` dengan UUID product yang ada
3. Update `image_url` dengan URL gambar valid
4. Klik **Send**

### 5. Contoh Request Body (Copy-Paste Ready)

#### Create Product - Minimal
```json
{
  "category_id": "550e8400-e29b-41d4-a716-446655440004",
  "name": "Test Product",
  "sku": "TEST-001",
  "price": 100000
}
```

#### Create Product - Full
```json
{
  "category_id": "550e8400-e29b-41d4-a716-446655440004",
  "name": "Laptop Gaming ASUS ROG Strix G16",
  "description": "Laptop gaming dengan processor Intel Core i7-13650HX",
  "sku": "LAP-ASUS-ROG-G16-001",
  "price": 18500000,
  "stock": 15,
  "weight": 2500,
  "thumbnail": "https://example.com/images/products/laptop-asus-rog-strix.jpg",
  "is_active": true,
  "is_featured": true
}
```

#### Update Product - Partial
```json
{
  "price": 14500000,
  "stock": 8
}
```

#### Add Product Image
```json
{
  "image_url": "https://example.com/images/product-detail-1.jpg",
  "sort_order": 1
}
```

### 6. UUID Dummy Data untuk Testing

**Categories:**
- Elektronik: `550e8400-e29b-41d4-a716-446655440001`
- Fashion: `550e8400-e29b-41d4-a716-446655440002`
- Makanan & Minuman: `550e8400-e29b-41d4-a716-446655440003`
- Laptop: `550e8400-e29b-41d4-a716-446655440004`
- Smartphone: `550e8400-e29b-41d4-a716-446655440005`

**Products:**
- Laptop ASUS ROG: `660e8400-e29b-41d4-a716-446655440001`
- MacBook Pro: `660e8400-e29b-41d4-a716-446655440002`
- iPhone 15 Pro Max: `660e8400-e29b-41d4-a716-446655440003`
- Samsung S24 Ultra: `660e8400-e29b-41d4-a716-446655440004`

### 7. Query Parameters Examples

#### Get Featured Products Only
```
GET /api/v1/products?featured=true&active_only=true
```

#### Get Products by Category
```
GET /api/v1/products?category_id=550e8400-e29b-41d4-a716-446655440004
```

#### Get Products with Pagination
```
GET /api/v1/products?page=1&limit=5
```

#### Combined Filters
```
GET /api/v1/products?page=1&limit=10&category_id=550e8400-e29b-41d4-a716-446655440004&featured=true&active_only=true
```

### 8. Troubleshooting

**Error: "category not found"**
- Pastikan category_id sudah ada di database
- Gunakan UUID dari dummy data atau create category terlebih dahulu

**Error: "SKU already exists"**
- Gunakan SKU yang berbeda saat create product
- SKU harus unique

**Error: Connection refused**
- Pastikan server sudah running
- Check apakah port 5000 sudah digunakan
- Cek `docker-compose ps` untuk melihat status container

**Error: Product not found**
- Pastikan UUID yang digunakan sudah benar
- Product mungkin sudah di-soft delete

### 9. Next Steps

1. Test semua endpoints secara berurutan
2. Coba error cases (invalid UUID, missing fields, dll)
3. Test pagination dengan data banyak
4. Test filter combinations
5. Integrate dengan frontend/mobile app

### 10. Tips

- Gunakan Postman Environment untuk switch antara dev/staging/prod
- Save responses sebagai examples untuk dokumentasi
- Gunakan Postman Collection Runner untuk automated testing
- Export collection untuk backup atau sharing dengan team
