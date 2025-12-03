package service

import (
	"context"

	"github.com/xendit/xendit-go/v7"
	"github.com/xendit/xendit-go/v7/invoice"
)

type PaymentService struct {
	x *xendit.APIClient;
}

func NewPaymentService(x *xendit.APIClient) *PaymentService {
	return &PaymentService{ x: x };
}

func(s *PaymentService) CreateInvoice(ctx context.Context, trxId string, amount float64) (*invoice.Invoice, error) {
	createInvoiceRequest := *invoice.NewCreateInvoiceRequest(trxId, amount);
	resp, _, err := s.x.InvoiceApi.CreateInvoice(ctx).
		CreateInvoiceRequest(createInvoiceRequest).
		Execute();

	if err != nil {
		return nil, err;
	}

	return resp, nil;
}
