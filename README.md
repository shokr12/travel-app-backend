🛫 Travel App Backend

A RESTful API backend built with Go and Gin for managing users, flights, hotels, visas, bookings, and support tickets. This API serves as the server for a complete travel booking system.

✔️ Includes user authentication, booking logic, admin routes, and trip management.

🚀 Features

User Authentication: Signup & login endpoints

Flights: List, filter (by city/date), book, cancel, and user bookings

Hotels: List, filter, book, cancel, and user bookings

Visas: Create, update, approve/reject (admin), and view per user

Support Tickets: Create, update, delete, and admin view all

Admin Routes: Manage users, visas, hotels, and support tickets

Health check endpoint

📁 Endpoints Overview
🔐 Authentication
POST /api/v1/signup
POST /api/v1/login

👤 Users
GET    /api/v1/users/:id
PUT    /api/v1/users/:id
DELETE /api/v1/users/:id
GET    /api/v1/admin/users

✈️ Flights
GET  /api/v1/flights
GET  /api/v1/flights/:id
GET  /api/v1/flights/city/:city
GET  /api/v1/flights/date/:date
POST /api/v1/flights/book
POST /api/v1/flights/cancel
GET  /api/v1/flights/user/:userId

🏨 Hotels
GET  /api/v1/hotels
GET  /api/v1/hotels/:id
GET  /api/v1/hotels/city/:city
GET  /api/v1/hotels/checkin/:date
GET  /api/v1/hotels/checkout/:date
POST /api/v1/hotels/book
POST /api/v1/hotels/cancel
GET  /api/v1/hotels/user/:userId

🛂 Visas
GET  /api/v1/visas
GET  /api/v1/visas/:id
GET  /api/v1/visas/user/:userId
POST /api/v1/visas
PUT  /api/v1/visas/:id
DELETE /api/v1/visas/:id

Admin Visa Management
GET  /api/v1/admin/visas
GET  /api/v1/admin/visas/approved
GET  /api/v1/admin/visas/pending
GET  /api/v1/admin/visas/rejected
POST /api/v1/admin/visas/:id/approve
POST /api/v1/admin/visas/:id/reject

🆘 Support
POST   /api/v1/support
GET    /api/v1/support/:id
GET    /api/v1/support/user/:userId
PUT    /api/v1/support/:id
DELETE /api/v1/support/:id
GET    /api/v1/admin/support

🩺 Health Check
GET /health

🧠 Tech Stack
Component	Technology
Backend	Go (Golang)
Framework	Gin Web Framework
Auth	JWT Tokens
Routing	REST API
Database	(Configurable — e.g., PostgreSQL/MySQL)*

* Database driver depends on your setup and environment variables.

⚙️ Setup & Development
1. Install Dependencies
go mod tidy

2. Set Environment Variables

Create a .env:

PORT=8080
DB_HOST=localhost
DB_USER=your_db_user
DB_PASSWORD=your_db_password
DB_NAME=travel_db
JWT_SECRET=your_jwt_secret

3. Run the Server
go run main.go


Server will run on:

http://localhost:8080

🧪 Testing

Use API tools like Postman or Insomnia to test endpoints.
Send authenticated requests using the JWT token from the login endpoint in headers:

Authorization: Bearer <token>

📦 Project Structure
/internal
  /handlers      → Route handlers
  /models        → Data models
  /middlewares   → Auth / middleware logic
  /db            → Database connection & config
/cmd             → Entry point
/main.go         → Start the server

📝 Contribution

You’re welcome to contribute!
Steps:

Fork this repo

Create a new branch

Make your changes

Submit a pull request

📄 License

Add your license details here (e.g., MIT, GPL, etc.)
