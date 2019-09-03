# Contributing

Since Armeria is open-source, anyone can run Armeria locally and contribute to
the game. All contributions are welcome and extremely appreciated. The only
thing not included is access to the production data. However, there is an
abundance of example data included with the repo.

## Appendix

* [Getting Started](#getting-started)
* [Local Development](#local-development)
* [Upgrading Dependencies](#upgrading-dependencies)
* [Publishing Features](#publishing-features)
* [CI/CD](#cicd)
* [Releasing](#releasing)
* [Discord](#discord)

## Getting Started

To begin contributing to Armeria, create a fork of the `armeria` repo. You can
work on your changes on a feature branch based off of your forked repo.

## Local Development

To get started, you will need the following installed:

* Golang (>= 1.12)
* Node.js (>= 10.0.0)
* Yarn (>= 1.6.0)

To run the server, you can either make  a copy of the `example-data` directory
(called `data`), or symlink to it. If you plan on contributing additional
example data upstream, it is best to make a symlink to it so the game can modify
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

You can load the client here: [http://localhost:8080/](http://localhost:8080/).

You should now be able to login using the command: `/login admin admin`.

While running `yarn serve`, you can make changes to the client files and the
changes will be hot reloaded immediately within the browser. Some of these
changes may terminate your connection to the game server and require you to
re-login.

## Upgrading Dependencies

This section outlines upgrading dependencies for both the client and the server.

### Client

To upgrade Vue.js, make sure you're using the latest version of the Vue CLI. If
you don't have it installed, install it with:

```bash
$ yarn global add @vue/cli
$ yarn global add @vue/cli-upgrade
```

If you already have the Vue CLI installed, upgrade it to the latest version
with:

```bash
$ yarn global upgrade @vue/cli
$ yarn global upgrade @vue/cli-upgrade
```

Next, use the Vue CLI UI to handle upgrades gracefully. To do this:

```bash
$ cd client
$ vue ui
```

The Vue CLI UI should pop up in your default browser. If the project is not
already imported into the UI, go ahead and import it now.

Use the **Plugins** and **Dependencies** tabs to upgrade all of the plugins and
dependencies that have updates available. Be especially careful when updating
the `vue` and `vuex` dependencies to make sure everything is working within the
client. It's a good idea to read through the Vue/Vuex release notes when doing
this as well.

Note that when performing these upgrades, the `package.json` and `yarn.lock`
files will be updated accordingly and, assuming everything is working, these
should be committed to the repo.

### Server

To upgrade the Golang version, you should first upgade Golang locally and make
sure the server binary can be built successfully. You should then upgrade the
build pipeline by modifying these two files:

* `.github/workflows/dev.yml`
* `.github/workflows/deploy.yml`

Modify the `go-version` setting under the `with` section of the build step.

To upgrade Armeria's dependencies, use:

```
$ go get -u
$ go mod tidy
```

You should also run `go mod tidy` any time you introduce or remove any new or
existing dependencies.

## Publishing Features

When you have finished working on a feature, be sure there are no broken unit
tests. You can check this via:

```bash
$ go test ./...
```

If you are adding new code, be sure you are also modifying or adding unit tests
to ensure test coverage. Furthermore, anything that should be documented
externally should be written in `/docs` as well.

Once ready, create a Pull Request (PR) from your forked repo's feature branch to
this repo's `master` branch.

## CI/CD

Once you create a Pull Request, a [GitHub
Actions](https://github.com/heyitsmdr/armeria/actions) build will kick off and
ensure that all unit tests are passing, the server binary can be built and the
client can be compiled without errors. If any of these steps fail, you will be
notified via the open PR.

## Releasing

A GitHub [release](https://github.com/heyitsmdr/armeria/releases) will be cut
once there is enough content on the `master` branch to warrant a server restart
and build update. There will always be a draft release that contributors to the
main repo will be able to see and they will edit it as they approve/merge PRs
into the `master` branch.

Once a release is published, a GitHub Actions workflow will deploy the changes
to the production version of Armeria and the game server will be restarted with
the new version.

## Discord

If you plan to develop for Armeria, it's strongly encouraged that you join our
community Discord at: https://discord.gg/hMzjH6n (room: #armeria-dev).