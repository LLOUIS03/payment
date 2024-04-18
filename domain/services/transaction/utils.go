package transaction

import (
	"context"
	"fmt"
	"strings"
	"time"

	"emperror.dev/errors"
	"github.com/deuna/payment/domain/domainerrors"
	"github.com/deuna/payment/domain/models"
	"github.com/deuna/payment/infraestructure/db/repos"
	"github.com/deuna/payment/utils/retry"
	"github.com/google/uuid"
)

const bankRefundMaxAttempts = 3
const bankRefundSleep = 100 * time.Millisecond

const bankPlaceMaxAttempts = 3
const bankPlaceSleep = 100 * time.Millisecond

func (t *transaction) placeReqValidator(amount float64, merchantID string) error {
	result := domainerrors.NewValidationError()

	if amount <= 0 {
		result.AddErrorMsg("Amount must be greater than 0")
	}

	if _, err := uuid.Parse(merchantID); err != nil {
		result.AddErrorMsg("Invalid merchant id")
	}

	if result.IsFailure() {
		return result
	}

	return nil
}

func (t *transaction) listTxReqValidator(merchantID string) error {
	result := domainerrors.NewValidationError()

	if _, err := uuid.Parse(merchantID); err != nil {
		result.AddErrorMsg("Invalid merchant id")
	}

	if result.IsFailure() {
		return result
	}

	return nil
}

func (t *transaction) refundReqValidator(id, merchantID string) error {
	result := domainerrors.NewValidationError()

	if _, err := uuid.Parse(id); err != nil {
		result.AddErrorMsg("Invalid transaction id")
	}

	if _, err := uuid.Parse(merchantID); err != nil {
		result.AddErrorMsg("Invalid merchant id")
	}

	if result.IsFailure() {
		return result
	}

	return nil
}

func (t *transaction) updateTranType(ctx context.Context, id uuid.UUID, tranType models.TransactionType, updatedAt time.Time) error {
	_, err := t.queries.UpdateTransactionType(ctx, repos.UpdateTransactionTypeParams{
		ID:        id,
		TxType:    tranType.Value(),
		UpdatedAt: updatedAt,
	})
	if err != nil {
		return domainerrors.NewInternalServerError(err)
	}
	return nil
}

func (t *transaction) bankRefund(ctx context.Context, id uuid.UUID) error {
	resp, retryErr := retry.NewRetry(func() error {
		if err := t.bankRepo.Refund(ctx, id); err != nil {
			return err
		}
		return nil
	}).WithMaxAttempts(bankRefundMaxAttempts).
		WithSleep(bankRefundSleep).
		Run()

	if retryErr != nil {
		strBuilder := strings.Builder{}
		fmt.Println("Value:", resp)
		for i := 0; i < resp.Attempts(); i++ {
			fmt.Println("index: ", i, "Value:", resp[i])
			strBuilder.WriteString(fmt.Sprintf("Attemp %d: %s\n", i, resp[i].Error()))
		}
		return errors.Wrapf(retryErr, "Calling %v", strBuilder.String())
	}
	return nil
}

func (t *transaction) bankPlace(ctx context.Context, id uuid.UUID, amount float64, merchantID uuid.UUID) error {
	resp, retryErr := retry.NewRetry(func() error {
		if err := t.bankRepo.Place(ctx, id, amount, merchantID); err != nil {
			return err
		}
		return nil
	}).WithMaxAttempts(bankPlaceMaxAttempts).
		WithSleep(bankPlaceSleep).
		Run()

	if retryErr != nil {
		fmt.Println("retryErr.Error():", retryErr.Error(), "Vas", errors.As(retryErr, &domainerrors.UnauthorizedTranError{}))
		strBuilder := strings.Builder{}
		for i := 0; i < resp.Attempts(); i++ {
			strBuilder.WriteString(fmt.Sprintf("Attemp %d: %s", i, resp[i].Error()))
		}
		return errors.Wrapf(retryErr, "Calling %v", strBuilder.String())
	}
	return nil
}
