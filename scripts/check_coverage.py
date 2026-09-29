#!/usr/bin/env python3
import sys

threshold = float(sys.argv[1]) if len(sys.argv) > 1 else 80.0
total = covered = 0
with open('coverage.out') as stream:
    next(stream)
    for line in stream:
        _, statements, count = line.rsplit(' ', 2)
        statements = int(statements)
        total += statements
        if int(count): covered += statements
if total == 0: raise SystemExit('No statements measured')
percent = covered / total * 100
print(f'Core coverage: {percent:.1f}% (minimum {threshold:.1f}%)')
if percent < threshold: raise SystemExit(1)
