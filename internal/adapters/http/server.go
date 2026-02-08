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

// @title			Tech Challenge - Wherehouse Service
// @version		1.0.0
// @description	Microsserviço de estoque e catálogo de serviços do Tech Challenge - Fase 4 do 12SOAT FIAP.
// @description
// @description				## 🚀 Funcionalidades
// @description				- **Consumíveis**: CRUD completo de peças e itens de estoque
// @description				- **Definições de Serviço**: Catálogo de serviços disponíveis com precificação
// @description				- **Integração via Eventos**: Consumo de eventos Kafka para atualização de estoque
//
// @contact.name				Kaique Moreira
// @contact.email				kaiquepng@gmail.com
// @contact.url				https://github.com/Logstay/tech-challange-wherehouse-service
//
// @license.name				MIT License
// @license.url				https://opensource.org/licenses/MIT
//
// @BasePath					/v1
// @securityDefinitions.apikey	BearerAuth
// @in							header
// @name						Authorization
// @description				Exemplo: 'eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9...'
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

		consumables := v1.Group("/consumables")
		{
			consumables.POST("", s.handler.SaveConsumable)
			consumables.GET("/:uid", s.handler.GetConsumable)
			consumables.GET("", s.handler.ListConsumable)
			consumables.PUT("/:uid", s.handler.UpdateConsumable)
			consumables.DELETE("/:uid", s.handler.DeleteConsumable)
		}

		serviceDefinitions := v1.Group("/service-definitions")
		{
			serviceDefinitions.POST("", s.handler.SaveServiceDefinition)
			serviceDefinitions.GET("/:uid", s.handler.GetServiceDefinition)
			serviceDefinitions.GET("", s.handler.ListServiceDefinition)
			serviceDefinitions.PUT("/:uid", s.handler.UpdateServiceDefinition)
			serviceDefinitions.DELETE("/:uid", s.handler.DeleteServiceDefinition)
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
