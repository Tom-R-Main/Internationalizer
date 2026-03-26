import { readFile } from "node:fs/promises";

const rootPackage = JSON.parse(await readFile(new URL("../package.json", import.meta.url), "utf8"));
const packagePaths = [
  "../packages/darwin-arm64/package.json",
  "../packages/darwin-x64/package.json",
  "../packages/linux-arm64/package.json",
  "../packages/linux-x64/package.json",
  "../packages/win32-x64/package.json",
];

for (const packagePath of packagePaths) {
  const pkg = JSON.parse(await readFile(new URL(packagePath, import.meta.url), "utf8"));
  if (pkg.version !== rootPackage.version) {
    throw new Error(`${pkg.name} version ${pkg.version} does not match root version ${rootPackage.version}`);
  }
}
