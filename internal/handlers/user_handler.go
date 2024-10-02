package handlers

import (
	"net/http"
	"taskmaster/internal/models"
	"taskmaster/internal/services"

	"github.com/labstack/echo/v4"
)

type UserHandler struct {
   userService services.UserService
}

func NewUserHandler(userService services.UserService) *UserHandler {
   return &UserHandler{
      userService: userService,
   }
}

// handles user registration
func (h *UserHandler) Register(c echo.Context) error {
   user := new(models.User)
   if err := c.Bind(user); err != nil {
      return c.JSON(http.StatusBadRequest, map[string]string{ "error": "Invalid request payload" })
   }

   err := h.userService.CreateUser(c.Request().Context(), user)
   if err != nil {
      return c.JSON(http.StatusInternalServerError, map[string]string{ "error": err.Error() })
   }

   return c.JSON(http.StatusCreated, user)
}
