package transaction

import (
	"context"
	"errors"
	"testing"

	"github.com/deuna/payment/test/mocks"
	"github.com/golang/mock/gomock"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
)

func Test_refundReqValidator(t *testing.T) {
	t.Run("Valid Transaction ID and Merchant ID", func(t *testing.T) {
		// Arrange
		id := uuid.New().String()
		merchantID := uuid.New().String()
		svc := transaction{}

		// Act
		err := svc.refundReqValidator(id, merchantID)

		// Assert
		assert.NoError(t, err)
	})

	t.Run("Invalid Transaction ID", func(t *testing.T) {
		// Arrange
		id := "invalid-id"
		merchantID := uuid.New().String()
		svc := transaction{}

		// Act
		err := svc.refundReqValidator(id, merchantID)

		// Assert
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "Invalid transaction id")
	})

	t.Run("Invalid Merchant ID", func(t *testing.T) {
		// Arrange
		id := uuid.New().String()
		merchantID := "invalid-id"
		svc := transaction{}

		// Act
		err := svc.refundReqValidator(id, merchantID)

		// Assert
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "Invalid merchant id")
	})
}

func Test_placeReqValidator(t *testing.T) {
	t.Run("Valid Amount and Merchant ID", func(t *testing.T) {
		// Arrange
		amount := 100.0
		merchantID := uuid.New().String()
		svc := transaction{}

		// Act
		err := svc.placeReqValidator(amount, merchantID)

		// Assert
		assert.NoError(t, err)
	})

	t.Run("Invalid Amount", func(t *testing.T) {
		// Arrange
		amount := -10.0
		merchantID := uuid.New().String()
		svc := transaction{}

		// Act
		err := svc.placeReqValidator(amount, merchantID)

		// Assert
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "Amount must be greater than 0")
	})

	t.Run("Invalid Merchant ID", func(t *testing.T) {
		// Arrange
		amount := 100.0
		merchantID := "invalid-id"
		svc := transaction{}

		// Act
		err := svc.placeReqValidator(amount, merchantID)

		// Assert
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "Invalid merchant id")
	})
}
func Test_listTxReqValidator(t *testing.T) {
	t.Run("Valid Merchant ID", func(t *testing.T) {
		// Arrange
		merchantID := uuid.New().String()
		svc := transaction{}

		// Act
		err := svc.listTxReqValidator(merchantID)

		// Assert
		assert.NoError(t, err)
	})

	t.Run("Invalid Merchant ID", func(t *testing.T) {
		// Arrange
		merchantID := "invalid-id"
		svc := transaction{}

		// Act
		err := svc.listTxReqValidator(merchantID)

		// Assert
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "Invalid merchant id")
	})
}
func Test_bankRefund(t *testing.T) {
	ctx := context.Background()
	id := uuid.New()

	t.Run("Successful refund", func(t *testing.T) {
		// Arrange
		ctrl := gomock.NewController(t)
		mockBankRepo := mocks.NewMockBank(ctrl)
		mockBankRepo.EXPECT().Refund(ctx, id).Return(nil).AnyTimes()

		svc := &transaction{
			bankRepo: mockBankRepo,
		}

		// Act
		err := svc.bankRefund(ctx, id)

		// Assert
		assert.NoError(t, err)
	})

	t.Run("Failed refund with retries", func(t *testing.T) {
		// Arrange
		exptErr := errors.New("refund failed")
		ctrl := gomock.NewController(t)
		mockBankRepo := mocks.NewMockBank(ctrl)
		mockBankRepo.EXPECT().Refund(ctx, id).Return(exptErr).AnyTimes()

		svc := &transaction{
			bankRepo: mockBankRepo,
		}

		// Act
		err := svc.bankRefund(ctx, id)

		// Assert
		assert.Error(t, err)
		assert.ErrorContains(t, err, exptErr.Error())
	})
}
func Test_bankPlace(t *testing.T) {
	ctx := context.Background()
	id := uuid.New()
	amount := 100.0
	merchantID := uuid.New()

	t.Run("Successful place", func(t *testing.T) {
		// Arrange
		ctrl := gomock.NewController(t)
		mockBankRepo := mocks.NewMockBank(ctrl)
		mockBankRepo.EXPECT().Place(ctx, id, amount, merchantID).Return(nil).AnyTimes()

		svc := &transaction{
			bankRepo: mockBankRepo,
		}

		// Act
		err := svc.bankPlace(ctx, id, amount, merchantID)

		// Assert
		assert.NoError(t, err)
	})

	t.Run("Failed place with retries", func(t *testing.T) {
		// Arrange
		exptErr := errors.New("place failed")
		ctrl := gomock.NewController(t)
		mockBankRepo := mocks.NewMockBank(ctrl)
		mockBankRepo.EXPECT().Place(ctx, id, amount, merchantID).Return(exptErr).AnyTimes()

		svc := &transaction{
			bankRepo: mockBankRepo,
		}

		// Act
		err := svc.bankPlace(ctx, id, amount, merchantID)

		// Assert
		assert.Error(t, err)
		assert.ErrorContains(t, err, exptErr.Error())
	})
}
