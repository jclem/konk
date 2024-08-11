const pkg = require("./package.json");
const fs = require("node:fs");
const path = require("node:path");

const archMap = {
  ia32: "386",
  x64: "amd64_v1",
  arm64: "arm64",
};

const actions = {
  install,
  uninstall,
};

try {
  main();
} catch (err) {
  console.error(err);
  process.exit(1);
}

function install() {
  const { goBinary } = pkg;

  const src = path.join(
    ".",
    "dist",
    `konk_${process.platform}_${archMap[process.arch]}`,
    goBinary.name,
  );

  const dest = path.join(goBinary.path, goBinary.name);

  fs.mkdirSync(goBinary.path, { recursive: true });
  fs.copyFileSync(src, dest);
}

function uninstall() {
  const dest = path.join(goBinary.path, goBinary.name);
  console.log("Removing konk binaries from", dest);

  fs.unlinkSync(dest);
}

function main() {
  const command = process.argv[2];
  const action = actions[command];

  if (!action) {
    console.error("Unknown command:", command);
    console.error("Usage: node tasks.cjs [install|uninstall]");
    process.exit(1);
  }

  action();
}
