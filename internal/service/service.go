package service

import (
	"pos-api/internal/store"

	"github.com/aws/aws-sdk-go-v2/service/s3"
)

type Service struct {
	UserService *UserService;
	AuthService *AuthService;
}

func New(q *store.Queries, awsClient *s3.Client, bucket string) *Service {
	return &Service{
		UserService: NewUserService(q, awsClient, bucket, "users"),
		AuthService: NewAuthService(q),
	}
}
