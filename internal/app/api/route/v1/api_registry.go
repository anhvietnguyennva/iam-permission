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
	common.RegisterAPI(rg, newPermissionAPI(), "/permissions")
}

func (a APIv1) RegisterAdminAPIs(rg *gin.RouterGroup) {
	common.RegisterAdminAPI(rg, newServiceAPI(), "/services")
	common.RegisterAdminAPI(rg, newRelationDefinitionAPI(), "/relation-definitions")
	common.RegisterAdminAPI(rg, newSubjectRelationTupleAPI(), "/subject-relation-tuples")
	common.RegisterAdminAPI(rg, newSubjectSetAPI(), "/subject-sets")
	common.RegisterAdminAPI(rg, newSubjectSetRelationTupleAPI(), "/subject-set-relation-tuples")
}
