package repository

import (
	"context"
	"errors"
	"fmt"

	"emailaddress.horse/thousand/repository/queries"
	"github.com/jackc/pgx/v4"
	"github.com/jackc/pgx/v4/log/zapadapter"
	"github.com/jackc/pgx/v4/pgxpool"
	"go.uber.org/zap"
)

type Repository struct {
	pool    *pgxpool.Pool
	txConn  txConnector
	queries *queries.Queries
}

type txConnector interface {
	Begin(context.Context) (pgx.Tx, error)
}

type Options struct {
	DatabaseURL string
	Logger      *zap.Logger
}

func New(opts Options) (*Repository, error) {
	if opts.Logger == nil {
		opts.Logger = zap.NewNop()
	}

	config, err := pgxpool.ParseConfig(opts.DatabaseURL)
	if err != nil {
		return nil, fmt.Errorf("error parsing database URL: %w", err)
	}

	config.ConnConfig.Logger = zapadapter.NewLogger(opts.Logger)
	config.LazyConnect = true

	pool, err := pgxpool.ConnectConfig(context.Background(), config)
	if err != nil {
		return nil, fmt.Errorf("error connecting to database: %w", err)
	}

	return &Repository{
		pool:    pool,
		txConn:  pool,
		queries: queries.New(pool),
	}, nil
}

func (r *Repository) Ping(ctx context.Context) error {
	if r.pool == nil {
		return errors.New("cannot ping a transaction")
	}
	return r.pool.Ping(ctx)
}

func (r *Repository) Stat() *pgxpool.Stat {
	return r.pool.Stat()
}

func (r *Repository) WithTx(ctx context.Context) (*Repository, pgx.Tx, error) {
	tx, err := r.txConn.Begin(ctx)
	if err != nil {
		return nil, nil, fmt.Errorf("error beginning transaction: %w", err)
	}

	return &Repository{
		txConn:  tx,
		queries: queries.New(tx),
	}, tx, nil
}

func (r *Repository) WithSavepoint(ctx context.Context) (*Repository, pgx.Tx, error) {
	tx, ok := r.txConn.(pgx.Tx)
	if !ok {
		return nil, nil, errors.New("can only create savepoint from a transaction")
	}

	spTx, err := tx.Begin(ctx)
	if err != nil {
		return nil, nil, fmt.Errorf("error beginning savepoint transaction: %w", err)
	}

	return &Repository{
		txConn:  spTx,
		queries: queries.New(spTx),
	}, spTx, nil
}
