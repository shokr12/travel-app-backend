# 🛫 Travel Companion Backend

A robust, high-performance RESTful API built with **Go** and the **Gin** framework. This backend serves as the core engine for a complete travel management system, handling everything from user authentication to flight/hotel bookings and visa processing.

---

## 🚀 Key Features

### 🔐 Security & Identity

- **JWT Authentication**: Secure token-based access for protected resources.
- **Role-Based Access Control**: Separate flows for Users and Administrators.
- **CORS Enabled**: Configured for seamless interaction with modern frontend frameworks (React/Vue/etc).

### ✈️ Travel Operations

- **Hotels**: Full CRUD for listings, search by city/date, and booking lifecycle management.
- **Flights**: Real-time listing, filtering by destination or date, and ticketing system.
- **Visas**: Digital visa application platform with status tracking and admin approval workflows.

### 🛠️ Support & Management

- **Support Tickets**: Integrated helpdesk system for user queries.
- **Admin Dashboard**: Specialized endpoints for managing users, approving visas, and overseeing support.
- **Health Monitoring**: Built-in health check endpoint for uptime tracking.

---

## 🧪 Tech Stack

| Category        | Technology                                  |
| :-------------- | :------------------------------------------ |
| **Language**    | [Go (Golang)](https://golang.org/)          |
| **Framework**   | [Gin Web Framework](https://gin-gonic.com/) |
| **Database**    | GORM (ORM) with MySQL/PostgreSQL support    |
| **Auth**        | JWT (JSON Web Tokens)                       |
| **Environment** | godotenv                                    |

---

## 📂 Project Structure

```text
backend/
├── cmd/
│   └── server/          # Application entry point (main.go)
├── internal/
│   ├── handlers/        # HTTP Request handling logic
│   ├── services/        # Core business logic layer
│   └── repos/           # Database abstractions (Repository pattern)
├── middleware/          # Auth and custom Gin middlewares
├── migration/           # Database schema migrations
├── models/              # GORM data models
├── pkg/                 # Shared utilities and helpers
├── config/              # Database & App configurations
└── .env                 # Environment variables (not tracked)
```

---

## 🛠️ Getting Started

### 1. Prerequisites

- Go 1.25+ installed.
- A running database (configured via `.env`).

### 2. Installation

```bash
# Clone the repository
git clone <your-repo-url>
cd backend

# Install dependencies
go mod tidy
```

### 3. Configuration

Create a `.env` file in the root directory:

```env
PORT=8080
DB_HOST=localhost
DB_PORT=3306
DB_USER=root
DB_PASSWORD=yourpassword
DB_NAME=travel_db
JWT_SECRET=your_super_secret_key
```

### 4. Running the Server

```bash
# Start with live reloading (if Air is installed)
air

# OR standard run
go run cmd/server/main.go
```

---

## 📡 API Endpoints (Quick Reference)

### 👤 Authentication

- `POST /api/v1/signup` - Register a new account
- `POST /api/v1/login` - Authenticate and receive JWT

### 🏨 Hotels & ✈️ Flights

- `GET /api/v1/hotels` - List all hotels
- `POST /api/v1/hotels/book` - Reserve a room (Auth required)
- `GET /api/v1/flights` - Search available flights
- `POST /api/v1/flights/book` - Book a flight (Auth required)

### 🛂 Visa Management

- `POST /api/v1/visas` - Submit a visa application
- `GET /api/v1/admin/visas/pending` - Review pending applications (Admin only)

---

## 📄 License

This project is licensed under the MIT License - see the LICENSE file for details.
