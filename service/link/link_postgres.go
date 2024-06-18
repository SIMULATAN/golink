package link

import (
	"context"
	"database/sql/driver"
	_ "embed"
	"fmt"
	"github.com/XSAM/otelsql"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/stdlib"
	"github.com/jackc/pgx/v5/tracelog"
	"github.com/jmoiron/sqlx"
	"go.opentelemetry.io/otel/attribute"
	semconv "go.opentelemetry.io/otel/semconv/v1.21.0"
	"go.opentelemetry.io/otel/trace"
	"golink/model"
	"golink/service"
	"log/slog"
)

type PostgresLinkService struct {
	conn   *sqlx.DB
	tracer trace.Tracer
}

func NewPostgresService(ctx context.Context, cfg PostgresConfig, schemaSql string) (service.LinkService, error) {
	driverName := "pgx"
	dbSystem := semconv.DBSystemPostgreSQL
	dataSourceConnection, err := cfg.PostgresDataSourceName()
	if err != nil {
		return nil, err
	}
	pgCfg, err := pgx.ParseConfig(dataSourceConnection)
	if err != nil {
		return nil, err
	}

	pgCfg.Tracer = &tracelog.TraceLog{
		Logger: tracelog.LoggerFunc(func(ctx context.Context, level tracelog.LogLevel, msg string, data map[string]any) {
			args := make([]any, 0, len(data))
			for k, v := range data {
				args = append(args, slog.Any(k, v))
			}
			slog.DebugContext(ctx, msg, slog.Group("data", args...))
		}),
		LogLevel: tracelog.LogLevelDebug,
	}
	dataSourceName := stdlib.RegisterConnConfig(pgCfg)

	sqlDB, err := otelsql.Open(driverName, dataSourceName,
		otelsql.WithAttributes(dbSystem),
		otelsql.WithSQLCommenter(true),
		otelsql.WithAttributesGetter(func(ctx context.Context, method otelsql.Method, query string, args []driver.NamedValue) []attribute.KeyValue {
			return []attribute.KeyValue{
				semconv.DBOperationKey.String(string(method)),
				semconv.DBStatementKey.String(query),
			}
		}),
	)
	if err != nil {
		return nil, fmt.Errorf("failed to open database: %w", err)
	}

	if err = otelsql.RegisterDBStatsMetrics(sqlDB, otelsql.WithAttributes(dbSystem)); err != nil {
		return nil, fmt.Errorf("failed to register database stats metrics: %w", err)
	}

	dbx := sqlx.NewDb(sqlDB, driverName)
	if err = dbx.PingContext(ctx); err != nil {
		return nil, fmt.Errorf("failed to ping database: %w", err)
	}
	// execute schema
	if _, err = dbx.ExecContext(ctx, schemaSql); err != nil {
		return nil, fmt.Errorf("failed to execute schema setup: %w", err)
	}

	db := &PostgresLinkService{
		conn: dbx,
	}

	return db, nil
}

func (s *PostgresLinkService) GetLink(code string) (*model.Link, error) {
	var link model.Link
	err := s.conn.Get(&link, "SELECT * FROM link WHERE code = $1", code)
	if err != nil {
		return nil, err
	}
	return &link, nil
}

func (s *PostgresLinkService) CreateLink(target string, short string) (*model.Link, error) {
	link := &model.Link{
		Target: target,
		Code:   short,
	}
	// check if target exists
	if _, err := s.conn.Exec("SELECT 1 FROM link WHERE target = $1", target); err == nil {
		return nil, service.CodeExistsError
	}
	if _, err := s.conn.NamedExec("INSERT INTO link (target, code) VALUES (:target, :code)", link); err != nil {
		return nil, err
	}

	return link, nil
}

func (s *PostgresLinkService) Init() {}

type PostgresConfig struct {
	Host     *string `cfg:"host"`
	Port     *int    `cfg:"port"`
	Username *string `cfg:"username"`
	Password *string `cfg:"password"`
	Database *string `cfg:"database"`
	SSLMode  *string `cfg:"ssl_mode"`

	Uri *string `cfg:"uri"`
}

func (c PostgresConfig) PostgresDataSourceName() (string, error) {
	if c.Uri != nil {
		return *c.Uri, nil
	}

	if c.Host == nil || c.Username == nil || c.Password == nil || c.Database == nil {
		return "", fmt.Errorf("missing required fields")
	}
	if c.Port == nil {
		port := 5432
		c.Port = &port
	}
	if c.SSLMode == nil {
		sslMode := "disable"
		c.SSLMode = &sslMode
	}

	return fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=%s",
		*c.Host,
		*c.Port,
		*c.Username,
		*c.Password,
		*c.Database,
		*c.SSLMode,
	), nil
}
