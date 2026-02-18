package database

import (
	"context"
	"errors"
	"strconv"

	"github.com/Logstay/tech-challange-execution-service/config"
	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
)

const (
	postgres = "postgres"
)

type ConnectionFactory interface {
	NewConnection() UnitOfWork

	NewConnectionWithTransaction() (UnitOfWork, error)

	Delegate() any

	Close() error

	GetDB() *sqlx.DB
}

type connectionFactory struct {
	db *sqlx.DB
}

func NewConnectionFactory(ctx context.Context, delegate *sqlx.DB) (ConnectionFactory, error) {
	if err := delegate.PingContext(ctx); err != nil {
		return nil, err
	}

	return &connectionFactory{db: delegate}, nil
}

func NewPostgresConnectionFactory(ctx context.Context, config config.DatabaseConfig) (ConnectionFactory, error) {
	conn := sqlx.MustOpen(postgres, buildConnString(config))
	if err := conn.PingContext(ctx); err != nil {
		return nil, err
	}

	conn.SetMaxIdleConns(config.PoolMin)
	conn.SetMaxOpenConns(config.PoolMax)
	return &connectionFactory{db: conn}, nil
}

func (f *connectionFactory) NewConnection() UnitOfWork {
	return NewUnitOfWork(f.db, nil)
}

func (f *connectionFactory) NewConnectionWithTransaction() (UnitOfWork, error) {
	tx := f.db.MustBegin()
	if tx == nil {
		return nil, errors.New("nenhuma transação foi iniciada")
	}

	return NewUnitOfWork(nil, tx), nil
}

func (f *connectionFactory) Delegate() any {
	return f.db
}

func (f *connectionFactory) Close() error {
	if f.db != nil {
		return f.db.Close()
	}
	return nil
}

func (f *connectionFactory) GetDB() *sqlx.DB {
	return f.db
}

func buildConnString(config config.DatabaseConfig) string {
	return "host=" + config.Host +
		" port=" + strconv.Itoa(config.Port) +
		" user=" + config.User +
		" password=" + config.Password +
		" dbname=" + config.Database +
		" sslmode=" + config.SSLMode +
		" TimeZone=America/Fortaleza"
}
