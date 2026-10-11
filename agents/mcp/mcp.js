import { McpServer } from '@modelcontextprotocol/server';
import { serveStdio } from '@modelcontextprotocol/server/stdio';
import * as z from 'zod/v4';
import { lstat, readdir, readFile, realpath } from 'node:fs/promises';
import { fileURLToPath } from 'node:url';
import path from 'node:path';

const repositoryRoot = path.resolve(path.dirname(fileURLToPath(import.meta.url)), '../..');
const featuresPath = path.join(repositoryRoot, 'specs/features');
const validFeatureId = /^[a-z0-9][a-z0-9-]*$/;
const allowedDocuments = new Set(['spec.md', 'plan.md', 'validation.md']);

export async function listFeatures() {
  const entries = await readdir(featuresPath, { withFileTypes: true });
  return entries.filter((entry) => entry.isDirectory() && validFeatureId.test(entry.name))
    .map((entry) => entry.name).sort();
}

export async function readFeatureDocument(featureId, document, directory = featuresPath) {
  if (!validFeatureId.test(featureId) || !allowedDocuments.has(document)) {
    throw new Error('Invalid feature ID or document');
  }
  const root = await realpath(directory);
  const featurePath = path.join(root, featureId);
  const documentPath = path.join(featurePath, document);
  if (!(await lstat(featurePath)).isDirectory() || !(await lstat(documentPath)).isFile()) {
    throw new Error('Feature or document is not a regular directory or file');
  }
  const candidate = await realpath(documentPath);
  if (path.dirname(path.dirname(candidate)) !== root) {
    throw new Error('Document is outside the feature directory');
  }
  return readFile(candidate, 'utf8');
}

export function createServer() {
  const server = new McpServer({ name: 'project-specs', version: '0.1.0' });
  server.registerTool('list_features', {
    description: 'List feature IDs in this repository.',
    inputSchema: z.object({}),
  }, async () => ({ content: [{ type: 'text', text: JSON.stringify(await listFeatures()) }] }));
  server.registerTool('read_feature_document', {
    description: 'Read one spec, plan or validation document for a feature.',
    inputSchema: z.object({
      featureId: z.string(),
      document: z.enum(['spec.md', 'plan.md', 'validation.md']),
    }),
  }, async ({ featureId, document }) => {
    try {
      return { content: [{ type: 'text', text: await readFeatureDocument(featureId, document) }] };
    } catch (error) {
      return { isError: true, content: [{ type: 'text', text: error.message }] };
    }
  });
  return server;
}

if (process.argv[1] && path.resolve(process.argv[1]) === fileURLToPath(import.meta.url)) {
  await serveStdio(() => createServer());
}
