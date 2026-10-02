package main

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
	"strings"
)

func runCoverage(args []string) error {
	if len(args) > 1 {
		return fmt.Errorf("usage: repoguard coverage [minimum-percent]")
	}
	minimum := 80.0
	if len(args) == 1 {
		value, err := strconv.ParseFloat(args[0], 64)
		if err != nil || value < 0 || value > 100 {
			return fmt.Errorf("invalid coverage minimum: %q", args[0])
		}
		minimum = value
	}
	file, err := os.Open("coverage.out")
	if err != nil {
		return err
	}
	defer file.Close()
	percent, err := coveragePercent(bufio.NewScanner(file))
	if err != nil {
		return err
	}
	fmt.Printf("Core coverage: %.1f%% (minimum %.1f%%)\n", percent, minimum)
	if percent < minimum {
		return fmt.Errorf("coverage below minimum")
	}
	return nil
}

func coveragePercent(scanner *bufio.Scanner) (float64, error) {
	if !scanner.Scan() || !strings.HasPrefix(scanner.Text(), "mode: ") {
		return 0, fmt.Errorf("invalid Go coverage profile")
	}
	var total, covered int
	for scanner.Scan() {
		fields := strings.Fields(scanner.Text())
		if len(fields) != 3 {
			return 0, fmt.Errorf("invalid Go coverage record")
		}
		statements, statementErr := strconv.Atoi(fields[1])
		count, countErr := strconv.Atoi(fields[2])
		if statementErr != nil || countErr != nil || statements < 0 || count < 0 {
			return 0, fmt.Errorf("invalid Go coverage counts")
		}
		total += statements
		if count > 0 {
			covered += statements
		}
	}
	if err := scanner.Err(); err != nil {
		return 0, err
	}
	if total == 0 {
		return 0, fmt.Errorf("no statements measured")
	}
	return float64(covered) / float64(total) * 100, nil
}
