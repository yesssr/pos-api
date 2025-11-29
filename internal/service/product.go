package service

import (
	"context"
	"fmt"
	"math/big"
	"net/http"
	"path"
	"pos-api/internal/lib"
	"pos-api/internal/middleware"
	"pos-api/internal/store"

	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/jackc/pgx/v5/pgtype"
)

type ProductService struct {
	q         *store.Queries
	awsClient *s3.Client
	bucket    string
	dir       string
}

func NewProductService(q *store.Queries, awsClient *s3.Client, bucket string, dir string) *ProductService {
	return &ProductService{
		q:         q,
		awsClient: awsClient,
		bucket:    bucket,
		dir:       dir,
	}
}

func (s *ProductService) CreateProduct(
	r *http.Request,
	name string,
	price *big.Int,
	stock int32,
) (store.Product, error) {
	ctx := r.Context()

	imageUrl, err := lib.UploadHandler(r, s.awsClient, s.bucket, s.dir, name, nil);
	if err != nil {
		return store.Product{}, err;
	}

	p := store.CreateProductParams{
		Name:     name,
		Price:    *lib.IntToPgNumeric(price),
		Stock:    stock,
		ImageUrl: imageUrl,
	}

	product, err := s.q.CreateProduct(ctx, p);
	if err != nil {
		return store.Product{}, err;
	}

	return product, nil;
}

func (s *ProductService) ListProducts(ctx context.Context, l, o int, oBy, oDir, search string,) (any, int, error) {
	var list any;
	if oDir != "desc" {
		args := store.ListProductsAscParams{
			Limit:  int32(l),
			Offset: int32(o),
			Column3: oBy,
			Column4: pgtype.Text{String: search, Valid: true},
		}
		p, err := s.q.ListProductsAsc(ctx, args);
		if err != nil {
			return nil, 0, err;
		}
		list = p;
	} else {
		args := store.ListProductsDescParams{
			Limit:  int32(l),
			Offset: int32(o),
			Column3: oBy,
			Column4: pgtype.Text{String: search, Valid: true},
		}
		p, err := s.q.ListProductsDesc(ctx, args);
		if err != nil {
			return nil, 0, err;
		}
		list = p;
	}

	c, _ := s.q.CountProducts(ctx, pgtype.Text{String: search, Valid: true});
	t := int(c);

	totalPages := lib.GetTotalPages(t, l);
	return list, totalPages, nil;
}

func (s *ProductService) UpdateProduct(
	r *http.Request,
	name string,
	price *big.Int,
	stock int32,
	isActive bool,
	submittedImgUrl string,
) (store.Product, error) {
	ctx := r.Context();
	id, _ := middleware.GetIdFromCtx(r);

	params := store.UpdateProductParams{
		ID:       id,
		Name:     name,
		Price:    *lib.IntToPgNumeric(price),
		Stock:    stock,
		ImageUrl: submittedImgUrl,
		IsActive: isActive,
	}

	kDel := path.Base(submittedImgUrl);
	imageUrl, err := lib.UploadHandler(r, s.awsClient, s.bucket, s.dir, name, &kDel)
	if err != nil {
		return store.Product{}, err;
	}
	if imageUrl != "" {
		params.ImageUrl = imageUrl;
	}

	p, err := s.q.UpdateProduct(ctx, params)
	if err != nil {
		return store.Product{}, err;
	}

	return p, nil;
}

func (s *ProductService) DeleteProduct(ctx context.Context, id pgtype.UUID) (store.Product, error) {
	p, err := s.q.DeleteProduct(ctx, id);
	if err != nil {
		return store.Product{}, err;
	}

	k := path.Base(p.ImageUrl)
	if err := lib.DeleteImageFromCloud(ctx, s.awsClient, s.bucket, s.dir, k); err != nil {
		fmt.Println("Warning: failed to delete image:", err);
	}

	return p, nil;
}

func (s *ProductService) GetProduct(ctx context.Context, id pgtype.UUID) (store.GetProductRow, error) {
	p, err := s.q.GetProduct(ctx, id);
	if err != nil {
		return store.GetProductRow{}, err;
	}
	return p, nil;
}

func (s *ProductService) GetTotalProduct(ctx context.Context) (int, error) {
	t, err := s.q.CountProducts(ctx, pgtype.Text{Valid: true});
	if err != nil {
		return 0, err;
	}
	return int(t), nil;
}
