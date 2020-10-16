package main

import (
	"time"

	"os"

	ginzap "github.com/gin-contrib/zap"
	"github.com/gin-gonic/gin"
	uuid "github.com/nu7hatch/gouuid"
	"go.uber.org/zap"
)

func main() {
	count := 1
	callUuidV4, _ := uuid.NewV4()

	nodeName := os.Getenv("NODE_NAME")
	podName := os.Getenv("POD_NAME")
	podNamespace := os.Getenv("POD_NAMESPACE")
	podIP := os.Getenv("POD_IP")

	port, ok := os.LookupEnv("PORT")
	if ok != true {
		port = "8080"
	}

	message, ok := os.LookupEnv("MESSAGE")
	if ok != true {
		port = "ok"
	}

	gin.SetMode(gin.ReleaseMode)

	r := gin.New()

	logger, _ := zap.NewProduction()
	r.Use(ginzap.Ginzap(logger, time.RFC3339, true))

	r.GET("/", func(c *gin.Context) {
		instanceUuidV4, _ := uuid.NewV4()

		c.JSON(200, gin.H{
			"message":       message,
			"time":          time.Now(),
			"count":         count,
			"uuid_call":     instanceUuidV4.String(),
			"uuid_instance": callUuidV4.String(),
			"client_ip":     c.ClientIP(),
			"node_name":     nodeName,
			"pod_name":      podName,
			"pod_namespace": podNamespace,
			"pod_ip":        podIP,
			"pod_port":      port,
		})
		count++
	})
	r.Run(":" + port)
}
