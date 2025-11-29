package handler

import (
	"pos-api/internal/service"
)

type Handler struct {
	Auth *AuthHandler;
	User *UserHandler;
	Product *ProductHandler;
	Customer *CustomerHandler;
}

func New(s *service.Service) *Handler {
	return &Handler{
		Auth: NewAuthHandler(s.AuthService),
		User: NewUserHandler(s.UserService),
		Product: NewProductHandler(s.ProductService),
		Customer: NewCustomerHandler(s.CustomerService),
	}
}
