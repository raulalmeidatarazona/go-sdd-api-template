package httpapi

import (
	"context"
	"crypto/sha256"
	"crypto/subtle"
	"encoding/json"
	"errors"
	"io"
	"net/http"
	"strings"
	"time"

	"github.com/raulalmeidatarazona/go-sdd-api-template/internal/application"
	"github.com/raulalmeidatarazona/go-sdd-api-template/internal/domain"
)

type Jobs interface {
	Create(context.Context, string, string) (domain.Job, bool, error)
	Get(context.Context, string) (domain.Job, error)
}

type ReadyCheck func(context.Context) error

func NewHandler(jobs Jobs, apiKey string, ready ReadyCheck) http.Handler {
	mux := http.NewServeMux()
	mux.HandleFunc("GET /live", func(w http.ResponseWriter, _ *http.Request) {
		respond(w, http.StatusOK, map[string]string{"status": "live"})
	})
	mux.HandleFunc("GET /ready", func(w http.ResponseWriter, r *http.Request) {
		ctx, cancel := context.WithTimeout(r.Context(), 2*time.Second)
		defer cancel()
		if err := ready(ctx); err != nil {
			respond(w, http.StatusServiceUnavailable, map[string]string{"status": "unavailable"})
			return
		}
		respond(w, http.StatusOK, map[string]string{"status": "ready"})
	})
	mux.Handle("POST /v1/jobs", authorize(apiKey, http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if !strings.HasPrefix(r.Header.Get("Content-Type"), "application/json") {
			respondError(w, http.StatusUnsupportedMediaType, "application/json required")
			return
		}
		body := http.MaxBytesReader(w, r.Body, 4096)
		decoder := json.NewDecoder(body)
		decoder.DisallowUnknownFields()
		var input struct {
			Name string `json:"name"`
		}
		if err := decoder.Decode(&input); err != nil {
			respondError(w, http.StatusBadRequest, "invalid JSON body")
			return
		}
		if err := decoder.Decode(new(any)); !errors.Is(err, io.EOF) {
			respondError(w, http.StatusBadRequest, "single JSON object required")
			return
		}
		job, created, err := jobs.Create(r.Context(), input.Name, r.Header.Get("Idempotency-Key"))
		if err != nil {
			handleError(w, err)
			return
		}
		status := http.StatusOK
		if created {
			status = http.StatusCreated
		}
		respond(w, status, job)
	})))
	mux.Handle("GET /v1/jobs/{id}", authorize(apiKey, http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		job, err := jobs.Get(r.Context(), r.PathValue("id"))
		if err != nil {
			handleError(w, err)
			return
		}
		respond(w, http.StatusOK, job)
	})))
	return mux
}

func authorize(apiKey string, next http.Handler) http.Handler {
	expected := sha256.Sum256([]byte(apiKey))
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		candidate := strings.TrimPrefix(r.Header.Get("Authorization"), "Bearer ")
		actual := sha256.Sum256([]byte(candidate))
		if candidate == "" || subtle.ConstantTimeCompare(actual[:], expected[:]) != 1 {
			respondError(w, http.StatusUnauthorized, "unauthorized")
			return
		}
		next.ServeHTTP(w, r)
	})
}

func handleError(w http.ResponseWriter, err error) {
	switch {
	case errors.Is(err, domain.ErrInvalidName), errors.Is(err, application.ErrInvalidIdempotencyKey):
		respondError(w, http.StatusBadRequest, err.Error())
	case errors.Is(err, application.ErrConflict):
		respondError(w, http.StatusConflict, err.Error())
	case errors.Is(err, application.ErrNotFound):
		respondError(w, http.StatusNotFound, err.Error())
	default:
		respondError(w, http.StatusInternalServerError, "internal error")
	}
}

func respondError(w http.ResponseWriter, status int, message string) {
	respond(w, status, map[string]string{"error": message})
}
func respond(w http.ResponseWriter, status int, value any) {
	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("X-Content-Type-Options", "nosniff")
	w.WriteHeader(status)
	_ = json.NewEncoder(w).Encode(value)
}
