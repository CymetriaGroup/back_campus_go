package database

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	"hexagonal-go-backend/internal/ent"

	entsql "entgo.io/ent/dialect/sql"

	_ "github.com/jackc/pgx/v5/stdlib"
)

// OpenPostgres creates an Ent client backed by PostgreSQL. When autoMigrate is
// enabled, it also synchronizes the database from the Ent schemas.
func OpenPostgres(ctx context.Context, databaseURL string, maxOpenConns, maxIdleConns int, connMaxLifetime time.Duration, autoMigrate bool) (*ent.Client, *sql.DB, error) {
	if databaseURL == "" {
		return nil, nil, fmt.Errorf("DATABASE_URL is required")
	}

	db, err := sql.Open("pgx", databaseURL)
	if err != nil {
		return nil, nil, fmt.Errorf("open postgres: %w", err)
	}
	db.SetMaxOpenConns(maxOpenConns)
	db.SetMaxIdleConns(maxIdleConns)
	db.SetConnMaxLifetime(connMaxLifetime)

	if err := db.PingContext(ctx); err != nil {
		_ = db.Close()
		return nil, nil, fmt.Errorf("ping postgres: %w", err)
	}

	driver := entsql.OpenDB("postgres", db)
	client := ent.NewClient(ent.Driver(driver))
	if autoMigrate {
		if err := client.Schema.Create(ctx); err != nil {
			_ = client.Close()
			return nil, nil, fmt.Errorf("create ent schema: %w", err)
		}
	}
	return client, db, nil
}
