package v1

import (
	errt "github.com/anhvietnguyennva/go-error/pkg/transformer"
	iamentity "github.com/anhvietnguyennva/iam-go-sdk/pkg/oauth/entity"
	"github.com/gin-gonic/gin"

	"iam-permission/internal/app/api/dto"
	"iam-permission/internal/app/api/middleware"
	"iam-permission/internal/app/api/service"
	"iam-permission/internal/pkg/constant"
	"iam-permission/pkg/context"
)

type serviceAPI struct {
	serviceService service.IServiceService
}

func newServiceAPI() *serviceAPI {
	return &serviceAPI{
		serviceService: service.ServiceServiceInstance(),
	}
}

func (a *serviceAPI) SetupAdminRoute(rg *gin.RouterGroup) {
	rg.POST("", middleware.RequireBearerAuthorizationJWT, a.createService)
}

func (a *serviceAPI) createService(c *gin.Context) {
	ctx := context.New(c)

	var req dto.CreateServiceRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		apiErr := errt.RestTransformerInstance().ValidationErrToRestAPIErr(err)
		dto.RespondError(c, apiErr)
		return
	}

	accessToken, _ := c.Value(constant.CtxAccessTokenKey).(iamentity.AccessToken)

	voReq := req.ToValueObject()
	voReq.CreatedBy = accessToken.Subject
	s, domainErr := a.serviceService.Create(ctx, voReq)
	if domainErr != nil {
		apiErr := errt.RestTransformerInstance().DomainErrToRestAPIErr(domainErr)
		dto.RespondError(c, apiErr)
		return
	}
	dto.RespondSuccess(c, new(dto.CreateServiceResponse).FromEntity(s))
}
