package db

import (
	"context"

	"github.com/jackc/pgx/v5/pgxpool"
)

// Provides all functions to execute database queries and Transactions
type Store interface {
	Querier
	TransferTx(ctx context.Context, arg TransferTxParams) (TransferTxResult, error)
	CreateUserTx(ctx context.Context, arg CreateUserTxParams) (CreateUserTxResult, error)
	VerifyEmailTx(ctx context.Context, arg VerifyEmailTxParams) (VerifyEmailTxResult, error)
}

// Provides all functions to execute SQL queries and Transactions
type SQLStore struct {
	*Queries
	connPool *pgxpool.Pool
}

func NewStore(connPool *pgxpool.Pool) Store {

	return &SQLStore{connPool: connPool, Queries: New(connPool)}
}
