package dto

import (
	"net/http"

	errc "github.com/anhvietnguyennva/go-error/pkg/constant"
	"github.com/anhvietnguyennva/go-error/pkg/errors"
	"github.com/gin-gonic/gin"

	"iam-permission/internal/pkg/constant"
	"iam-permission/internal/pkg/util/logger"
)

type Response struct {
	Code    int         `json:"code"`
	Message string      `json:"message"`
	Data    interface{} `json:"data,omitempty"`
}

func RespondSuccess(c *gin.Context, data interface{}) {
	c.AbortWithStatusJSON(http.StatusOK, Response{
		Code:    errc.ClientErrCodeOK,
		Message: errc.ClientErrMsgOK,
		Data:    data,
	})
}

func RespondError(c *gin.Context, apiErr *errors.RestAPIError) {
	traceID, ok := c.Value(constant.CtxTraceIDKey).(string)
	if ok && apiErr != nil {
		apiErr.Details = append(apiErr.Details, struct {
			TraceID string `json:"traceId"`
		}{
			TraceID: traceID,
		})
	}
	logger.Error(c, apiErr)
	c.AbortWithStatusJSON(apiErr.HttpStatus, apiErr)
}
