package multiinstruction

import (
	"context"
	"testing"

	"github.com/deuna/payment/infraestructure/db/repos"
	"github.com/deuna/payment/test/mocks"
	"github.com/stretchr/testify/assert"
)

func TestMultiInstruction(t *testing.T) {
	// Arrange
	ctx := context.Background()
	dbMock, sqlDB := mocks.SetUpMockDB(t)
	fn := func(q repos.Querier) error {
		return nil
	}

	dbMock.ExpectBegin()
	dbMock.ExpectCommit()

	// Act
	err := NewMultiInstruction(sqlDB).ExecTx(ctx, fn)

	// Assert

	assert.NoError(t, err)
}

func TestMultiInstruction_ExecTx(t *testing.T) {
	// Arrange
	ctx := context.Background()
	_, sqlDB := mocks.SetUpMockDB(t)
	fn := func(q repos.Querier) error {
		return nil
	}

	// Act
	err := NewMultiInstruction(sqlDB).ExecTx(ctx, fn)

	// Assert
	assert.Error(t, err)
}

func TestMultiInstruction_ExecTx_Rollback(t *testing.T) {
	// Arrange
	ctx := context.Background()
	dbMock, sqlDB := mocks.SetUpMockDB(t)
	fn := func(q repos.Querier) error {
		return assert.AnError
	}

	dbMock.ExpectBegin()
	dbMock.ExpectRollback().WillReturnError(assert.AnError)

	// Act
	err := NewMultiInstruction(sqlDB).ExecTx(ctx, fn)

	// Assert
	assert.ErrorContains(t, err, "tx err:")
}
