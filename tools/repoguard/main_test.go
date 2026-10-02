package main

import (
	"bufio"
	"strings"
	"testing"
)

func TestArchitectureViolations(t *testing.T) {
	module := "example.com/service"
	cases := []struct {
		name    string
		current packageInfo
		want    bool
	}{
		{"domain cannot import adapter", packageInfo{ImportPath: module + "/internal/domain", Imports: []string{module + "/internal/adapters/postgres"}}, true},
		{"domain cannot import third party", packageInfo{ImportPath: module + "/internal/domain", Imports: []string{"github.com/vendor/lib"}}, true},
		{"application may import domain", packageInfo{ImportPath: module + "/internal/application", Imports: []string{module + "/internal/domain"}}, false},
		{"application cannot import adapter", packageInfo{ImportPath: module + "/internal/application", Imports: []string{module + "/internal/adapters/postgres"}}, true},
		{"adapters cannot cross-import", packageInfo{ImportPath: module + "/internal/adapters/httpapi", Imports: []string{module + "/internal/adapters/postgres"}}, true},
	}
	for _, test := range cases {
		t.Run(test.name, func(t *testing.T) {
			got := len(architectureViolations(module, []packageInfo{test.current})) > 0
			if got != test.want {
				t.Fatalf("violation = %t, want %t", got, test.want)
			}
		})
	}
}

func TestSDDGateRejectsCodeWithoutSpec(t *testing.T) {
	if err := validateSDD([]string{"tools/repoguard/main.go"}); err == nil {
		t.Fatal("tooling change without a spec passed")
	}
	if err := validateSDD([]string{"internal/application/handler.go", "specs/features/0002/spec.md"}); err != nil {
		t.Fatal(err)
	}
}

func TestSecretDetection(t *testing.T) {
	token := "gh" + "p_" + strings.Repeat("A", 32)
	if len(detectSecrets([]byte(token))) == 0 {
		t.Fatal("token was not detected")
	}
	if len(detectSecrets([]byte("API_KEY=local-development-example"))) != 0 {
		t.Fatal("example value was rejected")
	}
}

func TestCoveragePercent(t *testing.T) {
	scanner := bufio.NewScanner(strings.NewReader("mode: set\nfile.go:1.1,1.2 3 1\nfile.go:2.1,2.2 1 0\n"))
	percent, err := coveragePercent(scanner)
	if err != nil || percent != 75 {
		t.Fatalf("coverage = %f, err = %v", percent, err)
	}
}
