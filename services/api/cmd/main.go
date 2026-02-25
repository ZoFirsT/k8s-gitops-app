package main

import (
    "log"
    "os"

    "github.com/gin-gonic/gin"
    "github.com/prometheus/client_golang/prometheus/promhttp"
    "github.com/redis/go-redis/v9"

    "github.com/ZoFirsT/k8s-gitops-app/api/internal/handler"
    "github.com/ZoFirsT/k8s-gitops-app/api/internal/middleware"
)

func main() {
    // Redis connection
    redisAddr := os.Getenv("REDIS_ADDR")
    if redisAddr == "" {
        redisAddr = "localhost:6379"
    }

    rdb := redis.NewClient(&redis.Options{
        Addr: redisAddr,
    })

    // Setup router
    r := gin.New()
    r.Use(gin.Recovery())
    r.Use(middleware.PrometheusMetrics())

    // Handlers
    taskHandler := handler.NewTaskHandler(rdb)

    // Routes
    r.GET("/health", taskHandler.HealthCheck)
    r.GET("/metrics", gin.WrapH(promhttp.Handler()))

    v1 := r.Group("/api/v1")
    {
        v1.POST("/tasks", taskHandler.CreateTask)
        v1.GET("/tasks/:id", taskHandler.GetTask)
    }

    port := os.Getenv("PORT")
    if port == "" {
        port = "8080"
    }

    log.Printf("🚀 API Service starting on port %s", port)
    if err := r.Run(":" + port); err != nil {
        log.Fatal(err)
    }
}
