package httpapi

import (
	"context"
	"errors"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/raulalmeidatarazona/go-sdd-api-template/internal/application"
	"github.com/raulalmeidatarazona/go-sdd-api-template/internal/domain"
)

type stubJobs struct {
	err     error
	created bool
}

func (stub stubJobs) Create(_ context.Context, name, key string) (domain.Job, bool, error) {
	if stub.err != nil {
		return domain.Job{}, false, stub.err
	}
	return domain.Job{ID: "123e4567-e89b-12d3-a456-426614174000", Name: name, Status: "pending"}, stub.created, nil
}
func (stub stubJobs) Get(_ context.Context, id string) (domain.Job, error) {
	if stub.err != nil {
		return domain.Job{}, stub.err
	}
	return domain.Job{ID: id, Name: "test", Status: "pending"}, nil
}

func TestHandlerContracts(t *testing.T) {
	key := "a-valid-local-test-key"
	cases := []struct {
		name, method, path, body, auth, content, idem string
		jobs                                          stubJobs
		ready                                         ReadyCheck
		expected                                      int
	}{
		{name: "live", method: "GET", path: "/live", expected: 200},
		{name: "ready", method: "GET", path: "/ready", ready: func(context.Context) error { return errors.New("db down") }, expected: 503},
		{name: "unauthorized", method: "GET", path: "/v1/jobs/1", expected: 401},
		{name: "created", method: "POST", path: "/v1/jobs", body: `{"name":"ok"}`, auth: key, content: "application/json", idem: "one", jobs: stubJobs{created: true}, expected: 201},
		{name: "retry", method: "POST", path: "/v1/jobs", body: `{"name":"ok"}`, auth: key, content: "application/json", idem: "one", expected: 200},
		{name: "content-type", method: "POST", path: "/v1/jobs", body: `{}`, auth: key, expected: 415},
		{name: "unknown field", method: "POST", path: "/v1/jobs", body: `{"oops":1}`, auth: key, content: "application/json", expected: 400},
		{name: "trailing json", method: "POST", path: "/v1/jobs", body: `{} {}`, auth: key, content: "application/json", expected: 400},
		{name: "conflict", method: "POST", path: "/v1/jobs", body: `{}`, auth: key, content: "application/json", jobs: stubJobs{err: application.ErrConflict}, expected: 409},
		{name: "invalid", method: "POST", path: "/v1/jobs", body: `{}`, auth: key, content: "application/json", jobs: stubJobs{err: domain.ErrInvalidName}, expected: 400},
		{name: "failure", method: "POST", path: "/v1/jobs", body: `{}`, auth: key, content: "application/json", jobs: stubJobs{err: errors.New("secret db error")}, expected: 500},
		{name: "not found", method: "GET", path: "/v1/jobs/1", auth: key, jobs: stubJobs{err: application.ErrNotFound}, expected: 404},
		{name: "get", method: "GET", path: "/v1/jobs/123e4567-e89b-12d3-a456-426614174000", auth: key, expected: 200},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			ready := tc.ready
			if ready == nil {
				ready = func(context.Context) error { return nil }
			}
			request := httptest.NewRequest(tc.method, tc.path, strings.NewReader(tc.body))
			if tc.auth != "" {
				request.Header.Set("Authorization", "Bearer "+tc.auth)
			}
			if tc.content != "" {
				request.Header.Set("Content-Type", tc.content)
			}
			if tc.idem != "" {
				request.Header.Set("Idempotency-Key", tc.idem)
			}
			recorder := httptest.NewRecorder()
			NewHandler(tc.jobs, key, ready).ServeHTTP(recorder, request)
			if recorder.Code != tc.expected {
				t.Fatalf("status %d expected %d; body %s", recorder.Code, tc.expected, recorder.Body.String())
			}
			if strings.Contains(recorder.Body.String(), "secret db error") {
				t.Fatal("internal error leaked")
			}
			if recorder.Code == http.StatusOK && tc.name == "get" && !strings.Contains(recorder.Body.String(), "123e4567") {
				t.Fatal("missing job")
			}
		})
	}
}
