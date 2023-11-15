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

type subjectRelationTupleAPI struct {
	subjectRelationTupleService service.ISubjectRelationTupleService
}

func newSubjectRelationTupleAPI() *subjectRelationTupleAPI {
	return &subjectRelationTupleAPI{
		subjectRelationTupleService: service.SubjectRelationTupleServiceInstance(),
	}
}

func (a *subjectRelationTupleAPI) SetupAdminRoute(rg *gin.RouterGroup) {
	rg.POST("", middleware.RequireBearerAuthorizationJWT, a.createSubjectRelationTuple)
}

func (a *subjectRelationTupleAPI) createSubjectRelationTuple(c *gin.Context) {
	ctx := context.New(c)

	var req dto.CreateSubjectRelationTupleReq
	if err := c.ShouldBindJSON(&req); err != nil {
		apiErr := errt.RestTransformerInstance().ValidationErrToRestAPIErr(err)
		dto.RespondError(c, apiErr)
		return
	}

	accessToken, _ := c.Value(constant.CtxAccessTokenKey).(*oentity.AccessToken)

	voReq := req.ToValueObject()
	voReq.CreatedBy = accessToken.Subject

	tuple, domainErr := a.subjectRelationTupleService.Create(ctx, voReq)
	if domainErr != nil {
		apiErr := errt.RestTransformerInstance().DomainErrToRestAPIErr(domainErr)
		dto.RespondError(c, apiErr)
		return
	}
	dto.RespondSuccess(c, new(dto.CreateSubjectRelationTupleRes).FromEntity(tuple))
}
