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

func (suite *TestSuite) TestAPI_CreateSubjectSet_Successfully() {
	name := "TestAPI_CreateSubjectSet_Successfully"
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
	endpoint := "/admin/api/v1/subject-sets"
	body := map[string]any{
		"namespace":   "iam-Permission",
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

	// check db
	var sets []*postgres.SubjectSet
	db.Instance().Find(&sets)
	suite.EqualValues(1, len(sets))
	subjectSet := sets[0]
	suite.EqualValues("410e54a2-2dcd-4675-99f1-df27a01e882d", subjectSet.ServiceID.String())
	suite.EqualValues("4fe1bfd3-7156-4b38-bf32-f90b57aff7fc", subjectSet.RelationDefinitionID.String())
	suite.EqualValues("iam-permission", subjectSet.Namespace)
	suite.EqualValues("object", subjectSet.Object)
	suite.EqualValues("viewer", subjectSet.Relation)
	suite.EqualValues("description 1", subjectSet.Description)
	suite.EqualValues(mockSubject, subjectSet.CreatedBy)
	suite.EqualValues(mockSubject, subjectSet.UpdatedBy)

	// assert API response
	suite.EqualValues(http.StatusOK, w.Code)
	expectedResBody := map[string]any{
		"code":    float64(errc.ClientErrCodeOK),
		"message": errc.ClientErrMsgOK,
		"data": map[string]any{
			"id": subjectSet.ID.String(),
		},
	}
	actualResBody := map[string]any{}
	_ = json.Unmarshal(responseData, &actualResBody)
	suite.Equal(expectedResBody, actualResBody)
}

func (suite *TestSuite) TestAPI_CreateSubjectSet_Failed_1() {
	name := "TestAPI_CreateSubjectSet_Failed: missing namespace"
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
	endpoint := "/admin/api/v1/subject-sets"
	body := map[string]any{
		"namespace":   "",
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
		"code":          float64(errc.ClientErrCodeRequired),
		"message":       fmt.Sprintf("%s: %s", errc.ClientErrMsgRequired, "Namespace"),
		"errorEntities": []any{"Namespace"},
		"details":       any(nil),
	}
	actualResBody := map[string]any{}
	_ = json.Unmarshal(responseData, &actualResBody)
	suite.Equal(expectedResBody, actualResBody)
}

func (suite *TestSuite) TestAPI_CreateSubjectSet_Failed_2() {
	name := "TestAPI_CreateSubjectSet_Failed: missing object"
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
	endpoint := "/admin/api/v1/subject-sets"
	body := map[string]any{
		"namespace":   "iam-Permission",
		"object":      "",
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
		"message":       fmt.Sprintf("%s: %s", errc.ClientErrMsgRequired, "Object"),
		"errorEntities": []any{"Object"},
		"details":       any(nil),
	}
	actualResBody := map[string]any{}
	_ = json.Unmarshal(responseData, &actualResBody)
	suite.Equal(expectedResBody, actualResBody)
}

func (suite *TestSuite) TestAPI_CreateSubjectSet_Failed_3() {
	name := "TestAPI_CreateSubjectSet_Failed: missing relation"
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
	endpoint := "/admin/api/v1/subject-sets"
	body := map[string]any{
		"namespace":   "iam-Permission",
		"object":      "object",
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

func (suite *TestSuite) TestAPI_CreateSubjectSet_Failed_4() {
	name := "TestAPI_CreateSubjectSet_Failed: namespace exceeds max length"
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
	endpoint := "/admin/api/v1/subject-sets"
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

func (suite *TestSuite) TestAPI_CreateSubjectSet_Failed_5() {
	name := "TestAPI_CreateSubjectSet_Failed: object exceeds max length"
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
	endpoint := "/admin/api/v1/subject-sets"
	body := map[string]any{
		"namespace":   "iam-permission",
		"object":      exceedLength255String,
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
		"message":       fmt.Sprintf("%s: %s", errc.ClientErrMsgOutOfRange, "Object"),
		"errorEntities": []any{"Object"},
		"details":       any(nil),
	}
	actualResBody := map[string]any{}
	_ = json.Unmarshal(responseData, &actualResBody)
	suite.Equal(expectedResBody, actualResBody)
}

func (suite *TestSuite) TestAPI_CreateSubjectSet_Failed_6() {
	name := "TestAPI_CreateSubjectSet_Failed: relation exceeds max length"
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
	endpoint := "/admin/api/v1/subject-sets"
	body := map[string]any{
		"namespace":   "iam-permission",
		"object":      "object",
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

func (suite *TestSuite) TestAPI_CreateSubjectSet_Failed_7() {
	name := "TestAPI_CreateSubjectSet_Failed: description exceeds max length"
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
	endpoint := "/admin/api/v1/subject-sets"
	body := map[string]any{
		"namespace":   "iam-permission",
		"object":      "object",
		"relation":    "Viewer",
		"description": exceedLength500String,
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

func (suite *TestSuite) TestAPI_CreateSubjectSet_Failed_8() {
	name := "TestAPI_CreateSubjectSet_Failed: not found namespace"
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
	endpoint := "/admin/api/v1/subject-sets"
	body := map[string]any{
		"namespace":   "iam-Permission1",
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

func (suite *TestSuite) TestAPI_CreateSubjectSet_Failed_9() {
	name := "TestAPI_CreateSubjectSet_Failed: not found relation"
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
	endpoint := "/admin/api/v1/subject-sets"
	body := map[string]any{
		"namespace":   "iam-Permission",
		"object":      "object",
		"relation":    "Viewer1",
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
		"message":       fmt.Sprintf("%s: %s", errc.ClientErrMsgNotFound, constant.FieldRelationDefinition),
		"errorEntities": []any{constant.FieldRelationDefinition},
		"details":       any(nil),
	}
	actualResBody := map[string]any{}
	_ = json.Unmarshal(responseData, &actualResBody)
	suite.Equal(expectedResBody, actualResBody)
}

func (suite *TestSuite) TestAPI_CreateSubjectSet_Failed_10() {
	name := "TestAPI_CreateSubjectSet_Failed: subject set already exists"
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
	if err := db.Instance().Exec("INSERT INTO subject_sets " +
		"(id, service_id, relation_definition_id, namespace, object, relation, description, created_at, updated_at, created_by, updated_by) VALUES " +
		"('3221d4a9-ff49-4ac3-a3c2-ec39b8eef6ca', '410e54a2-2dcd-4675-99f1-df27a01e882d', '4fe1bfd3-7156-4b38-bf32-f90b57aff7fc', 'iam-permission', 'object', 'viewer', 'description 1', 0, 0, '', '')").Error; err != nil {
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
	endpoint := "/admin/api/v1/subject-sets"
	body := map[string]any{
		"namespace":   "iam-Permission",
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
	suite.EqualValues(http.StatusConflict, w.Code)
	expectedResBody := map[string]any{
		"code":          float64(errc.ClientErrCodeDuplicate),
		"message":       fmt.Sprintf("%s: %s", errc.ClientErrMsgDuplicate, constant.FieldSubjectSet),
		"errorEntities": []any{constant.FieldSubjectSet},
		"details":       any(nil),
	}
	actualResBody := map[string]any{}
	_ = json.Unmarshal(responseData, &actualResBody)
	suite.Equal(expectedResBody, actualResBody)
}

func (suite *TestSuite) TestAPI_CreateSubjectSet_Failed_11() {
	name := "TestAPI_CreateSubjectSet_Failed: invalid access token"
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
	endpoint := "/admin/api/v1/subject-sets"
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
