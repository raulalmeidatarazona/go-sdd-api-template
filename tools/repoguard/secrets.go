package main

import (
	"bytes"
	"fmt"
	"os/exec"
	"regexp"
	"strings"
)

var secretPatterns = []struct {
	name    string
	pattern *regexp.Regexp
}{
	{"private key", regexp.MustCompile("-----BEGIN (RSA |EC |OPENSSH )?PRIVATE KEY-----")},
	{"GitHub token", regexp.MustCompile("gh[pousr]_[A-Za-z0-9_]{30,}")},
	{"AWS access key", regexp.MustCompile("AKIA[0-9A-Z]{16}")},
	{"credential assignment", regexp.MustCompile("(?i)(password|secret|api[_-]?key|token)[[:space:]]*[:=][[:space:]]*['\"]([A-Za-z0-9/+_=.-]{20,})['\"]")},
}

func runSecrets() error {
	output, err := exec.Command("git", "ls-files", "-z").Output()
	if err != nil {
		return fmt.Errorf("list indexed files: %w", err)
	}
	var violations []string
	for _, name := range bytes.Split(bytes.TrimSuffix(output, []byte{0}), []byte{0}) {
		if len(name) == 0 {
			continue
		}
		content, err := exec.Command("git", "show", ":"+string(name)).Output()
		if err != nil {
			return fmt.Errorf("read staged %s: %w", name, err)
		}
		if bytes.IndexByte(content, 0) >= 0 {
			continue
		}
		for _, label := range detectSecrets(content) {
			violations = append(violations, string(name)+": "+label)
		}
	}
	if len(violations) > 0 {
		return fmt.Errorf("potential secrets found; inspect these files:\n%s", strings.Join(violations, "\n"))
	}
	fmt.Println("Secret scan: OK")
	return nil
}

func detectSecrets(content []byte) []string {
	var found []string
	for _, candidate := range secretPatterns {
		if candidate.pattern.Match(content) {
			found = append(found, candidate.name)
		}
	}
	return found
}
