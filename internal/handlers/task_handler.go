package handlers

import (
	"net/http"
	"strconv"

	"taskmaster/internal/models"
	"taskmaster/internal/services"
	"taskmaster/pkg/auth"

	"github.com/golang-jwt/jwt/v5"
	"github.com/labstack/echo/v4"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

type TaskHandler struct {
   taskService services.TaskService
   logger      zerolog.Logger
}

func NewTaskHandler(taskService services.TaskService) *TaskHandler {
   handlerLogger := log.With().Str("component", "task_handler").Caller().Logger()
   return &TaskHandler{
      taskService: taskService,
      logger: handlerLogger,
   }
}

func (h *TaskHandler) CreateTask(c echo.Context) error {
   claims := c.Get("user").(*jwt.Token).Claims.(*auth.Claims)
   userID := claims.UserID
   
   var task models.Task
   if err := c.Bind(&task); err != nil {
      h.logger.Err(err).Msg("Invalid request payload")
      return c.JSON(http.StatusBadRequest, map[string]string{ "error": "Invalid request payload" })
   }

   task.UserID = userID

   err := h.taskService.CreateTask(c.Request().Context(), &task)
   if err != nil {
      h.logger.Err(err)
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

func (h *TaskHandler) DeleteTask(c echo.Context) error {
   claims := c.Get("user").(*jwt.Token).Claims.(*auth.Claims)
   userID := claims.UserID

   taskID, err := strconv.ParseInt(c.Param("id"), 10, 64)
   if err != nil {
      return c.JSON(http.StatusBadRequest, map[string]string{ "error": "Invalid task id"})
   }

   err = h.taskService.DeleteTask(c.Request().Context(), taskID, userID)
   if err != nil {
      return c.JSON(http.StatusInternalServerError, map[string]string{ "error": err.Error() })
   }

   return c.NoContent(http.StatusNoContent)
}
