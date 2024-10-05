package services

import (
	"context"
	"errors"

	"taskmaster/internal/models"
	"taskmaster/internal/repositories"
)

type TaskService interface {
   CreateTask(ctx context.Context, task *models.Task) error
   GetTasksByUserID(ctx context.Context, userID int64) ([]models.Task, error)
}

type taskService struct {
   taskRepo repositories.TaskRepository
}

func NewTaskService(repo repositories.TaskRepository) TaskService {
   return &taskService{taskRepo: repo}
}

func (s *taskService) CreateTask(ctx context.Context, task *models.Task) error {
   // Business rules go here, such as checking task validity
   if task.Title == "" {
      return errors.New("Task title is required")
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
