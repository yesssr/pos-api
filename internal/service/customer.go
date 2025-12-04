package service

import (
	"context"
	"pos-api/internal/lib"
	"pos-api/internal/store"

	"github.com/jackc/pgx/v5/pgtype"
)

type CustomerService struct {
	q *store.Queries;
}

func NewCustomerService(q *store.Queries) *CustomerService {
	return &CustomerService{ q: q };
}

func(s *CustomerService) CreateCustomer(ctx context.Context, name, phone, address string) (store.Customer, error) {
	c, err := s.q.CreateCustomer(ctx, store.CreateCustomerParams{
		Name:    name,
		Phone:   pgtype.Text{String: phone, Valid: true},
		Address: pgtype.Text{String: address, Valid: true},
	});
	if err != nil {
		return store.Customer{}, err;
	}
	return c, nil;
}

func(s *CustomerService) ListCustomers(ctx context.Context, l, o int, search string) ([]store.Customer, int, error) {
	args := store.ListCustomersParams{
		Limit:  int32(l),
		Offset: int32(o),
		Column3: pgtype.Text{String: search, Valid: true},
	}
	list, _ := s.q.ListCustomers(ctx, args);
	c, _ := s.q.CountCustomers(ctx, pgtype.Text{String: search, Valid: true});
	t := int(c);

	totalPages := lib.GetTotalPages(t, l);
	return list, totalPages, nil;
}

func(s *CustomerService) GetCustomerByID(ctx context.Context, id pgtype.UUID) (store.Customer, error) {
	c, err := s.q.GetCustomerByID(ctx, id);
	if err != nil {
		return store.Customer{}, err;
	}
	return c, nil;
}

func(s *CustomerService) UpdateCustomer(ctx context.Context, id pgtype.UUID, name, phone, address string) (store.Customer, error) {
	c, err := s.q.UpdateCustomer(ctx, store.UpdateCustomerParams{
		ID:      id,
		Name:    name,
		Phone:   pgtype.Text{String: phone, Valid: true},
		Address: pgtype.Text{String: address, Valid: true},
	});
	if err != nil {
		return store.Customer{}, err;
	}
	return c, nil;
}

func(s *CustomerService) DeleteCustomer(ctx context.Context, id pgtype.UUID) (store.Customer, error) {
	c, err := s.q.DeleteCustomer(ctx, id);
	if err != nil {
		return store.Customer{}, err;
	}
	return c, nil;
}

func (s *CustomerService) GetTotalCustomer(ctx context.Context) (int, error) {
	t, err := s.q.CountCustomers(ctx, pgtype.Text{Valid: true});
	if err != nil {
		return 0, err;
	}
	return int(t), nil;
}
