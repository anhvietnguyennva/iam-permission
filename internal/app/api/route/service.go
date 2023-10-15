package route

import (
	"strings"

	"github.com/gin-gonic/gin"

	"iam-permission/internal/app/api/dto"
	"iam-permission/internal/app/api/service"
	"iam-permission/pkg/context"
)

type serviceAPI struct {
	serviceService service.IServiceService
}

func newServiceAPI(serviceService service.IServiceService) *serviceAPI {
	return &serviceAPI{
		serviceService: serviceService,
	}
}

func (a *serviceAPI) setupRoute(rg *gin.RouterGroup) {
	rg.POST("", a.createService)
}

func (a *serviceAPI) createService(c *gin.Context) {
	ctx := context.New(c)

	var req dto.CreateWatchedWalletReq
	if err := c.ShouldBindJSON(&req); err != nil {
		apiErr := kst.RestTransformerInstance().ValidationErrToRestAPIErr(err)
		dto2.RespondError(c, apiErr)
		return
	}

	watchedWallet, domainErr := a.watchedWalletService.Create(ctx, strings.ToLower(req.WalletAddress))
	if domainErr != nil {
		apiErr := kst.RestTransformerInstance().DomainErrToRestAPIErr(domainErr)
		dto2.RespondError(c, apiErr)
		return
	}
	dto2.RespondSuccess(c, new(dto.CreateWatchedWalletRes).FromEntity(watchedWallet))
}
