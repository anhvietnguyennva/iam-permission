package builder

import (
	"fmt"
	"iam-permission/internal/app/api/middleware"
	"net/http"
	"testing"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

func TestBuildEngine(t *testing.T) {
	name := "TestBuildEngine"
	t.Log(name)

	engine := newEngine()
	handlers := engine.Handlers

	corsConfig := cors.DefaultConfig()
	corsConfig.AddAllowMethods(http.MethodOptions)
	corsConfig.AllowAllOrigins = true

	assert.EqualValues(t, 4, len(handlers))
	assert.EqualValues(t, fmt.Sprintf("%p", gin.LoggerWithWriter(gin.DefaultWriter, "/api/v1/health")), fmt.Sprintf("%p", handlers[0]))
	assert.EqualValues(t, fmt.Sprintf("%p", gin.Recovery()), fmt.Sprintf("%p", handlers[1]))
	assert.EqualValues(t, fmt.Sprintf("%p", middleware.Logger()), fmt.Sprintf("%p", handlers[2]))
	assert.EqualValues(t, fmt.Sprintf("%p", cors.New(corsConfig)), fmt.Sprintf("%p", handlers[3]))
}
