# 📈 TradePulse — Trading Dashboard

A full-stack trading dashboard built with **Go + React + PostgreSQL + Docker**.

![Go](https://img.shields.io/badge/Go-1.21-00ADD8?style=flat&logo=go)
![React](https://img.shields.io/badge/React-18-61DAFB?style=flat&logo=react)
![PostgreSQL](https://img.shields.io/badge/PostgreSQL-16-336791?style=flat&logo=postgresql)
![Docker](https://img.shields.io/badge/Docker-Compose-2496ED?style=flat&logo=docker)

---

## ✨ Features

- 🔐 JWT-based authentication (Register / Login)
- 🔍 Search any stock, ETF, or mutual fund (US + Indian markets)
- 📊 Live price quotes with change percentage
- 📈 Interactive area chart — 1 Day & 30 Day history
- 🗄️ Backend data pipeline — Yahoo Finance → PostgreSQL → API
- 🐳 One command Docker setup

---

## 🏗️ Architecture

```
Browser (React + Tailwind)
        │
        │  JWT in Authorization header
        ▼
  Go REST API (Gin)          ← Port 8080
        │
        ├── POST /api/auth/register
        ├── POST /api/auth/login
        ├── GET  /api/auth/me
        ├── GET  /api/market/search?q=
        ├── GET  /api/market/quote/:symbol
        └── GET  /api/market/history/:symbol?period=
                    │
                    │  Data Pipeline
                    ▼
          Yahoo Finance API  ──→  PostgreSQL   ← Port 5432
                                  (price_history,
                                   users,
                                   symbol_metadata)
```

---

## 📁 Project Structure

```
trading-dashboard/
├── backend/
│   ├── cmd/main.go                  ← Entry point, router setup
│   ├── internal/
│   │   ├── auth/
│   │   │   ├── handler.go           ← Register / Login / Me
│   │   │   └── middleware.go        ← JWT validation
│   │   ├── config/config.go         ← Env var loading
│   │   ├── db/db.go                 ← DB connect + migrations
│   │   ├── market/
│   │   │   ├── handler.go           ← HTTP handlers
│   │   │   ├── service.go           ← Data pipeline logic
│   │   │   └── repository.go        ← DB queries
│   │   └── models/models.go         ← All structs
│   ├── go.mod
│   ├── go.sum
│   └── Dockerfile
├── frontend/
│   ├── src/
│   │   ├── context/AuthContext.jsx  ← Global auth state
│   │   ├── services/api.js          ← Axios + interceptors
│   │   ├── pages/
│   │   │   ├── LoginPage.jsx
│   │   │   ├── RegisterPage.jsx
│   │   │   ├── DashboardPage.jsx    ← Search + quick picks
│   │   │   └── StockPage.jsx        ← Chart + OHLCV table
│   │   ├── App.jsx
│   │   └── main.jsx
│   ├── package.json
│   ├── vite.config.js
│   ├── Dockerfile
│   └── nginx.conf
├── docker-compose.yml
├── .gitignore
└── README.md
```

---

## 🚀 Quick Start — One Command (Docker)

### Prerequisites
- [Docker Desktop](https://www.docker.com/products/docker-desktop/) installed and running

```bash
# 1. Clone the repo
git clone https://github.com/YOUR_USERNAME/trading-dashboard.git
cd trading-dashboard

# 2. Generate go.sum (one time only — needs Go installed)
cd backend && go mod tidy && cd ..

# 3. Run everything
docker-compose up --build
```

| Service     | URL                    |
|-------------|------------------------|
| 🌐 Frontend  | http://localhost:3000  |
| ⚙️ Backend   | http://localhost:8080  |
| 🗄️ Database  | localhost:5432         |

```bash
# Stop everything
docker-compose down

# Stop and delete DB data too
docker-compose down -v
```

---

## 💻 Local Development (Without Docker)

### Step 1 — Start PostgreSQL via Docker
```bash
docker run --name trading_db \
  -e POSTGRES_USER=postgres \
  -e POSTGRES_PASSWORD=password \
  -e POSTGRES_DB=trading_db \
  -p 5432:5432 \
  -d postgres:16
```

### Step 2 — Start Backend
```bash
cd backend
go mod tidy
go run ./cmd/main.go
# http://localhost:8080
```

### Step 3 — Start Frontend
```bash
cd frontend
npm install
npm run dev
# http://localhost:3000
```

---

## 🧪 API Reference

### Health Check
```bash
curl http://localhost:8080/health
```

### Register
```bash
curl -X POST http://localhost:8080/api/auth/register \
  -H "Content-Type: application/json" \
  -d '{"email":"test@test.com","password":"pass123","name":"Test User"}'
```

### Login
```bash
curl -X POST http://localhost:8080/api/auth/login \
  -H "Content-Type: application/json" \
  -d '{"email":"test@test.com","password":"pass123"}'
```

### Get Current User
```bash
curl http://localhost:8080/api/auth/me \
  -H "Authorization: Bearer YOUR_TOKEN"
```

### Search
```bash
curl "http://localhost:8080/api/market/search?q=apple" \
  -H "Authorization: Bearer YOUR_TOKEN"
```

### Live Quote
```bash
curl http://localhost:8080/api/market/quote/AAPL \
  -H "Authorization: Bearer YOUR_TOKEN"

curl http://localhost:8080/api/market/quote/RELIANCE.NS \
  -H "Authorization: Bearer YOUR_TOKEN"
```

### Price History
```bash
curl "http://localhost:8080/api/market/history/AAPL?period=30d" \
  -H "Authorization: Bearer YOUR_TOKEN"

curl "http://localhost:8080/api/market/history/AAPL?period=1d" \
  -H "Authorization: Bearer YOUR_TOKEN"
```

---

## 🗄️ Database Schema

```sql
CREATE TABLE users (
    id            SERIAL PRIMARY KEY,
    email         VARCHAR(255) UNIQUE NOT NULL,
    password_hash VARCHAR(255) NOT NULL,
    name          VARCHAR(255) NOT NULL,
    created_at    TIMESTAMP DEFAULT NOW()
);

CREATE TABLE price_history (
    id          SERIAL PRIMARY KEY,
    symbol      VARCHAR(20)   NOT NULL,
    trade_date  DATE          NOT NULL,
    open_price  DECIMAL(14,4),
    high_price  DECIMAL(14,4),
    low_price   DECIMAL(14,4),
    close_price DECIMAL(14,4),
    volume      BIGINT,
    fetched_at  TIMESTAMP DEFAULT NOW(),
    UNIQUE(symbol, trade_date)
);

CREATE TABLE symbol_metadata (
    symbol       VARCHAR(20) PRIMARY KEY,
    company_name VARCHAR(255),
    currency     VARCHAR(10),
    exchange     VARCHAR(50),
    last_fetched TIMESTAMP DEFAULT NOW()
);
```

---

## 🔐 Security

- Passwords hashed with **bcrypt** (cost factor 12)
- JWT signed with **HS256**, 24-hour expiry
- Same error for wrong email/password (prevents user enumeration)
- JWT middleware validates signing algorithm

---

## 📦 Tech Stack

| Layer       | Technology                      |
|-------------|----------------------------------|
| Backend     | Go 1.21, Gin                    |
| Frontend    | React 18, Vite, Tailwind CSS    |
| Database    | PostgreSQL 16                   |
| Auth        | JWT (golang-jwt/jwt)            |
| Charts      | Recharts                        |
| Market Data | Yahoo Finance API               |
| Container   | Docker, Docker Compose          |