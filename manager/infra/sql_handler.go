package infra

import (
	"context"
	"database/sql"
	"fmt"
	_ "github.com/go-sql-driver/mysql"
	"github.com/jmoiron/sqlx"
	"touchgift-job-manager/config"
	"touchgift-job-manager/domain/repository"
)

// SQLHandler TODO: repository層にあるべきでは？？
type SQLHandler interface {
	PrepareNamedContext(ctx context.Context, query string) (*sqlx.NamedStmt, error)
	Begin(ctx context.Context) (repository.Transaction, error)
	PrepareContext(ctx context.Context, query string) (*sqlx.Stmt, error)
	Select(ctx context.Context, dest interface{}, query string, args ...interface{}) error
	In(query string, arg interface{}) (*string, []interface{}, error)
	Close()
}

type sqlHandler struct {
	DB     *sqlx.DB
	logger *Logger
}

func NewSQLHandler(logger *Logger) SQLHandler {
	connectionString := fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?charset=utf8mb4&parseTime=True",
		config.Env.Db.User, config.Env.Db.Password, config.Env.Db.Host, config.Env.Db.Port, config.Env.Db.Database)
	db, err := sqlx.Open(config.Env.Db.DriverName, connectionString)
	if err != nil {
		panic(err.Error)
	}
	db.SetMaxOpenConns(config.Env.MaxOpenConns)
	db.SetMaxIdleConns(config.Env.MaxIdleConns)
	db.SetConnMaxLifetime(config.Env.ConnMaxLifetime)
	handler := sqlHandler{
		DB:     db,
		logger: logger,
	}
	err = db.Ping()
	if err != nil {
		panic(err)
	}
	logger.Info().Msg("Connect database")
	return &handler
}

// Close is function
func (s *sqlHandler) Close() {
	s.DB.Close()
}

// PrepareNamedContext is function
func (s *sqlHandler) PrepareNamedContext(ctx context.Context, query string) (*sqlx.NamedStmt, error) {
	return s.DB.PrepareNamedContext(ctx, query)
}

// PrepareContext is function
func (s *sqlHandler) PrepareContext(ctx context.Context, query string) (*sqlx.Stmt, error) {
	return s.DB.PreparexContext(ctx, query)
}

// Select is function
func (s *sqlHandler) Select(ctx context.Context, dest interface{}, query string, args ...interface{}) error {
	err := s.DB.SelectContext(ctx, dest, query, args...)
	if err != nil {
		return err
	}
	return nil
}

// In is function
func (s *sqlHandler) In(query string, arg interface{}) (*string, []interface{}, error) {
	query, args, err := sqlx.Named(query, arg)
	if err != nil {
		return nil, nil, err
	}
	query, args, err = sqlx.In(query, args...)
	if err != nil {
		return nil, nil, err
	}
	query = s.DB.Rebind(query)
	return &query, args, nil
}

// Begin is function
func (s *sqlHandler) Begin(ctx context.Context) (repository.Transaction, error) {
	tx, err := s.DB.BeginTxx(ctx, &sql.TxOptions{})
	if err != nil {
		return nil, err
	}
	return &Transaction{tx}, nil
}

// Transaction is struct
type Transaction struct {
	Tx *sqlx.Tx
}

// Commit is function
func (t *Transaction) Commit() error {
	return t.Tx.Commit()
}

// Rollback is function
func (t *Transaction) Rollback() error {
	return t.Tx.Rollback()
}
