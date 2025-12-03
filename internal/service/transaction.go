package service

import (
	"context"
	"pos-api/internal/lib"
	"pos-api/internal/store"
	"time"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgtype"
	"github.com/jackc/pgx/v5/pgxpool"
)

type TransactionService struct {
	dbx *pgxpool.Pool;
	q *store.Queries;
	pay *PaymentService;
}

func NewTransactionService(q *store.Queries, dbx *pgxpool.Pool, pay *PaymentService) *TransactionService {
	return &TransactionService{q: q, dbx: dbx, pay: pay};
}

func(s *TransactionService) ListTransactions(ctx context.Context, start, end time.Time, l, o int) ([]store.ListTransactionsRow, error) {
	p := store.ListTransactionsParams{
		Date: pgtype.Timestamptz{Time: start, Valid: true},
		Date_2: pgtype.Timestamptz{Time: end, Valid: true},
		Limit: int32(l),
		Offset: int32(o),
	}
	list, err := s.q.ListTransactions(ctx, p);
	if err != nil {
		return []store.ListTransactionsRow{}, err;
	}
	return list, nil;
}

func(s *TransactionService) CreateTransaction(ctx context.Context, header store.CreateTransactionParams, detail []store.CreateDetailTransactionParams) (store.Transaction, *string, error) {
	tx, err := s.dbx.BeginTx(ctx, pgx.TxOptions{});
	if err != nil {
	 	return store.Transaction{}, nil, err;
	}
 	defer tx.Rollback(ctx);

  qtx := s.q.WithTx(tx);

  trx, err := qtx.CreateTransaction(ctx, header);
  if err != nil {
 		return store.Transaction{}, nil, err;
  }

  for _, d := range detail {
 		d.IDTransaction = trx.ID;
   	_, err := qtx.CreateDetailTransaction(ctx, d);
    if err != nil {
    	return store.Transaction{}, nil, err;
    }
  }

 	if err := tx.Commit(ctx);  err != nil {
 		return store.Transaction{}, nil, err;
 	}

  args := store.UpdateTransactionStatusParams{
		ID: trx.ID,
		PaymentStatus: trx.PaymentStatus,
		PaymentMethod: trx.PaymentMethod,
		IDTransactionGateway: pgtype.Text{Valid: false},
	};

  if header.PaymentMethod == "cash" {
  	args.PaymentStatus = "paid";
  	t, err := s.q.UpdateTransactionStatus(ctx, args);
   	if err != nil {
	 		return store.Transaction{}, nil, err;
	 	}
		return t, nil, nil;
  }

  amount := lib.NumericToFloat(trx.Total);
  res, err := s.pay.CreateInvoice(ctx, trx.ID.String(), amount);

  if err != nil {
		return store.Transaction{}, nil, err;
	}

	invoiceUrl := res.GetInvoiceUrl();
	args.PaymentStatus = "pending";
	args.IDTransactionGateway = pgtype.Text{
		String: res.GetId(),
		Valid:  true,
	};

	t, err := s.q.UpdateTransactionStatus(ctx, args);
 	if err != nil {
 		return store.Transaction{}, nil, err;
 	}

	trx = t;

  return trx, &invoiceUrl, nil;
}
