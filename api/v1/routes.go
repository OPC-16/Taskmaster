package v1

import (
	"taskmaster/internal/handlers"
	"taskmaster/pkg/auth"

	"github.com/golang-jwt/jwt/v5"
	"github.com/labstack/echo-jwt/v4"
	"github.com/labstack/echo/v4"
)

func SetupRoutes(e *echo.Echo, userHandler *handlers.UserHandler, jwtSecret string) {
   v1 := e.Group("/api/v1")

   // User routes (public)
   v1.POST("/register", userHandler.Register)
   v1.POST("/login", userHandler.Login)

   // Restricted routes (requires JWT)
   restricted_tasks := v1.Group("/tasks")
   restricted_tasks.Use(echojwt.WithConfig(echojwt.Config{
      SigningKey: []byte(jwtSecret),
      SigningMethod: "HS512",
      NewClaimsFunc: func(c echo.Context) jwt.Claims {
         return new(auth.Claims)
      },
   }))

   restricted_user := v1.Group("/user")
   restricted_user.Use(echojwt.WithConfig(echojwt.Config{
      SigningKey: []byte(jwtSecret),
      SigningMethod: "HS512",
      NewClaimsFunc: func(c echo.Context) jwt.Claims {
         return new(auth.Claims)
      },
   }))

   restricted_user.GET("/profile", userHandler.GetProfile)
}
