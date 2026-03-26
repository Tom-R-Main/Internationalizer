const test = require("node:test");
const assert = require("node:assert/strict");

const {
  getPlatformBinaryName,
  getPlatformPackageName,
  resolvePlatformBinaryPath,
} = require("../bin/platform-package");

test("maps supported platforms to npm package names", () => {
  assert.equal(getPlatformPackageName("darwin", "arm64"), "internationalizer-darwin-arm64");
  assert.equal(getPlatformPackageName("darwin", "x64"), "internationalizer-darwin-x64");
  assert.equal(getPlatformPackageName("linux", "arm64"), "internationalizer-linux-arm64");
  assert.equal(getPlatformPackageName("linux", "x64"), "internationalizer-linux-x64");
  assert.equal(getPlatformPackageName("win32", "x64"), "internationalizer-win32-x64");
});

test("maps platform to binary name", () => {
  assert.equal(getPlatformBinaryName("linux"), "internationalizer");
  assert.equal(getPlatformBinaryName("darwin"), "internationalizer");
  assert.equal(getPlatformBinaryName("win32"), "internationalizer.exe");
});

test("rejects unsupported platforms", () => {
  assert.throws(() => getPlatformPackageName("freebsd", "x64"), /unsupported platform/);
  assert.throws(() => getPlatformPackageName("linux", "ppc64"), /unsupported platform/);
});

test("resolves installed binary path from a platform package", () => {
  const fakeResolve = (request) => {
    if (request === "internationalizer-linux-x64/package.json") {
      return "/tmp/node_modules/internationalizer-linux-x64/package.json";
    }
    throw new Error(`unexpected request: ${request}`);
  };

  assert.equal(
    resolvePlatformBinaryPath("linux", "x64", fakeResolve),
    "/tmp/node_modules/internationalizer-linux-x64/bin/internationalizer"
  );
});
