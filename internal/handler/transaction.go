package handler

import (
	"net/http"
	"pos-api/internal/lib"
	"pos-api/internal/middleware"
	"pos-api/internal/service"
	"pos-api/internal/store"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgtype"
)

type Item struct {
	IDProduct   string  	`json:"id_product" validate:"required,uuid4"`
	Qty    int32    		  `json:"qty" validate:"required,gt=0"`
	Price  int  		  		`json:"price" validate:"required,gt=0"`
	Subtotal int 					`json:"subtotal" validate:"required,gt=0"`
}

type createTransactionInput struct {
	IDCustomer           *string    		`json:"id_customer,omitempty" validate:"omitempty,uuid4"`
	Total                int 						`json:"total" validate:"required,gt=0"`
	PaymentMethod        string  				`json:"payment_method" validate:"required,oneof=cash qris credit debit"`
	Items 							[]Item  				`json:"items" validate:"required,dive,required"`
}

type TransactionHandler struct {
	s *service.TransactionService;
}

func NewTransactionHandler(s *service.TransactionService) *TransactionHandler {
	return &TransactionHandler{s: s}
}

func(h *TransactionHandler) CreateTransaction(w http.ResponseWriter, r *http.Request) {
	var b createTransactionInput;
	if !lib.ValidateJSON(w, r, &b) {
		return;
	}

	if err := lib.ValidateStruct(&b); err != nil {
		lib.SendErrorResponse(w, err, b);
		return;
	}

	u, err := middleware.GetUserPayload(r);
	if err != nil {
		lib.SendErrorResponse(w, err, nil);
		return;
	}

	userId, err := uuid.Parse(u.Id);
	if err != nil {
		lib.SendErrorResponse(w, err, nil);
		return;
	}

	trxParams := store.CreateTransactionParams{
		IDUser: pgtype.UUID{Bytes: userId, Valid: true},
		IDCustomer: pgtype.UUID{Valid: false},
		Total: *lib.IntToPgNumeric(b.Total),
		PaymentMethod: store.PaymentMethod(b.PaymentMethod),
		PaymentStatus: store.PaymentStatus("pending"),
		IDTransactionGateway: pgtype.Text{Valid: false},
	}

	if b.IDCustomer != nil {
		custId, err := uuid.Parse(*b.IDCustomer);
		if err != nil {
			lib.SendErrorResponse(w, err, nil);
			return;
		}
		trxParams.IDCustomer = pgtype.UUID{Bytes: custId, Valid: true}
	}

	var items []store.CreateDetailTransactionParams;
	for _, item := range b.Items {
		productId, err := uuid.Parse(item.IDProduct);
		if err != nil {
			lib.SendErrorResponse(w, err, nil);
			return;
		}
		items = append(items, store.CreateDetailTransactionParams{
			IDTransaction: pgtype.UUID{Valid: false},
			IDProduct: pgtype.UUID{Bytes: productId, Valid: true},
			Qty: item.Qty,
			Price: *lib.IntToPgNumeric(item.Price),
			Subtotal: *lib.IntToPgNumeric(item.Subtotal),
		})
	}

	trx, invoiceUrl, err := h.s.CreateTransaction(r.Context(), trxParams, items);
	if err != nil {
		lib.SendErrorResponse(w, err, nil);
		return;
	}

	res := map[string]any{
		"invoice_url": invoiceUrl,
		"transaction": trx,
	}
	lib.SendResponse(w, http.StatusCreated, "Successfully created transaction", res, nil, nil);
}
