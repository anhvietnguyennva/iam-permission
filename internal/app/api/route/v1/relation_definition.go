package v1

import (
	errt "github.com/anhvietnguyennva/go-error/pkg/transformer"
	oentity "github.com/anhvietnguyennva/iam-go-sdk/oauth/entity"
	"github.com/gin-gonic/gin"

	"iam-permission/internal/app/api/dto"
	"iam-permission/internal/app/api/middleware"
	"iam-permission/internal/app/api/service"
	"iam-permission/internal/pkg/constant"
	"iam-permission/pkg/context"
)

type relationDefinitionAPI struct {
	relationDefinitionService service.IRelationDefinitionService
}

func newRelationDefinitionAPI() *relationDefinitionAPI {
	return &relationDefinitionAPI{
		relationDefinitionService: service.RelationDefinitionServiceInstance(),
	}
}

func (a *relationDefinitionAPI) SetupAdminRoute(rg *gin.RouterGroup) {
	rg.POST("", middleware.RequireBearerAuthorizationJWT, a.createRelationDefinition)
}

func (a *relationDefinitionAPI) createRelationDefinition(c *gin.Context) {
	ctx := context.New(c)

	var req dto.CreateRelationDefinitionReq
	if err := c.ShouldBindJSON(&req); err != nil {
		apiErr := errt.RestTransformerInstance().ValidationErrToRestAPIErr(err)
		dto.RespondError(c, apiErr)
		return
	}

	accessToken, _ := c.Value(constant.CtxAccessTokenKey).(*oentity.AccessToken)

	voReq := req.ToValueObject()
	voReq.CreatedBy = accessToken.Subject

	definition, domainErr := a.relationDefinitionService.Create(ctx, voReq)
	if domainErr != nil {
		apiErr := errt.RestTransformerInstance().DomainErrToRestAPIErr(domainErr)
		dto.RespondError(c, apiErr)
		return
	}
	dto.RespondSuccess(c, new(dto.CreateRelationDefinitionResponse).FromEntity(definition))
}
