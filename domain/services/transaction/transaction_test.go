package transaction

import (
	"context"
	"database/sql"
	"errors"
	"testing"
	"time"

	"github.com/deuna/payment/domain/domainerrors"
	"github.com/deuna/payment/infraestructure/db/repos"
	"github.com/deuna/payment/test/mocks"
	"github.com/golang/mock/gomock"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
)

func TestListTx_HappyPath(t *testing.T) {
	// Arrange
	ctx := context.Background()
	ctrl := gomock.NewController(t)

	bankRepo := mocks.NewMockBank(ctrl)
	queries := mocks.NewMockQuerier(ctrl)

	uuidID := uuid.New()

	transactions := []repos.Tx{
		{
			ID:         uuid.New(),
			MerchantID: uuid.New(),
			Amount:     100.00,
		}, {
			ID:         uuid.New(),
			MerchantID: uuid.New(),
			Amount:     100.00,
		},
	}

	queries.EXPECT().ListTransactionsByMechant(gomock.Any(), gomock.Any()).Return(transactions, nil).AnyTimes()

	service := NewService(queries, bankRepo)

	// Act
	resp, err := service.ListTx(ctx, uuidID.String())

	// Assert
	assert.NoError(t, err)
	assert.Equal(t, 2, len(resp))
}

func TestListTx_InvalidUUID(t *testing.T) {
	// Arrange
	ctx := context.Background()
	ctrl := gomock.NewController(t)

	bankRepo := mocks.NewMockBank(ctrl)
	queries := mocks.NewMockQuerier(ctrl)

	service := NewService(queries, bankRepo)

	// Act
	resp, err := service.ListTx(ctx, "testig")

	// Assert
	assert.Nil(t, resp)
	assert.ErrorContains(t, err, "Invalid merchant id")
}

func TestListTx_NoData(t *testing.T) {
	// Arrange
	ctx := context.Background()
	ctrl := gomock.NewController(t)

	bankRepo := mocks.NewMockBank(ctrl)
	queries := mocks.NewMockQuerier(ctrl)

	uuidID := uuid.New()

	queries.EXPECT().ListTransactionsByMechant(gomock.Any(), gomock.Any()).Return(nil, sql.ErrNoRows).AnyTimes()

	service := NewService(queries, bankRepo)

	// Act
	resp, err := service.ListTx(ctx, uuidID.String())

	// Assert
	assert.NoError(t, err)
	assert.Equal(t, 0, len(resp))
}

func TestListTx_Error(t *testing.T) {
	// Arrange
	ctx := context.Background()
	ctrl := gomock.NewController(t)

	bankRepo := mocks.NewMockBank(ctrl)
	queries := mocks.NewMockQuerier(ctrl)

	uuidID := uuid.New()
	expectedErr := errors.New("test error")

	queries.EXPECT().ListTransactionsByMechant(gomock.Any(), gomock.Any()).Return(nil, expectedErr).AnyTimes()

	service := NewService(queries, bankRepo)

	// Act
	resp, err := service.ListTx(ctx, uuidID.String())

	// Assert
	assert.Nil(t, resp)
	assert.ErrorContains(t, err, "internal server error")
}

func TestPlace_HappyPath(t *testing.T) {
	// Arrange
	ctx := context.Background()
	ctrl := gomock.NewController(t)

	bankRepo := mocks.NewMockBank(ctrl)
	queries := mocks.NewMockQuerier(ctrl)

	uuidID := uuid.New()

	bankRepo.EXPECT().Place(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).Return(nil).AnyTimes()
	queries.EXPECT().CreateTransaction(gomock.Any(), gomock.Any()).Return(repos.CreateTransactionRow{
		ID:        uuid.New(),
		UpdatedAt: time.Now(),
	}, nil).AnyTimes()

	queries.EXPECT().UpdateTransactionType(gomock.Any(), gomock.Any()).Return(repos.UpdateTransactionTypeRow{}, nil).AnyTimes()

	service := NewService(queries, bankRepo)

	// Act
	resp, err := service.Place(ctx, 100.00, uuidID.String())

	// Assert
	assert.NoError(t, err)
	assert.NotNil(t, resp)
}

func TestPlace_InvalidParameters(t *testing.T) {
	// Arrange
	ctx := context.Background()
	ctrl := gomock.NewController(t)

	bankRepo := mocks.NewMockBank(ctrl)
	queries := mocks.NewMockQuerier(ctrl)

	// uuidID := uuid.New()

	service := NewService(queries, bankRepo)

	// Act
	resp, err := service.Place(ctx, 0, "test")

	// Assert
	assert.Nil(t, resp)
	assert.ErrorContains(t, err, "Amount must be greater than 0;Invalid merchant id")

}

func TestPlace_ErrorCreatingTransaction(t *testing.T) {
	// Arrange
	ctx := context.Background()
	ctrl := gomock.NewController(t)

	bankRepo := mocks.NewMockBank(ctrl)
	queries := mocks.NewMockQuerier(ctrl)

	uuidID := uuid.New()
	expectedErr := errors.New("test error")

	bankRepo.EXPECT().Place(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).Return(nil).AnyTimes()
	queries.EXPECT().CreateTransaction(gomock.Any(), gomock.Any()).Return(repos.CreateTransactionRow{
		ID:        uuid.New(),
		UpdatedAt: time.Now(),
	}, expectedErr).AnyTimes()

	queries.EXPECT().UpdateTransactionType(gomock.Any(), gomock.Any()).Return(repos.UpdateTransactionTypeRow{}, nil).AnyTimes()

	service := NewService(queries, bankRepo)

	// Act
	resp, err := service.Place(ctx, 100.00, uuidID.String())

	// Assert
	assert.Nil(t, resp)
	assert.ErrorContains(t, err, "internal server error")
}

func TestPlace_BankPlaceErr(t *testing.T) {
	// Arrange
	ctx := context.Background()
	ctrl := gomock.NewController(t)

	bankRepo := mocks.NewMockBank(ctrl)
	queries := mocks.NewMockQuerier(ctrl)

	uuidID := uuid.New()

	expectedErr := errors.New("test error")
	bankRepo.EXPECT().Place(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).Return(expectedErr).AnyTimes()
	queries.EXPECT().CreateTransaction(gomock.Any(), gomock.Any()).Return(repos.CreateTransactionRow{
		ID:        uuid.New(),
		UpdatedAt: time.Now(),
	}, nil).AnyTimes()

	queries.EXPECT().UpdateTransactionType(gomock.Any(), gomock.Any()).Return(repos.UpdateTransactionTypeRow{}, nil).AnyTimes()

	service := NewService(queries, bankRepo)

	// Act
	resp, err := service.Place(ctx, 100.00, uuidID.String())

	// Assert
	assert.ErrorContains(t, err, "internal server error")
	assert.Nil(t, resp)
}

func TestPlace_Unauthorized(t *testing.T) {
	// Arrange
	ctx := context.Background()
	ctrl := gomock.NewController(t)

	bankRepo := mocks.NewMockBank(ctrl)
	queries := mocks.NewMockQuerier(ctrl)

	uuidID := uuid.New()
	expectedErr := domainerrors.NewUnauthorizedTranError(nil)

	bankRepo.EXPECT().Place(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).Return(expectedErr).AnyTimes()
	queries.EXPECT().CreateTransaction(gomock.Any(), gomock.Any()).Return(repos.CreateTransactionRow{
		ID:        uuid.New(),
		UpdatedAt: time.Now(),
	}, nil).AnyTimes()

	queries.EXPECT().UpdateTransactionType(gomock.Any(), gomock.Any()).Return(repos.UpdateTransactionTypeRow{}, nil).AnyTimes()

	service := NewService(queries, bankRepo)

	// Act
	resp, err := service.Place(ctx, 100.00, uuidID.String())

	// Assert
	assert.ErrorContains(t, err, "unable to authorize transaction")
	assert.Nil(t, resp)
}

func TestRefund_HappyPath(t *testing.T) {
	// Arrange
	ctx := context.Background()
	ctrl := gomock.NewController(t)

	bankRepo := mocks.NewMockBank(ctrl)
	queries := mocks.NewMockQuerier(ctrl)

	mercahntUUID := uuid.New()
	txUUID := uuid.New()

	queries.EXPECT().GetTransaction(gomock.Any(), gomock.Any()).Return(repos.Tx{
		MerchantID: mercahntUUID,
	}, nil).AnyTimes()

	queries.EXPECT().UpdateTransactionType(gomock.Any(), gomock.Any()).Return(repos.UpdateTransactionTypeRow{}, nil).AnyTimes()
	bankRepo.EXPECT().Refund(gomock.Any(), gomock.Any()).Return(nil).AnyTimes()

	service := NewService(queries, bankRepo)

	// Act
	err := service.Refund(ctx, txUUID.String(), mercahntUUID.String())

	// Assert
	assert.NoError(t, err)
}

func TestRefund_InvalidParam(t *testing.T) {
	// Arrange
	ctx := context.Background()
	ctrl := gomock.NewController(t)

	bankRepo := mocks.NewMockBank(ctrl)
	queries := mocks.NewMockQuerier(ctrl)

	service := NewService(queries, bankRepo)

	// Act
	err := service.Refund(ctx, "test", "test")

	// Assert
	assert.ErrorContains(t, err, "Invalid transaction id;Invalid merchant id")
}

func TestRefund_NoTransaction(t *testing.T) {
	// Arrange
	ctx := context.Background()
	ctrl := gomock.NewController(t)

	bankRepo := mocks.NewMockBank(ctrl)
	queries := mocks.NewMockQuerier(ctrl)

	mercahntUUID := uuid.New()
	txUUID := uuid.New()

	expectedErr := errors.New("test error")

	queries.EXPECT().GetTransaction(gomock.Any(), gomock.Any()).Return(repos.Tx{
		MerchantID: mercahntUUID,
	}, expectedErr).AnyTimes()

	service := NewService(queries, bankRepo)

	// Act
	err := service.Refund(ctx, txUUID.String(), mercahntUUID.String())

	// Assert
	assert.ErrorContains(t, err, "internal server error")
}

func TestRefund_NoTransactionFound(t *testing.T) {
	// Arrange
	ctx := context.Background()
	ctrl := gomock.NewController(t)

	bankRepo := mocks.NewMockBank(ctrl)
	queries := mocks.NewMockQuerier(ctrl)

	mercahntUUID := uuid.New()
	txUUID := uuid.New()

	expectedErr := sql.ErrNoRows

	queries.EXPECT().GetTransaction(gomock.Any(), gomock.Any()).Return(repos.Tx{
		MerchantID: mercahntUUID,
	}, expectedErr).AnyTimes()

	service := NewService(queries, bankRepo)

	// Act
	err := service.Refund(ctx, txUUID.String(), mercahntUUID.String())

	// Assert
	assert.ErrorContains(t, err, "no transaction found")
}

func TestRefund_DiffMerchandat(t *testing.T) {
	// Arrange
	ctx := context.Background()
	ctrl := gomock.NewController(t)

	bankRepo := mocks.NewMockBank(ctrl)
	queries := mocks.NewMockQuerier(ctrl)

	mercahntUUID := uuid.New()
	txUUID := uuid.New()

	queries.EXPECT().GetTransaction(gomock.Any(), gomock.Any()).Return(repos.Tx{
		MerchantID: uuid.New(),
	}, nil).AnyTimes()

	queries.EXPECT().UpdateTransactionType(gomock.Any(), gomock.Any()).Return(repos.UpdateTransactionTypeRow{}, nil).AnyTimes()
	bankRepo.EXPECT().Refund(gomock.Any(), gomock.Any()).Return(nil).AnyTimes()

	service := NewService(queries, bankRepo)

	// Act
	err := service.Refund(ctx, txUUID.String(), mercahntUUID.String())

	// Assert
	assert.ErrorContains(t, err, "no transaction found")
}
