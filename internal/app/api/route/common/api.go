package common

import "github.com/gin-gonic/gin"

func RegisterAPI(router gin.IRouter, api IAPI, apiPath string) {
	rg := router.Group(apiPath)
	api.SetupRoute(rg)
}

func RegisterAdminAPI(router gin.IRouter, api IAdminAPI, apiPath string) {
	rg := router.Group(apiPath)
	api.SetupAdminRoute(rg)
}

func RegisterAPIGroup(router gin.IRouter, apiGroup IAPIGroup, apiGroupPath string) {
	rg := router.Group(apiGroupPath)
	apiGroup.RegisterAPIs(rg)
}

func RegisterAdminAPIGroup(router gin.IRouter, apiGroup IAdminAPIGroup, apiGroupPath string) {
	rg := router.Group(apiGroupPath)
	apiGroup.RegisterAdminAPIs(rg)
}
