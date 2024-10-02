package v1

import (
   "taskmaster/internal/handlers"

	"github.com/labstack/echo/v4"
)

func SetupRoutes(e *echo.Echo, userHandler *handlers.UserHandler) {
   v1 := e.Group("/api/v1")

   // User routes
   v1.POST("/register", userHandler.Register)
}
