package store

import (
	"context"
	"database/sql"
	"fmt"
	"os"
	"path/filepath"
	"runtime"
	"strings"
	"testing"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/pressly/goose/v3"
	_ "github.com/jackc/pgx/v5/stdlib" // goose needs database/sql driver

	"github.com/heniums/habitx-server/db"
)

// TestDB creates a fresh, uniquely-named database on the dev Postgres
// instance, applies all migrations, and returns queries plus a cleanup
// function. Each test run is fully isolated: own database, no shared
// state with dev data or other test runs.
func TestDB(t *testing.T) (*db.Queries, func()) {
	t.Helper()

	adminDSN := os.Getenv("HABITX_TEST_DB_DSN")
	if adminDSN == "" {
		adminDSN = "postgres://habitx:habitx@localhost:5433/postgres?sslmode=disable"
	}

	dbName := fmt.Sprintf("habitx_test_%d_%d", os.Getpid(), time.Now().UnixNano())
	ctx := context.Background()

	// 1. connect to the maintenance database ("postgres") to CREATE DATABASE
	admin, err := pgxpool.New(ctx, adminDSN)
	if err != nil {
		t.Fatalf("connect admin: %v", err)
	}
	defer admin.Close()

	if _, err := admin.Exec(ctx, fmt.Sprintf(`CREATE DATABASE %s`, dbName)); err != nil {
		t.Fatalf("create test db: %v", err)
	}

	// 2. DSN pointing at the fresh database
	testDSN := strings.Replace(adminDSN, "/postgres?", "/"+dbName+"?", 1)

	// 3. migrate it with goose (library mode)
	if err := migrateUp(t, testDSN); err != nil {
		t.Fatalf("migrate: %v", err)
	}

	// 4. open the pool tests will use
	pool, err := pgxpool.New(ctx, testDSN)
	if err != nil {
		t.Fatalf("connect test db: %v", err)
	}

	cleanup := func() {
		pool.Close()
		drop, err := pgxpool.New(ctx, adminDSN)
		if err != nil {
			return
		}
		defer drop.Close()
		_, _ = drop.Exec(ctx, fmt.Sprintf(`DROP DATABASE %s WITH (FORCE)`, dbName))
	}

	return db.New(pool), cleanup
}

func migrateUp(t *testing.T, dsn string) error {
	t.Helper()

	// absolute path to migrations dir, regardless of test's working directory
	_, thisFile, _, _ := runtime.Caller(0)
	migrationsDir := filepath.Join(filepath.Dir(thisFile), "..", "..", "db", "migrations")

	// goose needs a *sql.DB (database/sql), not a DSN string — open one via pgx's stdlib driver
	sqldb, err := sql.Open("pgx", dsn)
	if err != nil {
		return err
	}
	defer sqldb.Close()

	goose.SetBaseFS(nil)
	return goose.UpContext(context.Background(), sqldb, migrationsDir)
}
