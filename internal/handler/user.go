package handler

import (
	"net/http"
	"pos-api/internal/lib"
	"pos-api/internal/middleware"
	"pos-api/internal/service"
	"pos-api/internal/store"
	"strconv"
)
const userDir = "users";
type CreateUserInput struct {
  Username string `json:"username" validate:"required,min=3,username"`
  Password string `json:"password" validate:"required,min=6"`
  Role     store.Roles `json:"role" validate:"required,oneof=admin cashier"`
}

type UpdateUserInput struct {
	Username string `json:"username" validate:"required,min=3,username"`
	Role     store.Roles `json:"role" validate:"required,oneof=admin cashier"`
	Password string `json:"password" validate:"omitempty,min=6"`
	ImageUrl string `json:"image_url" validate:"required"`
	IsActive *bool `json:"is_active,omitempty" validate:"omitempty"`
}

type UserHandler struct {
	s *service.UserService;
}

func NewUserHandler(s *service.UserService) *UserHandler {
	return &UserHandler{
		s: s,
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
		lib.SendErrorResponse(w, err, b);
		return;
	}

	u, err := h.s.CreateUser(r, b.Username, b.Password, b.Role);
	if err != nil {
		lib.SendErrorResponse(w, err, nil);
		return;
	}
	lib.SendResponse(w, http.StatusCreated, "Successfully added user", u, nil, nil);
}

func (h *UserHandler) ListUsers(w http.ResponseWriter, r *http.Request) {
	p := lib.GetPagination(r);
	offset := (p.CurrentPage - 1) * p.PerPage;
	l, t, err := h.s.ListUsers(r.Context(), p.PerPage, offset);
	if err != nil {
		lib.SendErrorResponse(w, err, nil);
		return;
	}
	p.TotalPages = &t;
	lib.SendResponse(w, http.StatusOK, "List of users", l, p, nil)
}

func (h *UserHandler) GetUser(w http.ResponseWriter, r *http.Request) {
	id, _ := middleware.GetIdFromCtx(r);
	u, err := h.s.GetUserById(r.Context(), id);
	if err != nil {
		lib.SendErrorResponse(w, err, nil);
		return;
	}
	lib.SendResponse(w, http.StatusOK, "User details", u, nil, nil);
}

func (h *UserHandler) UpdateUser(w http.ResponseWriter, r *http.Request) {
	username := r.FormValue("username");
	password := r.FormValue("password");
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
		Password: password,
	}

	if err := lib.ValidateStruct(b); err != nil {
		lib.SendErrorResponse(w, err, b);
		return;
	}

	u, err := h.s.UpdateUser(r, b.Username, b.Password, b.IsActive, b.Role, b.ImageUrl);
	if err != nil {
		lib.SendErrorResponse(w, err, nil);
		return;
	}
	lib.SendResponse(w, http.StatusOK, "Successfully updated user", u, nil, nil);
}

func (h *UserHandler) DeleteUser(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context();
	id, _ := middleware.GetIdFromCtx(r);

	u, err := h.s.DeleteUser(ctx, id);
	if err != nil {
		lib.SendErrorResponse(w, err, nil);
		return;
	}
	lib.SendResponse(w, http.StatusOK, "Successfully deleted user", u, nil, nil);
}

func (h *UserHandler) GetTotalUser(w http.ResponseWriter, r *http.Request) {
	t, err := h.s.GetTotalUser(r.Context());
	if err != nil {
		lib.SendErrorResponse(w, err, nil);
		return;
	}
	lib.SendResponse(w, http.StatusOK, "Total users", t, nil, nil);
}
