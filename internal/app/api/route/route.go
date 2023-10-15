package route

import (
	"github.com/gin-gonic/gin"

	"iam-permission/internal/app/api/route/common"
	"iam-permission/internal/app/api/route/v1"
)

func Register(engine *gin.Engine, isAdmin bool) {
	common.RegisterAPI(engine, newHealthAPI(), "health")
	if isAdmin {
		common.RegisterAdminAPIGroup(engine, v1.New(), "/admin/api/v1")
	} else {
		common.RegisterAPIGroup(engine, v1.New(), "/api/v1")
	}
}
