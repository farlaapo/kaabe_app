package main

import (
	"database/sql"
	"fmt"
	"kaabe-app/internal/api/controller"
	"kaabe-app/internal/api/gateway"
	"kaabe-app/internal/api/routes"
	"kaabe-app/internal/config"
	"kaabe-app/internal/domain/service"

	// utils "kaabe-app/pkg/config"

	"log"
	// "net/http"
	// "os"
	// "time"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
)

func init() {
	err := godotenv.Load(".env")
	if err != nil {
		log.Println(".env file not found, relying on system environment variables")
	}
}

func main() {
	// Load configuration
	config.LoadEnv()
	appCfg, err := config.LoadAppConfig()
	if err != nil {
		log.Fatalf("Failed to load config.yaml: %v", err)
	}

	dbCfg := config.LoadDBConfig()
	db := config.InitDB(dbCfg)
	if db == nil {
		log.Fatal("Failed to initialize the database")
	}
	defer func() {
		if err := db.Close(); err != nil {
			log.Printf("Error closing DB: %v", err)
		}
	}()

	fmt.Printf("Server running on port %s in %s mode\n", appCfg.App.Port, appCfg.App.Env)

	// // Setup debug routes using net/http
	// go func() {
	// 	http.HandleFunc("/debug/jwt", func(w http.ResponseWriter, r *http.Request) {
	// 		jwt := os.Getenv("JWT_SECRET")
	// 		refresh := os.Getenv("JWT_REFRESH_SECRET")

	// 		if jwt == "" || refresh == "" {
	// 			http.Error(w, "JWT secrets not found", http.StatusInternalServerError)
	// 			return
	// 		}

	// 		fmt.Fprintf(w, "JWT_SECRET: %s\nJWT_REFRESH_SECRET: %s\n", jwt, refresh)
	// 	})

	// 	http.HandleFunc("/debug/token", func(w http.ResponseWriter, r *http.Request) {
	// 		userID := "123"

	// 		accessToken, err := utils.GenereteToken(userID, time.Now().Add(15*time.Minute).Unix())
	// 		if err != nil {
	// 			http.Error(w, "Access token error: "+err.Error(), http.StatusInternalServerError)
	// 			return
	// 		}

	// 		refreshToken, err := utils.GenerateRefreshToken(userID, time.Now().Add(7*24*time.Hour).Unix())
	// 		if err != nil {
	// 			http.Error(w, "Refresh token error: "+err.Error(), http.StatusInternalServerError)
	// 			return
	// 		}

	// 		fmt.Fprintf(w, "Access Token:\n%s\n\nRefresh Token:\n%s\n", accessToken, refreshToken)
	// 	})

	// 	port := os.Getenv("DEBUG_PORT")
	// 	if port == "" {
	// 		port = "9090"
	// 	}
	// 	log.Printf("Debug server running on port %s", port)
	// 	log.Fatal(http.ListenAndServe(":"+port, nil))
	// }()
	// Initialize Repositories with *sql.DB

	var dbConn *sql.DB = db

	userRepo := gateway.NewUserRepository(dbConn)
	tokenRepo := gateway.NewTokenRepository(dbConn)
	courseRepo := gateway.NewCourseRepository(dbConn)
	lessonRepo := gateway.NewLessonRepository(dbConn)
	ratingRepo := gateway.NewRatingRepository(dbConn)
	SubscriptionRepo := gateway.NewSubscriptionImpl(dbConn)
	WithdrawalRepo := gateway.NewWithdrawalRepositoryImpl(dbConn)
	paymentRepo := gateway.NewPaymentRepository(dbConn)


	// Initialize Services
	userService := service.NewUserService(userRepo, tokenRepo)
	courseService := service.NewCourseService(courseRepo, tokenRepo)
	lessonService := service.NewLessonService(lessonRepo, tokenRepo)
	ratingService := service.NewRatingService(ratingRepo, tokenRepo)
	subscriptionService := service.NewSubscriptionService(SubscriptionRepo, tokenRepo)
	withdrawalService := service.NewWithdrawalService(WithdrawalRepo, tokenRepo)
	paymentService := service.NewPaymentService(paymentRepo)


	// Initialize Controllers
	userController := controller.NewUserController(userService)
	courseController := controller.NewCourseController(courseService)
	lessonController := controller.NewLessonController(lessonService)
	ratingController := controller.NewRatingController(ratingService)
	subscriptionController := controller.NewSubscriptionController(subscriptionService)
	withdrawalController := controller.NewWithdrawalController(withdrawalService)
	paymentController := controller.NewPaymentController(paymentService)

	// Setup Gin HTTP Server
	r := gin.Default()
	r.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"http://localhost:3000"},
		AllowMethods:     []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "Content-Type", "Authorization"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true,
	}))

	// Register API Routes
	routes.RegisterUserRoutes(r, userController, tokenRepo)
	routes.RegisterCoursesRoutes(r, courseController, tokenRepo)
	routes.RegisterLessonRoutes(r, lessonController, tokenRepo)
	routes.RegisterRatingRoutes(r, ratingController, tokenRepo)
	routes.RegisterSubscriptionRoutes(r, subscriptionController, tokenRepo)
	routes.RegisterWithdrawalRoutes(r, withdrawalController, tokenRepo)
	routes.RegisterPaymentRoutes(r, paymentController, tokenRepo)
	// Start main API server
	if err := r.Run(":" + appCfg.App.Port); err != nil {
		log.Fatal("Failed to start API server:", err)
	}
}
