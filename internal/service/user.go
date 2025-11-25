package service

import (
	"context"
	"fmt"
	"net/http"
	"path"
	"pos-api/internal/lib"
	"pos-api/internal/middleware"
	"pos-api/internal/store"

	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/jackc/pgx/v5/pgtype"
)

type UserService struct {
	q *store.Queries;
	awsClient *s3.Client;
	bucket string;
	dir string;
}

func NewUserService(q *store.Queries, awsClient *s3.Client, bucket string, dir string) *UserService {
	return &UserService{
		q: q,
		awsClient: awsClient,
		bucket: bucket,
		dir: dir,
	}
}

func(s *UserService) CreateUser(
	r *http.Request,
	username string,
	password string,
	role store.Roles,
) (store.User, error) {
	ctx := r.Context();
	if _, err := s.q.GetUserByUsername(ctx, username); err == nil {
		return store.User{}, &lib.AppError{
			Message: "Username already exists",
			StatusCode: http.StatusBadRequest,
		};
	}

	hash, _ := lib.HashPassword(password);
	imageUrl, err := lib.UploadHandler(r, s.awsClient, s.bucket, s.dir, username, nil);
	if err != nil {
		return store.User{}, err;
	}

	u, err := s.q.CreateUser(ctx, store.CreateUserParams{
		Username: username,
		Password: hash,
		Role:     role,
		ImageUrl: imageUrl,
	});

	if err != nil {
		return store.User{}, err;
	}

	return u, nil;
}

func(s *UserService) ListUsers(ctx context.Context, l, o int) ([]store.ListUsersRow, int, error) {
	limit := l;

	args := store.ListUsersParams{
		Limit:  int32(limit),
		Offset: int32(o),
	}
	list, _ := s.q.ListUsers(ctx, args);
	t, _ := s.q.CountUsers(ctx);
	totalPages := (int(t) + limit - 1) / limit;
	return list, totalPages, nil;
}

func(s *UserService) GetUserById(ctx context.Context, id pgtype.UUID) (store.GetUserByIdRow, error) {
	u, _ := s.q.GetUserById(ctx, id);
	return u, nil
}

func(s *UserService) UpdateUser(
	r *http.Request,
	username string,
	password string,
	isActive *bool,
	role store.Roles,
	submittedImgUrl string,
) (store.User, error) {
	ctx := r.Context();
	id, err := middleware.GetIdFromCtx(r);

	if u, err := s.q.GetUserByUsername(ctx, username); err == nil && u.ID != id {
		return store.User{}, &lib.AppError{
			Message: "Username already exists",
			StatusCode: http.StatusBadRequest,
		};
	}

	pgBool := lib.BoolPtrToPgBool(isActive);
	params := store.UpdateUserParams{
		ID: id,
		Username: username,
		Role: role,
		ImageUrl: submittedImgUrl,
		IsActive: pgBool,
	}

	kDel := path.Base(submittedImgUrl);
	imageUrl, err := lib.UploadHandler(r, s.awsClient, s.bucket, s.dir, username, &kDel);
	if err != nil {
		return store.User{}, err;
	}
	if imageUrl != "" {
		params.ImageUrl = imageUrl;
	}
	u, _ := s.q.UpdateUser(ctx, params);
	if password != "" {
		pass, err := lib.HashPassword(password);
		if err != nil {
			return store.User{}, err;
		}
		args := store.UpdatePassParams{
			ID: id,
			Password: pass,
		}
		s.q.UpdatePass(ctx, args);
	}

	return u, nil;
}

func (s *UserService) DeleteUser(ctx context.Context, id pgtype.UUID) (store.User, error) {
	u, err := s.q.DeleteUser(ctx, id)
	if err != nil {
		return store.User{}, err
	}

	k := path.Base(u.ImageUrl)
	if err := lib.DeleteImageFromCloud(ctx, s.awsClient, s.bucket, s.dir, k); err != nil {
		fmt.Println("Warning: failed to delete image:", err)
	}
	return u, nil
}

func (s *UserService) GetTotalUser(ctx context.Context) (int, error) {
	t, err := s.q.CountUsers(ctx);
	if err != nil {
		return 0, err;
	}
	return int(t), nil;
}
