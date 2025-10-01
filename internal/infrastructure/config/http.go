package config

import (
	"github.com/gin-gonic/gin"
	webrouter "katseye/internal/infrastructure/web/router"
)

func buildHTTPServer(httpCfg HTTPConfig, handlers webrouter.Handlers, middlewares []gin.HandlerFunc) *webrouter.Server {
	return webrouter.New(webrouter.Config{
		Port:        httpCfg.Port,
		Mode:        httpCfg.GinMode,
		Handlers:    handlers,
		Middlewares: middlewares,
	})
}
