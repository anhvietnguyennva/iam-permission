package route

import (
	"github.com/gin-gonic/gin"
)

type IAPI interface {
	setupRoute(rg *gin.RouterGroup)
}

type IAdminAPI interface {
	setupAdminRoute(rg *gin.RouterGroup)
}

func AddPublicRouterV1(engine *gin.Engine) {
	// API v1
	v1 := engine.Group("/api/v1")
	addPublicApi(newHealthApi(), "/health", v1)
}

func addPublicApi(api IAPI, path string, rg *gin.RouterGroup) {
	publicApiRg := rg.Group(path)
	api.setupRoute(publicApiRg)
}

func AddAdminRouterV1(engine *gin.Engine) {
	// admin API v1
	v1 := engine.Group("/api/v1/admin")
	addAdminApi(newHealthApi(), "/health", v1)
}

func addAdminApi(api IAdminAPI, path string, rg *gin.RouterGroup) {
	apiRg := rg.Group(path)
	api.setupAdminRoute(apiRg)
}
