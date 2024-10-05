package repositories

import (
	"context"
	"database/sql"

	"taskmaster/internal/models"
)

type TaskRepository interface {
   CreateTask(ctx context.Context, task *models.Task) error
   GetTasksByUserID(ctx context.Context, userID int64) ([]models.Task, error)
}

type taskRepository struct {
   db *sql.DB
}

func NewTaskRepository(db *sql.DB) TaskRepository {
   return &taskRepository{db: db}
}

func (r *taskRepository) CreateTask(ctx context.Context, task *models.Task) error {
   query := "INSERT INTO tasks (title, description, status, priority, deadline, user_id) VALUES (?, ?, ?, ?, ?, ?)"
   _, err := r.db.ExecContext(ctx, query, task.Title, task.Description, task.Status, task.Priority, task.Deadline, task.UserID)
   return err
}

func (r *taskRepository) GetTasksByUserID(ctx context.Context, userID int64) ([]models.Task, error) {
   var tasks []models.Task
   query := "SELECT id, title, description, category, deadline, priority, status, created_at, updated_at FROM tasks WHERE user_id = ?"
   rows, err := r.db.QueryContext(ctx, query, userID)
   if err != nil {
      return nil, err
   }
   defer rows.Close()

   for rows.Next() {
      var task models.Task
      if err := rows.Scan(&task.ID, &task.Title, &task.Description, &task.Category, &task.Deadline,
         &task.Priority, &task.Status, &task.CreatedAt, &task.UpdatedAt); err != nil {
         return nil, err
      }
      tasks = append(tasks, task)
   }

   return tasks, nil
}
