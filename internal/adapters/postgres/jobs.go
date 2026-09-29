package postgres

import (
	"context"
	"encoding/json"
	"errors"
	"time"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/raulalmeidatarazona/go-sdd-api-template/internal/application"
	"github.com/raulalmeidatarazona/go-sdd-api-template/internal/domain"
)

type JobStore struct{ pool *pgxpool.Pool }

func NewJobStore(pool *pgxpool.Pool) *JobStore { return &JobStore{pool: pool} }

func (store *JobStore) Create(ctx context.Context, job domain.Job, key, hash, eventID string) (domain.Job, bool, error) {
	tx, err := store.pool.Begin(ctx)
	if err != nil {
		return domain.Job{}, false, err
	}
	defer tx.Rollback(ctx)
	var insertedID string
	err = tx.QueryRow(ctx, `INSERT INTO jobs(id,name,status,idempotency_key,request_hash,created_at)
        VALUES($1,$2,$3,$4,$5,$6) ON CONFLICT (idempotency_key) DO NOTHING RETURNING id`,
		job.ID, job.Name, job.Status, key, hash, job.CreatedAt).Scan(&insertedID)
	if errors.Is(err, pgx.ErrNoRows) {
		var existing domain.Job
		var existingHash string
		err = tx.QueryRow(ctx, `SELECT id,name,status,created_at,request_hash FROM jobs WHERE idempotency_key=$1`, key).
			Scan(&existing.ID, &existing.Name, &existing.Status, &existing.CreatedAt, &existingHash)
		if err != nil {
			return domain.Job{}, false, err
		}
		if existingHash != hash {
			return domain.Job{}, false, application.ErrConflict
		}
		existing.CreatedAt = existing.CreatedAt.UTC()
		if err = tx.Commit(ctx); err != nil {
			return domain.Job{}, false, err
		}
		return existing, false, nil
	}
	if err != nil {
		return domain.Job{}, false, err
	}
	payload, err := json.Marshal(map[string]any{"eventId": eventID, "version": 1, "jobId": job.ID, "name": job.Name, "occurredAt": job.CreatedAt})
	if err != nil {
		return domain.Job{}, false, err
	}
	_, err = tx.Exec(ctx, `INSERT INTO outbox(id,event_type,payload) VALUES($1,$2,$3)`, eventID, "JobCreated.v1", payload)
	if err != nil {
		return domain.Job{}, false, err
	}
	if err = tx.Commit(ctx); err != nil {
		return domain.Job{}, false, err
	}
	return job, true, nil
}

func (store *JobStore) Get(ctx context.Context, id string) (domain.Job, error) {
	var job domain.Job
	err := store.pool.QueryRow(ctx, `SELECT id,name,status,created_at FROM jobs WHERE id=$1`, id).
		Scan(&job.ID, &job.Name, &job.Status, &job.CreatedAt)
	if errors.Is(err, pgx.ErrNoRows) {
		return domain.Job{}, application.ErrNotFound
	}
	job.CreatedAt = job.CreatedAt.UTC()
	return job, err
}

func Ready(ctx context.Context, pool *pgxpool.Pool, maxLag time.Duration) error {
	if err := pool.Ping(ctx); err != nil {
		return err
	}
	var stale bool
	err := pool.QueryRow(ctx, `SELECT EXISTS (SELECT 1 FROM outbox WHERE delivered_at IS NULL AND created_at < now() - ($1::integer * interval '1 second'))`, int(maxLag.Seconds())).Scan(&stale)
	if err != nil {
		return err
	}
	if stale {
		return errors.New("outbox delivery lag exceeded")
	}
	return nil
}
