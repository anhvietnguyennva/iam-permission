package common

import "github.com/gin-gonic/gin"

type IAPI interface {
	SetupRoute(*gin.RouterGroup)
}

type IAdminAPI interface {
	SetupAdminRoute(group *gin.RouterGroup)
}

type IAPIGroup interface {
	RegisterAPIs(*gin.RouterGroup)
}

type IAdminAPIGroup interface {
	RegisterAdminAPIs(*gin.RouterGroup)
}
