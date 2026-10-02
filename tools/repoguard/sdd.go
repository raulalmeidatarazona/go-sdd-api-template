package main

import (
	"fmt"
	"os/exec"
	"strings"
)

func runSDD(args []string) error {
	diffArgs := []string{"diff", "--name-only", "-z", "--diff-filter=ACMR"}
	switch {
	case len(args) == 0 || (len(args) == 1 && args[0] == "--staged"):
		diffArgs = append(diffArgs, "--cached")
	case len(args) == 2 && args[0] == "--base":
		diffArgs = append(diffArgs, args[1]+"...HEAD")
	default:
		return fmt.Errorf("usage: repoguard sdd [--staged | --base <ref>]")
	}
	output, err := exec.Command("git", diffArgs...).Output()
	if err != nil {
		return fmt.Errorf("read changed files: %w", err)
	}
	if err := validateSDD(strings.Split(strings.TrimSuffix(string(output), "\x00"), "\x00")); err != nil {
		return err
	}
	fmt.Println("SDD gate: OK")
	return nil
}

func validateSDD(changed []string) error {
	code, spec := false, false
	for _, path := range changed {
		if path == "" {
			continue
		}
		if hasPrefix(path, "internal/", "cmd/", "tools/", "migrations/", "agents/mcp/") {
			code = true
		}
		if hasPrefix(path, "specs/features/", "specs/bugs/") && strings.HasSuffix(path, ".md") {
			spec = true
		}
	}
	if code && !spec {
		return fmt.Errorf("SDD gate: code, migration or MCP changes require a feature/bug spec or validation update in the same change")
	}
	return nil
}

func hasPrefix(value string, prefixes ...string) bool {
	for _, prefix := range prefixes {
		if strings.HasPrefix(value, prefix) {
			return true
		}
	}
	return false
}
