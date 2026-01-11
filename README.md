

# POS Backend – Setup & API Documentation

Backend service for a **Point of Sale (POS)** system built with **Golang**, using **PostgreSQL**, **Redis**, and **WebSocket** for realtime updates.
Designed to be consumed by **Flutter**, **Web**, or other POS clients.

---

## 🚀 Features

* Authentication (Register & Login)
* Role-based access (Kasir, Kitchen, etc.)
* Order Management
* Table Session (QR Ordering)
* Cashier Shift Management
* Realtime updates via WebSocket
* Docker-based infrastructure

---

## 🛠 Tech Stack

* **Golang** ≥ 1.22
* **PostgreSQL** (Docker)
* **Redis** (Docker)
* **Docker Compose**
* **WebSocket**
* **JWT Authentication**

---

## 📦 Prerequisites

Make sure you have the following installed:

* **Go (Golang)** ≥ 1.22
  👉 [https://go.dev/dl/](https://go.dev/dl/)
* **Docker Desktop**
  👉 [https://www.docker.com/products/docker-desktop/](https://www.docker.com/products/docker-desktop/)
* **Git**
  👉 [https://git-scm.com/](https://git-scm.com/)

---

## 🐳 Start Infrastructure (Docker)

We use Docker Compose to run PostgreSQL and Redis.

1. Open terminal (CMD / PowerShell) in the project root
2. Run:

```bash
docker-compose up -d
```

> `-d` means detached mode (runs in background)

3. Verify containers:

```bash
docker ps
```

You should see:

* `pos_postgres`
* `pos_redis`

---

## 🗄 Run Database Migrations

Create database tables using the provided migration script.

```bash
./migrate-up.bat
```

✅ If successful, you’ll see:

```
Migrasi berhasil dijalankan
```

### Troubleshooting

* Make sure PostgreSQL container is running and healthy
* Check Docker logs if migration fails

---

## ▶️ Run the Application

### Download dependencies (first time only)

```bash
go mod tidy
```

### Run server

```bash
go run cmd/main.go
```

Expected log:

```text
2026/01/11 17:00:00 Starting server on 0.0.0.0:8080
```

---

## ✅ Verify Server

Open browser or Postman:

```
http://localhost:8080/api/v1/ping
```

If configured correctly, you should receive a response or see logs confirming the server is alive.

---

## 📚 API Documentation

**Base URL**

```
http://localhost:8080/api/v1
```

---

### 1️⃣ Authentication

#### Register User

**POST** `/auth/register`
Auth: ❌ None

```json
{
  "username": "cashier1",
  "password": "strongpassword",
  "role": "KASIR"
}
```

**Response – 201 Created**

```json
{
  "id": "uuid...",
  "username": "cashier1",
  "role": "KASIR"
}
```

---

#### Login

**POST** `/auth/login`
Auth: ❌ None

```json
{
  "username": "cashier1",
  "password": "strongpassword"
}
```

**Response – 200 OK**

```json
{
  "access_token": "ey...",
  "user": { }
}
```

📌 **Note:**
Use this token for protected routes:

```
Authorization: Bearer <access_token>
```

---

### 2️⃣ Table Sessions (QR Ordering)

#### Create Session

**POST** `/table-sessions`

```json
{
  "table_id": "uuid..."
}
```

**Response – 201 Created**

```json
{
  "token": "session_token_..."
}
```

---

#### Validate Session

**POST** `/table-sessions/validate`

```json
{
  "token": "session_token_..."
}
```

---

### 3️⃣ Orders

#### Create Order

**POST** `/orders`
Auth: ✅ Required

```json
{
  "outlet_id": "uuid...",
  "table_session_id": "uuid... (optional)",
  "note": "No onions",
  "items": [
    {
      "product_id": "uuid...",
      "quantity": 2,
      "note": "Extra cheese"
    }
  ]
}
```

**Response – 201 Created**

```json
{
  "id": "uuid...",
  "status": "NEW",
  "total_amount": 150000
}
```

---

#### Get Order

**GET** `/orders/:id`
Auth: ✅ Required

---

#### Update Order Status

**PATCH** `/orders/:id/status`
Auth: ✅ Required

```json
{
  "status": "COOKING"
}
```

**Valid Statuses:**

```
NEW | ACCEPTED | COOKING | READY | DONE | VOIDED
```

---

### 4️⃣ Shifts (Cashier)

#### Open Shift

**POST** `/shifts/open`
Auth: ✅ Cashier

```json
{
  "outlet_id": "uuid...",
  "start_cash": 200000,
  "user_id": "uuid..."
}
```

---

#### Close Shift

**POST** `/shifts/close`
Auth: ✅ Cashier

```json
{
  "shift_id": "uuid...",
  "end_cash": 1500000
}
```

---

#### Get Current Shift

**GET** `/shifts/current`
Auth: ✅ Required

---

### 5️⃣ Realtime (WebSocket)

**URL**

```
ws://localhost:8080/api/v1/ws
```

#### Events

**NEW_ORDER**

```json
{
  "type": "NEW_ORDER",
  "payload": {
    "id": "...",
    "status": "NEW"
  }
}
```

---

## 🧪 Postman Testing Guide

### Prerequisites

* Postman installed
  👉 [https://www.postman.com/](https://www.postman.com/)
* Backend server running
* Import `postman_collection.json`

---

### Step 1: Import & Configure

1. Open Postman → Import
2. Import `postman_collection.json`
3. Open **POS Backend API → Variables**
4. Set:

   * `base_url` → `http://localhost:8080`
   * `outlet_id` → valid UUID
   * `product_id` → valid UUID

---

### Step 2: Authentication

1. Open **Auth / Login**
2. Send request
3. ✅ Token is auto-saved to `{{token}}`

---

### Step 3: Shift Management

1. Open **Shifts / Open Shift**
2. Send request
3. Verify shift created successfully

---

### Step 4: Order & Realtime Test

#### WebSocket

1. New → WebSocket Request
2. URL:

   ```
   ws://localhost:8080/api/v1/ws
   ```
3. Click **Connect**

#### Create Order

1. Open **Orders / Create Order**
2. Send request
3. ✅ WebSocket should receive:

```json
{
  "type": "NEW_ORDER"
}
```

---

## 🧾 Cheatsheet

| Action           | Command                |
| ---------------- | ---------------------- |
| Start DB & Redis | `docker-compose up -d` |
| Stop DB & Redis  | `docker-compose down`  |
| Run App          | `go run cmd/main.go`   |
| Run Tests        | `go test ./...`        |
| Run Migration    | `./migrate-up.bat`     |
| Reset Database   | `./migrate-down.bat`   |


