package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"os/exec"
	"strings"
)

type packageInfo struct {
	ImportPath string
	Imports    []string
}

func runArchitecture() error {
	moduleOutput, err := exec.Command("go", "list", "-m").Output()
	if err != nil {
		return fmt.Errorf("read Go module: %w", err)
	}
	module := strings.TrimSpace(string(moduleOutput))
	output, err := exec.Command("go", "list", "-json", "./...").Output()
	if err != nil {
		return fmt.Errorf("list Go packages: %w", err)
	}
	decoder := json.NewDecoder(bytes.NewReader(output))
	var packages []packageInfo
	for {
		var current packageInfo
		if err := decoder.Decode(&current); err == io.EOF {
			break
		} else if err != nil {
			return fmt.Errorf("decode Go packages: %w", err)
		}
		packages = append(packages, current)
	}
	violations := architectureViolations(module, packages)
	if len(violations) > 0 {
		return fmt.Errorf("architecture violations:\n%s", strings.Join(violations, "\n"))
	}
	fmt.Println("Architecture: OK")
	return nil
}

func architectureViolations(module string, packages []packageInfo) []string {
	var violations []string
	for _, current := range packages {
		layer, local := strings.CutPrefix(current.ImportPath, module+"/")
		if !local {
			continue
		}
		for _, target := range current.Imports {
			imported, internal := strings.CutPrefix(target, module+"/")
			external := strings.Contains(strings.Split(target, "/")[0], ".")
			invalid := false
			switch {
			case inLayer(layer, "internal/domain"):
				invalid = internal || external
			case inLayer(layer, "internal/application"):
				invalid = (internal && !inLayer(imported, "internal/domain")) || (!internal && external)
			case strings.HasPrefix(layer, "internal/adapters/") && internal && strings.HasPrefix(imported, "internal/adapters/"):
				invalid = strings.Split(layer, "/")[2] != strings.Split(imported, "/")[2]
			}
			if invalid {
				violations = append(violations, layer+" -> "+target)
			}
		}
	}
	return violations
}

func inLayer(path, layer string) bool {
	return path == layer || strings.HasPrefix(path, layer+"/")
}
