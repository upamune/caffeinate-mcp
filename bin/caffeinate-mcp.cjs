#!/usr/bin/env node

// This is a CommonJS wrapper for the ESM module
// It helps with npx compatibility issues

async function main() {
  await import('../dist/index.js');
}

main().catch(console.error);