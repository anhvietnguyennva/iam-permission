package route

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

type healthApi struct{}

func newHealthApi() *healthApi {
	return &healthApi{}
}

func (a *healthApi) setupRoute(rg *gin.RouterGroup) {
	rg.GET("", func(c *gin.Context) {
		c.AbortWithStatusJSON(http.StatusOK, "OK")
	})
}

func (a *healthApi) setupAdminRoute(rg *gin.RouterGroup) {
	rg.GET("", func(c *gin.Context) {
		c.AbortWithStatusJSON(http.StatusOK, "OK")
	})
}
