package handler

import (
	"net/http"
	"pos-api/internal/lib"
	"pos-api/internal/middleware"
	"pos-api/internal/service"
)

type CustomerHandler struct {
	s *service.CustomerService;
}

func NewCustomerHandler(s *service.CustomerService) *CustomerHandler {
	return &CustomerHandler{s: s}
}

type createCustomerInput struct {
	Name  string `json:"name" validate:"required,min=3"`;
	Phone string `json:"phone" validate:"omitempty,min=10,max=15"`;
	Address string `json:"address" validate:"omitempty,min=10"`;
}

type updateCustomerInput struct {
	Name  string `json:"name" validate:"required,min=3"`;
	Phone string `json:"phone" validate:"required,min=10,max=15"`;
	Address string `json:"address" validate:"required,min=10"`;
}

func(h *CustomerHandler) CreateCustomer(w http.ResponseWriter, r *http.Request) {
	b := &createCustomerInput{};
	if !lib.ValidateJSON(w, r, b) {
		return;
	}

	if err := lib.ValidateStruct(b); err != nil {
		lib.SendErrorResponse(w, err, b);
		return;
	}

	c, err := h.s.CreateCustomer(r.Context(), b.Name, b.Phone, b.Address);
	if err != nil {
		lib.SendErrorResponse(w, err, nil);
		return;
	}

	lib.SendResponse(w, http.StatusCreated, "Successfully added customer", c, nil, nil);
}

func(h *CustomerHandler) ListCustomers(w http.ResponseWriter, r *http.Request) {
	q, _ := middleware.GetQueryFromCtx(r);
	p := lib.GetPagination(r);
	offset := (p.CurrentPage - 1) * p.PerPage;
	c, t, err := h.s.ListCustomers(r.Context(), p.PerPage, offset, q.Search);
	if err != nil {
		lib.SendErrorResponse(w, err, nil);
		return;
	}

	totalPages := lib.GetTotalPages(t, p.PerPage);
	p.TotalPages = &totalPages;

	lib.SendResponse(w, http.StatusOK, "List of customers", c, p, nil);
}

func (h *CustomerHandler) GetCustomer(w http.ResponseWriter, r *http.Request) {
	id, err := middleware.GetIdFromCtx(r);
	if err != nil {
		lib.SendErrorResponse(w, err, nil);
		return;
	}

	c, err := h.s.GetCustomerByID(r.Context(), id);
	if err != nil {
		lib.SendErrorResponse(w, err, nil);
		return;
	}

	lib.SendResponse(w, http.StatusOK, "Get customer", c, nil, nil);
}

func (h *CustomerHandler) UpdateCustomer(w http.ResponseWriter, r *http.Request) {
	id, err := middleware.GetIdFromCtx(r);
	if err != nil {
		lib.SendErrorResponse(w, err, nil);
		return;
	}

	var b updateCustomerInput;
	if !lib.ValidateJSON(w, r, &b) {
		return;
	}

	if err := lib.ValidateStruct(&b); err != nil {
		lib.SendErrorResponse(w, err, b);
		return;
	}

	c, err := h.s.UpdateCustomer(r.Context(), id, b.Name, b.Phone, b.Address);
	if err != nil {
		lib.SendErrorResponse(w, err, nil);
		return;
	}

	lib.SendResponse(w, http.StatusOK, "Successfully updated customer", c, nil, nil);
}

func (h *CustomerHandler) DeleteCustomer(w http.ResponseWriter, r *http.Request) {
	id, err := middleware.GetIdFromCtx(r);
	if err != nil {
		lib.SendErrorResponse(w, err, nil);
		return;
	}

	c, err := h.s.DeleteCustomer(r.Context(), id);
	if err != nil {
		lib.SendErrorResponse(w, err, nil);
		return;
	}

	lib.SendResponse(w, http.StatusOK, "Successfully deleted customer", c, nil, nil);
}
