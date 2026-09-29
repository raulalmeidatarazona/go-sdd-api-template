#!/usr/bin/env python3
import re
import subprocess
import sys

PATTERNS = {
    'private key': re.compile(r'-----BEGIN (?:RSA |EC |OPENSSH )?PRIVATE KEY-----'),
    'GitHub token': re.compile(r'gh[pousr]_[A-Za-z0-9_]{30,}'),
    'AWS access key': re.compile(r'AKIA[0-9A-Z]{16}'),
    'credential assignment': re.compile(r"(?i)(?:password|secret|api[_-]?key|token)\s*[:=]\s*['\"](?!test-|example-|dummy-|placeholder-|local-development-)[A-Za-z0-9/+_=.-]{20,}['\"]"),
}

files = subprocess.check_output(['git', 'ls-files', '-z'], text=True).split('\0')
violations = []
for name in files:
    if not name or name.endswith(('.lock', '.sum')) or name == '.env.example':
        continue
    try:
        content = subprocess.check_output(['git', 'show', ':' + name]).decode('utf-8')
    except (UnicodeError, OSError):
        continue
    for label, pattern in PATTERNS.items():
        if pattern.search(content):
            violations.append(f'{name}: {label}')
if violations:
    print('Potential secrets found; inspect these files:', *violations, sep='\n', file=sys.stderr)
    raise SystemExit(1)
print('Secret scan: OK')
