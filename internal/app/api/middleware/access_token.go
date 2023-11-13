package middleware

import (
	"github.com/anhvietnguyennva/go-error/pkg/errors"
	sdk "github.com/anhvietnguyennva/iam-go-sdk"
	"github.com/gin-gonic/gin"

	"iam-permission/internal/app/api/dto"
	"iam-permission/internal/pkg/constant"
)

func RequireBearerAuthorizationJWT(c *gin.Context) {
	authorization := c.GetHeader("Authorization")
	token, err := sdk.SDK().ParseBearerJWT(authorization)
	if err != nil {
		dto.RespondError(c, errors.NewRestAPIErrUnauthenticated(err))
		return
	}
	c.Set(constant.CtxAccessTokenKey, token)

	c.Next()
}
