package main

import (
	"log"
	"net/http"

	"taskmaster/config"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

func main() {
   cfg, err := config.Load()
   if err != nil {
      log.Fatalf("Failed to load configuration: %v", err)
   }

   // initialize Echo instance
   e := echo.New()

   // middleware
   e.Use(middleware.Logger())
   e.Use(middleware.Recover())

   // setting up a hello world route for "/"
   e.GET("/", func(c echo.Context) error {
      return c.String(http.StatusOK, "Hello, World!\n")
   })

   // start the server
   e.Logger.Fatal(e.Start(cfg.ServerAddress))
}
