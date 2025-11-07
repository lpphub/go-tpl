package dbs

import (
	"context"
	"errors"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

func setupTestDB(t *testing.T) (*gorm.DB, sqlmock.Sqlmock) {
	db, mock, err := sqlmock.New()
	require.NoError(t, err)

	gormDB, err := gorm.Open(mysql.New(mysql.Config{
		Conn: db,
		SkipInitializeWithVersion: true,
	}), &gorm.Config{})
	require.NoError(t, err)

	return gormDB, mock
}

func TestWithTransaction(t *testing.T) {
	db, mock := setupTestDB(t)
	mock.ExpectBegin()
	tx := db.Begin()
	defer func() {
		mock.ExpectRollback()
		tx.Rollback()
	}()

	ctx := context.Background()
	ctxWithTx := WithTransaction(ctx, tx)

	retrievedTx := TransactionFromContext(ctxWithTx)
	assert.NotNil(t, retrievedTx)
	assert.Equal(t, tx, retrievedTx)
}

func TestTransactionFromContext(t *testing.T) {
	t.Run("no transaction in context", func(t *testing.T) {
		ctx := context.Background()
		tx := TransactionFromContext(ctx)
		assert.Nil(t, tx)
	})

	t.Run("transaction exists in context", func(t *testing.T) {
		db, mock := setupTestDB(t)
		mock.ExpectBegin()
		tx := db.Begin()
		defer func() {
			mock.ExpectRollback()
			tx.Rollback()
		}()

		ctx := WithTransaction(context.Background(), tx)
		retrievedTx := TransactionFromContext(ctx)
		assert.NotNil(t, retrievedTx)
		assert.Equal(t, tx, retrievedTx)
	})
}

func TestInTransaction(t *testing.T) {
	t.Run("success execution", func(t *testing.T) {
		db, mock := setupTestDB(t)
		mock.ExpectBegin()
		mock.ExpectCommit()

		ctx := context.Background()
		called := false

		err := InTransaction(ctx, db, func(txCtx context.Context) error {
			called = true
			tx := TransactionFromContext(txCtx)
			assert.NotNil(t, tx)
			return nil
		})

		assert.NoError(t, err)
		assert.True(t, called)
	})

	t.Run("function returns error", func(t *testing.T) {
		db, mock := setupTestDB(t)
		mock.ExpectBegin()
		mock.ExpectRollback()

		ctx := context.Background()
		testErr := errors.New("test error")

		err := InTransaction(ctx, db, func(txCtx context.Context) error {
			return testErr
		})

		assert.Error(t, err)
		assert.Equal(t, testErr, err)
	})
}