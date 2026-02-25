package model

import "time"

type Task struct {
    ID          string    `json:"id"`
    Title       string    `json:"title"`
    Description string    `json:"description"`
    Status      string    `json:"status"` // pending, processing, done
    CreatedAt   time.Time `json:"created_at"`
}

type CreateTaskRequest struct {
    Title       string `json:"title" binding:"required"`
    Description string `json:"description"`
}
