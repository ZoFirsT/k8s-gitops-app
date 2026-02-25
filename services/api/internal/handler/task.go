package handler

import (
    "context"
    "fmt"
    "net/http"
    "time"

    "github.com/gin-gonic/gin"
    "github.com/google/uuid"
    "github.com/redis/go-redis/v9"

    "github.com/ZoFirsT/k8s-gitops-app/api/internal/model"
)

type TaskHandler struct {
    redis *redis.Client
}

func NewTaskHandler(rdb *redis.Client) *TaskHandler {
    return &TaskHandler{redis: rdb}
}

func (h *TaskHandler) CreateTask(c *gin.Context) {
    var req model.CreateTaskRequest
    if err := c.ShouldBindJSON(&req); err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
        return
    }

    task := model.Task{
        ID:          uuid.New().String(),
        Title:       req.Title,
        Description: req.Description,
        Status:      "pending",
        CreatedAt:   time.Now(),
    }

    // save to redis
    ctx := context.Background()
    key := fmt.Sprintf("task:%s", task.ID)
    h.redis.HSet(ctx, key, map[string]interface{}{
        "id":          task.ID,
        "title":       task.Title,
        "description": task.Description,
        "status":      task.Status,
        "created_at":  task.CreatedAt.Format(time.RFC3339),
    })
    h.redis.Expire(ctx, key, 24*time.Hour)

    c.JSON(http.StatusCreated, task)
}

func (h *TaskHandler) GetTask(c *gin.Context) {
    id := c.Param("id")
    ctx := context.Background()
    key := fmt.Sprintf("task:%s", id)

    data, err := h.redis.HGetAll(ctx, key).Result()
    if err != nil || len(data) == 0 {
        c.JSON(http.StatusNotFound, gin.H{"error": "task not found"})
        return
    }

    c.JSON(http.StatusOK, data)
}

func (h *TaskHandler) HealthCheck(c *gin.Context) {
    ctx := context.Background()
    if err := h.redis.Ping(ctx).Err(); err != nil {
        c.JSON(http.StatusServiceUnavailable, gin.H{
            "status": "unhealthy",
            "redis":  "unreachable",
        })
        return
    }
    c.JSON(http.StatusOK, gin.H{
        "status":  "healthy",
        "service": "api",
        "version": "1.0.0",
    })
}
