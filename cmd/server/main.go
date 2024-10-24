package main

import (
	"fmt"
	"net/http"
	"os"
	"time"

	"taskmaster/api/v1"
	"taskmaster/config"
	"taskmaster/internal/handlers"
	"taskmaster/internal/repositories"
	"taskmaster/internal/services"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

func main() {
   fmt.Println("Started....")
   // Open or create a log file
   logFile, err := os.OpenFile("logs.log", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0666)
   if err != nil {
      fmt.Printf("Failed to open log file: %v", err)
      os.Exit(1)
   }

   // setup zerolog (default is JSON structed logging)
   zerolog.TimeFieldFormat = time.RFC1123
   log.Logger = zerolog.New(logFile).With().
      Timestamp().
      Logger()

   defer logFile.Close()   // Ensure the log file is closed when the application exits

   // Initialize configuration
   cfg, err := config.Load()
   if err != nil {
      log.Fatal().Err(err).Msg("Failed to load configuration")
   }

   // Initialize database
   db, err := repositories.NewSQLiteDB(cfg.DatabasePath)
   if err != nil {
      log.Fatal().Err(err).Msg("Failed to connect to database")
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
   e.HideBanner = true
   e.HidePort = true

   // Middleware for recover
   e.Use(middleware.Recover())

   // Configure Echo's logger to write to the log file
   e.Use(zerologRequestLogger())
   e.Use(middleware.RequestID())

   // setting up a hello world route for "/"
   e.GET("/", func(c echo.Context) error {
      log.Info().Msg("Handling root endpoint")
      return c.String(http.StatusOK, "Hello, World!\n")
   })

   // Routes
   v1.SetupRoutes(e, userHandler, taskHandler, cfg.JWTSecret)

   // start the server
   // e.Logger.Fatal(e.Start(cfg.ServerAddress))
   err = e.Start(cfg.ServerAddress)
   if err != nil {
      log.Fatal().Err(err).Msg("Failed to start the server")
   }
}

// creates a middleware for json structured request logging
func zerologRequestLogger() echo.MiddlewareFunc {
   return func(next echo.HandlerFunc) echo.HandlerFunc {
      return func(c echo.Context) error {
         start := time.Now()

         // execute the request handler
         err := next(c)

         latency := time.Since(start)

         // capture the response status after the handler has been executed
         status := c.Response().Status

         // if there's an error and the status is still 200, update it
         if err != nil && status == http.StatusOK {
            he, ok := err.(*echo.HTTPError)
            if ok {
               status = he.Code
            } else {
               status = http.StatusInternalServerError
            }
         }

         // Determine the message based on the status code
         var message string
         switch {
         case status >= 500:
            message = "Server error"
         case status >= 400:
            message = "Client error"
         case status >= 300:
            message = "Redirection"
         case status >= 200:
            message = "Success"
         default:
            message = "Unknown status"
         }

         method := c.Request().Method
         uri := c.Request().RequestURI
         reqID := c.Response().Header().Get(echo.HeaderXRequestID)
         userAgent := c.Request().UserAgent()
         bytesIn := c.Request().Header.Get(echo.HeaderContentLength)
         bytesOut := c.Response().Size

         log.Info().
            Str("method", method).
            Str("uri", uri).
            Str("request_id", reqID).
            Int("status", status).
            Dur("latency", latency).
            Str("user_agent", userAgent).
            Str("bytes_in", bytesIn).
            Int64("bytes_out", bytesOut).
            Msg(message)

         return err
      }
   }
}
