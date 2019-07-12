# Armeria

Armeria is an open-source web-based online role-playing game being written in [Golang](https://golang.org/)
and [Vue.js](https://vuejs.org/).

## Contributing

Since Armeria is open-source, anyone can run Armeria locally and contribute to the game. The only thing not included
is access to the production data. However, there is an abundance of example data included with the repo.

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

You can load the client here: [http://localhost:8080/](http://localhost:8080/)