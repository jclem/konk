const pkg = require('./package.json')
const fs = require('node:fs')
const path = require('node:path')
const {execSync} = require('node:child_process')

const archMap = {
  ia32: '386',
  x64: 'amd64',
  arm64: 'arm64'
}

const actions = {
  install,
  uninstall
}

try {
  main()
} catch (err) {
  console.error(err)
  process.exit(1)
}

function install() {
  const {goBinary} = pkg

  const src = path.join(
    '.',
    'dist',
    `konk_${process.platform}_${archMap[process.arch]}`,
    goBinary.name
  )

  const binDir = getInstallationDir()
  const dest = getInstallationPath(binDir)

  fs.mkdirSync(binDir, {recursive: true})
  fs.copyFileSync(src, dest)
}

function getInstallationDir() {
  return execSync('npm bin').toString().trim()
}

function getInstallationPath(binDir) {
  const {goBinary} = pkg
  const dest = path.join(binDir, goBinary.name)
  return dest
}

function uninstall() {
  const binDir = getInstallationDir()
  const dest = getInstallationPath(binDir)

  fs.unlinkSync(dest)
}

function main() {
  const command = process.argv[2]
  const action = actions[command]

  if (!action) {
    console.error('Unknown command:', command)
    console.error('Usage: node tasks.cjs [install|uninstall]')
    process.exit(1)
  }

  action()
}
