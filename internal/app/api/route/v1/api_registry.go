package v1

import (
	"github.com/gin-gonic/gin"

	"iam-permission/internal/app/api/route/common"
)

type APIv1 struct{}

func New() *APIv1 {
	return &APIv1{}
}

func (a APIv1) RegisterAPIs(rg *gin.RouterGroup) {
}

func (a APIv1) RegisterAdminAPIs(rg *gin.RouterGroup) {
	common.RegisterAdminAPI(rg, newServiceAPI(), "/services")
	common.RegisterAdminAPI(rg, newSubjectRelationTupleAPI(), "/subject-relation-tuples")
}
