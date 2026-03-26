#!/usr/bin/env node

const { existsSync } = require("node:fs");
const { spawn } = require("node:child_process");
const { resolvePlatformBinaryPath } = require("./platform-package");

let binaryPath = process.env.INTERNATIONALIZER_BINARY_PATH;

if (!binaryPath) {
  try {
    binaryPath = resolvePlatformBinaryPath();
  } catch (error) {
    console.error(error.message);
    console.error(
      "Reinstall the package with npm so the matching optional platform package is installed."
    );
    process.exit(1);
  }
}

if (!existsSync(binaryPath)) {
  console.error(
    "internationalizer binary is missing. Reinstall the package so npm can install the matching platform package."
  );
  process.exit(1);
}

const child = spawn(binaryPath, process.argv.slice(2), {
  stdio: "inherit",
});

child.on("error", (error) => {
  console.error(error.message);
  process.exit(1);
});

child.on("exit", (code, signal) => {
  if (signal) {
    process.kill(process.pid, signal);
    return;
  }
  process.exit(code ?? 1);
});
