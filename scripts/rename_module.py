#!/usr/bin/env python3
import re
import subprocess
import sys
from pathlib import Path

if len(sys.argv) != 2 or not re.fullmatch(r"[A-Za-z0-9][A-Za-z0-9._~/-]+/[A-Za-z0-9._~-]+", sys.argv[1]):
    raise SystemExit("usage: python3 scripts/rename_module.py github.com/owner/service")

new_module = sys.argv[1]
module_file = Path("go.mod")
old_module = re.search(r"^module (\S+)$", module_file.read_text(), re.MULTILINE).group(1)
if new_module == old_module:
    raise SystemExit("new module path must differ from template path")

for path in [module_file, *Path("cmd").rglob("*.go"), *Path("internal").rglob("*.go")]:
    content = path.read_text()
    path.write_text(content.replace(old_module, new_module))

subprocess.run(["go", "mod", "tidy"], check=True)
print(f"Renamed Go module from {old_module} to {new_module}")
