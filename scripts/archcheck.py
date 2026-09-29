#!/usr/bin/env python3
import json
import subprocess
import sys

module = subprocess.check_output(['go','list','-m'], text=True).strip()
output = subprocess.check_output(['go','list','-json','./...'], text=True)
decoder = json.JSONDecoder()
position = 0
violations = []
while position < len(output):
    package, end = decoder.raw_decode(output, position)
    position = end
    while position < len(output) and output[position].isspace(): position += 1
    path = package['ImportPath']
    imports = package.get('Imports', [])
    layer = path.removeprefix(module + '/')
    for target in imports:
        internal = target.removeprefix(module + '/')
        if layer.startswith('internal/domain') and (target.startswith(module + '/') or '.' in target.split('/')[0]):
            violations.append(f'{layer} -> {target}')
        if layer.startswith('internal/application') and ((target.startswith(module + '/') and not internal.startswith('internal/domain')) or ('.' in target.split('/')[0] and not target.startswith(module + '/'))):
            violations.append(f'{layer} -> {target}')
        if layer.startswith('internal/adapters/') and internal.startswith('internal/adapters/') and target.startswith(module + '/') and internal.split('/')[2] != layer.split('/')[2]:
            violations.append(f'{layer} -> {target}')
if violations:
    print('Architecture violations:', *violations, sep='\n', file=sys.stderr)
    raise SystemExit(1)
print('Architecture: OK')
