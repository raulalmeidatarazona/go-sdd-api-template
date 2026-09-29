package application

import (
	"context"
	"errors"
	"testing"

	"github.com/raulalmeidatarazona/go-sdd-api-template/internal/domain"
)

type memoryStore struct {
	jobs   map[string]domain.Job
	hashes map[string]string
	events int
}

func (store *memoryStore) Create(_ context.Context, job domain.Job, key, hash, _ string) (domain.Job, bool, error) {
	if existing, ok := store.jobs[key]; ok {
		if store.hashes[key] != hash {
			return domain.Job{}, false, ErrConflict
		}
		return existing, false, nil
	}
	store.jobs[key], store.hashes[key] = job, hash
	store.events++
	return job, true, nil
}

func (store *memoryStore) Get(_ context.Context, id string) (domain.Job, error) {
	for _, job := range store.jobs {
		if job.ID == id {
			return job, nil
		}
	}
	return domain.Job{}, ErrNotFound
}

func TestCreateIsIdempotentAndKeepsOneEvent(t *testing.T) {
	store := &memoryStore{jobs: map[string]domain.Job{}, hashes: map[string]string{}}
	service := NewJobService(store)
	first, created, err := service.Create(context.Background(), " hello ", "key-1")
	if err != nil || !created || first.Name != "hello" {
		t.Fatalf("first: %#v %v %v", first, created, err)
	}
	second, created, err := service.Create(context.Background(), "hello", "key-1")
	if err != nil || created || second.ID != first.ID || store.events != 1 {
		t.Fatalf("retry: %#v %v %v events=%d", second, created, err, store.events)
	}
	_, _, err = service.Create(context.Background(), "different", "key-1")
	if !errors.Is(err, ErrConflict) {
		t.Fatalf("expected conflict, got %v", err)
	}
	got, err := service.Get(context.Background(), first.ID)
	if err != nil || got.ID != first.ID {
		t.Fatalf("get: %#v %v", got, err)
	}
}

func TestValidation(t *testing.T) {
	service := NewJobService(&memoryStore{jobs: map[string]domain.Job{}, hashes: map[string]string{}})
	if _, _, err := service.Create(context.Background(), "", "key"); !errors.Is(err, domain.ErrInvalidName) {
		t.Fatalf("expected invalid name: %v", err)
	}
	if _, _, err := service.Create(context.Background(), "ok", ""); !errors.Is(err, ErrInvalidIdempotencyKey) {
		t.Fatalf("expected invalid key: %v", err)
	}
	if _, err := service.Get(context.Background(), "bad"); !errors.Is(err, ErrNotFound) {
		t.Fatalf("expected not found: %v", err)
	}
}
