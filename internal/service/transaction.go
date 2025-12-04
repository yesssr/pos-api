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

func(s *TransactionService) ListTransactions(ctx context.Context, start, end time.Time, l, o int) ([]store.ListTransactionsRow, int, error) {
	p := store.ListTransactionsParams{
		Date: pgtype.Date{Time: start, Valid: true},
		Date_2: pgtype.Date{Time: end, Valid: true},
		Limit: int32(l),
		Offset: int32(o),
	}
	list, err := s.q.ListTransactions(ctx, p);
	if err != nil {
		return []store.ListTransactionsRow{}, 0, err;
	}
	c, _ := s.q.CountTransactions(ctx);
	t := int(c);

	totalPages := lib.GetTotalPages(t, l);
	return list, totalPages, nil;
}

func(s *TransactionService) CreateTransaction(ctx context.Context, header store.CreateTransactionParams, detail []store.CreateDetailTransactionParams) (store.Transaction, *string, error) {
	tx, err := s.dbx.BeginTx(ctx, pgx.TxOptions{});
	if err != nil {
	 	return store.Transaction{}, nil, err;
	}
 	defer tx.Rollback(ctx);

  total := 0;
  qtx := s.q.WithTx(tx);

  trx, err := qtx.CreateTransaction(ctx, header);
  if err != nil {
 		return store.Transaction{}, nil, err;
  }

  for _, d := range detail {
  	p, e := qtx.GetProductForUpdate(ctx, d.IDProduct);
   	if e != nil {
    	return store.Transaction{}, nil, e;
    }

    if p.Stock < d.Qty {
			return store.Transaction{}, nil, &lib.AppError{
				Message: "Stok product " + p.Name + " tidak cukup",
				StatusCode: 400,
			};
	 	}

		if _, err := qtx.UpdateProductStock(ctx, store.UpdateProductStockParams{
			ID: d.IDProduct,
			Stock: p.Stock - d.Qty,
		}); err != nil {
			return store.Transaction{}, nil, err;
		}

 		d.IDTransaction = trx.ID;
   	d.Price = p.Price;
    raw, _ := p.Price.Int64Value()
    price := raw.Int64
    subtotal := price * int64(d.Qty)

    d.Subtotal = *lib.IntToPgNumeric(int(subtotal));
   	_, err := qtx.CreateDetailTransaction(ctx, d);
    if err != nil {
    	return store.Transaction{}, nil, err;
    }
    total += int(subtotal);
  }

  args := store.UpdateTransactionStatusParams{
		ID: trx.ID,
		PaymentStatus: trx.PaymentStatus,
		PaymentMethod: trx.PaymentMethod,
		IDTransactionGateway: pgtype.Text{Valid: false},
		Total: *lib.IntToPgNumeric(total),
	};

  if header.PaymentMethod == "cash" {
  	args.PaymentStatus = "paid";
  } else {
  	args.PaymentStatus = "pending";
  }

  t, err := qtx.UpdateTransactionStatus(ctx, args);
	if err != nil {
		return store.Transaction{}, nil, err;
	}

 	if err := tx.Commit(ctx);  err != nil {
 		return store.Transaction{}, nil, err;
 	}

  var invoiceUrl string;
  amount := lib.NumericToFloat(t.Total);
  if header.PaymentMethod != "cash" {
	  res, err := s.pay.CreateInvoice(ctx, t.ID.String(), amount);
	  if err != nil {
			return store.Transaction{}, nil, err;
		}
		invoiceUrl = res.GetInvoiceUrl();
		args.IDTransactionGateway = pgtype.Text{
			String: res.GetId(),
			Valid:  true,
		};
		uTr, _ := s.q.UpdateTransactionStatus(ctx, args);
		t = uTr;
  }

	trx = t;
  return trx, &invoiceUrl, nil;
}

func (s *TransactionService) UpdateTrxStatus(ctx context.Context, args store.UpdateStatusByPaymentIdParams) (store.Transaction, error) {
	trx, err := s.q.UpdateStatusByPaymentId(ctx, args);
	if err != nil {
		return store.Transaction{}, err;
	}
	return trx, nil;
}
