package router

import (
	"strings"

	"github.com/gin-gonic/gin"
	"katseye/internal/infrastructure/web/handlers"
)

type Handlers struct {
	Product *handlers.ProductHandler
	Partner *handlers.PartnerHandler
	Address *handlers.AddressHandler
	Auth    *handlers.AuthHandler
}

type Server struct {
	engine *gin.Engine
	port   string
}

func New(cfg Config) *Server {
	if mode := strings.TrimSpace(cfg.Mode); mode != "" {
		gin.SetMode(mode)
	}

	engine := gin.Default()

	if len(cfg.Middlewares) > 0 {
		engine.Use(cfg.Middlewares...)
	}

	ConfigureRoutes(engine, cfg.Handlers)

	return &Server{engine: engine, port: cfg.Port}
}

func (s *Server) Run() error {
	address := s.port

	switch {
	case address == "":
		address = ":8080"
	case !strings.HasPrefix(address, ":"):
		address = ":" + address
	}

	return s.engine.Run(address)
}

func (s *Server) Engine() *gin.Engine {
	return s.engine
}
