---
schema_version: 1
id: npm
title: npm
summary: Install JavaScript packages and run a project's scripts.
aliases: []
tags: [developer, javascript, packages]
implementations: [upstream]
updated: 2026-11-10
references:
  - https://docs.npmjs.com/cli/commands
  - https://docs.npmjs.com/
---
## Common
### Overview
`npm` is the package manager that comes with Node.js. It downloads the libraries that a project depends on into a `node_modules` folder, records them in `package.json`, and runs the scripts that the project defines, such as `npm test`. The examples here use only your own folder, and the commands that download packages are reviewed only.

### Background
`npm` is an open-source command with documentation at docs.npmjs.com. It talks to a public registry for downloads. The tool tested here is the npm that comes with the long-term-support Node.js on the host.

### Examples
Print the version; it starts with a digit:

```sh
npm --version | grep -c '^[0-9]'
```

Start a project: `npm init -y` writes a `package.json` with default answers, and `node` reads its version:

```sh
npm init -y > /dev/null 2>&1 && node -p "require('./package.json').version"
```

Add a script to the project and run it; `--silent` hides the lines npm adds:

```sh
npm init -y > /dev/null 2>&1 && npm pkg set scripts.hi="echo hi" > /dev/null 2>&1 && npm run hi --silent
```

Download a package into this project; this needs the network and is reviewed only:

```sh
npm install left-pad
```

List what is installed; reviewed only because it needs an install first:

```sh
npm ls --depth=0
```

### Subcommands
`npm install` downloads the packages that the project needs. `npm run NAME` runs a script from `package.json`. `npm init` creates a `package.json`.

### Common mistakes
Run `npm` in the folder that holds `package.json`, or it uses a different project or none. `npm install` writes to `node_modules` and `package-lock.json`; commit the lock file and do not commit `node_modules`. A package can run code when it is installed, so install only what you trust. A global install (`-g`) changes your whole system; prefer `npx` or a project install. Do not use `sudo` with `npm`: it leaves files that you cannot change.

### Portability
The commands are the same on every system, and details change between npm versions. The examples that use the network are reviewed only: they were not run by the tests, and the ones that call a public service were checked by hand on 2026-11-10. The first three examples were executed on the npm of the host. The Common examples were executed on npm (the version that comes with the Node.js on the host) (see the development guide for the exact check).

## Upstream
### About upstream npm
This is the upstream `npm` tool, installed with Node.js. Its version follows the Node.js release.
