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
  id uuid [pk, default: `gen_random_uuid()`]
  email varchar(255) [not null, unique]
  username varchar(100) [unique]
  phone varchar(20)
  full_name varchar(255) [not null]
  password_hash varchar(255)
  user_type varchar(50) [default: 'member']
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
  created_at timestamp [default: `now()`]
  updated_at timestamp [default: `now()`]
  deleted_at timestamp
  
  indexes {
    email
    username
    google_id
    deleted_at
  }
}

Table categories {
  id uuid [pk, default: `gen_random_uuid()`]
  name varchar(255) [not null]
  description text
  slug varchar(255) [unique, not null]
  image_url text
  parent_id uuid
  is_active boolean [default: true]
  created_at timestamp [default: `now()`]
  updated_at timestamp [default: `now()`]
  
  indexes {
    slug
    parent_id
  }
}

Table products {
  id uuid [pk, default: `gen_random_uuid()`]
  category_id uuid [not null]
  name varchar(255) [not null]
  description text
  sku varchar(100) [unique, not null]
  price int [not null]
  stock int [default: 0]
  weight int [note: 'in grams']
  thumbnail text
  is_active boolean [default: true]
  is_featured boolean [default: false]
  created_at timestamp [default: `now()`]
  updated_at timestamp [default: `now()`]
  deleted_at timestamp
  
  indexes {
    category_id
    sku
    is_active
    is_featured
    deleted_at
  }
}

Table product_images {
  id uuid [pk, default: `gen_random_uuid()`]
  product_id uuid [not null]
  image_url text [not null]
  sort_order int [default: 0]
  created_at timestamp [default: `now()`]
  
  indexes {
    product_id
  }
}

Table carts {
  id uuid [pk, default: `gen_random_uuid()`]
  user_id uuid [not null, unique]
  created_at timestamp [default: `now()`]
  updated_at timestamp [default: `now()`]
  
  indexes {
    user_id
  }
}

Table cart_items {
  id uuid [pk, default: `gen_random_uuid()`]
  cart_id uuid [not null]
  product_id uuid [not null]
  quantity int [not null, default: 1]
  price int [not null]
  created_at timestamp [default: `now()`]
  updated_at timestamp [default: `now()`]
  
  indexes {
    cart_id
    product_id
    (cart_id, product_id) [unique]
  }
}

Table addresses {
  id uuid [pk, default: `gen_random_uuid()`]
  user_id uuid [not null]
  label varchar(100) [note: 'e.g., Home, Office']
  recipient_name varchar(255) [not null]
  phone varchar(20) [not null]
  address_line1 text [not null]
  address_line2 text
  city varchar(100) [not null]
  province varchar(100) [not null]
  postal_code varchar(10) [not null]
  is_default boolean [default: false]
  created_at timestamp [default: `now()`]
  updated_at timestamp [default: `now()`]
  
  indexes {
    user_id
    is_default
  }
}

Table orders {
  id uuid [pk, default: `gen_random_uuid()`]
  order_number varchar(50) [unique, not null]
  user_id uuid [not null]
  shipping_address_id uuid [not null]
  subtotal int [not null]
  shipping_cost int [default: 0]
  total_amount int [not null]
  status varchar(50) [not null, default: 'pending', note: 'pending, processing, shipped, delivered, cancelled']
  notes text
  created_at timestamp [default: `now()`]
  updated_at timestamp [default: `now()`]
  
  indexes {
    order_number
    user_id
    status
    created_at
  }
}

Table order_items {
  id uuid [pk, default: `gen_random_uuid()`]
  order_id uuid [not null]
  product_id uuid [not null]
  product_name varchar(255) [not null, note: 'snapshot of product name']
  quantity int [not null]
  price int [not null, note: 'price at time of order']
  subtotal int [not null]
  created_at timestamp [default: `now()`]
  
  indexes {
    order_id
    product_id
  }
}

Table payments {
  id uuid [pk, default: `gen_random_uuid()`]
  order_id varchar(50) [unique, not null, note: 'order number from orders table']
  order_uuid uuid [not null]
  midtrans_transaction_id varchar(255)
  amount int [not null]
  total_amount int [not null]
  status varchar(50) [not null, default: 'pending', note: 'pending, success, failed, cancelled, expired']
  payment_method varchar(50) [not null, note: 'bank_transfer, gopay, credit_card, qris']
  payment_type varchar(50) [default: 'midtrans']
  fraud_status varchar(50)
  va_number varchar(50)
  bank_type varchar(50)
  qr_code_url text
  expiry_time timestamp
  midtrans_response text [note: 'raw JSON response from Midtrans']
  created_at timestamp [default: `now()`]
  updated_at timestamp [default: `now()`]
  
  indexes {
    order_id
    order_uuid
    midtrans_transaction_id
    status
  }
}

Table reviews {
  id uuid [pk, default: `gen_random_uuid()`]
  product_id uuid [not null]
  user_id uuid [not null]
  order_id uuid [not null]
  rating int [not null, note: '1-5 stars']
  review_text text
  review_images text [note: 'JSON array of image URLs']
  created_at timestamp [default: `now()`]
  updated_at timestamp [default: `now()`]
  
  indexes {
    product_id
    user_id
    order_id
    rating
    (user_id, product_id, order_id) [unique, note: 'one review per product per order']
  }
}

// Relationships
Ref: categories.parent_id > categories.id 
Ref: products.category_id > categories.id [delete: cascade]
Ref: product_images.product_id > products.id [delete: cascade]
Ref: carts.user_id - users.id [delete: cascade]
Ref: cart_items.cart_id > carts.id [delete: cascade]
Ref: cart_items.product_id > products.id [delete: cascade]
Ref: addresses.user_id > users.id [delete: cascade]
Ref: orders.user_id > users.id
Ref: orders.shipping_address_id > addresses.id
Ref: order_items.order_id > orders.id [delete: cascade]
Ref: order_items.product_id > products.id
Ref: payments.order_uuid > orders.id
Ref: reviews.product_id > products.id [delete: cascade]
Ref: reviews.user_id > users.id [delete: cascade]
Ref: reviews.order_id > orders.id [delete: cascade]