package service

import (
	"pos-api/internal/store"

	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/xendit/xendit-go/v7"
)

type Service struct {
	AuthService *AuthService;
	UserService *UserService;
	ProductService *ProductService;
	CustomerService *CustomerService;
	TransactionService *TransactionService;
	Payment *PaymentService;
}

func New(q *store.Queries, awsClient *s3.Client, bucket string, dbx *pgxpool.Pool, x *xendit.APIClient) *Service {
	pay := NewPaymentService(x);
	return &Service{
		AuthService: NewAuthService(q),
		UserService: NewUserService(q, awsClient, bucket, "users"),
		ProductService: NewProductService(q, awsClient, bucket, "products"),
		CustomerService: NewCustomerService(q),
		Payment: pay,
		TransactionService: NewTransactionService(q, dbx, pay),
	}
}
