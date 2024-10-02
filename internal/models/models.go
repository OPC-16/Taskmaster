package models

import "time"

type User struct {
   ID       int64  `json:"id"`
   Username string `json:"username"`
   Email    string `json:"email"`
   Password string `json:"-"` // Password is not exposed in JSON
}

type Task struct {
   ID          int64     `json:"id"`
   UserID      int64     `json:"user_id"`
   Title       string    `json:"title"`
   Description string    `json:"description"`
   Category    string    `json:"category"`
   Deadline    time.Time `json:"deadline"`
   Priority    int       `json:"priority"`
   Status      string    `json:"status"`
   CreatedAt   time.Time `json:"created_at"`
   UpdatedAt   time.Time `json:"updated_at"`
}
