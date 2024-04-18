package transaction

import (
	"context"
	"database/sql"
	"errors"

	"github.com/deuna/payment/domain/clients"
	"github.com/deuna/payment/domain/domainerrors"
	"github.com/deuna/payment/domain/models"

	"github.com/deuna/payment/infraestructure/db/repos"
	"github.com/google/uuid"
)

// Transaction is the interface that wraps the transaction methods
type Transaction interface {
	Place(context.Context, float64, string) (*uuid.UUID, error)
	Refund(context.Context, string, string) error
	ListTx(context.Context, string) ([]repos.Tx, error)
}

type transaction struct {
	bankRepo clients.Bank
	queries  repos.Querier
}

// NewService creates a new transaction service
func NewService(queries repos.Querier, bankRepo clients.Bank) Transaction {
	return &transaction{
		bankRepo: bankRepo,
		queries:  queries,
	}
}

// ListTx returns a list of transactions
func (t *transaction) ListTx(ctx context.Context, merchantID string) ([]repos.Tx, error) {
	if err := t.listTxReqValidator(merchantID); err != nil {
		return nil, err
	}

	merchantIDUUID := uuid.MustParse(merchantID)

	resp, err := t.queries.ListTransactionsByMechant(ctx, merchantIDUUID)
	if err != nil {
		if err == sql.ErrNoRows {
			return []repos.Tx{}, nil
		}

		return nil, domainerrors.NewInternalServerError(err)
	}

	return resp, nil
}

// Place creates a transaction and places the amount in the bank
func (t *transaction) Place(ctx context.Context, amount float64, merchantID string) (*uuid.UUID, error) {
	if err := t.placeReqValidator(amount, merchantID); err != nil {
		return nil, err
	}

	merchantIDUUID := uuid.MustParse(merchantID)

	resp, err := t.queries.CreateTransaction(ctx, repos.CreateTransactionParams{
		ID:         uuid.New(),
		Amount:     amount,
		MerchantID: merchantIDUUID,
		TxType:     models.Initiate.Value(),
	})
	if err != nil {
		return nil, domainerrors.NewInternalServerError(err)
	}

	err = t.bankPlace(ctx, resp.ID, amount, merchantIDUUID)
	if err != nil {
		if errors.As(err, &domainerrors.UnauthorizedTranError{}) {
			if err := t.updateTranType(ctx, resp.ID, models.Rejected, resp.UpdatedAt); err != nil {
				return nil, err
			}
			return nil, err
		}
		if err := t.updateTranType(ctx, resp.ID, models.Cancelled, resp.UpdatedAt); err != nil {
			return nil, err
		}
		return nil, domainerrors.NewInternalServerError(err)
	}

	if err := t.updateTranType(ctx, resp.ID, models.Complete, resp.UpdatedAt); err != nil {
		return nil, err
	}

	return &resp.ID, nil
}

func (t *transaction) Refund(ctx context.Context, id, merchantID string) error {
	if err := t.refundReqValidator(id, merchantID); err != nil {
		return err
	}

	newUUID := uuid.MustParse(id)

	tx, err := t.queries.GetTransaction(ctx, newUUID)
	if err != nil {
		if err == sql.ErrNoRows {
			return domainerrors.NewNoTransactionError(err)
		}
		return domainerrors.NewInternalServerError(err)
	}

	if tx.MerchantID.String() != merchantID {
		return domainerrors.NewNoTransactionError(nil)
	}

	if err := t.bankRefund(ctx, tx.ID); err != nil {
		return domainerrors.NewInternalServerError(err)
	}

	if err := t.updateTranType(ctx, tx.ID, models.Refunded, tx.UpdatedAt); err != nil {
		return err
	}

	return nil
}
