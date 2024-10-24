package services

import (
	"context"
	"errors"

	"taskmaster/internal/models"
	"taskmaster/internal/repositories"

   "github.com/rs/zerolog"
   "github.com/rs/zerolog/log"
)

type TaskService interface {
   CreateTask(ctx context.Context, task *models.Task) error
   GetTasksByUserID(ctx context.Context, userID int64) ([]models.Task, error)
   DeleteTask(ctx context.Context, taskID int64, userID int64) error
}

type taskService struct {
   taskRepo repositories.TaskRepository
   logger   zerolog.Logger
}

func NewTaskService(repo repositories.TaskRepository) TaskService {
   serviceLogger := log.With().Str("component", "task_service").Caller().Logger()
   return &taskService{
      taskRepo: repo,
      logger: serviceLogger,
   }
}

func (s *taskService) CreateTask(ctx context.Context, task *models.Task) error {
   // Business rules go here, such as checking task validity
   if task.Title == "" {
      err := "Task title is required"
      s.logger.Error().Msg(err)
      return errors.New(err)
   }

   // Call the repository to save the task
   return s.taskRepo.CreateTask(ctx, task)
}

func (s *taskService) GetTasksByUserID(ctx context.Context, userID int64) ([]models.Task, error) {
   tasks, err := s.taskRepo.GetTasksByUserID(ctx, userID)
   if err != nil {
      return nil, err
   }

   return tasks, nil
}

func (s *taskService) DeleteTask(ctx context.Context, taskID, userID int64) error {
   err := s.taskRepo.DeleteTask(ctx, taskID, userID)
   if err != nil {
      return err
   }

   return nil
}
