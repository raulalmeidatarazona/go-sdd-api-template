package main

import (
	"bytes"
	"fmt"
	"io/fs"
	"os"
	"os/exec"
	"path/filepath"
	"regexp"
	"strings"
)

var modulePathPattern = regexp.MustCompile("^[A-Za-z0-9][A-Za-z0-9._~/-]+/[A-Za-z0-9._~-]+$")

func renameModule(args []string) error {
	if len(args) != 1 || !modulePathPattern.MatchString(args[0]) {
		return fmt.Errorf("usage: repoguard rename-module github.com/owner/service")
	}
	newModule := args[0]
	goMod, err := os.ReadFile("go.mod")
	if err != nil {
		return err
	}
	oldModule := ""
	for _, line := range strings.Split(string(goMod), "\n") {
		if strings.HasPrefix(line, "module ") {
			oldModule = strings.TrimSpace(strings.TrimPrefix(line, "module "))
			break
		}
	}
	if oldModule == "" || oldModule == newModule {
		return fmt.Errorf("new module path must differ from the current module")
	}
	paths := []string{"go.mod"}
	for _, root := range []string{"cmd", "internal", "tools"} {
		err := filepath.WalkDir(root, func(path string, entry fs.DirEntry, walkErr error) error {
			if walkErr != nil {
				return walkErr
			}
			if !entry.IsDir() && strings.HasSuffix(path, ".go") {
				paths = append(paths, path)
			}
			return nil
		})
		if err != nil {
			return err
		}
	}
	for _, path := range paths {
		content, err := os.ReadFile(path)
		if err != nil {
			return err
		}
		updated := bytes.ReplaceAll(content, []byte(oldModule), []byte(newModule))
		if bytes.Equal(content, updated) {
			continue
		}
		info, err := os.Stat(path)
		if err != nil {
			return err
		}
		if err := os.WriteFile(path, updated, info.Mode().Perm()); err != nil {
			return err
		}
	}
	if output, err := exec.Command("go", "mod", "tidy").CombinedOutput(); err != nil {
		return fmt.Errorf("go mod tidy: %w\n%s", err, output)
	}
	fmt.Printf("Renamed Go module from %s to %s\n", oldModule, newModule)
	return nil
}
