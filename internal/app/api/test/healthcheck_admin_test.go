package test

import (
	"io"
	"net/http"
	"net/http/httptest"
)

func (suite *TestSuite) TestAPI_HealthcheckAdmin_Successfully() {
	name := "TestAPI_HealthcheckAdmin_Successfully"
	suite.T().Log(name)

	// call API
	method := "GET"
	endpoint := "/api/v1/admin/health"
	req, _ := http.NewRequest(method, endpoint, nil)
	w := httptest.NewRecorder()
	testAdminRouter.ServeHTTP(w, req)
	responseData, _ := io.ReadAll(w.Body)

	// assert response
	suite.Equal(`"OK"`, string(responseData))
}
