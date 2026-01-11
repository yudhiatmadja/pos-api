# POS Backend - Point of Sale System

Backend service untuk sistem **Point of Sale (POS)** yang dibangun dengan **Golang**, menggunakan **PostgreSQL**, **Redis**, dan **WebSocket** untuk update realtime. Dirancang untuk dikonsumsi oleh client **Flutter**, **Web**, atau aplikasi POS lainnya.

---

## 🚀 Fitur

- **Authentication & Authorization** - Register, Login, Role-based Access Control (RBAC)
- **Order Management** - Create, Update, Track orders dengan state machine
- **Table Sessions** - QR Code ordering untuk pelanggan
- **Cashier Shift Management** - Open/Close shift dengan cash tracking
- **Realtime Updates** - WebSocket untuk notifikasi live
- **Idempotency** - Mencegah duplicate orders
- **Audit Logging** - Track semua perubahan data penting

---

## 🛠 Tech Stack

- **Golang** ≥ 1.22
- **PostgreSQL** (via Docker)
- **Redis** (via Docker)
- **Docker Compose**
- **WebSocket** untuk realtime communication
- **JWT** untuk authentication

---

## 📦 Prerequisites

Pastikan tools berikut sudah terinstall:

- **Go (Golang)** ≥ 1.22  
  👉 [https://go.dev/dl/](https://go.dev/dl/)
- **Docker Desktop**  
  👉 [https://www.docker.com/products/docker-desktop/](https://www.docker.com/products/docker-desktop/)
- **Git**  
  👉 [https://git-scm.com/](https://git-scm.com/)

---

## 🚀 Quick Start

### 1. Clone Repository

```bash
git clone <repository-url>
cd pos-backend
```

### 2. Start Infrastructure (Docker)

Jalankan PostgreSQL dan Redis menggunakan Docker Compose:

```bash
docker-compose up -d
```

Verifikasi containers berjalan:

```bash
docker ps
```

Anda harus melihat:
- `pos_postgres`
- `pos_redis`

### 3. Run Database Migrations

Buat database tables dengan migration script:

```bash
./migrate-up.bat
```

✅ Jika berhasil, akan muncul pesan:
```
Migrasi berhasil dijalankan
```

### 4. Download Dependencies

```bash
go mod tidy
```

### 5. Run Application

```bash
go run cmd/main.go
```

Expected output:
```
2026/01/11 17:00:00 Starting server on 0.0.0.0:8080
```

### 6. Verify Server

Buka browser atau Postman dan akses:

```
http://localhost:8080/api/v1/ping
```

---

## 📚 API Documentation

### Base URL

```
http://localhost:8080/api/v1
```

### Authentication

Semua endpoint (kecuali `/auth/register` dan `/auth/login`) memerlukan **Bearer Token**.

**Header:**
```
Authorization: Bearer <access_token>
```

---

## 🔐 1. Authentication

### Register User

Membuat user baru (staff, cashier, owner).

**Endpoint:** `POST /auth/register`  
**Auth:** ❌ None

**Request Body:**
```json
{
  "email": "owner@tokopi.com",
  "password": "strongpassword",
  "full_name": "John Doe",
  "role": "STORE_OWNER",
  "store_id": "uuid-store-123"
}
```

**Response:** `201 Created`
```json
{
  "id": "uuid...",
  "email": "owner@tokopi.com",
  "full_name": "John Doe",
  "role": "STORE_OWNER",
  "store_id": "uuid-store-123"
}
```

### Login

Authenticate dan mendapatkan JWT access token.

**Endpoint:** `POST /auth/login`  
**Auth:** ❌ None

**Request Body:**
```json
{
  "email": "owner@tokopi.com",
  "password": "strongpassword"
}
```

**Response:** `200 OK`
```json
{
  "access_token": "eyJhbGciOiJIUzI1NiIs...",
  "user": {
    "id": "uuid...",
    "email": "owner@tokopi.com",
    "role": "STORE_OWNER"
  }
}
```

---

## 👥 Roles & Permissions (RBAC)

| Role | Permissions |
|------|-------------|
| `SUPER_ADMIN` | Full platform access (all stores & data) |
| `STORE_OWNER` | Manage store, products, staff, reports, approve VOID |
| `KASIR` | Open/Close shift, create order, process payment |
| `KITCHEN` | View orders, update status (COOKING, READY) |
| `STAFF` | Create order only (waiter) |

---

## 🍽 2. Table Sessions (QR Ordering)

Untuk customer self-ordering via QR code.

### Create Table Session

**Endpoint:** `POST /table-sessions`  
**Auth:** ✅ Required

**Request Body:**
```json
{
  "table_id": "uuid-table-123"
}
```

**Response:** `201 Created`
```json
{
  "id": "uuid-session-456",
  "token": "session_token_xxx",
  "expires_at": "2026-01-11T18:00:00Z"
}
```

### Validate Table Session

**Endpoint:** `POST /table-sessions/validate`  
**Auth:** ❌ None

**Request Body:**
```json
{
  "token": "session_token_xxx"
}
```

**Response:** `200 OK`
```json
{
  "table_id": "uuid-table-123",
  "store_id": "uuid-store-123",
  "valid": true
}
```

---

## 🧾 3. Orders

### Create Order

**Endpoint:** `POST /orders`  
**Auth:** ✅ Required (STAFF, KASIR, STORE_OWNER)

**Request Body:**
```json
{
  "store_id": "uuid-store-123",
  "table_session_id": "uuid-session-456",
  "items": [
    {
      "product_id": "uuid-product-789",
      "quantity": 2,
      "note": "Less ice"
    }
  ],
  "note": "Table 5",
  "idempotency_key": "order-abc-123"
}
```

**Response:** `201 Created`
```json
{
  "id": "uuid-order-999",
  "order_number": "ORD-00021",
  "status": "NEW",
  "total_amount": 150000,
  "payment_status": "UNPAID"
}
```

### Get Order Detail

**Endpoint:** `GET /orders/:id`  
**Auth:** ✅ Required

### Update Order Status

**Endpoint:** `PATCH /orders/:id/status`  
**Auth:** ✅ Required (KITCHEN, STORE_OWNER)

**Request Body:**
```json
{
  "status": "COOKING"
}
```

**Valid Status Flow:**
```
NEW → ACCEPTED → COOKING → READY → DONE
       ↓
    VOIDED (terminal, requires STORE_OWNER)
```

---

## 💰 4. Shifts (Cashier)

### Open Shift

**Endpoint:** `POST /shifts/open`  
**Auth:** ✅ KASIR

**Request Body:**
```json
{
  "store_id": "uuid-store-123",
  "opening_cash": 500000
}
```

### Close Shift

**Endpoint:** `POST /shifts/close`  
**Auth:** ✅ KASIR

**Request Body:**
```json
{
  "shift_id": "uuid-shift-123",
  "end_cash": 1500000
}
```

### Get Current Shift

**Endpoint:** `GET /shifts/current`  
**Auth:** ✅ Required

---

## 🔴 5. Realtime (WebSocket)

Menerima live events untuk orders.

**WebSocket URL:**
```
ws://localhost:8080/api/v1/ws
```

### Event: NEW_ORDER

Triggered ketika order baru dibuat.

```json
{
  "type": "NEW_ORDER",
  "payload": {
    "id": "uuid-order-999",
    "order_number": "ORD-00021",
    "status": "NEW"
  }
}
```

---

## 🧠 6. Business Logic Rules

### Order State Machine

- Strict transition enforcement (tidak bisa skip atau revert)
- Terminal states: `DONE`, `VOIDED`
- Flow: `NEW → ACCEPTED → COOKING → READY → DONE`

### Idempotency (Create Order)

Untuk mencegah duplicate orders:
- Kirim `idempotency_key` pada request
- Jika key sudah ada → server return response original
- Tidak ada duplicate order yang dibuat

### RBAC Enforcement

- **VOID** requires `STORE_OWNER`
- **KASIR** cannot VOID
- **STAFF** cannot update order status
- **KITCHEN** cannot create orders

---

## 🗄 7. Database Schema

### Core Tables

**stores**
```
id, name, location, created_at
```

**profiles**
```
id, email, full_name, role, store_id
```

**roles**
```
code, name, description
```

### Order Management

**orders**
```
id, store_id, order_number, status,
total_amount, payment_status,
table_session_id, created_at
```

**order_items**
```
id, order_id, product_id,
quantity, price, total_price, note
```

### Products

**products**
```
id, store_id, name, description,
price, stock, category_id, is_available
```

**categories**
```
id, store_id, name
```

### Payment & Shift

**payments**
```
id, order_id, amount, payment_method, status, qris_url
```

**shifts**
```
id, store_id, user_id, start_time, end_time, start_cash, end_cash
```

### Security & Audit

**table_sessions**
```
id, table_id, token, expires_at
```

**audit_logs**
```
id, user_id, action, entity, entity_id, before_state, after_state
```

**idempotency_keys**
```
key, response_status, response_body, created_at
```

---

## 🐛 Troubleshooting

### Database Migration Gagal

- Pastikan PostgreSQL container berjalan dan healthy
- Check Docker logs: `docker logs pos_postgres`
- Verify connection string di config

### Server Tidak Bisa Start

- Cek apakah port 8080 sudah digunakan
- Pastikan Redis dan PostgreSQL sudah running
- Verify environment variables

### WebSocket Connection Failed

- Pastikan server sudah running
- Check firewall settings
- Verify WebSocket URL format

---

## 📝 License

[Add your license here]

---

## 🤝 Contributing

[Add contributing guidelines here]

---

## 📧 Contact

[Add contact information here]