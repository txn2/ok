package main

import (
	"net/http"
	"os"
	"time"

	ginzap "github.com/gin-contrib/zap"
	"github.com/gin-gonic/gin"
	uuid "github.com/nu7hatch/gouuid"
	"go.uber.org/zap"
)

var (
	nodeName     = getEnv("NODE_NAME", "undefined")
	podName      = getEnv("POD_NAME", "undefined")
	podNamespace = getEnv("POD_NAMESPACE", "undefined")
	podIP        = getEnv("POD_IP", "undefined")
	port         = getEnv("PORT", "8080")
	ip           = getEnv("IP", "127.0.0.1")
	message      = getEnv("MESSAGE", "ok")
)

var Version = "0.0.0"
var Service = "ok"

func main() {
	count := 1
	callUuidV4, _ := uuid.NewV4()

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

	logger.Info("Starting "+Service+" API Server",
		zap.String("version", Version),
		zap.String("type", "server_startup"),
		zap.String("port", port),
		zap.String("ip", ip),
	)

	s := &http.Server{
		Addr:           ip + ":" + port,
		Handler:        r,
		ReadTimeout:    10 * time.Second,
		WriteTimeout:   10 * time.Second,
		MaxHeaderBytes: 1 << 20, // 1 MB
	}

	err := s.ListenAndServe()
	if err != nil {
		logger.Fatal(err.Error())
	}
}

// getEnv gets an environment variable or sets a default if
// one does not exist.
func getEnv(key, fallback string) string {
	value := os.Getenv(key)
	if len(value) == 0 {
		return fallback
	}

	return value
}
