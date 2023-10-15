package builder

import (
	"github.com/gin-gonic/gin"

	"iam-permission/internal/app/api/route"
	"iam-permission/internal/pkg/config"
)

type apiServer struct {
	engine *gin.Engine
}

func NewAPIServer(admin bool) IRunner {
	engine := newEngine()

	if admin {
		route.AddAdminRouterV1(engine)
	} else {
		route.AddPublicRouterV1(engine)
	}

	return &apiServer{engine: engine}
}

func (f *apiServer) Run() error {
	return f.engine.Run(config.Instance().Http.BindAddress)
}
