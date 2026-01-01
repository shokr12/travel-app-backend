package main

import (
	"Visa/config"
	"Visa/internal/handlers"
	"Visa/internal/repos"
	"Visa/internal/services"
	"Visa/migration"
	"log"
	"time"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
)

func main() {
	migration.Migrate()
	config.ConnectToDB()
	godotenv.Load()
	// Initialize services
	userService := services.NewUserService(repos.NewUserRepo(config.Db))
	visaService := services.NewVisaService(repos.NewVisaRepo(config.Db))
	hotelService := services.NewHotelService(repos.NewHotelRepo(config.Db), repos.NewReservationRepo(config.Db))
	flightService := services.NewFlightService(repos.NewFlightRepo(config.Db), repos.NewReservationRepo(config.Db))
	supportService := services.NewSupportService(repos.NewSupportRepo(config.Db))

	// Initialize handlers
	userHandler := handlers.NewUserHandler(userService)
	visaHandler := handlers.NewVisaHandler(visaService)
	hotelHandler := handlers.NewHotelHandler(hotelService)
	flightHandler := handlers.NewFlightHandler(flightService)
	supportHandler := handlers.NewSupportHandler(supportService)

	// Setup router
	r := gin.Default()

	// CORS configuration
	r.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"http://localhost:5173", "https://yourdomain.com", "http://localhost:8080", "http://localhost:8081"},
		AllowMethods:     []string{"GET", "POST", "PUT", "PATCH", "DELETE", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "Content-Type", "Authorization"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true,
		MaxAge:           12 * time.Hour,
	}))

	// Rate limiting middleware (basic example)
	// You should use a proper rate limiter like github.com/ulule/limiter/v3

	// Public routes

	public := r.Group("/api/v1/")
	{
		// Auth

		public.POST("/login", userHandler.LoginUser)
		public.POST("/signup", userHandler.CreateUser)
		// Public flight and hotel viewing
		public.GET("/flights", flightHandler.GetFlights)
		public.GET("/flights/:id", flightHandler.GetFlightById)
		public.GET("/flights/city/:city", flightHandler.GetFlightsByCity)
		public.GET("/flights/date/:date", flightHandler.GetFlightsByDepartDate)

		public.GET("/hotels", hotelHandler.GetAllHotels)
		public.GET("/hotels/:id", hotelHandler.GetHotelById)
		public.GET("/hotels/city/:city", hotelHandler.GetHotelsByCity)
		public.GET("/hotels/checkin/:date", hotelHandler.GetHotelsByCheckInDate)
		public.GET("/hotels/checkout/:date", hotelHandler.GetHotelsByCheckOutDate)
		public.GET("/visas", visaHandler.GetAllVisa)
	}

	// Protected routes (require authentication)
	protected := r.Group("/api/v1")
	protected.Use(handlers.AuthMiddleware())
	{
		// User routes
		protected.GET("/users/:id", userHandler.GetUserById)
		protected.PUT("/users/:id", userHandler.UpdateUser)
		protected.DELETE("/users/:id", userHandler.DeleteUser)

		// Visa routes
		protected.POST("/visas", visaHandler.CreateVisa)
		protected.GET("/visas/:id", visaHandler.GetVisaById)
		protected.GET("/visas/user/:userId", visaHandler.GetVisasByUser)
		protected.PUT("/visas/:id", visaHandler.UpdateVisa)
		protected.DELETE("/visas/:id", visaHandler.DeleteVisa)

		// Hotel booking routes
		protected.POST("/hotels/book", hotelHandler.BookHotel)
		protected.POST("/hotels/cancel", hotelHandler.CancelHotel)
		protected.GET("/hotels/user/:userId", hotelHandler.GetHotelsByUser)

		// Flight booking routes
		protected.POST("/flights/book", flightHandler.BookFlight)
		protected.POST("/flights/cancel", flightHandler.CancelFlight)
		protected.GET("/flights/user/:userId", flightHandler.GetFlightsByUser)

		// Support ticket routes
		protected.POST("/support", supportHandler.CreateTicket)
		protected.GET("/support/:id", supportHandler.GetTicketById)
		protected.GET("/support/user/:userId", supportHandler.GetTicketsByUser)
		protected.PUT("/support/:id", supportHandler.UpdateTicket)
		protected.DELETE("/support/:id", supportHandler.DeleteTicket)
	}

	// Admin routes (require admin role)
	admin := r.Group("/api/v1/admin")
	admin.Use(handlers.AuthMiddleware(), handlers.AdminMiddleware())
	{
		// User management
		admin.GET("/users", userHandler.GetAllUsers)

		// Visa management
		admin.GET("/visas", visaHandler.GetAllVisa)
		admin.GET("/visas/approved", visaHandler.GetApprovedVisas)
		admin.GET("/visas/pending", visaHandler.GetPendingVisas)
		admin.GET("/visas/rejected", visaHandler.GetRejectedVisas)
		admin.POST("/visas/:id/approve", visaHandler.ApproveVisa)
		admin.POST("/visas/:id/reject", visaHandler.RejectVisa)

		// Hotel management
		admin.POST("/hotels", hotelHandler.CreateHotel)
		admin.PUT("/hotels/:id", hotelHandler.UpdateHotel)
		admin.DELETE("/hotels/:id", hotelHandler.DeleteHotel)

		// Support ticket management
		admin.GET("/support", supportHandler.GetAllTickets)
	}

	// Health check
	r.GET("/health", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"status": "ok",
			"time":   time.Now(),
		})
	})

	log.Println("Server starting on :8080")
	if err := r.Run(":8080"); err != nil {
		log.Fatal("Failed to start server:", err)
	}
}
