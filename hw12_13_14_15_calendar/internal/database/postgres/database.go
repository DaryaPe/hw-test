package postgres

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/DaryaPe/hw-test/hw12_13_14_15_calendar/common"
	"github.com/DaryaPe/hw-test/hw12_13_14_15_calendar/internal/config"
	"github.com/DaryaPe/hw-test/hw12_13_14_15_calendar/pkg/errors"
	"github.com/jmoiron/sqlx"
)

const (
	ctxSQL  = "sql"
	ctxArgs = "args"

	ctxQuery    = "query"
	ctxQueryRow = "query row"
	ctxExec     = "exec"
)

type DB struct {
	config *config.DatabaseConfig
	conn   *sqlx.DB
	logger common.Logger
}

func New(cfg *config.DatabaseConfig, log common.Logger) *DB {
	return &DB{
		config: cfg,
		logger: log,
	}
}

func (p *DB) Open(ctx context.Context) error {
	source := fmt.Sprintf(
		"dbname=%v user=%v password=%v host=%v port=%v sslmode=disable",
		p.config.Database,
		p.config.Username,
		p.config.Password,
		p.config.Host,
		p.config.Port,
	)
	db, err := sqlx.ConnectContext(ctx, p.config.Source, source)
	if err != nil {
		return errors.Wrap(err, "sqlx.ConnectContext")
	}
	p.conn = db
	p.logger.Debugf("connect pg db to url='%v:%v'", p.config.Host, p.config.Port)
	return nil
}

func (p *DB) Close() error {
	p.logger.Debugf("close connection pg db to url='%v:%v'", p.config.Host, p.config.Port)
	return p.conn.Close()
}

func (p *DB) Query(ctx context.Context, sql string, args ...interface{}) (context.CancelFunc, *sql.Rows, error) {
	ctx, cancel := context.WithTimeout(ctx, p.config.Timeout)
	p.log(ctx, ctxQuery, sql, args...)
	rows, err := p.conn.QueryContext(ctx, sql, args...)
	if err != nil {
		return cancel, nil, errors.Wrap(err, "QueryContext")
	}
	return cancel, rows, nil
}

func (p *DB) QueryRow(ctx context.Context, sql string, args ...interface{}) (context.CancelFunc, *sql.Row) {
	ctx, cancel := context.WithTimeout(ctx, p.config.Timeout)
	p.log(ctx, ctxQueryRow, sql, args...)
	return cancel, p.conn.QueryRowContext(ctx, sql, args...)
}

func (p *DB) Exec(ctx context.Context, sql string, args ...interface{}) error {
	ctx, cancel := context.WithTimeout(ctx, p.config.Timeout)
	defer cancel()
	p.log(ctx, ctxExec, sql, args...)
	_, err := p.conn.ExecContext(ctx, sql, args...)
	if err != nil {
		return errors.Wrap(err, "ExecContext")
	}
	return nil
}

func (p *DB) Rebind(query string) string {
	return p.conn.Rebind(query)
}

func (p *DB) log(ctx context.Context, action, sql string, args ...interface{}) {
	p.logger.Debugw(action, common.CtxActionID, ctx.Value(common.CtxActionID), ctxSQL, sql, ctxArgs, args)
}
