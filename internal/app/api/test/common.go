package test

import "github.com/anhvietnguyennva/iam-go-sdk/oauth/entity"

//go:generate mockgen -destination=iam_sdk_mock.go -package test . IIAMSDK
type IIAMSDK interface {
	ParseBearerJWT(bearerTokenJWTString string) (*entity.AccessToken, error)
	ParseJWT(tokenJWTString string) (*entity.AccessToken, error)
}

const (
	mockValidBearerAccessToken   = "Bearer valid"
	mockInvalidBearerAccessToken = "Bearer invalid"
	mockSubject                  = "3aa6ae4f-2f4b-4ef9-a6cb-6a33c52b98f5"
)

func mockAccessToken() *entity.AccessToken {
	return &entity.AccessToken{
		Subject: mockSubject,
	}
}
