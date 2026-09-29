#!/usr/bin/env python3
import subprocess
import sys

def changed_files():
    if len(sys.argv) == 1 or sys.argv[1] == '--staged':
        command = ['git', 'diff', '--cached', '--name-only', '--diff-filter=ACMR']
    elif len(sys.argv) == 3 and sys.argv[1] == '--base':
        command = ['git', 'diff', '--name-only', '--diff-filter=ACMR', f'{sys.argv[2]}...HEAD']
    else:
        raise SystemExit('usage: check_sdd.py [--staged | --base <ref>]')
    return set(subprocess.check_output(command, text=True).splitlines())

changed = changed_files()
code = [name for name in changed if name.startswith(('internal/', 'cmd/', 'src/', 'migrations/'))]
docs = [name for name in changed if name.startswith(('specs/features/', 'specs/bugs/')) and name.endswith('.md')]
if code and not docs:
    raise SystemExit('SDD gate: code or migration changes require a staged feature/bug spec or validation update in the same change')
print('SDD gate: OK')
