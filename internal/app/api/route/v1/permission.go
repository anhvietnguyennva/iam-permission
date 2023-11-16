package v1

import (
	errt "github.com/anhvietnguyennva/go-error/pkg/transformer"
	"github.com/gin-gonic/gin"

	"iam-permission/internal/app/api/dto"
	"iam-permission/internal/app/api/service"
	"iam-permission/pkg/context"
)

type permissionAPI struct {
	permissionService service.IPermissionService
}

func newPermissionAPI() *permissionAPI {
	return &permissionAPI{
		permissionService: service.PermissionServiceInstance(),
	}
}

func (a *permissionAPI) SetupRoute(rg *gin.RouterGroup) {
	rg.GET("/check", a.checkPermission)
}

func (a *permissionAPI) checkPermission(c *gin.Context) {
	ctx := context.New(c)

	var req dto.CheckPermissionReq
	if err := c.ShouldBindQuery(&req); err != nil {
		apiErr := errt.RestTransformerInstance().ValidationErrToRestAPIErr(err)
		dto.RespondError(c, apiErr)
		return
	}

	allowed, domainErr := a.permissionService.CheckPermission(ctx, req.ToValueObject())
	if domainErr != nil {
		apiErr := errt.RestTransformerInstance().DomainErrToRestAPIErr(domainErr)
		dto.RespondError(c, apiErr)
		return
	}
	dto.RespondSuccess(c, &dto.CheckPermissionRes{Allowed: allowed})
}
