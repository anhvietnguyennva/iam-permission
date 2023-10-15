package builder

import (
	"net/http"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"

	"iam-permission/internal/app/api/middleware"
	"iam-permission/internal/pkg/config"
)

func newEngine() *gin.Engine {
	gin.SetMode(config.Instance().Http.Mode)
	newEngine := gin.New()
	newEngine.Use(gin.LoggerWithWriter(gin.DefaultWriter, "/health"))
	newEngine.Use(gin.Recovery())
	newEngine.Use(middleware.Logger())
	setCORS(newEngine)

	return newEngine
}

func setCORS(engine *gin.Engine) {
	corsConfig := cors.DefaultConfig()
	corsConfig.AddAllowMethods(http.MethodOptions)
	corsConfig.AllowAllOrigins = true
	engine.Use(cors.New(corsConfig))
}
