package main

import (
	"log"
	"net/http"

   "taskmaster/api/v1"
	"taskmaster/config"
   "taskmaster/internal/handlers"
   "taskmaster/internal/repositories"
   "taskmaster/internal/services"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

func main() {
   // Initialize configuration
   cfg, err := config.Load()
   if err != nil {
      log.Fatalf("Failed to load configuration: %v", err)
   }

   // Initialize database
   db, err := repositories.NewSQLiteDB(cfg.DatabasePath)
   if err != nil {
      log.Fatalf("Failed to connect to database: %v", err)
   }
   defer db.Close()

   // Initialize repositories
   userRepo := repositories.NewUserRepository(db)
   taskRepo := repositories.NewTaskRepository(db)

   // Initialize services
   userService := services.NewUserService(userRepo)
   taskService := services.NewTaskService(taskRepo)

   // Initialize handlers
   userHandler := handlers.NewUserHandler(userService, cfg.JWTSecret)
   taskHandler := handlers.NewTaskHandler(taskService)

   // Initialize Echo instance
   e := echo.New()

   // Middleware
   e.Use(middleware.Logger())
   e.Use(middleware.Recover())

   // setting up a hello world route for "/"
   e.GET("/", func(c echo.Context) error {
      return c.String(http.StatusOK, "Hello, World!\n")
   })

   // Routes
   v1.SetupRoutes(e, userHandler, taskHandler, cfg.JWTSecret)


   // start the server
   e.Logger.Fatal(e.Start(cfg.ServerAddress))
}
