const { dirname, join } = require("node:path");

function getPlatformBinaryName(platform = process.platform) {
  return platform === "win32" ? "internationalizer.exe" : "internationalizer";
}

function getPlatformPackageName(platform = process.platform, arch = process.arch) {
  if (platform === "darwin" && arch === "arm64") {
    return "internationalizer-darwin-arm64";
  }
  if (platform === "darwin" && arch === "x64") {
    return "internationalizer-darwin-x64";
  }
  if (platform === "linux" && arch === "arm64") {
    return "internationalizer-linux-arm64";
  }
  if (platform === "linux" && arch === "x64") {
    return "internationalizer-linux-x64";
  }
  if (platform === "win32" && arch === "x64") {
    return "internationalizer-win32-x64";
  }

  throw new Error(`unsupported platform: ${platform}/${arch}`);
}

function resolvePlatformBinaryPath(platform = process.platform, arch = process.arch, resolvePackage = require.resolve) {
  const packageName = getPlatformPackageName(platform, arch);
  const packageJsonPath = resolvePackage(`${packageName}/package.json`);
  return join(dirname(packageJsonPath), "bin", getPlatformBinaryName(platform));
}

module.exports = {
  getPlatformBinaryName,
  getPlatformPackageName,
  resolvePlatformBinaryPath,
};
