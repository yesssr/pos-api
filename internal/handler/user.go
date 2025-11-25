package handler

import (
	"net/http"
	"os"
	"path"
	"pos-api/internal/lib"
	"pos-api/internal/middleware"
	"pos-api/internal/store"
	"strconv"

	"github.com/aws/aws-sdk-go-v2/service/s3"
)
const userDir = "users";
type CreateUserInput struct {
  Username string `json:"username" validate:"required,min=3"`
  Password string `json:"password" validate:"required,min=6"`
  Role     store.Roles `json:"role" validate:"required,oneof=admin cashier"`
}

type UpdateUserInput struct {
	Username string `json:"username" validate:"required,min=3"`
	Role     store.Roles `json:"role" validate:"required,oneof=admin cashier"`
	ImageUrl string `json:"image_url" validate:"required"`
	IsActive *bool `json:"is_active,omitempty" validate:"omitempty"`
}

type UserHandler struct {
	queries *store.Queries;
	r2Client *s3.Client;
}

func NewUserHandler(q *store.Queries, r *s3.Client) *UserHandler {
	return &UserHandler{
		queries: q,
		r2Client: r,
	}
}

func (h *UserHandler) CreateUser(w http.ResponseWriter, r *http.Request) {
	username := r.FormValue("username");
  password := r.FormValue("password");
  role := r.FormValue("role");

	b := &CreateUserInput{
		Username: username,
		Password: password,
		Role: store.Roles(role),
	};

	if err := lib.ValidateStruct(b); err != nil {
		lib.SendErrorResponse(w, err);
		return;
	}

	imageUrl, err := lib.UploadHandler(r, h.r2Client, userDir, b.Username, nil);

	ctx := r.Context()
	if _, err := h.queries.GetUserByUsername(ctx, b.Username); err == nil {
		lib.SendErrorResponse(w, &lib.AppError{
			Message: "Username already exists",
			StatusCode: http.StatusBadRequest,
		});
		return;
	}

	pass, _ := lib.HashPassword(b.Password);

	args := store.CreateUserParams{
		Username: b.Username,
		Password: pass,
		Role:     b.Role,
		ImageUrl: imageUrl,
	}

	u, err := h.queries.CreateUser(ctx, args);
	if err != nil {
		lib.SendErrorResponse(w, err);
		return;
	}

	lib.SendResponse(w, http.StatusCreated, "Successfully added user", u, nil, nil)
}

func (h *UserHandler) ListUsers(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context();
	p := lib.GetPagination(r);
	limit := p.PerPage;
	offset := (p.CurrentPage - 1) * p.PerPage;

	args := store.ListUsersParams{
		Limit:  int32(limit),
		Offset: int32(offset),
	}

	list, _ := h.queries.ListUsers(ctx, args);
	t, _ := h.queries.CountUsers(ctx);
	totalPages := (int(t) + limit - 1) / limit;
	p.TotalPages = &totalPages;
	lib.SendResponse(w, http.StatusOK, "List of users", list, p, nil)
}

func (h *UserHandler) GetUser(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context();
	id, err := middleware.GetIdFromCtx(r);
	if err != nil {
		lib.SendErrorResponse(w, err);
		return;
	}
	u, _ := h.queries.GetUserById(ctx, id);
	lib.SendResponse(w, http.StatusOK, "User details", u, nil, nil);
}

func (h *UserHandler) UpdateUser(w http.ResponseWriter, r *http.Request) {
	username := r.FormValue("username");
  strIsActive := r.FormValue("is_active");
  role := r.FormValue("role");
  submittedImageUrl := r.FormValue("image_url");
  var isAct *bool
	if temp, err := strconv.ParseBool(strIsActive); err != nil {
		isAct = nil;
	} else {
		isAct = &temp;
	}

 	b := &UpdateUserInput{
		Username: username,
		Role: store.Roles(role),
		ImageUrl: submittedImageUrl,
		IsActive: isAct,
	}

	if err := lib.ValidateStruct(b); err != nil {
		lib.SendErrorResponse(w, err);
		return;
	}

	id, err := middleware.GetIdFromCtx(r);
	if err != nil {
		lib.SendErrorResponse(w, err);
		return;
	}

	pgBool := lib.BoolPtrToPgBool(b.IsActive);
	kDel := path.Base(submittedImageUrl);
	imageUrl, err := lib.UploadHandler(r, h.r2Client, userDir, b.Username, &kDel);
	if err != nil {
		lib.SendErrorResponse(w, err);
		return;
	}

	if imageUrl != "http://default.png" {
		b.ImageUrl = imageUrl;
	}

	args := store.UpdateUserParams{
		ID: id,
		Username: b.Username,
		Role: b.Role,
		IsActive: pgBool,
		ImageUrl: b.ImageUrl,
	}

	ctx := r.Context();
	u, _ := h.queries.UpdateUser(ctx, args);
	lib.SendResponse(w, http.StatusOK, "Successfully updated user", u, nil, nil);
}

func (h *UserHandler) DeleteUser(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context();
	id, err := middleware.GetIdFromCtx(r);
	if err != nil {
		lib.SendErrorResponse(w, err);
		return;
	}
	u, err := h.queries.DeleteUser(ctx, id);
	if err != nil {
		lib.SendErrorResponse(w, err);
		return;
	}
	bucket := os.Getenv("R2_BUCKET_NAME");
	k := path.Base(u.ImageUrl);
	if err := lib.DeleteImageFromCloud(ctx, h.r2Client, bucket, userDir, k); err != nil {
		lib.SendErrorResponse(w, err);
		return;
	}
	lib.SendResponse(w, http.StatusOK, "Successfully deleted user", u, nil, nil);
}
