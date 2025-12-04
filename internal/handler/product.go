package handler

import (
	"net/http"
	"pos-api/internal/lib"
	"pos-api/internal/middleware"
	"pos-api/internal/service"
	"strconv"

	"github.com/jackc/pgx/v5"
)

type createProductInput struct {
	Name  string `json:"name" validate:"required,min=2"`
	Price string `json:"price" validate:"required,gt=0"`
	Stock int32  `json:"stock" validate:"required,gte=0"`
	ImageUrl    string `json:"image_url,omitempty" validate:"omitempty"`
}

type updateProductInput struct {
	Name        string `json:"name" validate:"required,min=2"`
	Price       string `json:"price" validate:"required"`
	Stock       int32  `json:"stock" validate:"required"`
	ImageUrl    string `json:"image_url" validate:"required"`
	IsActiveStr string `json:"is_active" validate:"required"`
}

type ProductHandler struct {
	s *service.ProductService
}

func NewProductHandler(s *service.ProductService) *ProductHandler {
	return &ProductHandler{
		s: s,
	}
}

func (h *ProductHandler) CreateProduct(w http.ResponseWriter, r *http.Request) {
	name := r.FormValue("name");
	priceStr := r.FormValue("price");
	stockStr := r.FormValue("stock");
	imageUrl := "default.png";

	stock64, _ := strconv.ParseInt(stockStr, 10, 32);
	b := &createProductInput{
		Name:  name,
		Price: priceStr,
		Stock: int32(stock64),
		ImageUrl: imageUrl,
	}

	if err := lib.ValidateStruct(b); err != nil {
		lib.SendErrorResponse(w, err, b);
		return;
	}

	priceInt, err := strconv.ParseInt(priceStr, 10, 64);
	if err != nil {
		lib.SendErrorResponse(w, &lib.AppError{
			Message:    "Invalid price format",
			StatusCode: http.StatusBadRequest,
		}, nil);
		return;
	}

	p, err := h.s.CreateProduct(r, b.Name, int(priceInt), b.Stock)
	if err != nil {
		lib.SendErrorResponse(w, err, nil);
		return;
	}

	lib.SendResponse(w, http.StatusCreated, "Successfully added product", p, nil, nil)
}

func (h *ProductHandler) ListProducts(w http.ResponseWriter, r *http.Request) {
	p := lib.GetPagination(r);
	q, _ := middleware.GetQueryFromCtx(r);
	offset := (p.CurrentPage - 1) * p.PerPage;

	list, totalPages, err := h.s.ListProducts(r.Context(), p.PerPage, offset, q.OrderBy, q.OrderDir, q.Search);
	if err != nil {
		lib.SendErrorResponse(w, err, nil);
		return;
	}

	p.TotalPages = &totalPages;
	lib.SendResponse(w, http.StatusOK, "List of products", list, p, nil);
}

func (h *ProductHandler) ListProductsActive(w http.ResponseWriter, r *http.Request) {
	p := lib.GetPagination(r);
	q, _ := middleware.GetQueryFromCtx(r);
	offset := (p.CurrentPage - 1) * p.PerPage;

	list, totalPages, err := h.s.ListProductsActive(r.Context(), p.PerPage, offset, q.OrderBy, q.OrderDir, q.Search);
	if err != nil {
		lib.SendErrorResponse(w, err, nil);
		return;
	}

	p.TotalPages = &totalPages;
	lib.SendResponse(w, http.StatusOK, "List of products active", list, p, nil);
}

func (h *ProductHandler) GetProduct(w http.ResponseWriter, r *http.Request) {
	id, _ := middleware.GetIdFromCtx(r)
	pd, err := h.s.GetProduct(r.Context(), id)
	if err != nil && err != pgx.ErrNoRows{
		lib.SendErrorResponse(w, err, nil)
		return
	}
	lib.SendResponse(w, http.StatusOK, "Product details", pd, nil, nil)
}

func (h *ProductHandler) UpdateProduct(w http.ResponseWriter, r *http.Request) {
	name := r.FormValue("name")
	priceStr := r.FormValue("price")
	stockStr := r.FormValue("stock")
	isActiveStr := r.FormValue("is_active")
	submittedImageUrl := r.FormValue("image_url")

	stock64, _ := strconv.ParseInt(stockStr, 10, 32)
	b := &updateProductInput{
		Name:        name,
		Price:       priceStr,
		Stock:       int32(stock64),
		ImageUrl:    submittedImageUrl,
		IsActiveStr: isActiveStr,
	}

	if err := lib.ValidateStruct(b); err != nil {
		lib.SendErrorResponse(w, err, b)
		return
	}

	priceInt, err := strconv.ParseInt(priceStr, 10, 64);
	if err != nil {
		lib.SendErrorResponse(w, &lib.AppError{
			Message:    "Invalid price format",
			StatusCode: http.StatusBadRequest,
		}, nil);
		return;
	}

	isActive, err := strconv.ParseBool(b.IsActiveStr)
	if err != nil {
		lib.SendErrorResponse(w, &lib.AppError{
			Message:    "Invalid is_active value",
			StatusCode: http.StatusBadRequest,
		}, nil)
		return
	}

	pd, err := h.s.UpdateProduct(r, b.Name, int(priceInt), b.Stock, isActive, b.ImageUrl)
	if err != nil {
		lib.SendErrorResponse(w, err, nil)
		return
	}

	lib.SendResponse(w, http.StatusOK, "Successfully updated product", pd, nil, nil)
}

func (h *ProductHandler) DeleteProduct(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	id, _ := middleware.GetIdFromCtx(r)

	pd, err := h.s.DeleteProduct(ctx, id)
	if err != nil && err != pgx.ErrNoRows {
		lib.SendErrorResponse(w, err, nil)
		return
	}

	lib.SendResponse(w, http.StatusOK, "Successfully deleted product", pd, nil, nil)
}

func (h *ProductHandler) GetTotalProduct(w http.ResponseWriter, r *http.Request) {
	t, err := h.s.GetTotalProduct(r.Context())
	if err != nil {
		lib.SendErrorResponse(w, err, nil)
		return
	}
	lib.SendResponse(w, http.StatusOK, "Total products", t, nil, nil)
}
