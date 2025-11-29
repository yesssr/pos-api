package service

import (
	"pos-api/internal/store"

	"github.com/aws/aws-sdk-go-v2/service/s3"
)

type Service struct {
	AuthService *AuthService;
	UserService *UserService;
	ProductService *ProductService;
	CustomerService *CustomerService;
}

func New(q *store.Queries, awsClient *s3.Client, bucket string) *Service {
	return &Service{
		AuthService: NewAuthService(q),
		UserService: NewUserService(q, awsClient, bucket, "users"),
		ProductService: NewProductService(q, awsClient, bucket, "products"),
		CustomerService: NewCustomerService(q),
	}
}
