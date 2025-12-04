package handler

import (
	"math/big"
	"net/http"
	"pos-api/internal/configuration"
	"pos-api/internal/lib"
	"pos-api/internal/middleware"
	"pos-api/internal/service"
	"pos-api/internal/store"
	"strings"
	"time"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgtype"
)

type Item struct {
	IDProduct   string  	`json:"id_product" validate:"required,uuid4"`
	Qty    int32    		  `json:"qty" validate:"required,gt=0"`
}

type createTransactionInput struct {
	IDCustomer           *string    		`json:"id_customer,omitempty" validate:"omitempty,uuid4"`
	PaymentMethod        string  				`json:"payment_method" validate:"required,oneof=cash qris credit debit"`
	Items 							[]Item  				`json:"items" validate:"required,dive"`
}

type updateTransactionInput struct {
	ID string 							`json:"id" validate:"required"`;
	Status string 					`json:"status" validate:"required"`;
}

type TransactionHandler struct {
	s *service.TransactionService;
	ws *configuration.Hub;
}

func NewTransactionHandler(s *service.TransactionService, ws *configuration.Hub) *TransactionHandler {
	return &TransactionHandler{s: s, ws: ws};
}

func(h *TransactionHandler) CreateTransaction(w http.ResponseWriter, r *http.Request) {
	b := &createTransactionInput{};
	if !lib.ValidateJSON(w, r, b) {
		return;
	}

	if err := lib.ValidateStruct(b); err != nil {
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

	jakarta, _ := time.LoadLocation("Asia/Jakarta");
	t := time.Now().In(jakarta);

	trxParams := store.CreateTransactionParams{
		IDUser: pgtype.UUID{Bytes: userId, Valid: true},
		IDCustomer: pgtype.UUID{Valid: false},
		Total: *lib.IntToPgNumeric(0),
		PaymentMethod: store.PaymentMethod(b.PaymentMethod),
		IDTransactionGateway: pgtype.Text{Valid: false},
		Date: pgtype.Date{Time: t, Valid: true},
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
			Price: pgtype.Numeric{Int: big.NewInt(0), Valid: false},
			Subtotal: pgtype.Numeric{Int: big.NewInt(0), Valid: false},
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

func(h *TransactionHandler) ListTransactions(w http.ResponseWriter, r *http.Request) {
	q, _ := middleware.GetQueryFromCtx(r);
	p := lib.GetPagination(r);
	offset := (p.CurrentPage - 1) * p.PerPage;

	list, t, err := h.s.ListTransactions(r.Context(), q.StartAt, q.EndAt, p.PerPage, offset);
	if err != nil {
		lib.SendErrorResponse(w, err, nil);
		return;
	}

	p.TotalPages = &t;
	lib.SendResponse(w, http.StatusOK, "List of transactions", list, p, nil);
}

func(h *TransactionHandler) WebHookXendit(w http.ResponseWriter, r *http.Request) {
	b := &updateTransactionInput{};
	if !lib.ValidateJSON(w, r, b) {
		return;
	}

	if err := lib.ValidateStruct(b); err != nil {
		lib.SendErrorResponse(w, err, b);
		return;
	}
	args := store.UpdateStatusByPaymentIdParams{
		IDTransactionGateway: pgtype.Text{String: b.ID, Valid: true},
		PaymentStatus:store.PaymentStatus(strings.ToLower(b.Status)),
	}
	trx, err := h.s.UpdateTrxStatus(r.Context(), args);
	if err != nil {
		h.ws.NotifyUser(trx.IDUser.String(), map[string]any{
			"type": "transaction_update",
			"error": err.Error(),
		});
		return;
	}
	h.ws.NotifyUser(trx.IDUser.String(), map[string]any{
		"type": "transaction_update",
		"data": trx,
	});
	w.WriteHeader(http.StatusOK);
}

func (h *TransactionHandler) SalesByPeriods(w http.ResponseWriter, r *http.Request) {
	q, _ := middleware.GetQueryFromCtx(r);

	list, err := h.s.ListSalesByPeriods(r.Context(), q.Period, q.StartAt, q.EndAt);
	if err != nil {
		lib.SendErrorResponse(w, err, nil);
		return;
	}

	lib.SendResponse(w, http.StatusOK, "Sales report", list, nil, nil);
}
