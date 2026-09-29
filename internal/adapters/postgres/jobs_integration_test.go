//go:build integration

package postgres

import (
	"context"
	"errors"
	"os"
	"testing"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/raulalmeidatarazona/go-sdd-api-template/internal/application"
)

func TestTransactionalJobAndOutbox(t *testing.T) {
	ctx := context.Background()
	pool, err := pgxpool.New(ctx, os.Getenv("TEST_DATABASE_URL"))
	if err != nil {
		t.Fatal(err)
	}
	defer pool.Close()
	service := application.NewJobService(NewJobStore(pool))
	key := "integration-key-" + time.Now().Format(time.RFC3339Nano)
	job, created, err := service.Create(ctx, "test", key)
	if err != nil || !created {
		t.Fatalf("create: %v", err)
	}
	retry, created, err := service.Create(ctx, "test", key)
	if err != nil || created || retry.ID != job.ID {
		t.Fatalf("retry: %v", err)
	}
	_, _, err = service.Create(ctx, "different", key)
	if !errors.Is(err, application.ErrConflict) {
		t.Fatalf("conflict: %v", err)
	}
	var jobs, events int
	if err := pool.QueryRow(ctx, `SELECT count(*) FROM jobs WHERE id=$1`, job.ID).Scan(&jobs); err != nil {
		t.Fatal(err)
	}
	if err := pool.QueryRow(ctx, `SELECT count(*) FROM outbox WHERE payload->>'jobId'=$1`, job.ID).Scan(&events); err != nil {
		t.Fatal(err)
	}
	if jobs != 1 || events != 1 {
		t.Fatalf("jobs=%d events=%d", jobs, events)
	}
}
