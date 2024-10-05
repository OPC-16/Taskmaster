package handlers

import (
	"net/http"
   "time"

	"taskmaster/internal/models"
	"taskmaster/internal/services"
   "taskmaster/pkg/auth"

	"github.com/labstack/echo/v4"
)

type UserHandler struct {
   userService services.UserService
   jwtSecret   string
}

func NewUserHandler(userService services.UserService, jwtSecret string) *UserHandler {
   return &UserHandler{
      userService: userService,
      jwtSecret: jwtSecret,
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

// authenticates the user and returns the JWT
func (h *UserHandler) Login(c echo.Context) error {
   var loginRequest struct {
      Username string `json:"username"`
      Password string `json:"password"`
   }

   if err := c.Bind(&loginRequest); err != nil {
      return c.JSON(http.StatusBadRequest, map[string]string{ "error": "Invalid request payload" })
   }

   // Authenticate the user
   user, err := h.userService.Authenticate(c.Request().Context(), loginRequest.Username, loginRequest.Password)
   if err != nil {
      return c.JSON(http.StatusUnauthorized, map[string]string{ "error": "Invalid credentials" })
   }

   // Generate the JWT
   token, err := auth.GenerateToken(user.ID, h.jwtSecret, 24*time.Hour)    // set the token expiry for 24 hrs
   if err != nil {
      return c.JSON(http.StatusInternalServerError, map[string]string{ "error": "Failed to generate token" })
   }

   return c.JSON(http.StatusOK, map[string]string{ "token": token })
}

// retrieves the user's profile
func (h *UserHandler) GetProfile(c echo.Context) error {
   // TODO: Implement get profile logic
   return c.JSON(http.StatusNotImplemented, map[string]string{ "message": "get profile not implemented yet" })
}

// updates the user's profile
func (h *UserHandler) UpdateProfile(c echo.Context) error {
   // TODO: Implement update profile logic
   return c.JSON(http.StatusNotImplemented, map[string]string{ "message": "update profile not implemented yet" })
}

// deletes a user account
func (h *UserHandler) DeleteUser(c echo.Context) error {
   // TODO: Implement delete user logic
   return c.JSON(http.StatusNotImplemented, map[string]string{ "message": "delete user not implemented yet" })
}
