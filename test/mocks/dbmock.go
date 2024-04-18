package mocks

import (
	"database/sql"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/stretchr/testify/assert"
)

// SetUpMockDB sets up a mock database for testing
func SetUpMockDB(t *testing.T) (dbmock sqlmock.Sqlmock, db *sql.DB) {
	db, dbmock, err := sqlmock.New()
	assert.NoError(t, err)
	return dbmock, db
}
