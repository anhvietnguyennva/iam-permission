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

func (suite *TestSuite) TestAPI_CreateSubjectSetRelationTuple_Successfully() {
	name := "TestAPI_CreateSubjectSetRelationTuple_Successfully"
	suite.T().Log(name)

	// mock
	if err := db.Instance().Exec("INSERT INTO services (id, namespace, description, created_at, updated_at, created_by, updated_by) VALUES " +
		"('410e54a2-2dcd-4675-99f1-df27a01e882d', 'iam-permission', 'IAM permission service', 0, 0, '', '')").Error; err != nil {
		panic(err)
	}
	if err := db.Instance().Exec("INSERT INTO relation_definitions " +
		"(id, service_id, namespace, relation, description, created_at, updated_at, created_by, updated_by) VALUES " +
		"('4fe1bfd3-7156-4b38-bf32-f90b57aff7fc', '410e54a2-2dcd-4675-99f1-df27a01e882d', 'iam-permission', 'viewer', 'view permission', 0, 0, '', '')," +
		"('344f382d-90aa-42ef-9876-5697f8150ba0', '410e54a2-2dcd-4675-99f1-df27a01e882d', 'iam-permission', 'operator', 'operator role', 0, 0, '', '')").
		Error; err != nil {
		panic(err)
	}
	if err := db.Instance().Exec("INSERT INTO subject_sets " +
		"(id, service_id, relation_definition_id, namespace, object, relation, description, created_at, updated_at, created_by, updated_by) VALUES " +
		"('3221d4a9-ff49-4ac3-a3c2-ec39b8eef6ca', '410e54a2-2dcd-4675-99f1-df27a01e882d', '344f382d-90aa-42ef-9876-5697f8150ba0', 'iam-permission', 'services/410e54a2-2dcd-4675-99f1-df27a01e882d', 'operator', 'description 1', 0, 0, '', '')").Error; err != nil {
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
	endpoint := "/admin/api/v1/subject-set-relation-tuples"
	body := map[string]any{
		"namespace":           "iam-Permission",
		"object":              "object",
		"relation":            "Viewer",
		"subjectSetNamespace": "iam-permission",
		"subjectSetObject":    "services/410e54a2-2dcd-4675-99f1-df27a01e882d",
		"subjectSetRelation":  "operator",
	}
	marshalledBody, _ := json.Marshal(body)
	req, _ := http.NewRequest(method, endpoint, bytes.NewReader(marshalledBody))
	req.Header.Add("Authorization", mockValidBearerAccessToken)
	w := httptest.NewRecorder()
	testAdminRouter.ServeHTTP(w, req)
	responseData, _ := io.ReadAll(w.Body)

	// check db
	var tuples []*postgres.SubjectSetRelationTuple
	db.Instance().Find(&tuples)
	suite.EqualValues(1, len(tuples))
	tuple := tuples[0]
	suite.EqualValues("410e54a2-2dcd-4675-99f1-df27a01e882d", tuple.ServiceID.String())
	suite.EqualValues("4fe1bfd3-7156-4b38-bf32-f90b57aff7fc", tuple.RelationDefinitionID.String())
	suite.EqualValues("3221d4a9-ff49-4ac3-a3c2-ec39b8eef6ca", tuple.SubjectSetID.String())
	suite.EqualValues("iam-permission", tuple.Namespace)
	suite.EqualValues("object", tuple.Object)
	suite.EqualValues("viewer", tuple.Relation)
	suite.EqualValues("iam-permission", tuple.SubjectSetNamespace)
	suite.EqualValues("services/410e54a2-2dcd-4675-99f1-df27a01e882d", tuple.SubjectSetObject)
	suite.EqualValues("operator", tuple.SubjectSetRelation)
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

func (suite *TestSuite) TestAPI_CreateSubjectSetRelationTuple_Failed_1() {
	name := "TestAPI_CreateSubjectSetRelationTuple_Failed: missing namespace"
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
	endpoint := "/admin/api/v1/subject-set-relation-tuples"
	body := map[string]any{
		"namespace":           "",
		"object":              "object",
		"relation":            "Viewer",
		"subjectSetNamespace": "iam-permission",
		"subjectSetObject":    "services/410e54a2-2dcd-4675-99f1-df27a01e882d",
		"subjectSetRelation":  "operator",
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

func (suite *TestSuite) TestAPI_CreateSubjectSetRelationTuple_Failed_2() {
	name := "TestAPI_CreateSubjectSetRelationTuple_Failed: missing object"
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
	endpoint := "/admin/api/v1/subject-set-relation-tuples"
	body := map[string]any{
		"namespace":           "iam-Permission",
		"object":              "",
		"relation":            "Viewer",
		"subjectSetNamespace": "iam-permission",
		"subjectSetObject":    "services/410e54a2-2dcd-4675-99f1-df27a01e882d",
		"subjectSetRelation":  "operator",
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

func (suite *TestSuite) TestAPI_CreateSubjectSetRelationTuple_Failed_3() {
	name := "TestAPI_CreateSubjectSetRelationTuple_Failed: missing relation"
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
	endpoint := "/admin/api/v1/subject-set-relation-tuples"
	body := map[string]any{
		"namespace":           "iam-Permission",
		"object":              "object",
		"relation":            "",
		"subjectSetNamespace": "iam-permission",
		"subjectSetObject":    "services/410e54a2-2dcd-4675-99f1-df27a01e882d",
		"subjectSetRelation":  "operator",
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

func (suite *TestSuite) TestAPI_CreateSubjectSetRelationTuple_Failed_4() {
	name := "TestAPI_CreateSubjectSetRelationTuple_Failed: missing subjectSetNamespace"
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
	endpoint := "/admin/api/v1/subject-set-relation-tuples"
	body := map[string]any{
		"namespace":           "iam-Permission",
		"object":              "object",
		"relation":            "Viewer",
		"subjectSetNamespace": "",
		"subjectSetObject":    "services/410e54a2-2dcd-4675-99f1-df27a01e882d",
		"subjectSetRelation":  "operator",
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
		"message":       fmt.Sprintf("%s: %s", errc.ClientErrMsgRequired, "SubjectSetNamespace"),
		"errorEntities": []any{"SubjectSetNamespace"},
		"details":       any(nil),
	}
	actualResBody := map[string]any{}
	_ = json.Unmarshal(responseData, &actualResBody)
	suite.Equal(expectedResBody, actualResBody)
}

func (suite *TestSuite) TestAPI_CreateSubjectSetRelationTuple_Failed_5() {
	name := "TestAPI_CreateSubjectSetRelationTuple_Failed: missing subjectSetObject"
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
	endpoint := "/admin/api/v1/subject-set-relation-tuples"
	body := map[string]any{
		"namespace":           "iam-Permission",
		"object":              "object",
		"relation":            "Viewer",
		"subjectSetNamespace": "iam-permission",
		"subjectSetObject":    "",
		"subjectSetRelation":  "operator",
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
		"message":       fmt.Sprintf("%s: %s", errc.ClientErrMsgRequired, "SubjectSetObject"),
		"errorEntities": []any{"SubjectSetObject"},
		"details":       any(nil),
	}
	actualResBody := map[string]any{}
	_ = json.Unmarshal(responseData, &actualResBody)
	suite.Equal(expectedResBody, actualResBody)
}

func (suite *TestSuite) TestAPI_CreateSubjectSetRelationTuple_Failed_6() {
	name := "TestAPI_CreateSubjectSetRelationTuple_Failed: missing subjectSetRelation"
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
	endpoint := "/admin/api/v1/subject-set-relation-tuples"
	body := map[string]any{
		"namespace":           "iam-Permission",
		"object":              "object",
		"relation":            "Viewer",
		"subjectSetNamespace": "iam-permission",
		"subjectSetObject":    "services/410e54a2-2dcd-4675-99f1-df27a01e882d",
		"subjectSetRelation":  "",
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
		"message":       fmt.Sprintf("%s: %s", errc.ClientErrMsgRequired, "SubjectSetRelation"),
		"errorEntities": []any{"SubjectSetRelation"},
		"details":       any(nil),
	}
	actualResBody := map[string]any{}
	_ = json.Unmarshal(responseData, &actualResBody)
	suite.Equal(expectedResBody, actualResBody)
}

func (suite *TestSuite) TestAPI_CreateSubjectSetRelationTuple_Failed_7() {
	name := "TestAPI_CreateSubjectSetRelationTuple_Failed: namespace exceeds max length"
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
	endpoint := "/admin/api/v1/subject-set-relation-tuples"
	body := map[string]any{
		"namespace":           exceedLength255String,
		"object":              "object",
		"relation":            "Viewer",
		"subjectSetNamespace": "iam-permission",
		"subjectSetObject":    "services/410e54a2-2dcd-4675-99f1-df27a01e882d",
		"subjectSetRelation":  "operator",
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

func (suite *TestSuite) TestAPI_CreateSubjectSetRelationTuple_Failed_8() {
	name := "TestAPI_CreateSubjectSetRelationTuple_Failed: object exceeds max length"
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
	endpoint := "/admin/api/v1/subject-set-relation-tuples"
	body := map[string]any{
		"namespace":           "iam-Permission",
		"object":              exceedLength255String,
		"relation":            "Viewer",
		"subjectSetNamespace": "iam-permission",
		"subjectSetObject":    "services/410e54a2-2dcd-4675-99f1-df27a01e882d",
		"subjectSetRelation":  "operator",
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

func (suite *TestSuite) TestAPI_CreateSubjectSetRelationTuple_Failed_9() {
	name := "TestAPI_CreateSubjectSetRelationTuple_Failed: relation exceeds max length"
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
	endpoint := "/admin/api/v1/subject-set-relation-tuples"
	body := map[string]any{
		"namespace":           "iam-Permission",
		"object":              "object",
		"relation":            exceedLength255String,
		"subjectSetNamespace": "iam-permission",
		"subjectSetObject":    "services/410e54a2-2dcd-4675-99f1-df27a01e882d",
		"subjectSetRelation":  "operator",
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

func (suite *TestSuite) TestAPI_CreateSubjectSetRelationTuple_Failed_10() {
	name := "TestAPI_CreateSubjectSetRelationTuple_Failed: subjectSetNamespace exceeds max length"
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
	endpoint := "/admin/api/v1/subject-set-relation-tuples"
	body := map[string]any{
		"namespace":           "iam-Permission",
		"object":              "object",
		"relation":            "Viewer",
		"subjectSetNamespace": exceedLength255String,
		"subjectSetObject":    "services/410e54a2-2dcd-4675-99f1-df27a01e882d",
		"subjectSetRelation":  "operator",
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
		"message":       fmt.Sprintf("%s: %s", errc.ClientErrMsgOutOfRange, "SubjectSetNamespace"),
		"errorEntities": []any{"SubjectSetNamespace"},
		"details":       any(nil),
	}
	actualResBody := map[string]any{}
	_ = json.Unmarshal(responseData, &actualResBody)
	suite.Equal(expectedResBody, actualResBody)
}

func (suite *TestSuite) TestAPI_CreateSubjectSetRelationTuple_Failed_11() {
	name := "TestAPI_CreateSubjectSetRelationTuple_Failed: subjectSetObject exceeds max length"
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
	endpoint := "/admin/api/v1/subject-set-relation-tuples"
	body := map[string]any{
		"namespace":           "iam-Permission",
		"object":              "object",
		"relation":            "Viewer",
		"subjectSetNamespace": "iam-permission",
		"subjectSetObject":    exceedLength255String,
		"subjectSetRelation":  "operator",
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
		"message":       fmt.Sprintf("%s: %s", errc.ClientErrMsgOutOfRange, "SubjectSetObject"),
		"errorEntities": []any{"SubjectSetObject"},
		"details":       any(nil),
	}
	actualResBody := map[string]any{}
	_ = json.Unmarshal(responseData, &actualResBody)
	suite.Equal(expectedResBody, actualResBody)
}

func (suite *TestSuite) TestAPI_CreateSubjectSetRelationTuple_Failed_12() {
	name := "TestAPI_CreateSubjectSetRelationTuple_Failed: subjectSetRelation exceeds max length"
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
	endpoint := "/admin/api/v1/subject-set-relation-tuples"
	body := map[string]any{
		"namespace":           "iam-Permission",
		"object":              "object",
		"relation":            "Viewer",
		"subjectSetNamespace": "iam-permission",
		"subjectSetObject":    "services/410e54a2-2dcd-4675-99f1-df27a01e882d",
		"subjectSetRelation":  exceedLength255String,
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
		"message":       fmt.Sprintf("%s: %s", errc.ClientErrMsgOutOfRange, "SubjectSetRelation"),
		"errorEntities": []any{"SubjectSetRelation"},
		"details":       any(nil),
	}
	actualResBody := map[string]any{}
	_ = json.Unmarshal(responseData, &actualResBody)
	suite.Equal(expectedResBody, actualResBody)
}

func (suite *TestSuite) TestAPI_CreateSubjectSetRelationTuple_Failed_13() {
	name := "TestAPI_CreateSubjectSetRelationTuple_Failed: not found namespace"
	suite.T().Log(name)

	// mock
	if err := db.Instance().Exec("INSERT INTO services (id, namespace, description, created_at, updated_at, created_by, updated_by) VALUES " +
		"('410e54a2-2dcd-4675-99f1-df27a01e882d', 'iam-permission', 'IAM permission service', 0, 0, '', '')").Error; err != nil {
		panic(err)
	}
	if err := db.Instance().Exec("INSERT INTO relation_definitions " +
		"(id, service_id, namespace, relation, description, created_at, updated_at, created_by, updated_by) VALUES " +
		"('4fe1bfd3-7156-4b38-bf32-f90b57aff7fc', '410e54a2-2dcd-4675-99f1-df27a01e882d', 'iam-permission', 'viewer', 'view permission', 0, 0, '', '')," +
		"('344f382d-90aa-42ef-9876-5697f8150ba0', '410e54a2-2dcd-4675-99f1-df27a01e882d', 'iam-permission', 'operator', 'operator role', 0, 0, '', '')").
		Error; err != nil {
		panic(err)
	}
	if err := db.Instance().Exec("INSERT INTO subject_sets " +
		"(id, service_id, relation_definition_id, namespace, object, relation, description, created_at, updated_at, created_by, updated_by) VALUES " +
		"('3221d4a9-ff49-4ac3-a3c2-ec39b8eef6ca', '410e54a2-2dcd-4675-99f1-df27a01e882d', '344f382d-90aa-42ef-9876-5697f8150ba0', 'iam-permission', 'services/410e54a2-2dcd-4675-99f1-df27a01e882d', 'operator', 'description 1', 0, 0, '', '')").Error; err != nil {
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
	endpoint := "/admin/api/v1/subject-set-relation-tuples"
	body := map[string]any{
		"namespace":           "iam-Permission1",
		"object":              "object",
		"relation":            "Viewer",
		"subjectSetNamespace": "iam-permission",
		"subjectSetObject":    "services/410e54a2-2dcd-4675-99f1-df27a01e882d",
		"subjectSetRelation":  "operator",
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

func (suite *TestSuite) TestAPI_CreateSubjectSetRelationTuple_Failed_14() {
	name := "TestAPI_CreateSubjectSetRelationTuple_Failed: not found relation"
	suite.T().Log(name)

	// mock
	// mock
	if err := db.Instance().Exec("INSERT INTO services (id, namespace, description, created_at, updated_at, created_by, updated_by) VALUES " +
		"('410e54a2-2dcd-4675-99f1-df27a01e882d', 'iam-permission', 'IAM permission service', 0, 0, '', '')").Error; err != nil {
		panic(err)
	}
	if err := db.Instance().Exec("INSERT INTO relation_definitions " +
		"(id, service_id, namespace, relation, description, created_at, updated_at, created_by, updated_by) VALUES " +
		"('4fe1bfd3-7156-4b38-bf32-f90b57aff7fc', '410e54a2-2dcd-4675-99f1-df27a01e882d', 'iam-permission', 'viewer', 'view permission', 0, 0, '', '')," +
		"('344f382d-90aa-42ef-9876-5697f8150ba0', '410e54a2-2dcd-4675-99f1-df27a01e882d', 'iam-permission', 'operator', 'operator role', 0, 0, '', '')").
		Error; err != nil {
		panic(err)
	}
	if err := db.Instance().Exec("INSERT INTO subject_sets " +
		"(id, service_id, relation_definition_id, namespace, object, relation, description, created_at, updated_at, created_by, updated_by) VALUES " +
		"('3221d4a9-ff49-4ac3-a3c2-ec39b8eef6ca', '410e54a2-2dcd-4675-99f1-df27a01e882d', '344f382d-90aa-42ef-9876-5697f8150ba0', 'iam-permission', 'services/410e54a2-2dcd-4675-99f1-df27a01e882d', 'operator', 'description 1', 0, 0, '', '')").Error; err != nil {
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
	endpoint := "/admin/api/v1/subject-set-relation-tuples"
	body := map[string]any{
		"namespace":           "iam-Permission",
		"object":              "object",
		"relation":            "Viewer1",
		"subjectSetNamespace": "iam-permission",
		"subjectSetObject":    "services/410e54a2-2dcd-4675-99f1-df27a01e882d",
		"subjectSetRelation":  "operator",
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

func (suite *TestSuite) TestAPI_CreateSubjectSetRelationTuple_Failed_15() {
	name := "TestAPI_CreateSubjectSetRelationTuple_Failed: not found subjectSet"
	suite.T().Log(name)

	// mock
	// mock
	if err := db.Instance().Exec("INSERT INTO services (id, namespace, description, created_at, updated_at, created_by, updated_by) VALUES " +
		"('410e54a2-2dcd-4675-99f1-df27a01e882d', 'iam-permission', 'IAM permission service', 0, 0, '', '')").Error; err != nil {
		panic(err)
	}
	if err := db.Instance().Exec("INSERT INTO relation_definitions " +
		"(id, service_id, namespace, relation, description, created_at, updated_at, created_by, updated_by) VALUES " +
		"('4fe1bfd3-7156-4b38-bf32-f90b57aff7fc', '410e54a2-2dcd-4675-99f1-df27a01e882d', 'iam-permission', 'viewer', 'view permission', 0, 0, '', '')," +
		"('344f382d-90aa-42ef-9876-5697f8150ba0', '410e54a2-2dcd-4675-99f1-df27a01e882d', 'iam-permission', 'operator', 'operator role', 0, 0, '', '')").
		Error; err != nil {
		panic(err)
	}
	if err := db.Instance().Exec("INSERT INTO subject_sets " +
		"(id, service_id, relation_definition_id, namespace, object, relation, description, created_at, updated_at, created_by, updated_by) VALUES " +
		"('3221d4a9-ff49-4ac3-a3c2-ec39b8eef6ca', '410e54a2-2dcd-4675-99f1-df27a01e882d', '344f382d-90aa-42ef-9876-5697f8150ba0', 'iam-permission', 'services/410e54a2-2dcd-4675-99f1-df27a01e882d', 'operator', 'description 1', 0, 0, '', '')").Error; err != nil {
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
	endpoint := "/admin/api/v1/subject-set-relation-tuples"
	body := map[string]any{
		"namespace":           "iam-Permission",
		"object":              "object",
		"relation":            "Viewer",
		"subjectSetNamespace": "iam-permission1",
		"subjectSetObject":    "services/410e54a2-2dcd-4675-99f1-df27a01e882d",
		"subjectSetRelation":  "operator",
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
		"message":       fmt.Sprintf("%s: %s", errc.ClientErrMsgNotFound, constant.FieldSubjectSet),
		"errorEntities": []any{constant.FieldSubjectSet},
		"details":       any(nil),
	}
	actualResBody := map[string]any{}
	_ = json.Unmarshal(responseData, &actualResBody)
	suite.Equal(expectedResBody, actualResBody)
}

func (suite *TestSuite) TestAPI_CreateSubjectSetRelationTuple_Failed_16() {
	name := "TestAPI_CreateSubjectSetRelationTuple_Failed: tuple already exists"
	suite.T().Log(name)

	// mock
	if err := db.Instance().Exec("INSERT INTO services (id, namespace, description, created_at, updated_at, created_by, updated_by) VALUES " +
		"('410e54a2-2dcd-4675-99f1-df27a01e882d', 'iam-permission', 'IAM permission service', 0, 0, '', '')").Error; err != nil {
		panic(err)
	}
	if err := db.Instance().Exec("INSERT INTO relation_definitions " +
		"(id, service_id, namespace, relation, description, created_at, updated_at, created_by, updated_by) VALUES " +
		"('4fe1bfd3-7156-4b38-bf32-f90b57aff7fc', '410e54a2-2dcd-4675-99f1-df27a01e882d', 'iam-permission', 'viewer', 'view permission', 0, 0, '', '')," +
		"('344f382d-90aa-42ef-9876-5697f8150ba0', '410e54a2-2dcd-4675-99f1-df27a01e882d', 'iam-permission', 'operator', 'operator role', 0, 0, '', '')").
		Error; err != nil {
		panic(err)
	}
	if err := db.Instance().Exec("INSERT INTO subject_sets " +
		"(id, service_id, relation_definition_id, namespace, object, relation, description, created_at, updated_at, created_by, updated_by) VALUES " +
		"('3221d4a9-ff49-4ac3-a3c2-ec39b8eef6ca', '410e54a2-2dcd-4675-99f1-df27a01e882d', '344f382d-90aa-42ef-9876-5697f8150ba0', 'iam-permission', 'services/410e54a2-2dcd-4675-99f1-df27a01e882d', 'operator', 'description 1', 0, 0, '', '')").Error; err != nil {
		panic(err)
	}
	if err := db.Instance().Exec("INSERT INTO subject_set_relation_tuples " +
		"(id, service_id, relation_definition_id, subject_set_id, namespace, object, relation, subject_set_namespace, subject_set_object, subject_set_relation, created_at, updated_at, created_by, updated_by) VALUES " +
		"('35bf4599-c634-4abc-af0c-b87513816486', '410e54a2-2dcd-4675-99f1-df27a01e882d', '4fe1bfd3-7156-4b38-bf32-f90b57aff7fc', '3221d4a9-ff49-4ac3-a3c2-ec39b8eef6ca', 'iam-permission', 'object', 'viewer', 'iam-permission', 'services/410e54a2-2dcd-4675-99f1-df27a01e882d', 'operator', 0, 0, '', '')").Error; err != nil {
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
	endpoint := "/admin/api/v1/subject-set-relation-tuples"
	body := map[string]any{
		"namespace":           "iam-Permission",
		"object":              "object",
		"relation":            "Viewer",
		"subjectSetNamespace": "iam-permission",
		"subjectSetObject":    "services/410e54a2-2dcd-4675-99f1-df27a01e882d",
		"subjectSetRelation":  "operator",
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
		"message":       fmt.Sprintf("%s: %s", errc.ClientErrMsgDuplicate, constant.FieldSubjectSetRelationTuple),
		"errorEntities": []any{constant.FieldSubjectSetRelationTuple},
		"details":       any(nil),
	}
	actualResBody := map[string]any{}
	_ = json.Unmarshal(responseData, &actualResBody)
	suite.Equal(expectedResBody, actualResBody)
}

func (suite *TestSuite) TestAPI_CreateSubjectSetRelationTuple_Failed_17() {
	name := "TestAPI_CreateSubjectSetRelationTuple_Failed: invalid access token"
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
	endpoint := "/admin/api/v1/subject-set-relation-tuples"
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
