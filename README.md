# 🛫 Travel Companion Backend

A production-ready RESTful API built with Go and the Gin framework.  
This backend powers a complete travel management platform, handling authentication, bookings, visa workflows, and admin operations.

## 🎯 What This Project Solves

Travel platforms require secure authentication, complex booking workflows, and admin moderation.  
This project demonstrates how to design and build a scalable backend system that supports real-world travel operations.

---

## 🚀 Key Features

### 🔐 Security & Identity
- JWT-based authentication
- Role-Based Access Control (User / Admin)
- Secure middleware-protected routes
- CORS enabled for frontend integration

### ✈️ Travel Operations
- **Hotels**
  - CRUD operations
  - Search by city and date
  - Booking lifecycle management
- **Flights**
  - Flight listing and filtering
  - Ticket booking system
- **Visas**
  - Digital visa application submission
  - Status tracking and admin approval workflow

### 🛠️ Support & Management
- Support ticket system for user issues
- Admin endpoints for user and visa management
- Health check endpoint for monitoring uptime

---

## 🧪 Tech Stack

| Category | Technology |
|-------|-----------|
| Language | Go (Golang) |
| Framework | Gin Web Framework |
| Database | GORM (MySQL / PostgreSQL) |
| Authentication | JWT |
| Configuration | godotenv |

---

## 📂 Project Structure

backend/
├── cmd/server/ # Application entry point
├── internal/
│ ├── handlers/ # HTTP handlers
│ ├── services/ # Business logic
│ └── repos/ # Repository layer
├── middleware/ # Auth & custom middlewares
├── migration/ # Database migrations
├── models/ # GORM models
├── pkg/ # Shared utilities
├── config/ # App configuration
└── .env # Environment variables

yaml
Copy code

---

## 🛠️ Getting Started

### Prerequisites
- Go 1.25+
- MySQL or PostgreSQL database

### Installation

```bash
git clone https://github.com/shokr12/travel-app-backend.git
cd backend
go mod tidy
Configuration
Create a .env file:

env
Copy code
PORT=8080
DB_HOST=localhost
DB_PORT=3306
DB_USER=root
DB_PASSWORD=yourpassword
DB_NAME=travel_db
JWT_SECRET=your_super_secret_key
Run the Server
bash
Copy code
go run cmd/server/main.go
Server runs at:

arduino
Copy code
http://localhost:8080
📡 API Endpoints (Quick Reference)
👤 Authentication
POST /api/v1/signup – Register a new user

POST /api/v1/login – Login and receive JWT

🏨 Hotels & ✈️ Flights
GET /api/v1/hotels – List hotels

POST /api/v1/hotels/book – Book a hotel (Auth required)

GET /api/v1/flights – Search flights

POST /api/v1/flights/book – Book flight (Auth required)

🛂 Visa Management
POST /api/v1/visas – Submit visa application

GET /api/v1/admin/visas/pending – Review pending visas (Admin)

📈 What This Project Demonstrates
Designing REST APIs with Go and Gin

Secure authentication & RBAC

Booking and approval workflows

Clean architecture and separation of concerns

Database modeling and migrations

👤 Author
Mahmoud Shokr
GitHub: https://github.com/shokr12
LinkedIn: https://www.linkedin.com/in/mahmoud-shokr12

⭐ If you find this project useful, feel free to star the repository.
