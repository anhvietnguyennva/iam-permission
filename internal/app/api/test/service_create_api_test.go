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

func (suite *TestSuite) TestAPI_CreateService_Successfully() {
	name := "TestAPI_CreateService_Successfully"
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
	endpoint := "/admin/api/v1/services"
	body := map[string]any{
		"namespace":   "iam-peRmisSion",
		"description": "iam service",
	}
	marshalledBody, _ := json.Marshal(body)
	req, _ := http.NewRequest(method, endpoint, bytes.NewReader(marshalledBody))
	req.Header.Add("Authorization", mockValidBearerAccessToken)
	w := httptest.NewRecorder()
	testAdminRouter.ServeHTTP(w, req)
	responseData, _ := io.ReadAll(w.Body)

	// check database
	var services []*postgres.Service
	db.Instance().Find(&services)
	suite.EqualValues(1, len(services))
	service := services[0]
	suite.EqualValues("iam-permission", service.Namespace)
	suite.EqualValues("iam service", service.Description)
	suite.EqualValues(mockSubject, service.CreatedBy)
	suite.EqualValues(mockSubject, service.UpdatedBy)

	var relationDefinitions []*postgres.RelationDefinition
	db.Instance().Find(&relationDefinitions)
	suite.EqualValues(4, len(relationDefinitions))

	relationViewer := relationDefinitions[0]
	suite.EqualValues(service.ID, relationViewer.ServiceID)
	suite.EqualValues(service.Namespace, relationViewer.Namespace)
	suite.EqualValues(constant.RelationViewer, relationViewer.Relation)
	suite.EqualValues(constant.RelationViewer, relationViewer.Description)
	suite.EqualValues(mockSubject, relationViewer.CreatedBy)
	suite.EqualValues(mockSubject, relationViewer.UpdatedBy)

	relationEditor := relationDefinitions[1]
	suite.EqualValues(service.ID, relationEditor.ServiceID)
	suite.EqualValues(service.Namespace, relationEditor.Namespace)
	suite.EqualValues(constant.RelationEditor, relationEditor.Relation)
	suite.EqualValues(constant.RelationEditor, relationEditor.Description)
	suite.EqualValues(mockSubject, relationEditor.CreatedBy)
	suite.EqualValues(mockSubject, relationEditor.UpdatedBy)

	relationOwner := relationDefinitions[2]
	suite.EqualValues(service.ID, relationOwner.ServiceID)
	suite.EqualValues(service.Namespace, relationOwner.Namespace)
	suite.EqualValues(constant.RelationOwner, relationOwner.Relation)
	suite.EqualValues(constant.RelationOwner, relationOwner.Description)
	suite.EqualValues(mockSubject, relationOwner.CreatedBy)
	suite.EqualValues(mockSubject, relationOwner.UpdatedBy)

	relationConsumer := relationDefinitions[3]
	suite.EqualValues(service.ID, relationConsumer.ServiceID)
	suite.EqualValues(service.Namespace, relationConsumer.Namespace)
	suite.EqualValues(constant.RelationConsumer, relationConsumer.Relation)
	suite.EqualValues(constant.RelationConsumer, relationConsumer.Description)
	suite.EqualValues(mockSubject, relationConsumer.CreatedBy)
	suite.EqualValues(mockSubject, relationConsumer.UpdatedBy)

	var relationConfigurations []*postgres.RelationConfiguration
	db.Instance().Find(&relationConfigurations)
	suite.EqualValues(2, len(relationConfigurations))

	configurationViewerEditor := relationConfigurations[0]
	suite.EqualValues(service.ID, configurationViewerEditor.ServiceID)
	suite.EqualValues(service.Namespace, configurationViewerEditor.Namespace)
	suite.EqualValues(relationViewer.ID, configurationViewerEditor.ParentRelationDefinitionID)
	suite.EqualValues(relationViewer.Relation, configurationViewerEditor.ParentRelation)
	suite.EqualValues(relationEditor.ID, configurationViewerEditor.ChildRelationDefinitionID)
	suite.EqualValues(relationEditor.Relation, configurationViewerEditor.ChildRelation)
	suite.EqualValues(mockSubject, configurationViewerEditor.CreatedBy)
	suite.EqualValues(mockSubject, configurationViewerEditor.UpdatedBy)

	configurationEditorOwner := relationConfigurations[1]
	suite.EqualValues(service.ID, configurationEditorOwner.ServiceID)
	suite.EqualValues(service.Namespace, configurationEditorOwner.Namespace)
	suite.EqualValues(relationEditor.ID, configurationEditorOwner.ParentRelationDefinitionID)
	suite.EqualValues(relationEditor.Relation, configurationEditorOwner.ParentRelation)
	suite.EqualValues(relationOwner.ID, configurationEditorOwner.ChildRelationDefinitionID)
	suite.EqualValues(relationOwner.Relation, configurationEditorOwner.ChildRelation)
	suite.EqualValues(mockSubject, configurationEditorOwner.CreatedBy)
	suite.EqualValues(mockSubject, configurationEditorOwner.UpdatedBy)

	// assert API response
	suite.EqualValues(http.StatusOK, w.Code)
	expectedResBody := map[string]any{
		"code":    float64(errc.ClientErrCodeOK),
		"message": errc.ClientErrMsgOK,
		"data": map[string]any{
			"id": service.ID.String(),
		},
	}
	actualResBody := map[string]any{}
	_ = json.Unmarshal(responseData, &actualResBody)
	suite.Equal(expectedResBody, actualResBody)
}

func (suite *TestSuite) TestAPI_CreateService_Failed_1() {
	name := "TestAPI_CreateService_Failed: missing namespace"
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
	endpoint := "/admin/api/v1/services"
	body := map[string]any{
		"namespace":   "",
		"description": "iam service",
	}
	marshalledBody, _ := json.Marshal(body)
	req, _ := http.NewRequest(method, endpoint, bytes.NewReader(marshalledBody))
	req.Header.Add("Authorization", mockValidBearerAccessToken)
	w := httptest.NewRecorder()
	testAdminRouter.ServeHTTP(w, req)
	responseData, _ := io.ReadAll(w.Body)

	// check database
	var services []*postgres.Service
	db.Instance().Find(&services)
	suite.EqualValues(0, len(services))

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

func (suite *TestSuite) TestAPI_CreateService_Failed_2() {
	name := "TestAPI_CreateService_Failed: missing description"
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
	endpoint := "/admin/api/v1/services"
	body := map[string]any{
		"namespace":   "iam-permission",
		"description": "",
	}
	marshalledBody, _ := json.Marshal(body)
	req, _ := http.NewRequest(method, endpoint, bytes.NewReader(marshalledBody))
	req.Header.Add("Authorization", mockValidBearerAccessToken)
	w := httptest.NewRecorder()
	testAdminRouter.ServeHTTP(w, req)
	responseData, _ := io.ReadAll(w.Body)

	// check database
	var services []*postgres.Service
	db.Instance().Find(&services)
	suite.EqualValues(0, len(services))

	// assert API response
	suite.EqualValues(http.StatusBadRequest, w.Code)
	expectedResBody := map[string]any{
		"code":          float64(errc.ClientErrCodeRequired),
		"message":       fmt.Sprintf("%s: %s", errc.ClientErrMsgRequired, "Description"),
		"errorEntities": []any{"Description"},
		"details":       any(nil),
	}
	actualResBody := map[string]any{}
	_ = json.Unmarshal(responseData, &actualResBody)
	suite.Equal(expectedResBody, actualResBody)
}

func (suite *TestSuite) TestAPI_CreateService_Failed_3() {
	name := "TestAPI_CreateService_Failed: namespace exceeds max length"
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
	endpoint := "/admin/api/v1/services"
	body := map[string]any{
		"namespace":   "hbmijewagbtraqgvvaqgwcujcgwgjzxxmtpcunkgckynxrkkctwnefppqztjijrvhzyparajjcjvtfykfkvawgnyvibnbmrnifugmyqvygvyfyykvythgpxfupnjxknnaqdjqcrkqpwkdpqihkjzxrqbybedmfanuzcbttxcevyyiwmmiwuurdwqrbdqwtfajytjackbhtkdaavmhwhidwfimqxipvkjgwemyqmwcvrhthreicruwccjgdtbjzkr",
		"description": "iam service",
	}
	marshalledBody, _ := json.Marshal(body)
	req, _ := http.NewRequest(method, endpoint, bytes.NewReader(marshalledBody))
	req.Header.Add("Authorization", mockValidBearerAccessToken)
	w := httptest.NewRecorder()
	testAdminRouter.ServeHTTP(w, req)
	responseData, _ := io.ReadAll(w.Body)

	// check database
	var services []*postgres.Service
	db.Instance().Find(&services)
	suite.EqualValues(0, len(services))

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

func (suite *TestSuite) TestAPI_CreateService_Failed_4() {
	name := "TestAPI_CreateService_Failed: description exceeds max length"
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
	endpoint := "/admin/api/v1/services"
	body := map[string]any{
		"namespace":   "iam-permission",
		"description": "hbmijewagbtraqgvvaqgwcujcgwgjzxxmtpcunkgckynxrkkctwnefppqztjijrvhzyparajjcjvtfykfkvawgnyvibnbmrnifugmyqvygvyfyykvythgpxfupnjxknnaqdjqcrkqpwkdpqihkjzxrqbybedmfanuzcbttxcevyyiwmmiwuurdwqrbdqwtfajytjackbhtkdaavmhwhidwfimqxipvkjgwemyqmwcvrhthreicruwccjgdtbjzkr",
	}
	marshalledBody, _ := json.Marshal(body)
	req, _ := http.NewRequest(method, endpoint, bytes.NewReader(marshalledBody))
	req.Header.Add("Authorization", mockValidBearerAccessToken)
	w := httptest.NewRecorder()
	testAdminRouter.ServeHTTP(w, req)
	responseData, _ := io.ReadAll(w.Body)

	// check database
	var services []*postgres.Service
	db.Instance().Find(&services)
	suite.EqualValues(0, len(services))

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

func (suite *TestSuite) TestAPI_CreateService_Failed_5() {
	name := "TestAPI_CreateService_Failed: namespace already exists"
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

	if err := db.Instance().Exec("INSERT INTO services (id, namespace, description, created_by, updated_by, created_at, updated_at) VALUES " +
		"('dd8cb2fc-a196-4071-b42e-cd656d934ad6', 'iam-permission', 'iam service', 'created_by', 'updated_by', 0, 0)").Error; err != nil {
		panic(err)
	}

	// call API
	method := "POST"
	endpoint := "/admin/api/v1/services"
	body := map[string]any{
		"namespace":   "iam-permission",
		"description": "iam service",
	}
	marshalledBody, _ := json.Marshal(body)
	req, _ := http.NewRequest(method, endpoint, bytes.NewReader(marshalledBody))
	req.Header.Add("Authorization", mockValidBearerAccessToken)
	w := httptest.NewRecorder()
	testAdminRouter.ServeHTTP(w, req)
	responseData, _ := io.ReadAll(w.Body)

	// check database
	var services []*postgres.Service
	db.Instance().Find(&services)
	suite.EqualValues(1, len(services))

	// assert API response
	suite.EqualValues(http.StatusConflict, w.Code)
	expectedResBody := map[string]any{
		"code":          float64(errc.ClientErrCodeDuplicate),
		"message":       fmt.Sprintf("%s: %s", errc.ClientErrMsgDuplicate, constant.FieldService),
		"errorEntities": []any{constant.FieldService},
		"details":       any(nil),
	}
	actualResBody := map[string]any{}
	_ = json.Unmarshal(responseData, &actualResBody)
	suite.Equal(expectedResBody, actualResBody)
}

func (suite *TestSuite) TestAPI_CreateService_Failed_6() {
	name := "TestAPI_CreateService_Failed: namespace already exists - case insensitive"
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

	if err := db.Instance().Exec("INSERT INTO services (id, namespace, description, created_by, updated_by, created_at, updated_at) VALUES " +
		"('dd8cb2fc-a196-4071-b42e-cd656d934ad6', 'iam-permission', 'iam service', 'created_by', 'updated_by', 0, 0)").Error; err != nil {
		panic(err)
	}

	// call API
	method := "POST"
	endpoint := "/admin/api/v1/services"
	body := map[string]any{
		"namespace":   "iam-permIssIon",
		"description": "iam service",
	}
	marshalledBody, _ := json.Marshal(body)
	req, _ := http.NewRequest(method, endpoint, bytes.NewReader(marshalledBody))
	req.Header.Add("Authorization", mockValidBearerAccessToken)
	w := httptest.NewRecorder()
	testAdminRouter.ServeHTTP(w, req)
	responseData, _ := io.ReadAll(w.Body)

	// check database
	var services []*postgres.Service
	db.Instance().Find(&services)
	suite.EqualValues(1, len(services))

	// assert API response
	suite.EqualValues(http.StatusConflict, w.Code)
	expectedResBody := map[string]any{
		"code":          float64(errc.ClientErrCodeDuplicate),
		"message":       fmt.Sprintf("%s: %s", errc.ClientErrMsgDuplicate, constant.FieldService),
		"errorEntities": []any{constant.FieldService},
		"details":       any(nil),
	}
	actualResBody := map[string]any{}
	_ = json.Unmarshal(responseData, &actualResBody)
	suite.Equal(expectedResBody, actualResBody)
}

func (suite *TestSuite) TestAPI_CreateService_Failed_7() {
	name := "TestAPI_CreateService_Failed: access token is invalid"
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
	endpoint := "/admin/api/v1/services"
	body := map[string]any{
		"namespace":   "iam-permission",
		"description": "iam service",
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
