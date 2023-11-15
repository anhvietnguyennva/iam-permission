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

func (suite *TestSuite) TestAPI_CreateSubjectRelationTuple_Successfully() {
	name := "TestAPI_CreateSubjectRelationTuple_Successfully"
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
	endpoint := "/admin/api/v1/subject-relation-tuples"
	body := map[string]any{
		"namespace": "iam-Permission",
		"object":    "object",
		"relation":  "Viewer",
		"subjectId": "1",
	}
	marshalledBody, _ := json.Marshal(body)
	req, _ := http.NewRequest(method, endpoint, bytes.NewReader(marshalledBody))
	req.Header.Add("Authorization", mockValidBearerAccessToken)
	w := httptest.NewRecorder()
	testAdminRouter.ServeHTTP(w, req)
	responseData, _ := io.ReadAll(w.Body)

	// check db
	var tuples []*postgres.SubjectRelationTuple
	db.Instance().Find(&tuples)
	suite.EqualValues(1, len(tuples))
	tuple := tuples[0]
	suite.EqualValues("410e54a2-2dcd-4675-99f1-df27a01e882d", tuple.ServiceID.String())
	suite.EqualValues("4fe1bfd3-7156-4b38-bf32-f90b57aff7fc", tuple.RelationDefinitionID.String())
	suite.EqualValues("iam-permission", tuple.Namespace)
	suite.EqualValues("object", tuple.Object)
	suite.EqualValues("viewer", tuple.Relation)
	suite.EqualValues("1", tuple.SubjectID)
	suite.EqualValues(mockSubject, tuple.CreatedBy)
	suite.EqualValues(mockSubject, tuple.UpdatedBy)

	// assert API response
	suite.EqualValues(http.StatusOK, w.Code)
	expectedResBody := map[string]any{
		"code":    float64(errc.ClientErrCodeOK),
		"message": errc.ClientErrMsgOK,
		"data": map[string]any{
			"id": tuple.ID.String(),
		},
	}
	actualResBody := map[string]any{}
	_ = json.Unmarshal(responseData, &actualResBody)
	suite.Equal(expectedResBody, actualResBody)
}

func (suite *TestSuite) TestAPI_CreateSubjectRelationTuple_Failed_1() {
	name := "TestAPI_CreateSubjectRelationTuple_Failed: missing namespace"
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
	endpoint := "/admin/api/v1/subject-relation-tuples"
	body := map[string]any{
		"namespace": "",
		"object":    "object",
		"relation":  "Viewer",
		"subjectId": "1",
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

func (suite *TestSuite) TestAPI_CreateSubjectRelationTuple_Failed_2() {
	name := "TestAPI_CreateSubjectRelationTuple_Failed: missing object"
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
	endpoint := "/admin/api/v1/subject-relation-tuples"
	body := map[string]any{
		"namespace": "iam-Permission",
		"object":    "",
		"relation":  "Viewer",
		"subjectId": "1",
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

func (suite *TestSuite) TestAPI_CreateSubjectRelationTuple_Failed_3() {
	name := "TestAPI_CreateSubjectRelationTuple_Failed: missing relation"
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
	endpoint := "/admin/api/v1/subject-relation-tuples"
	body := map[string]any{
		"namespace": "iam-Permission",
		"object":    "object",
		"relation":  "",
		"subjectId": "1",
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

func (suite *TestSuite) TestAPI_CreateSubjectRelationTuple_Failed_4() {
	name := "TestAPI_CreateSubjectRelationTuple_Failed: missing subjectId"
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
	endpoint := "/admin/api/v1/subject-relation-tuples"
	body := map[string]any{
		"namespace": "iam-Permission",
		"object":    "object",
		"relation":  "Viewer",
		"subjectId": "",
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
		"message":       fmt.Sprintf("%s: %s", errc.ClientErrMsgRequired, "SubjectID"),
		"errorEntities": []any{"SubjectID"},
		"details":       any(nil),
	}
	actualResBody := map[string]any{}
	_ = json.Unmarshal(responseData, &actualResBody)
	suite.Equal(expectedResBody, actualResBody)
}

func (suite *TestSuite) TestAPI_CreateSubjectRelationTuple_Failed_5() {
	name := "TestAPI_CreateSubjectRelationTuple_Failed: namespace exceeds max length"
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
	endpoint := "/admin/api/v1/subject-relation-tuples"
	body := map[string]any{
		"namespace": exceedLength255String,
		"object":    "object",
		"relation":  "Viewer",
		"subjectId": "1",
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

func (suite *TestSuite) TestAPI_CreateSubjectRelationTuple_Failed_6() {
	name := "TestAPI_CreateSubjectRelationTuple_Failed: object exceeds max length"
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
	endpoint := "/admin/api/v1/subject-relation-tuples"
	body := map[string]any{
		"namespace": "iam-permission",
		"object":    exceedLength255String,
		"relation":  "Viewer",
		"subjectId": "1",
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

func (suite *TestSuite) TestAPI_CreateSubjectRelationTuple_Failed_7() {
	name := "TestAPI_CreateSubjectRelationTuple_Failed: relation exceeds max length"
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
	endpoint := "/admin/api/v1/subject-relation-tuples"
	body := map[string]any{
		"namespace": "iam-permission",
		"object":    "object",
		"relation":  exceedLength255String,
		"subjectId": "1",
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

func (suite *TestSuite) TestAPI_CreateSubjectRelationTuple_Failed_8() {
	name := "TestAPI_CreateSubjectRelationTuple_Failed: subjectId exceeds max length"
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
	endpoint := "/admin/api/v1/subject-relation-tuples"
	body := map[string]any{
		"namespace": "iam-permission",
		"object":    "object",
		"relation":  "Viewer",
		"subjectId": exceedLength255String,
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
		"message":       fmt.Sprintf("%s: %s", errc.ClientErrMsgOutOfRange, "SubjectID"),
		"errorEntities": []any{"SubjectID"},
		"details":       any(nil),
	}
	actualResBody := map[string]any{}
	_ = json.Unmarshal(responseData, &actualResBody)
	suite.Equal(expectedResBody, actualResBody)
}

func (suite *TestSuite) TestAPI_CreateSubjectRelationTuple_Failed_9() {
	name := "TestAPI_CreateSubjectRelationTuple_Failed: not found namespace"
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
	endpoint := "/admin/api/v1/subject-relation-tuples"
	body := map[string]any{
		"namespace": "iam-Permission1",
		"object":    "object",
		"relation":  "Viewer",
		"subjectId": "1",
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

func (suite *TestSuite) TestAPI_CreateSubjectRelationTuple_Failed_10() {
	name := "TestAPI_CreateSubjectRelationTuple_Failed: not found relation"
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
	endpoint := "/admin/api/v1/subject-relation-tuples"
	body := map[string]any{
		"namespace": "iam-Permission",
		"object":    "object",
		"relation":  "Viewer1",
		"subjectId": "1",
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

func (suite *TestSuite) TestAPI_CreateSubjectRelationTuple_Failed_11() {
	name := "TestAPI_CreateSubjectRelationTuple_Failed: tuple already exists"
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
	if err := db.Instance().Exec("INSERT INTO subject_relation_tuples " +
		"(id, service_id, relation_definition_id, namespace, object, relation, subject_id, created_at, updated_at, created_by, updated_by) VALUES " +
		"('3221d4a9-ff49-4ac3-a3c2-ec39b8eef6ca', '410e54a2-2dcd-4675-99f1-df27a01e882d', '4fe1bfd3-7156-4b38-bf32-f90b57aff7fc', 'iam-permission', 'object', 'viewer', '1', 0, 0, '', '')").Error; err != nil {
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
	endpoint := "/admin/api/v1/subject-relation-tuples"
	body := map[string]any{
		"namespace": "iam-Permission",
		"object":    "object",
		"relation":  "Viewer",
		"subjectId": "1",
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
		"message":       fmt.Sprintf("%s: %s", errc.ClientErrMsgDuplicate, constant.FieldSubjectRelationTuple),
		"errorEntities": []any{constant.FieldSubjectRelationTuple},
		"details":       any(nil),
	}
	actualResBody := map[string]any{}
	_ = json.Unmarshal(responseData, &actualResBody)
	suite.Equal(expectedResBody, actualResBody)
}

func (suite *TestSuite) TestAPI_CreateSubjectRelationTuple_Failed_12() {
	name := "TestAPI_CreateSubjectRelationTuple_Failed: invalid access token"
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
	endpoint := "/admin/api/v1/subject-relation-tuples"
	body := map[string]any{
		"namespace": "iam-Permission",
		"object":    "object",
		"relation":  "Viewer",
		"subjectId": "1",
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
