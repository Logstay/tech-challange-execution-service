package http

import (
	"context"
	"fmt"
	"net/http"
	"time"

	"github.com/Logstay/tech-challange-wherehouse-service/config"
	_ "github.com/Logstay/tech-challange-wherehouse-service/docs"
	services "github.com/Logstay/tech-challange-wherehouse-service/internal/interfaces/services"
	"github.com/Logstay/tech-challange-wherehouse-service/pkg/constant"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

type Server struct {
	handler *Handler
	router  *gin.Engine
	cfg     *config.Config
}

func NewServer(ctx context.Context, services services.ServiceFactory, cfg *config.Config) *Server {
	handler := NewHandler(services)

	server := &Server{
		handler: handler,
		cfg:     cfg,
	}

	server.setupRouter()
	return server
}

func (s *Server) setupRouter() {
	gin.SetMode(gin.ReleaseMode)
	router := gin.New()

	router.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"*"},
		AllowMethods:     []string{"GET", "POST", "PUT", "PATCH", "DELETE", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "Content-Type", "Accept", "Authorization", "X-Request-ID"},
		ExposeHeaders:    []string{"Content-Length", "X-Request-ID"},
		AllowCredentials: true,
	}))

	router.Use(gin.Recovery())

	v1 := router.Group("/v1")
	v1.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))
	v1.GET("/health", s.handler.HealthCheck)

	{

		quotations := v1.Group("/quotations")
		{
			quotations.GET("/:id", s.handler.GetQuotation)
		}
	}

	s.router = router
}

func (s *Server) Start() error {
	addr := fmt.Sprintf("%s:%s", s.cfg.Server.Host, s.cfg.Server.Port)

	server := &http.Server{
		Addr:         addr,
		Handler:      s.router,
		WriteTimeout: time.Second * constant.ServerWriteReadTimeout,
		ReadTimeout:  time.Second * constant.ServerWriteReadTimeout,
		IdleTimeout:  time.Second * constant.ServerIdleTimeout,
	}

	return server.ListenAndServe()
}

func (s *Server) GetRouter() *gin.Engine {
	return s.router
}
