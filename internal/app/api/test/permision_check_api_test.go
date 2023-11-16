package test

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"

	errc "github.com/anhvietnguyennva/go-error/pkg/constant"

	"iam-permission/internal/pkg/constant"
	"iam-permission/internal/pkg/db"
	"iam-permission/internal/pkg/repository/postgres"
)

func (suite *TestSuite) TestAPI_CheckPermission_Successfully_1() {
	name := "TestAPI_CheckPermission_Successfully: subjectID is allowed via subject_relation_tuples"
	suite.T().Log(name)

	// mock
	service := &postgres.Service{
		Namespace:   "iam-permission",
		Description: "iam-permission",
	}
	if err := db.Instance().Create(service).Error; err != nil {
		panic(err)
	}

	relationDefinition := &postgres.RelationDefinition{
		ServiceID:   service.ID,
		Namespace:   service.Namespace,
		Relation:    constant.RelationViewer,
		Description: constant.RelationViewer,
	}
	if err := db.Instance().Create(relationDefinition).Error; err != nil {
		panic(err)
	}

	subjectRelationTuple := &postgres.SubjectRelationTuple{
		ServiceID:            service.ID,
		RelationDefinitionID: relationDefinition.ID,
		Namespace:            service.Namespace,
		Object:               "object 1",
		Relation:             relationDefinition.Relation,
		SubjectID:            "subject 1",
	}
	if err := db.Instance().Create(subjectRelationTuple).Error; err != nil {
		panic(err)
	}

	// call API
	method := "GET"
	endpoint := "/api/v1/permissions/check"
	req, _ := http.NewRequest(method, endpoint, nil)
	q := req.URL.Query()
	q.Add("namespace", "iam-permission")
	q.Add("object", "object 1")
	q.Add("relation", "viewer")
	q.Add("subjectId", "subject 1")
	req.URL.RawQuery = q.Encode()
	w := httptest.NewRecorder()
	testPublicRouter.ServeHTTP(w, req)
	responseData, _ := io.ReadAll(w.Body)

	// assert API response
	suite.EqualValues(http.StatusOK, w.Code)
	expectedResBody := map[string]any{
		"code":    float64(errc.ClientErrCodeOK),
		"message": errc.ClientErrMsgOK,
		"data": map[string]any{
			"allowed": true,
		},
	}
	actualResBody := map[string]any{}
	_ = json.Unmarshal(responseData, &actualResBody)
	suite.Equal(expectedResBody, actualResBody)
}

func (suite *TestSuite) TestAPI_CheckPermission_Successfully_2() {
	name := "TestAPI_CheckPermission_Successfully: subjectID is allowed via subject_sets"
	suite.T().Log(name)

	// mock
	service := &postgres.Service{
		Namespace:   "iam-permission",
		Description: "iam-permission",
	}
	if err := db.Instance().Create(service).Error; err != nil {
		panic(err)
	}

	relationViewer := &postgres.RelationDefinition{
		ServiceID:   service.ID,
		Namespace:   service.Namespace,
		Relation:    constant.RelationViewer,
		Description: constant.RelationViewer,
	}
	if err := db.Instance().Create(relationViewer).Error; err != nil {
		panic(err)
	}
	relationOperator := &postgres.RelationDefinition{
		ServiceID:   service.ID,
		Namespace:   service.Namespace,
		Relation:    "operator",
		Description: "operator",
	}
	if err := db.Instance().Create(relationOperator).Error; err != nil {
		panic(err)
	}

	subjectSet := &postgres.SubjectSet{
		ServiceID:            service.ID,
		RelationDefinitionID: relationOperator.ID,
		Namespace:            service.Namespace,
		Object:               "services/iam-permission",
		Relation:             relationOperator.Relation,
	}
	if err := db.Instance().Create(subjectSet).Error; err != nil {
		panic(err)
	}

	subjectSetRelationTuple := &postgres.SubjectSetRelationTuple{
		ServiceID:            service.ID,
		RelationDefinitionID: relationViewer.ID,
		SubjectSetID:         subjectSet.ID,
		Namespace:            service.Namespace,
		Object:               "object 1",
		Relation:             relationViewer.Relation,
		SubjectSetNamespace:  subjectSet.Namespace,
		SubjectSetObject:     subjectSet.Object,
		SubjectSetRelation:   subjectSet.Relation,
	}
	if err := db.Instance().Create(subjectSetRelationTuple).Error; err != nil {
		panic(err)
	}

	subjectRelationTuple := &postgres.SubjectRelationTuple{
		ServiceID:            service.ID,
		RelationDefinitionID: subjectSet.RelationDefinitionID,
		Namespace:            service.Namespace,
		Object:               subjectSet.Object,
		Relation:             subjectSet.Relation,
		SubjectID:            "subject 1",
	}
	if err := db.Instance().Create(subjectRelationTuple).Error; err != nil {
		panic(err)
	}

	// call API
	method := "GET"
	endpoint := "/api/v1/permissions/check"
	req, _ := http.NewRequest(method, endpoint, nil)
	q := req.URL.Query()
	q.Add("namespace", "iam-permission")
	q.Add("object", "object 1")
	q.Add("relation", "viewer")
	q.Add("subjectId", "subject 1")
	req.URL.RawQuery = q.Encode()
	w := httptest.NewRecorder()
	testPublicRouter.ServeHTTP(w, req)
	responseData, _ := io.ReadAll(w.Body)

	// assert API response
	suite.EqualValues(http.StatusOK, w.Code)
	expectedResBody := map[string]any{
		"code":    float64(errc.ClientErrCodeOK),
		"message": errc.ClientErrMsgOK,
		"data": map[string]any{
			"allowed": true,
		},
	}
	actualResBody := map[string]any{}
	_ = json.Unmarshal(responseData, &actualResBody)
	suite.Equal(expectedResBody, actualResBody)
}

func (suite *TestSuite) TestAPI_CheckPermission_Successfully_3() {
	name := "TestAPI_CheckPermission_Successfully: subjectID is allowed via parent relation"
	suite.T().Log(name)

	// mock
	service := &postgres.Service{
		Namespace:   "iam-permission",
		Description: "iam-permission",
	}
	if err := db.Instance().Create(service).Error; err != nil {
		panic(err)
	}

	relationViewer := &postgres.RelationDefinition{
		ServiceID:   service.ID,
		Namespace:   service.Namespace,
		Relation:    constant.RelationViewer,
		Description: constant.RelationViewer,
	}
	if err := db.Instance().Create(relationViewer).Error; err != nil {
		panic(err)
	}
	relationEditor := &postgres.RelationDefinition{
		ServiceID:   service.ID,
		Namespace:   service.Namespace,
		Relation:    constant.RelationEditor,
		Description: constant.RelationEditor,
	}
	if err := db.Instance().Create(relationEditor).Error; err != nil {
		panic(err)
	}

	relationConfiguration := &postgres.RelationConfiguration{
		ServiceID:                  service.ID,
		Namespace:                  service.Namespace,
		ParentRelationDefinitionID: relationViewer.ID,
		ParentRelation:             relationViewer.Relation,
		ChildRelationDefinitionID:  relationEditor.ID,
		ChildRelation:              relationEditor.Relation,
	}
	if err := db.Instance().Create(relationConfiguration).Error; err != nil {
		panic(err)
	}

	subjectRelationTuple := &postgres.SubjectRelationTuple{
		ServiceID:            service.ID,
		RelationDefinitionID: relationEditor.ID,
		Namespace:            service.Namespace,
		Object:               "object 1",
		Relation:             relationEditor.Relation,
		SubjectID:            "subject 1",
	}
	if err := db.Instance().Create(subjectRelationTuple).Error; err != nil {
		panic(err)
	}

	// call API
	method := "GET"
	endpoint := "/api/v1/permissions/check"
	req, _ := http.NewRequest(method, endpoint, nil)
	q := req.URL.Query()
	q.Add("namespace", "iam-permission")
	q.Add("object", "object 1")
	q.Add("relation", "viewer")
	q.Add("subjectId", "subject 1")
	req.URL.RawQuery = q.Encode()
	w := httptest.NewRecorder()
	testPublicRouter.ServeHTTP(w, req)
	responseData, _ := io.ReadAll(w.Body)

	// assert API response
	suite.EqualValues(http.StatusOK, w.Code)
	expectedResBody := map[string]any{
		"code":    float64(errc.ClientErrCodeOK),
		"message": errc.ClientErrMsgOK,
		"data": map[string]any{
			"allowed": true,
		},
	}
	actualResBody := map[string]any{}
	_ = json.Unmarshal(responseData, &actualResBody)
	suite.Equal(expectedResBody, actualResBody)
}

func (suite *TestSuite) TestAPI_CheckPermission_Successfully_4() {
	name := "TestAPI_CheckPermission_Successfully: subjectID is allowed via subject_sets and children_relation"
	suite.T().Log(name)

	// mock
	service := &postgres.Service{
		Namespace:   "iam-permission",
		Description: "iam-permission",
	}
	if err := db.Instance().Create(service).Error; err != nil {
		panic(err)
	}

	relationViewer := &postgres.RelationDefinition{
		ServiceID:   service.ID,
		Namespace:   service.Namespace,
		Relation:    constant.RelationViewer,
		Description: constant.RelationViewer,
	}
	if err := db.Instance().Create(relationViewer).Error; err != nil {
		panic(err)
	}
	relationOperator := &postgres.RelationDefinition{
		ServiceID:   service.ID,
		Namespace:   service.Namespace,
		Relation:    "operator",
		Description: "operator",
	}
	if err := db.Instance().Create(relationOperator).Error; err != nil {
		panic(err)
	}
	relationAdmin := &postgres.RelationDefinition{
		ServiceID:   service.ID,
		Namespace:   service.Namespace,
		Relation:    "admin",
		Description: "admin",
	}
	if err := db.Instance().Create(relationAdmin).Error; err != nil {
		panic(err)
	}

	relationConfiguration := &postgres.RelationConfiguration{
		ServiceID:                  service.ID,
		Namespace:                  service.Namespace,
		ParentRelationDefinitionID: relationOperator.ID,
		ParentRelation:             relationOperator.Relation,
		ChildRelationDefinitionID:  relationAdmin.ID,
		ChildRelation:              relationAdmin.Relation,
	}
	if err := db.Instance().Create(relationConfiguration).Error; err != nil {
		panic(err)
	}

	subjectSet := &postgres.SubjectSet{
		ServiceID:            service.ID,
		RelationDefinitionID: relationOperator.ID,
		Namespace:            service.Namespace,
		Object:               "services/iam-permission",
		Relation:             relationOperator.Relation,
	}
	if err := db.Instance().Create(subjectSet).Error; err != nil {
		panic(err)
	}

	subjectSetRelationTuple := &postgres.SubjectSetRelationTuple{
		ServiceID:            service.ID,
		RelationDefinitionID: relationViewer.ID,
		SubjectSetID:         subjectSet.ID,
		Namespace:            service.Namespace,
		Object:               "object 1",
		Relation:             relationViewer.Relation,
		SubjectSetNamespace:  subjectSet.Namespace,
		SubjectSetObject:     subjectSet.Object,
		SubjectSetRelation:   subjectSet.Relation,
	}
	if err := db.Instance().Create(subjectSetRelationTuple).Error; err != nil {
		panic(err)
	}

	subjectRelationTuple := &postgres.SubjectRelationTuple{
		ServiceID:            service.ID,
		RelationDefinitionID: relationAdmin.ID,
		Namespace:            service.Namespace,
		Object:               "services/iam-permission",
		Relation:             relationAdmin.Relation,
		SubjectID:            "subject 1",
	}
	if err := db.Instance().Create(subjectRelationTuple).Error; err != nil {
		panic(err)
	}

	// call API
	method := "GET"
	endpoint := "/api/v1/permissions/check"
	req, _ := http.NewRequest(method, endpoint, nil)
	q := req.URL.Query()
	q.Add("namespace", "iam-permission")
	q.Add("object", "object 1")
	q.Add("relation", "viewer")
	q.Add("subjectId", "subject 1")
	req.URL.RawQuery = q.Encode()
	w := httptest.NewRecorder()
	testPublicRouter.ServeHTTP(w, req)
	responseData, _ := io.ReadAll(w.Body)

	// assert API response
	suite.EqualValues(http.StatusOK, w.Code)
	expectedResBody := map[string]any{
		"code":    float64(errc.ClientErrCodeOK),
		"message": errc.ClientErrMsgOK,
		"data": map[string]any{
			"allowed": true,
		},
	}
	actualResBody := map[string]any{}
	_ = json.Unmarshal(responseData, &actualResBody)
	suite.Equal(expectedResBody, actualResBody)
}

func (suite *TestSuite) TestAPI_CheckPermission_Successfully_5() {
	name := "TestAPI_CheckPermission_Successfully: subjectID is allowed via children_relation and subject_sets"
	suite.T().Log(name)

	// mock
	service := &postgres.Service{
		Namespace:   "iam-permission",
		Description: "iam-permission",
	}
	if err := db.Instance().Create(service).Error; err != nil {
		panic(err)
	}

	relationViewer := &postgres.RelationDefinition{
		ServiceID:   service.ID,
		Namespace:   service.Namespace,
		Relation:    constant.RelationViewer,
		Description: constant.RelationViewer,
	}
	if err := db.Instance().Create(relationViewer).Error; err != nil {
		panic(err)
	}
	relationEditor := &postgres.RelationDefinition{
		ServiceID:   service.ID,
		Namespace:   service.Namespace,
		Relation:    constant.RelationEditor,
		Description: constant.RelationEditor,
	}
	if err := db.Instance().Create(relationEditor).Error; err != nil {
		panic(err)
	}
	relationOperator := &postgres.RelationDefinition{
		ServiceID:   service.ID,
		Namespace:   service.Namespace,
		Relation:    "operator",
		Description: "operator",
	}
	if err := db.Instance().Create(relationOperator).Error; err != nil {
		panic(err)
	}

	relationConfiguration := &postgres.RelationConfiguration{
		ServiceID:                  service.ID,
		Namespace:                  service.Namespace,
		ParentRelationDefinitionID: relationViewer.ID,
		ParentRelation:             relationViewer.Relation,
		ChildRelationDefinitionID:  relationEditor.ID,
		ChildRelation:              relationEditor.Relation,
	}
	if err := db.Instance().Create(relationConfiguration).Error; err != nil {
		panic(err)
	}

	subjectSet := &postgres.SubjectSet{
		ServiceID:            service.ID,
		RelationDefinitionID: relationOperator.ID,
		Namespace:            service.Namespace,
		Object:               "services/iam-permission",
		Relation:             relationOperator.Relation,
	}
	if err := db.Instance().Create(subjectSet).Error; err != nil {
		panic(err)
	}

	subjectSetRelationTuple := &postgres.SubjectSetRelationTuple{
		ServiceID:            service.ID,
		RelationDefinitionID: relationEditor.ID,
		SubjectSetID:         subjectSet.ID,
		Namespace:            service.Namespace,
		Object:               "object 1",
		Relation:             relationEditor.Relation,
		SubjectSetNamespace:  subjectSet.Namespace,
		SubjectSetObject:     subjectSet.Object,
		SubjectSetRelation:   subjectSet.Relation,
	}
	if err := db.Instance().Create(subjectSetRelationTuple).Error; err != nil {
		panic(err)
	}

	subjectRelationTuple := &postgres.SubjectRelationTuple{
		ServiceID:            service.ID,
		RelationDefinitionID: subjectSet.RelationDefinitionID,
		Namespace:            service.Namespace,
		Object:               subjectSet.Object,
		Relation:             subjectSet.Relation,
		SubjectID:            "subject 1",
	}
	if err := db.Instance().Create(subjectRelationTuple).Error; err != nil {
		panic(err)
	}

	// call API
	method := "GET"
	endpoint := "/api/v1/permissions/check"
	req, _ := http.NewRequest(method, endpoint, nil)
	q := req.URL.Query()
	q.Add("namespace", "iam-permission")
	q.Add("object", "object 1")
	q.Add("relation", "viewer")
	q.Add("subjectId", "subject 1")
	req.URL.RawQuery = q.Encode()
	w := httptest.NewRecorder()
	testPublicRouter.ServeHTTP(w, req)
	responseData, _ := io.ReadAll(w.Body)

	// assert API response
	suite.EqualValues(http.StatusOK, w.Code)
	expectedResBody := map[string]any{
		"code":    float64(errc.ClientErrCodeOK),
		"message": errc.ClientErrMsgOK,
		"data": map[string]any{
			"allowed": true,
		},
	}
	actualResBody := map[string]any{}
	_ = json.Unmarshal(responseData, &actualResBody)
	suite.Equal(expectedResBody, actualResBody)
}

func (suite *TestSuite) TestAPI_CheckPermission_Successfully_6() {
	name := "TestAPI_CheckPermission_Successfully: subjectID is not allowed"
	suite.T().Log(name)

	// mock
	service := &postgres.Service{
		Namespace:   "iam-permission",
		Description: "iam-permission",
	}
	if err := db.Instance().Create(service).Error; err != nil {
		panic(err)
	}

	relationDefinition := &postgres.RelationDefinition{
		ServiceID:   service.ID,
		Namespace:   service.Namespace,
		Relation:    constant.RelationViewer,
		Description: constant.RelationViewer,
	}
	if err := db.Instance().Create(relationDefinition).Error; err != nil {
		panic(err)
	}

	subjectRelationTuple := &postgres.SubjectRelationTuple{
		ServiceID:            service.ID,
		RelationDefinitionID: relationDefinition.ID,
		Namespace:            service.Namespace,
		Object:               "object 1",
		Relation:             relationDefinition.Relation,
		SubjectID:            "subject 1",
	}
	if err := db.Instance().Create(subjectRelationTuple).Error; err != nil {
		panic(err)
	}

	// call API
	method := "GET"
	endpoint := "/api/v1/permissions/check"
	req, _ := http.NewRequest(method, endpoint, nil)
	q := req.URL.Query()
	q.Add("namespace", "iam-permission")
	q.Add("object", "object 1")
	q.Add("relation", "viewer")
	q.Add("subjectId", "subject 2")
	req.URL.RawQuery = q.Encode()
	w := httptest.NewRecorder()
	testPublicRouter.ServeHTTP(w, req)
	responseData, _ := io.ReadAll(w.Body)

	// assert API response
	suite.EqualValues(http.StatusOK, w.Code)
	expectedResBody := map[string]any{
		"code":    float64(errc.ClientErrCodeOK),
		"message": errc.ClientErrMsgOK,
		"data": map[string]any{
			"allowed": false,
		},
	}
	actualResBody := map[string]any{}
	_ = json.Unmarshal(responseData, &actualResBody)
	suite.Equal(expectedResBody, actualResBody)
}

func (suite *TestSuite) TestAPI_CheckPermission_Successfully_7() {
	name := "TestAPI_CheckPermission_Successfully: subjectID is not allowed because maxDepth is too small"
	suite.T().Log(name)

	// mock
	service := &postgres.Service{
		Namespace:   "iam-permission",
		Description: "iam-permission",
	}
	if err := db.Instance().Create(service).Error; err != nil {
		panic(err)
	}

	relationViewer := &postgres.RelationDefinition{
		ServiceID:   service.ID,
		Namespace:   service.Namespace,
		Relation:    constant.RelationViewer,
		Description: constant.RelationViewer,
	}
	if err := db.Instance().Create(relationViewer).Error; err != nil {
		panic(err)
	}
	relationEditor := &postgres.RelationDefinition{
		ServiceID:   service.ID,
		Namespace:   service.Namespace,
		Relation:    constant.RelationEditor,
		Description: constant.RelationEditor,
	}
	if err := db.Instance().Create(relationEditor).Error; err != nil {
		panic(err)
	}
	relationOperator := &postgres.RelationDefinition{
		ServiceID:   service.ID,
		Namespace:   service.Namespace,
		Relation:    "operator",
		Description: "operator",
	}
	if err := db.Instance().Create(relationOperator).Error; err != nil {
		panic(err)
	}

	relationConfiguration := &postgres.RelationConfiguration{
		ServiceID:                  service.ID,
		Namespace:                  service.Namespace,
		ParentRelationDefinitionID: relationViewer.ID,
		ParentRelation:             relationViewer.Relation,
		ChildRelationDefinitionID:  relationEditor.ID,
		ChildRelation:              relationEditor.Relation,
	}
	if err := db.Instance().Create(relationConfiguration).Error; err != nil {
		panic(err)
	}

	subjectSet := &postgres.SubjectSet{
		ServiceID:            service.ID,
		RelationDefinitionID: relationOperator.ID,
		Namespace:            service.Namespace,
		Object:               "services/iam-permission",
		Relation:             relationOperator.Relation,
	}
	if err := db.Instance().Create(subjectSet).Error; err != nil {
		panic(err)
	}

	subjectSetRelationTuple := &postgres.SubjectSetRelationTuple{
		ServiceID:            service.ID,
		RelationDefinitionID: relationEditor.ID,
		SubjectSetID:         subjectSet.ID,
		Namespace:            service.Namespace,
		Object:               "object 1",
		Relation:             relationEditor.Relation,
		SubjectSetNamespace:  subjectSet.Namespace,
		SubjectSetObject:     subjectSet.Object,
		SubjectSetRelation:   subjectSet.Relation,
	}
	if err := db.Instance().Create(subjectSetRelationTuple).Error; err != nil {
		panic(err)
	}

	subjectRelationTuple := &postgres.SubjectRelationTuple{
		ServiceID:            service.ID,
		RelationDefinitionID: subjectSet.RelationDefinitionID,
		Namespace:            service.Namespace,
		Object:               subjectSet.Object,
		Relation:             subjectSet.Relation,
		SubjectID:            "subject 1",
	}
	if err := db.Instance().Create(subjectRelationTuple).Error; err != nil {
		panic(err)
	}

	// call API
	method := "GET"
	endpoint := "/api/v1/permissions/check"
	req, _ := http.NewRequest(method, endpoint, nil)
	q := req.URL.Query()
	q.Add("namespace", "iam-permission")
	q.Add("object", "object 1")
	q.Add("relation", "viewer")
	q.Add("subjectId", "subject 1")
	q.Add("maxDepth", "1")
	req.URL.RawQuery = q.Encode()
	w := httptest.NewRecorder()
	testPublicRouter.ServeHTTP(w, req)
	responseData, _ := io.ReadAll(w.Body)

	// assert API response
	suite.EqualValues(http.StatusOK, w.Code)
	expectedResBody := map[string]any{
		"code":    float64(errc.ClientErrCodeOK),
		"message": errc.ClientErrMsgOK,
		"data": map[string]any{
			"allowed": false,
		},
	}
	actualResBody := map[string]any{}
	_ = json.Unmarshal(responseData, &actualResBody)
	suite.Equal(expectedResBody, actualResBody)
}

func (suite *TestSuite) TestAPI_CheckPermission_Failed_1() {
	name := "TestAPI_CheckPermission_Failed: missing namespace"
	suite.T().Log(name)

	// call API
	method := "GET"
	endpoint := "/api/v1/permissions/check"
	req, _ := http.NewRequest(method, endpoint, nil)
	q := req.URL.Query()
	q.Add("namespace", "")
	q.Add("object", "object 1")
	q.Add("relation", "viewer")
	q.Add("subjectId", "subject 2")
	req.URL.RawQuery = q.Encode()
	w := httptest.NewRecorder()
	testPublicRouter.ServeHTTP(w, req)
	responseData, _ := io.ReadAll(w.Body)

	// assert API response
	suite.EqualValues(http.StatusBadRequest, w.Code)
	expectedResBody := map[string]any{
		"code":          float64(errc.ClientErrCodeRequired),
		"message":       fmt.Sprintf("%s: %s", errc.ClientErrMsgRequired, "Namespace"),
		"errorEntities": []any{"Namespace"},
		"details":       any(nil),
	}
	actualResBody := map[string]any{}
	_ = json.Unmarshal(responseData, &actualResBody)
	suite.Equal(expectedResBody, actualResBody)
}

func (suite *TestSuite) TestAPI_CheckPermission_Failed_2() {
	name := "TestAPI_CheckPermission_Failed: missing object"
	suite.T().Log(name)

	// call API
	method := "GET"
	endpoint := "/api/v1/permissions/check"
	req, _ := http.NewRequest(method, endpoint, nil)
	q := req.URL.Query()
	q.Add("namespace", "iam-permission")
	q.Add("object", "")
	q.Add("relation", "viewer")
	q.Add("subjectId", "subject 2")
	req.URL.RawQuery = q.Encode()
	w := httptest.NewRecorder()
	testPublicRouter.ServeHTTP(w, req)
	responseData, _ := io.ReadAll(w.Body)

	// assert API response
	suite.EqualValues(http.StatusBadRequest, w.Code)
	expectedResBody := map[string]any{
		"code":          float64(errc.ClientErrCodeRequired),
		"message":       fmt.Sprintf("%s: %s", errc.ClientErrMsgRequired, "Object"),
		"errorEntities": []any{"Object"},
		"details":       any(nil),
	}
	actualResBody := map[string]any{}
	_ = json.Unmarshal(responseData, &actualResBody)
	suite.Equal(expectedResBody, actualResBody)
}

func (suite *TestSuite) TestAPI_CheckPermission_Failed_3() {
	name := "TestAPI_CheckPermission_Failed: missing relation"
	suite.T().Log(name)

	// call API
	method := "GET"
	endpoint := "/api/v1/permissions/check"
	req, _ := http.NewRequest(method, endpoint, nil)
	q := req.URL.Query()
	q.Add("namespace", "iam-permission")
	q.Add("object", "object 1")
	q.Add("relation", "")
	q.Add("subjectId", "subject 2")
	req.URL.RawQuery = q.Encode()
	w := httptest.NewRecorder()
	testPublicRouter.ServeHTTP(w, req)
	responseData, _ := io.ReadAll(w.Body)

	// assert API response
	suite.EqualValues(http.StatusBadRequest, w.Code)
	expectedResBody := map[string]any{
		"code":          float64(errc.ClientErrCodeRequired),
		"message":       fmt.Sprintf("%s: %s", errc.ClientErrMsgRequired, "Relation"),
		"errorEntities": []any{"Relation"},
		"details":       any(nil),
	}
	actualResBody := map[string]any{}
	_ = json.Unmarshal(responseData, &actualResBody)
	suite.Equal(expectedResBody, actualResBody)
}

func (suite *TestSuite) TestAPI_CheckPermission_Failed_4() {
	name := "TestAPI_CheckPermission_Failed: missing subjectId"
	suite.T().Log(name)

	// call API
	method := "GET"
	endpoint := "/api/v1/permissions/check"
	req, _ := http.NewRequest(method, endpoint, nil)
	q := req.URL.Query()
	q.Add("namespace", "iam-permission")
	q.Add("object", "object 1")
	q.Add("relation", "viewer")
	q.Add("subjectId", "")
	req.URL.RawQuery = q.Encode()
	w := httptest.NewRecorder()
	testPublicRouter.ServeHTTP(w, req)
	responseData, _ := io.ReadAll(w.Body)

	// assert API response
	suite.EqualValues(http.StatusBadRequest, w.Code)
	expectedResBody := map[string]any{
		"code":          float64(errc.ClientErrCodeRequired),
		"message":       fmt.Sprintf("%s: %s", errc.ClientErrMsgRequired, "SubjectID"),
		"errorEntities": []any{"SubjectID"},
		"details":       any(nil),
	}
	actualResBody := map[string]any{}
	_ = json.Unmarshal(responseData, &actualResBody)
	suite.Equal(expectedResBody, actualResBody)
}

func (suite *TestSuite) TestAPI_CheckPermission_Failed_5() {
	name := "TestAPI_CheckPermission_Failed: namespace exceeds max length"
	suite.T().Log(name)

	// call API
	method := "GET"
	endpoint := "/api/v1/permissions/check"
	req, _ := http.NewRequest(method, endpoint, nil)
	q := req.URL.Query()
	q.Add("namespace", exceedLength255String)
	q.Add("object", "object 1")
	q.Add("relation", "viewer")
	q.Add("subjectId", "subject 2")
	req.URL.RawQuery = q.Encode()
	w := httptest.NewRecorder()
	testPublicRouter.ServeHTTP(w, req)
	responseData, _ := io.ReadAll(w.Body)

	// assert API response
	suite.EqualValues(http.StatusBadRequest, w.Code)
	expectedResBody := map[string]any{
		"code":          float64(errc.ClientErrCodeOutOfRange),
		"message":       fmt.Sprintf("%s: %s", errc.ClientErrMsgOutOfRange, "Namespace"),
		"errorEntities": []any{"Namespace"},
		"details":       any(nil),
	}
	actualResBody := map[string]any{}
	_ = json.Unmarshal(responseData, &actualResBody)
	suite.Equal(expectedResBody, actualResBody)
}

func (suite *TestSuite) TestAPI_CheckPermission_Failed_6() {
	name := "TestAPI_CheckPermission_Failed: object exceeds max length"
	suite.T().Log(name)

	// call API
	method := "GET"
	endpoint := "/api/v1/permissions/check"
	req, _ := http.NewRequest(method, endpoint, nil)
	q := req.URL.Query()
	q.Add("namespace", "iam-permission")
	q.Add("object", exceedLength255String)
	q.Add("relation", "viewer")
	q.Add("subjectId", "subject 2")
	req.URL.RawQuery = q.Encode()
	w := httptest.NewRecorder()
	testPublicRouter.ServeHTTP(w, req)
	responseData, _ := io.ReadAll(w.Body)

	// assert API response
	suite.EqualValues(http.StatusBadRequest, w.Code)
	expectedResBody := map[string]any{
		"code":          float64(errc.ClientErrCodeOutOfRange),
		"message":       fmt.Sprintf("%s: %s", errc.ClientErrMsgOutOfRange, "Object"),
		"errorEntities": []any{"Object"},
		"details":       any(nil),
	}
	actualResBody := map[string]any{}
	_ = json.Unmarshal(responseData, &actualResBody)
	suite.Equal(expectedResBody, actualResBody)
}

func (suite *TestSuite) TestAPI_CheckPermission_Failed_7() {
	name := "TestAPI_CheckPermission_Failed: relation exceeds max length"
	suite.T().Log(name)

	// call API
	method := "GET"
	endpoint := "/api/v1/permissions/check"
	req, _ := http.NewRequest(method, endpoint, nil)
	q := req.URL.Query()
	q.Add("namespace", "iam-permission")
	q.Add("object", "object 1")
	q.Add("relation", exceedLength255String)
	q.Add("subjectId", "subject 2")
	req.URL.RawQuery = q.Encode()
	w := httptest.NewRecorder()
	testPublicRouter.ServeHTTP(w, req)
	responseData, _ := io.ReadAll(w.Body)

	// assert API response
	suite.EqualValues(http.StatusBadRequest, w.Code)
	expectedResBody := map[string]any{
		"code":          float64(errc.ClientErrCodeOutOfRange),
		"message":       fmt.Sprintf("%s: %s", errc.ClientErrMsgOutOfRange, "Relation"),
		"errorEntities": []any{"Relation"},
		"details":       any(nil),
	}
	actualResBody := map[string]any{}
	_ = json.Unmarshal(responseData, &actualResBody)
	suite.Equal(expectedResBody, actualResBody)
}

func (suite *TestSuite) TestAPI_CheckPermission_Failed_8() {
	name := "TestAPI_CheckPermission_Failed: subjectId exceeds max length"
	suite.T().Log(name)

	// call API
	method := "GET"
	endpoint := "/api/v1/permissions/check"
	req, _ := http.NewRequest(method, endpoint, nil)
	q := req.URL.Query()
	q.Add("namespace", "iam-permission")
	q.Add("object", "object 1")
	q.Add("relation", "viewer")
	q.Add("subjectId", exceedLength255String)
	req.URL.RawQuery = q.Encode()
	w := httptest.NewRecorder()
	testPublicRouter.ServeHTTP(w, req)
	responseData, _ := io.ReadAll(w.Body)

	// assert API response
	suite.EqualValues(http.StatusBadRequest, w.Code)
	expectedResBody := map[string]any{
		"code":          float64(errc.ClientErrCodeOutOfRange),
		"message":       fmt.Sprintf("%s: %s", errc.ClientErrMsgOutOfRange, "SubjectID"),
		"errorEntities": []any{"SubjectID"},
		"details":       any(nil),
	}
	actualResBody := map[string]any{}
	_ = json.Unmarshal(responseData, &actualResBody)
	suite.Equal(expectedResBody, actualResBody)
}

func (suite *TestSuite) TestAPI_CheckPermission_Failed_9() {
	name := "TestAPI_CheckPermission_Failed: maxDepth exceeds max length"
	suite.T().Log(name)

	// call API
	method := "GET"
	endpoint := "/api/v1/permissions/check"
	req, _ := http.NewRequest(method, endpoint, nil)
	q := req.URL.Query()
	q.Add("namespace", "iam-permission")
	q.Add("object", "object 1")
	q.Add("relation", "viewer")
	q.Add("subjectId", "subject 1")
	q.Add("maxDepth", "4")
	req.URL.RawQuery = q.Encode()
	w := httptest.NewRecorder()
	testPublicRouter.ServeHTTP(w, req)
	responseData, _ := io.ReadAll(w.Body)

	// assert API response
	suite.EqualValues(http.StatusBadRequest, w.Code)
	expectedResBody := map[string]any{
		"code":          float64(errc.ClientErrCodeOutOfRange),
		"message":       fmt.Sprintf("%s: %s", errc.ClientErrMsgOutOfRange, "MaxDepth"),
		"errorEntities": []any{"MaxDepth"},
		"details":       any(nil),
	}
	actualResBody := map[string]any{}
	_ = json.Unmarshal(responseData, &actualResBody)
	suite.Equal(expectedResBody, actualResBody)
}

func (suite *TestSuite) TestAPI_CheckPermission_Failed_10() {
	name := "TestAPI_CheckPermission_Failed: maxDepth is invalid format"
	suite.T().Log(name)

	// call API
	method := "GET"
	endpoint := "/api/v1/permissions/check"
	req, _ := http.NewRequest(method, endpoint, nil)
	q := req.URL.Query()
	q.Add("namespace", "iam-permission")
	q.Add("object", "object 1")
	q.Add("relation", "viewer")
	q.Add("subjectId", "subject 1")
	q.Add("maxDepth", "a")
	req.URL.RawQuery = q.Encode()
	w := httptest.NewRecorder()
	testPublicRouter.ServeHTTP(w, req)
	responseData, _ := io.ReadAll(w.Body)

	// assert API response
	suite.EqualValues(http.StatusBadRequest, w.Code)
	expectedResBody := map[string]any{
		"code":          float64(errc.ClientErrCodeInvalidFormat),
		"message":       errc.ClientErrMsgInvalidFormat,
		"errorEntities": any(nil),
		"details":       any(nil),
	}
	actualResBody := map[string]any{}
	_ = json.Unmarshal(responseData, &actualResBody)
	suite.Equal(expectedResBody, actualResBody)
}
