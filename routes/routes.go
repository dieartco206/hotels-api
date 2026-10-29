package routes

import (
	"hotels-api/controllers"
	"hotels-api/middleware"
	"hotels-api/models"

	"github.com/gin-gonic/gin"
)

func SetupRouter() *gin.Engine {
	r := gin.Default()

	// Global CORS middleware
	r.Use(middleware.CORSMiddleware())

	// Health check
	r.GET("/health", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"status":  "healthy",
			"service": "hotels-api",
			"engine":  "Golang / Gin / GORM / MySQL",
		})
	})

	api := r.Group("/api")
	{
		// --- Public Auth ---
		auth := api.Group("/auth")
		{
			auth.POST("/register", controllers.Register)
			auth.POST("/login", controllers.Login)
			auth.GET("/me", middleware.AuthMiddleware(), controllers.Me)
		}

		// --- Public Hotels Catalog ---
		hotels := api.Group("/hotels")
		{
			hotels.GET("", controllers.GetHotels)
			hotels.GET("/:id", controllers.GetHotelByID)
		}

		// --- Public Reservations / Consumer Booking ---
		bookings := api.Group("/bookings")
		{
			bookings.POST("", controllers.CreateBooking)
			bookings.GET("/code/:code", controllers.GetBookingByCode)
		}

		// --- Public Partner Partnership Submission ---
		api.POST("/partners/apply", controllers.ApplyPartner)

		// --- Authenticated Guest Routes ---
		authenticated := api.Group("", middleware.AuthMiddleware())
		{
			authenticated.GET("/my-bookings", controllers.GetMyBookings)
			authenticated.POST("/bookings/:id/cancel", controllers.CancelBooking)
		}

		// --- Hotel Partner Routes (Extranet) ---
		partner := api.Group("/partner", middleware.AuthMiddleware(), middleware.RequireRole(models.RolePartner, models.RoleAdmin))
		{
			partner.GET("/dashboard", controllers.GetPartnerDashboard)
			partner.GET("/hotels", controllers.GetPartnerHotels)
			partner.POST("/hotels", controllers.CreateHotel)
			partner.PUT("/hotels/:id", controllers.UpdateHotel)

			partner.POST("/rooms", controllers.CreateRoom)
			partner.PUT("/rooms/:id", controllers.UpdateRoom)
			partner.PATCH("/rooms/:id/allotment", controllers.UpdateAllotment)

			partner.GET("/bookings", controllers.GetPartnerBookings)
			partner.POST("/checkin", controllers.CheckInGuest)
			partner.GET("/settlements", controllers.GetPartnerSettlements)
		}

		// --- Superadmin / Platform Owner Routes ---
		admin := api.Group("/admin", middleware.AuthMiddleware(), middleware.RequireRole(models.RoleAdmin))
		{
			admin.GET("/dashboard", controllers.GetAdminDashboard)
			admin.GET("/applications", controllers.GetPartnerApplications)
			admin.POST("/applications/:id/review", controllers.ReviewPartnerApplication)
			admin.GET("/bookings", controllers.GetAllBookings)
			admin.GET("/settlements", controllers.GetAllSettlements)
			admin.POST("/settlements/:id/payout", controllers.ApproveSettlementPayout)
			admin.PATCH("/hotels/:id/status", controllers.ToggleHotelStatus)
		}
	}

	return r
}
