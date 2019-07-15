# Armeria

Armeria is an open-source web-based online role-playing game being written in [Golang](https://golang.org/)
and [Vue.js](https://vuejs.org/).

## Appendix

* [Contributing](#contributing)
    * [Getting Started](#getting-started)
    * [Local Development](#local-development)
    * [Upgrading Dependencies](#upgrading-dependencies)
    * [Publishing Features](#publishing-features)
    * [CI/CD](#cicd)
    * [Releasing](#releasing)

## Contributing

Since Armeria is open-source, anyone can run Armeria locally and contribute to the game. All contributions are welcome
and extremely appreciated. The only thing not included is access to the production data. However, there is an
abundance of example data included with the repo.

### Getting Started

To begin contributing to Armeria, create a fork of the `armeria` repo. You can work on your changes on a feature
branch based off of your forked repo.

### Local Development

To get started, you will need the following installed:

* Golang (>= 1.11)
* Node.js (>= 10.0.0)
* Yarn (>= 1.6.0)

To run the server, you can either make  a copy of the `example-data` directory (called `data`), or symlink to it. If
you plan on contributing additional example data upstream, it is best to make a symlink to it so the game can modify
the example-data as you play locally.

To copy the data directory and not contribute data changes upstream:

```bash
$ cp -R ./example-data ./data
```

To symlink the data directory and contribute data changes upstream:

```bash
$ ln -s ./example-data ./data
```

To build and run the game:

```bash
$ go run cmd/armeria/main.go
```

To build and run the web client:

```bash
$ cd ./client
$ yarn install
$ yarn serve
```

You can load the client here: [http://localhost:8080/](http://localhost:8080/). While running `yarn serve`, you can
make changes to the client files and the changes will be hot reloaded immediately within the browser. Some of these
changes may terminate your connection to the game server and require you to re-login.

### Upgrading Dependencies

To upgrade Vue.js, you can use the Vue CLI. Be sure you have the Vue CLI and the Vue CLI Upgrade packages installed:

```bash
$ yarn global add @vue/cli
$ yarn global add @vue/cli-upgrade
```

You can upgrade the client by running:

```bash
$ cd client
$ vue upgrade major
```

Press `Y` to confirm updating `package.json` and re-running `yarn install` automatically. Make sure the client builds
successfully and commit the changes.

To upgrade another Node dependency, you can do this by running:

```bash
$ cd client
$ yarn add <dependency name>
```

### Publishing Features

When you have finished working on a feature, be sure there are no broken unit tests. You can check this via:

```bash
$ go test ./...
```

If you are adding new code, be sure you are also modifying or adding unit tests to ensure test coverage. Furthermore,
anything that should be documented externally should be written in `/docs` as well.

Once ready, create a Pull Request (PR) from your forked repo's feature branch to this repo's `master` branch.

### CI/CD

Once you've created a Pull Request, our CircleCI build will kick off and ensure all unit tests are passing. The build
will additionally confirm that the binary can be built and ran successfully. If any of these steps fail, you will be
notified via the open PR.

### Releasing

The main repo has a `master` branch and `release` branch. The former is a development branch that will contain squashed
merges from contributors that are feature-complete (and ready to be released at any time). Once we've deemed it
necessary to cut a release, we will create a Pull Request to merge `master` into `release`. Once merged, CircleCI will
deploy the code to the production server and trigger a server restart. Characters on the game will be kicked off for
a few seconds while the new server instance starts back up.