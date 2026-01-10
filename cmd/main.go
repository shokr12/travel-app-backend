package main

import (
	"Visa/config"
	"Visa/internal/handlers"
	"Visa/internal/repos"
	"Visa/internal/services"
	"Visa/middleware"
	"Visa/migration"
	"log"
	"os"
	"time"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

func GetMe(c *gin.Context) {
	user, exists := c.Get("user")
	if !exists {
		c.JSON(401, gin.H{"error": "unauthorized"})
		return
	}
	c.JSON(200, user)
}

func main() {
	config.ConnectToDB()
	migration.Migrate()
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
		AllowOrigins:     []string{"https://travel-visa-app1.netlify.app"},
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
	protected.Use(middleware.AuthMiddleware())
	{
		// User routes
		protected.GET("/users/:id", userHandler.GetUserById)
		protected.PUT("/users/:id", userHandler.UpdateUser)
		protected.DELETE("/users/:id", userHandler.DeleteUser)
		protected.GET("/auth/me", GetMe)

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
	admin.Use(middleware.AuthMiddleware(), middleware.AdminMiddleware())
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
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}
	addr := "0.0.0.0:" + port
	if err := r.Run(addr); err != nil {
		log.Fatal("Failed to start server:", err)
	}

}
