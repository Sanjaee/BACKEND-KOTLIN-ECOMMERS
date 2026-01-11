# Zacode Go Backend

Backend server untuk aplikasi Zacode menggunakan Go dengan arsitektur clean architecture.

## Struktur Proyek

```
/yourapp
│
├── cmd/
│   └── server/
│       └── main.go
│
├── internal/
│   ├── config/          # load env, config global
│   │   └── config.go
│   │
│   ├── app/             # HTTP handler + routing (Gin)
│   │   ├── router.go
│   │   ├── auth_handler.go
│   │   ├── chat_handler.go
│   │   └── ...
│   │
│   ├── service/         # business logic / usecase
│   │   ├── auth_service.go
│   │   ├── chat_service.go
│   │   └── ...
│   │
│   ├── repository/      # DB access (gorm / raw SQL)
│   │   ├── user_repo.go
│   │   ├── chat_repo.go
│   │   └── ...
│   │
│   ├── model/           # struct model untuk DB
│   │   ├── user.go
│   │   ├── chat.go
│   │   └── ...
│   │
│   ├── websocket/       # ws hub, manager, client
│   │   ├── hub.go
│   │   ├── client.go
│   │   ├── ws_handler.go
│   │   └── ...
│   │
│   └── util/            # helper: jwt, hash, error, response
│       ├── jwt.go
│       ├── hash.go
│       └── response.go
│
├── pkg/                 # library reusable (optional)
│   └── logger/
│       └── logger.go
│
├── go.mod
├── .env
├── Dockerfile
└── docker-compose.yml
```

## Deskripsi Folder

### `cmd/server/`
Entry point aplikasi. Berisi `main.go` yang menginisialisasi dan menjalankan server.

### `internal/config/`
Konfigurasi aplikasi, termasuk loading environment variables dan setup global config.

### `internal/app/`
Layer HTTP handler dan routing menggunakan Gin framework.
- `router.go`: Setup routing dan middleware
- `*_handler.go`: HTTP handlers untuk setiap endpoint

### `internal/service/`
Business logic layer (use case layer). Berisi logika bisnis aplikasi.

### `internal/repository/`
Data access layer. Interface dan implementasi untuk akses database (GORM atau raw SQL).

### `internal/model/`
Struct model untuk database. Definisi struct yang digunakan untuk mapping database.

### `internal/websocket/`
WebSocket implementation untuk real-time communication.
- `hub.go`: WebSocket hub untuk manage connections
- `client.go`: WebSocket client implementation
- `ws_handler.go`: WebSocket handler

### `internal/util/`
Utility functions dan helpers:
- `jwt.go`: JWT token generation dan validation
- `hash.go`: Password hashing utilities
- `response.go`: Standard response formatter

### `pkg/logger/`
Reusable logger library yang bisa digunakan di seluruh aplikasi.

## Setup

### Prerequisites
- Go 1.21+
- Docker & Docker Compose
- PostgreSQL
- Redis (optional)
- RabbitMQ (optional)

### Installation

1. Clone repository
```bash
git clone <repository-url>
cd /go
```

2. Copy environment file
```bash
cp .env.example .env
```

3. Update `.env` dengan konfigurasi yang sesuai

4. Install dependencies
```bash
go mod download
```

5. Run dengan Docker Compose
```bash
docker-compose up -d
```

6. Atau run secara lokal
```bash
go run cmd/server/main.go
```

## Environment Variables

Buat file `.env` dengan variabel berikut:

```env
# Server
PORT=5000
SERVER_HOST=0.0.0.0
CLIENT_URL=http://localhost:3000

# Database
POSTGRES_HOST=localhost
POSTGRES_PORT=5432
POSTGRES_USER=your_user
POSTGRES_PASSWORD=your_password
POSTGRES_DB=your_database
POSTGRES_SSLMODE=disable

# JWT
JWT_SECRET=your_jwt_secret_key

# Redis
REDIS_HOST=localhost
REDIS_PORT=6379
REDIS_PASSWORD=

# RabbitMQ
RABBITMQ_HOST=localhost
RABBITMQ_PORT=5672
RABBITMQ_USER=your_user
RABBITMQ_PASSWORD=your_password
```

## Development

### Run Development Server
```bash
go run cmd/server/main.go
```

### Build
```bash
go build -o bin/server cmd/server/main.go
```

### Run Tests
```bash
go test ./...
```

## Docker

### Build Image
```bash
docker build -t 
```

### Run with Docker Compose
```bash
docker-compose up -d
```

### Stop Services
```bash
docker-compose down
```

## Services & Ports

Setelah menjalankan `docker-compose up -d`, services berikut akan tersedia:

- **API Server**: http://localhost:5000
- **PostgreSQL**: localhost:5432
- **Redis**: localhost:6379
- **RabbitMQ Management UI**: http://localhost:15672
  - Username: `yourapp` (default)
  - Password: `password123` (default)
- **pgweb (Database UI)**: http://localhost:8081
  - Web-based PostgreSQL client untuk melihat dan mengelola database
  - Otomatis terhubung ke database yang dikonfigurasi

## Architecture

Aplikasi ini menggunakan **Clean Architecture** dengan layer separation:

1. **Handler Layer** (`internal/app/`): HTTP handlers dan routing
2. **Service Layer** (`internal/service/`): Business logic
3. **Repository Layer** (`internal/repository/`): Data access
4. **Model Layer** (`internal/model/`): Domain models

## License

MIT

// E-Commerce Database Schema
// Copy this code to https://dbdiagram.io/

Table users {
  id uuid [pk]
  email varchar(255) [not null, unique]
  username varchar(100) [unique]
  phone varchar(20)
  full_name varchar(255) [not null]
  password_hash varchar(255)
  user_type varchar(50) [default: 'member', note: 'member, seller, admin']
  profile_photo text
  date_of_birth date
  gender varchar(20)
  is_active boolean [default: true]
  is_verified boolean [default: false]
  last_login timestamp
  login_type varchar(50) [default: 'credential']
  google_id varchar(255) [unique]
  otp_code varchar(6)
  otp_expires_at timestamp
  reset_token text
  reset_expires_at timestamp
  created_at timestamp
  updated_at timestamp
  deleted_at timestamp
}

Table sellers {
  id uuid [pk]
  user_id uuid [not null, unique, ref: - users.id]
  shop_name varchar(255) [not null, unique]
  shop_slug varchar(255) [not null, unique]
  shop_description text
  shop_logo text
  shop_banner text
  shop_address text
  shop_city varchar(100)
  shop_province varchar(100)
  shop_phone varchar(20)
  shop_email varchar(255)
  is_verified boolean [default: false]
  is_active boolean [default: true]
  total_products int [default: 0]
  total_sales int [default: 0]
  rating_average decimal(3,2) [default: 0.00]
  total_reviews int [default: 0]
  created_at timestamp
  updated_at timestamp
}

Table categories {
  id uuid [pk]
  name varchar(255) [not null]
  description text
  slug varchar(255) [unique, not null]
  image_url text
  parent_id uuid [ref: > categories.id]
  is_active boolean [default: true]
  created_at timestamp
  updated_at timestamp
}

Table products {
  id uuid [pk]
  seller_id uuid [not null, ref: > sellers.id]
  category_id uuid [not null, ref: > categories.id]
  name varchar(255) [not null]
  description text
  sku varchar(100) [unique, not null]
  price int [not null]
  stock int [default: 0]
  weight int
  thumbnail text
  is_active boolean [default: true]
  is_featured boolean [default: false]
  created_at timestamp
  updated_at timestamp
  deleted_at timestamp
}

Table product_images {
  id uuid [pk]
  product_id uuid [not null, ref: > products.id]
  image_url text [not null]
  sort_order int [default: 0]
  created_at timestamp
}

Table carts {
  id uuid [pk]
  user_id uuid [not null, unique, ref: - users.id]
  created_at timestamp
  updated_at timestamp
}

Table cart_items {
  id uuid [pk]
  cart_id uuid [not null, ref: > carts.id]
  product_id uuid [not null, ref: > products.id]
  quantity int [not null, default: 1]
  price int [not null]
  created_at timestamp
  updated_at timestamp
}

Table addresses {
  id uuid [pk]
  user_id uuid [not null, ref: > users.id]
  label varchar(100)
  recipient_name varchar(255) [not null]
  phone varchar(20) [not null]
  address_line1 text [not null]
  address_line2 text
  city varchar(100) [not null]
  province varchar(100) [not null]
  postal_code varchar(10) [not null]
  is_default boolean [default: false]
  created_at timestamp
  updated_at timestamp
}

Table orders {
  id uuid [pk]
  order_number varchar(50) [unique, not null]
  user_id uuid [not null, ref: > users.id]
  shipping_address_id uuid [not null, ref: > addresses.id]
  subtotal int [not null]
  shipping_cost int [default: 0]
  total_amount int [not null]
  status varchar(50) [not null, default: 'pending']
  notes text
  created_at timestamp
  updated_at timestamp
}

Table order_items {
  id uuid [pk]
  order_id uuid [not null, ref: > orders.id]
  product_id uuid [not null, ref: > products.id]
  seller_id uuid [not null, ref: > sellers.id]
  product_name varchar(255) [not null]
  quantity int [not null]
  price int [not null]
  subtotal int [not null]
  created_at timestamp
}

Table payments {
  id uuid [pk]
  order_id varchar(50) [unique, not null]
  order_uuid uuid [not null, ref: - orders.id]
  midtrans_transaction_id varchar(255)
  amount int [not null]
  total_amount int [not null]
  status varchar(50) [not null, default: 'pending']
  payment_method varchar(50) [not null]
  payment_type varchar(50) [default: 'midtrans']
  fraud_status varchar(50)
  va_number varchar(50)
  bank_type varchar(50)
  qr_code_url text
  expiry_time timestamp
  midtrans_response text
  created_at timestamp
  updated_at timestamp
}

Table reviews {
  id uuid [pk]
  product_id uuid [not null, ref: > products.id]
  user_id uuid [not null, ref: > users.id]
  order_id uuid [not null, ref: > orders.id]
  rating int [not null]
  review_text text
  review_images text
  created_at timestamp
  updated_at timestamp
}

Ref: "sellers"."id" < "sellers"."shop_description"