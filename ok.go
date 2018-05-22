package main

import (
	"time"

	"os"

	"github.com/gin-gonic/gin"
	"github.com/nu7hatch/gouuid"
	"github.com/gin-contrib/zap"
	"go.uber.org/zap"
)

func main() {
	count := 1
	callUuidV4, _ := uuid.NewV4()

	nodeName := os.Getenv("NODE_NAME")
	podName := os.Getenv("POD_NAME")
	podNamespace := os.Getenv("POD_NAMESPACE")
	podIP := os.Getenv("POD_IP")
	serviceAccountName := os.Getenv("SERVICE_ACCOUNT")
	port, ok := os.LookupEnv("PORT")
	if ok != true {
		port = "8080"
	}

	gin.SetMode(gin.ReleaseMode)

	r := gin.New()

	logger, _ := zap.NewProduction()
	r.Use(ginzap.Ginzap(logger, time.RFC3339, true))

	r.GET("/", func(c *gin.Context) {
		instanceUuidV4, _ := uuid.NewV4()

		c.JSON(200, gin.H{
			"message":         "ok",
			"time":            time.Now(),
			"count":           count,
			"uuid_call":       instanceUuidV4.String(),
			"uuid_instance":   callUuidV4.String(),
			"client_ip":       c.ClientIP(),
			"version":         2,
			"version_msg":     "version 2",
			"node_name":       nodeName,
			"pod_name":        podName,
			"pod_namespace":   podNamespace,
			"pod_ip":          podIP,
			"service_account": serviceAccountName,
		})
		count++
	})
	r.Run(":" + port)
}
