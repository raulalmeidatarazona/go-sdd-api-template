package main

import (
	"fmt"
	"os"
)

func main() {
	if len(os.Args) < 2 {
		fail(fmt.Errorf("usage: go run ./tools/repoguard <arch|sdd|secrets|coverage|rename-module> [args]"))
	}

	var err error
	switch os.Args[1] {
	case "arch":
		err = runArchitecture()
	case "sdd":
		err = runSDD(os.Args[2:])
	case "secrets":
		err = runSecrets()
	case "coverage":
		err = runCoverage(os.Args[2:])
	case "rename-module":
		err = renameModule(os.Args[2:])
	default:
		err = fmt.Errorf("unknown guard: %s", os.Args[1])
	}
	if err != nil {
		fail(err)
	}
}

func fail(err error) {
	fmt.Fprintln(os.Stderr, err)
	os.Exit(1)
}
