package services

import (
	"context"
	"github.com/AntNoHuabei/Remo/pkg/api"
	"github.com/AntNoHuabei/Remo/pkg/persist"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/wailsapp/wails/v3/pkg/application"
	"net"
	"net/http"
	"time"
)

// GinService implements a Wails service that uses Gin for HTTP handling
type GinService struct {
	ginEngine   *gin.Engine
	app         *application.App
	netListener net.Listener
}

// NewGinService creates a new GinService instance
func NewGinService() *GinService {

	persist.InitDB()

	// Create a new Gin router
	ginEngine := gin.New()

	// Add middlewares
	ginEngine.Use(gin.Recovery())
	ginEngine.Use(cors.Default())
	ginEngine.Use(LoggingMiddleware())

	service := &GinService{
		ginEngine: ginEngine,
	}

	// Define routes
	service.setupRoutes()

	return service
}

// ServiceName returns the name of the service
func (s *GinService) ServiceName() string {
	return "Gin API Service"
}

// ServiceStartup is called when the service starts
func (s *GinService) ServiceStartup(ctx context.Context, options application.ServiceOptions) error {
	// You can access the application instance via ctx
	s.app = application.Get()

	s.setupHttpServe()
	return nil
}

// ServiceShutdown is called when the service shuts down
func (s *GinService) ServiceShutdown() error {
	if s.netListener != nil {
		s.netListener.Close()
	}
	// Clean up event handler to prevent memory leaks
	return nil
}

// ServeHTTP implements the http.Handler interface
func (s *GinService) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	// All other requests go to the Gin router
	s.ginEngine.ServeHTTP(w, r)
}

// setupRoutes configures the API routes
func (s *GinService) setupRoutes() {
	sessionGroup := s.ginEngine.Group("/session")
	sessionGroup.POST("/create", api.SessionCreate)
	sessionGroup.POST("/delete", api.SessionDelete)
	sessionGroup.POST("/list", api.SessionList)
	sessionGroup.POST("/messages", api.SessionMessages)
	s.ginEngine.POST("/chat", api.Chat)
}

// setupHttpServe 由于wails里面无法正常使用sse
func (s *GinService) setupHttpServe() {

	// 创建 TCP listener
	listener, err := net.Listen("tcp", ":9980")
	if err != nil {
		s.app.Logger.Error("Error creating listener", "error", err)
		return
	}
	s.netListener = listener
	go func() {
		s.ginEngine.RunListener(listener)
	}()
}

// LoggingMiddleware is a Gin middleware that logs request details
func LoggingMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		// Start timer
		start := time.Now()

		// Process request
		c.Next()

		// Calculate latency
		latency := time.Since(start)

		// Log request details
		statusCode := c.Writer.Status()
		clientIP := c.ClientIP()
		method := c.Request.Method
		path := c.Request.URL.Path

		// Get the application instance
		app := application.Get()
		if app != nil {
			app.Logger.Info("HTTP Request",
				"status", statusCode,
				"method", method,
				"path", path,
				"ip", clientIP,
				"latency", latency,
			)
		}
	}
}
