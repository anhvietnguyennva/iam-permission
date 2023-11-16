package test

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"

	errc "github.com/anhvietnguyennva/go-error/pkg/constant"
	sdk "github.com/anhvietnguyennva/iam-go-sdk"
	"github.com/golang/mock/gomock"

	"iam-permission/internal/pkg/constant"
	"iam-permission/internal/pkg/db"
	"iam-permission/internal/pkg/repository/postgres"
)

func (suite *TestSuite) TestAPI_CreateRelationDefinition_Successfully() {
	name := "TestAPI_CreateRelationDefinition_Successfully"
	suite.T().Log(name)

	// mock
	if err := db.Instance().Exec("INSERT INTO services (id, namespace, description, created_at, updated_at, created_by, updated_by) VALUES " +
		"('410e54a2-2dcd-4675-99f1-df27a01e882d', 'iam-permission', 'IAM permission service', 0, 0, '', '')").Error; err != nil {
		panic(err)
	}
	ctrl := gomock.NewController(suite.T())
	defer ctrl.Finish()
	mockSDK := NewMockIIAMSDK(ctrl)
	sdk.SetSDK(mockSDK)
	defer sdk.SetSDK(nil)
	mockSDK.EXPECT().
		ParseBearerJWT(mockValidBearerAccessToken).
		Return(mockAccessToken(), nil)

	// call API
	method := "POST"
	endpoint := "/admin/api/v1/relation-definitions"
	body := map[string]any{
		"namespace":   "iam-Permission",
		"relation":    "Viewer",
		"description": "description 1",
	}
	marshalledBody, _ := json.Marshal(body)
	req, _ := http.NewRequest(method, endpoint, bytes.NewReader(marshalledBody))
	req.Header.Add("Authorization", mockValidBearerAccessToken)
	w := httptest.NewRecorder()
	testAdminRouter.ServeHTTP(w, req)
	responseData, _ := io.ReadAll(w.Body)

	// check db
	var definitions []*postgres.RelationDefinition
	db.Instance().Find(&definitions)
	suite.EqualValues(1, len(definitions))
	definition := definitions[0]
	suite.EqualValues("410e54a2-2dcd-4675-99f1-df27a01e882d", definition.ServiceID.String())
	suite.EqualValues("iam-permission", definition.Namespace)
	suite.EqualValues("viewer", definition.Relation)
	suite.EqualValues("description 1", definition.Description)
	suite.EqualValues(mockSubject, definition.CreatedBy)
	suite.EqualValues(mockSubject, definition.UpdatedBy)

	// assert API response
	suite.EqualValues(http.StatusOK, w.Code)
	expectedResBody := map[string]any{
		"code":    float64(errc.ClientErrCodeOK),
		"message": errc.ClientErrMsgOK,
		"data": map[string]any{
			"id": definition.ID.String(),
		},
	}
	actualResBody := map[string]any{}
	_ = json.Unmarshal(responseData, &actualResBody)
	suite.Equal(expectedResBody, actualResBody)
}

func (suite *TestSuite) TestAPI_CreateRelationDefinition_Failed_1() {
	name := "TestAPI_CreateRelationDefinition_Failed: missing namespace"
	suite.T().Log(name)

	// mock
	ctrl := gomock.NewController(suite.T())
	defer ctrl.Finish()
	mockSDK := NewMockIIAMSDK(ctrl)
	sdk.SetSDK(mockSDK)
	defer sdk.SetSDK(nil)
	mockSDK.EXPECT().
		ParseBearerJWT(mockValidBearerAccessToken).
		Return(mockAccessToken(), nil)

	// call API
	method := "POST"
	endpoint := "/admin/api/v1/relation-definitions"
	body := map[string]any{
		"namespace":   "",
		"relation":    "Viewer",
		"description": "description 1",
	}
	marshalledBody, _ := json.Marshal(body)
	req, _ := http.NewRequest(method, endpoint, bytes.NewReader(marshalledBody))
	req.Header.Add("Authorization", mockValidBearerAccessToken)
	w := httptest.NewRecorder()
	testAdminRouter.ServeHTTP(w, req)
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

func (suite *TestSuite) TestAPI_CreateRelationDefinition_Failed_2() {
	name := "TestAPI_CreateRelationDefinition_Failed: missing relation"
	suite.T().Log(name)

	// mock
	ctrl := gomock.NewController(suite.T())
	defer ctrl.Finish()
	mockSDK := NewMockIIAMSDK(ctrl)
	sdk.SetSDK(mockSDK)
	defer sdk.SetSDK(nil)
	mockSDK.EXPECT().
		ParseBearerJWT(mockValidBearerAccessToken).
		Return(mockAccessToken(), nil)

	// call API
	method := "POST"
	endpoint := "/admin/api/v1/relation-definitions"
	body := map[string]any{
		"namespace":   "iam-Permission",
		"relation":    "",
		"description": "description 1",
	}
	marshalledBody, _ := json.Marshal(body)
	req, _ := http.NewRequest(method, endpoint, bytes.NewReader(marshalledBody))
	req.Header.Add("Authorization", mockValidBearerAccessToken)
	w := httptest.NewRecorder()
	testAdminRouter.ServeHTTP(w, req)
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

func (suite *TestSuite) TestAPI_CreateRelationDefinition_Failed_3() {
	name := "TestAPI_CreateRelationDefinition_Failed: namespace exceeds max length"
	suite.T().Log(name)

	// mock
	ctrl := gomock.NewController(suite.T())
	defer ctrl.Finish()
	mockSDK := NewMockIIAMSDK(ctrl)
	sdk.SetSDK(mockSDK)
	defer sdk.SetSDK(nil)
	mockSDK.EXPECT().
		ParseBearerJWT(mockValidBearerAccessToken).
		Return(mockAccessToken(), nil)

	// call API
	method := "POST"
	endpoint := "/admin/api/v1/relation-definitions"
	body := map[string]any{
		"namespace":   exceedLength255String,
		"object":      "object",
		"relation":    "Viewer",
		"description": "description 1",
	}
	marshalledBody, _ := json.Marshal(body)
	req, _ := http.NewRequest(method, endpoint, bytes.NewReader(marshalledBody))
	req.Header.Add("Authorization", mockValidBearerAccessToken)
	w := httptest.NewRecorder()
	testAdminRouter.ServeHTTP(w, req)
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

func (suite *TestSuite) TestAPI_CreateRelationDefinition_Failed_4() {
	name := "TestAPI_CreateRelationDefinition_Failed: relation exceeds max length"
	suite.T().Log(name)

	// mock
	ctrl := gomock.NewController(suite.T())
	defer ctrl.Finish()
	mockSDK := NewMockIIAMSDK(ctrl)
	sdk.SetSDK(mockSDK)
	defer sdk.SetSDK(nil)
	mockSDK.EXPECT().
		ParseBearerJWT(mockValidBearerAccessToken).
		Return(mockAccessToken(), nil)

	// call API
	method := "POST"
	endpoint := "/admin/api/v1/relation-definitions"
	body := map[string]any{
		"namespace":   "iam-permission",
		"relation":    exceedLength255String,
		"description": "description 1",
	}
	marshalledBody, _ := json.Marshal(body)
	req, _ := http.NewRequest(method, endpoint, bytes.NewReader(marshalledBody))
	req.Header.Add("Authorization", mockValidBearerAccessToken)
	w := httptest.NewRecorder()
	testAdminRouter.ServeHTTP(w, req)
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

func (suite *TestSuite) TestAPI_CreateRelationDefinition_Failed_5() {
	name := "TestAPI_CreateRelationDefinition_Failed: description exceeds max length"
	suite.T().Log(name)

	// mock
	ctrl := gomock.NewController(suite.T())
	defer ctrl.Finish()
	mockSDK := NewMockIIAMSDK(ctrl)
	sdk.SetSDK(mockSDK)
	defer sdk.SetSDK(nil)
	mockSDK.EXPECT().
		ParseBearerJWT(mockValidBearerAccessToken).
		Return(mockAccessToken(), nil)

	// call API
	method := "POST"
	endpoint := "/admin/api/v1/relation-definitions"
	body := map[string]any{
		"namespace":   "iam-permission",
		"relation":    "Viewer",
		"description": exceedLength255String,
	}
	marshalledBody, _ := json.Marshal(body)
	req, _ := http.NewRequest(method, endpoint, bytes.NewReader(marshalledBody))
	req.Header.Add("Authorization", mockValidBearerAccessToken)
	w := httptest.NewRecorder()
	testAdminRouter.ServeHTTP(w, req)
	responseData, _ := io.ReadAll(w.Body)

	// assert API response
	suite.EqualValues(http.StatusBadRequest, w.Code)
	expectedResBody := map[string]any{
		"code":          float64(errc.ClientErrCodeOutOfRange),
		"message":       fmt.Sprintf("%s: %s", errc.ClientErrMsgOutOfRange, "Description"),
		"errorEntities": []any{"Description"},
		"details":       any(nil),
	}
	actualResBody := map[string]any{}
	_ = json.Unmarshal(responseData, &actualResBody)
	suite.Equal(expectedResBody, actualResBody)
}

func (suite *TestSuite) TestAPI_CreateRelationDefinition_Failed_6() {
	name := "TestAPI_CreateRelationDefinition_Failed: not found namespace"
	suite.T().Log(name)

	// mock
	if err := db.Instance().Exec("INSERT INTO services (id, namespace, description, created_at, updated_at, created_by, updated_by) VALUES " +
		"('410e54a2-2dcd-4675-99f1-df27a01e882d', 'iam-permission', 'IAM permission service', 0, 0, '', '')").Error; err != nil {
		panic(err)
	}
	ctrl := gomock.NewController(suite.T())
	defer ctrl.Finish()
	mockSDK := NewMockIIAMSDK(ctrl)
	sdk.SetSDK(mockSDK)
	defer sdk.SetSDK(nil)
	mockSDK.EXPECT().
		ParseBearerJWT(mockValidBearerAccessToken).
		Return(mockAccessToken(), nil)

	// call API
	method := "POST"
	endpoint := "/admin/api/v1/relation-definitions"
	body := map[string]any{
		"namespace":   "iam-Permission1",
		"relation":    "Viewer",
		"description": "description 1",
	}
	marshalledBody, _ := json.Marshal(body)
	req, _ := http.NewRequest(method, endpoint, bytes.NewReader(marshalledBody))
	req.Header.Add("Authorization", mockValidBearerAccessToken)
	w := httptest.NewRecorder()
	testAdminRouter.ServeHTTP(w, req)
	responseData, _ := io.ReadAll(w.Body)

	// assert API response
	suite.EqualValues(http.StatusNotFound, w.Code)
	expectedResBody := map[string]any{
		"code":          float64(errc.ClientErrCodeNotFound),
		"message":       fmt.Sprintf("%s: %s", errc.ClientErrMsgNotFound, constant.FieldService),
		"errorEntities": []any{constant.FieldService},
		"details":       any(nil),
	}
	actualResBody := map[string]any{}
	_ = json.Unmarshal(responseData, &actualResBody)
	suite.Equal(expectedResBody, actualResBody)
}

func (suite *TestSuite) TestAPI_CreateRelationDefinition_Failed_7() {
	name := "TestAPI_CreateRelationDefinition_Failed: relation definition already exists"
	suite.T().Log(name)

	// mock
	if err := db.Instance().Exec("INSERT INTO services (id, namespace, description, created_at, updated_at, created_by, updated_by) VALUES " +
		"('410e54a2-2dcd-4675-99f1-df27a01e882d', 'iam-permission', 'IAM permission service', 0, 0, '', '')").Error; err != nil {
		panic(err)
	}
	if err := db.Instance().Exec("INSERT INTO relation_definitions " +
		"(id, service_id, namespace, relation, description, created_at, updated_at, created_by, updated_by) VALUES " +
		"('4fe1bfd3-7156-4b38-bf32-f90b57aff7fc', '410e54a2-2dcd-4675-99f1-df27a01e882d', 'iam-permission', 'viewer', 'view permission', 0, 0, '', '')").Error; err != nil {
		panic(err)
	}
	ctrl := gomock.NewController(suite.T())
	defer ctrl.Finish()
	mockSDK := NewMockIIAMSDK(ctrl)
	sdk.SetSDK(mockSDK)
	defer sdk.SetSDK(nil)
	mockSDK.EXPECT().
		ParseBearerJWT(mockValidBearerAccessToken).
		Return(mockAccessToken(), nil)

	// call API
	method := "POST"
	endpoint := "/admin/api/v1/relation-definitions"
	body := map[string]any{
		"namespace":   "iam-Permission",
		"relation":    "Viewer",
		"description": "description 1",
	}
	marshalledBody, _ := json.Marshal(body)
	req, _ := http.NewRequest(method, endpoint, bytes.NewReader(marshalledBody))
	req.Header.Add("Authorization", mockValidBearerAccessToken)
	w := httptest.NewRecorder()
	testAdminRouter.ServeHTTP(w, req)
	responseData, _ := io.ReadAll(w.Body)

	// assert API response
	suite.EqualValues(http.StatusConflict, w.Code)
	expectedResBody := map[string]any{
		"code":          float64(errc.ClientErrCodeDuplicate),
		"message":       fmt.Sprintf("%s: %s", errc.ClientErrMsgDuplicate, constant.FieldRelationDefinition),
		"errorEntities": []any{constant.FieldRelationDefinition},
		"details":       any(nil),
	}
	actualResBody := map[string]any{}
	_ = json.Unmarshal(responseData, &actualResBody)
	suite.Equal(expectedResBody, actualResBody)
}

func (suite *TestSuite) TestAPI_CreateRelationDefinition_Failed_8() {
	name := "TestAPI_CreateRelationDefinition_Failed: invalid access token"
	suite.T().Log(name)

	// mock
	ctrl := gomock.NewController(suite.T())
	defer ctrl.Finish()
	mockSDK := NewMockIIAMSDK(ctrl)
	sdk.SetSDK(mockSDK)
	defer sdk.SetSDK(nil)
	mockSDK.EXPECT().
		ParseBearerJWT(mockInvalidBearerAccessToken).
		Return(nil, errors.New("invalid token"))

	// call API
	method := "POST"
	endpoint := "/admin/api/v1/relation-definitions"
	body := map[string]any{
		"namespace":   "iam-Permission",
		"object":      "object",
		"relation":    "Viewer",
		"description": "description 1",
	}
	marshalledBody, _ := json.Marshal(body)
	req, _ := http.NewRequest(method, endpoint, bytes.NewReader(marshalledBody))
	req.Header.Add("Authorization", mockInvalidBearerAccessToken)
	w := httptest.NewRecorder()
	testAdminRouter.ServeHTTP(w, req)
	responseData, _ := io.ReadAll(w.Body)

	// check database
	var services []*postgres.Service
	db.Instance().Find(&services)
	suite.EqualValues(0, len(services))

	// assert API response
	suite.EqualValues(http.StatusUnauthorized, w.Code)
	expectedResBody := map[string]any{
		"code":          float64(errc.ClientErrCodeUnauthenticated),
		"message":       errc.ClientErrMsgUnauthenticated,
		"errorEntities": any(nil),
		"details":       any(nil),
	}
	actualResBody := map[string]any{}
	_ = json.Unmarshal(responseData, &actualResBody)
	suite.Equal(expectedResBody, actualResBody)
}
