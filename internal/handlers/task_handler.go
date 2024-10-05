package handlers

import (
	"net/http"

	"taskmaster/internal/models"
	"taskmaster/internal/services"
	"taskmaster/pkg/auth"

	"github.com/golang-jwt/jwt/v5"
	"github.com/labstack/echo/v4"
)

type TaskHandler struct {
   taskService services.TaskService
}

func NewTaskHandler(taskService services.TaskService) *TaskHandler {
   return &TaskHandler{taskService: taskService}
}

func (h *TaskHandler) CreateTask(c echo.Context) error {
   claims := c.Get("user").(*jwt.Token).Claims.(*auth.Claims)
   userID := claims.UserID
   
   var task models.Task
   if err := c.Bind(&task); err != nil {
      return c.JSON(http.StatusBadRequest, map[string]string{ "error": "Invalid request payload" })
   }

   task.UserID = userID

   err := h.taskService.CreateTask(c.Request().Context(), &task)
   if err != nil {
      return c.JSON(http.StatusInternalServerError, map[string]string{ "error": err.Error() })
   }

   return c.JSON(http.StatusCreated, task)
}

func (h *TaskHandler) GetTasks(c echo.Context) error {
   claims := c.Get("user").(*jwt.Token).Claims.(*auth.Claims)
   userID := claims.UserID

   tasks, err := h.taskService.GetTasksByUserID(c.Request().Context(), userID)
   if err != nil {
      return c.JSON(http.StatusInternalServerError, map[string]string{ "error": err.Error() })
   }

   return c.JSON(http.StatusOK, tasks)
}
