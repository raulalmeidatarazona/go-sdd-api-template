package application

import (
	"context"
	"crypto/rand"
	"crypto/sha256"
	"encoding/hex"
	"errors"
	"regexp"
	"strings"
	"time"

	"github.com/raulalmeidatarazona/go-sdd-api-template/internal/domain"
)

var (
	ErrInvalidIdempotencyKey = errors.New("idempotency key must contain 1 to 128 characters")
	ErrConflict              = errors.New("idempotency key already used for a different request")
	ErrNotFound              = errors.New("job not found")
	uuidPattern              = regexp.MustCompile(`(?i)^[0-9a-f]{8}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{12}$`)
)

type JobStore interface {
	Create(ctx context.Context, job domain.Job, idempotencyKey, requestHash, eventID string) (domain.Job, bool, error)
	Get(ctx context.Context, id string) (domain.Job, error)
}

type JobService struct{ store JobStore }

func NewJobService(store JobStore) *JobService { return &JobService{store: store} }

func (service *JobService) Create(ctx context.Context, nameInput, keyInput string) (domain.Job, bool, error) {
	name, err := domain.NormalizeName(nameInput)
	if err != nil {
		return domain.Job{}, false, err
	}
	key := strings.TrimSpace(keyInput)
	if len(key) < 1 || len(key) > 128 {
		return domain.Job{}, false, ErrInvalidIdempotencyKey
	}
	jobID, err := newID()
	if err != nil {
		return domain.Job{}, false, err
	}
	eventID, err := newID()
	if err != nil {
		return domain.Job{}, false, err
	}
	hash := sha256.Sum256([]byte(name))
	job := domain.Job{ID: jobID, Name: name, Status: "pending", CreatedAt: time.Now().UTC()}
	return service.store.Create(ctx, job, key, hex.EncodeToString(hash[:]), eventID)
}

func (service *JobService) Get(ctx context.Context, id string) (domain.Job, error) {
	if !uuidPattern.MatchString(id) {
		return domain.Job{}, ErrNotFound
	}
	return service.store.Get(ctx, id)
}

func newID() (string, error) {
	var bytes [16]byte
	if _, err := rand.Read(bytes[:]); err != nil {
		return "", err
	}
	bytes[6] = (bytes[6] & 0x0f) | 0x40
	bytes[8] = (bytes[8] & 0x3f) | 0x80
	return hex.EncodeToString(bytes[0:4]) + "-" + hex.EncodeToString(bytes[4:6]) + "-" + hex.EncodeToString(bytes[6:8]) + "-" + hex.EncodeToString(bytes[8:10]) + "-" + hex.EncodeToString(bytes[10:16]), nil
}
