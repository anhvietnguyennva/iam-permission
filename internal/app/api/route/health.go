package route

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

type HealthAPI struct{}

func newHealthAPI() *HealthAPI {
	return &HealthAPI{}
}

func (a *HealthAPI) SetupRoute(rg *gin.RouterGroup) {
	rg.GET("", func(c *gin.Context) {
		c.AbortWithStatusJSON(http.StatusOK, "OK")
	})
}
