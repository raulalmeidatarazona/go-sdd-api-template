import assert from 'node:assert/strict';
import { test } from 'node:test';
import { Client } from '@modelcontextprotocol/client';
import { StdioClientTransport } from '@modelcontextprotocol/client/stdio';
import { fileURLToPath } from 'node:url';
import { mkdtemp, mkdir, symlink, writeFile, rm } from 'node:fs/promises';
import { tmpdir } from 'node:os';
import path from 'node:path';
import { listFeatures, readFeatureDocument } from './mcp.js';

test('feature reader limits paths and reads a real spec', async () => {
  const features = await listFeatures();
  assert.ok(features.includes('0001-job-example'));
  assert.match(await readFeatureDocument('0001-job-example', 'spec.md'), /Job/i);
  await assert.rejects(readFeatureDocument('../..', 'spec.md'), /Invalid/);
  await assert.rejects(readFeatureDocument('0001-job-example', '../../AGENTS.md'), /Invalid/);
});

test('feature reader refuses symlink escapes', async () => {
  const root = await mkdtemp(path.join(tmpdir(), 'project-specs-'));
  try {
    await mkdir(path.join(root, 'safe'));
    await writeFile(path.join(root, 'secret.md'), 'private');
    await symlink(path.join(root, 'secret.md'), path.join(root, 'safe', 'spec.md'));
    await assert.rejects(readFeatureDocument('safe', 'spec.md', root), /regular/);
  } finally {
    await rm(root, { recursive: true, force: true });
  }
});

test('MCP server connects over stdio and exposes read-only tools', async () => {
  const client = new Client({ name: 'project-specs-test', version: '0.1.0' });
  const transport = new StdioClientTransport({ command: process.execPath, args: [fileURLToPath(new URL('./mcp.js', import.meta.url))] });
  await client.connect(transport);
  try {
    const { tools } = await client.listTools();
    assert.deepEqual(tools.map((tool) => tool.name).sort(), ['list_features', 'read_feature_document']);
    const result = await client.callTool({ name: 'list_features', arguments: {} });
    assert.match(result.content[0].text, /0001-job-example/);
    const refused = await client.callTool({ name: 'read_feature_document', arguments: { featureId: '../..', document: 'spec.md' } });
    assert.equal(refused.isError, true);
  } finally {
    await client.close();
  }
});
