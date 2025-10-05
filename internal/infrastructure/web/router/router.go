package router

import (
	"context"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"katseye/internal/infrastructure/web/handlers"
)

type Handlers struct {
	Product  *handlers.ProductHandler
	Partner  *handlers.PartnerHandler
	Address  *handlers.AddressHandler
	Consumer *handlers.ConsumerHandler
	Auth     *handlers.AuthHandler
}

type Server struct {
	engine *gin.Engine
	httpServer *http.Server
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

	address := cfg.Port

	switch {
	case address == "":
		address = ":8080"
	case !strings.HasPrefix(address, ":"):
		address = ":" + address
	}

	httpServer := &http.Server{
		Addr: address,
		Handler: engine,
	}

	return &Server{engine: engine, httpServer: httpServer}
}

func (s *Server) Run() error {
	return s.httpServer.ListenAndServe()
}

func (s *Server) Shutdown(ctx context.Context) error {
	return s.httpServer.Shutdown(ctx)
}

func (s *Server) Engine() *gin.Engine {
	return s.engine
}
