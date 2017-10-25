# currency-fetcher
Script that fetches latest currency rates from XE.com and puts them on the Kafka queue.

## Build Instructions

* Install Go 1.8.x
* Clone this repository into: $GOPATH/src/github.com/dailymotion-leo
* To build, run: **make**
    * You may need to pull in packages: **go get {package path}**
* To run tests, run: **make tests**
* To create a linux build, run: **make dist**
* Copy the following files to the service machine:
    * currency-fetcher (the binary)
    * bin/start.sh, bin/stop.sh, bin/status.sh (startup scripts)
    * ${ENV}.config.yaml (the config file appropriate to the environment)

## Running the service

$ ./discomotionslack --config dev.config.yaml

### Run on Docker

An alternative is to run the API in an isolated environment, using docker containers.

* Install [Docker](https://www.docker.com/) and [Docker Compose](https://docs.docker.com/compose/)
  * `brew cask install docker` on MacOS, and then launch the `Docker` application to start the docker daemon
* Clone this repository
  * `git clone https://github.com/dailymotion-leo/currency-fetcher.git`
* Start the application (and all its dependencies)
  * `make run-in-docker`
* The API is exposed on the port `8080`. And check the available HTTP endpoints.
* At the end, just hit `ctrl+c` in the docker output to stop the containers

### Dependencies

Dependencies are managed with [Dep](https://github.com/golang/dep), and are already "vendored" (stored in the `vendor` dir, so that you don't need to retrieve them yourself).

We vendor our dependencies to have reproductible builds and be sure of what we have in the binary.

#### Adding a new dependency

* Get [Dep](https://github.com/golang/dep): `go get -u github.com/golang/dep/cmd/dep`
* Use the new dependency in your source code (for example `import github.com/pkg/errors` in your `.go` files, and do something with it)
  * If your IDE runs `goimports` on save, and complains about the new pkg being unknown, `go get`-it before (for example `go get github.com/pkg/errors` for example)
* Run `dep ensure pkg@version` (for example `dep ensure github.com/pkg/errors@0.8.0`) - the version can be a git tag or branch
  * This will update the `Gopkg.toml` file
  * If you ommit the version, `dep` will not add the new dependency to the `Gopkg.toml` file, only to the `Gopkg.lock` file (which is not what we want)
  * The `Gopkg.lock` file and the `vendor` dir should have changed too. If not, that's because you didn't use the new dependency in your source code. `dep` is smart enough to not add a dependency that you are not using ;-)
* Review the changes, run the tests, and commit

#### Updating dependencies

If you want to update one or more dependencies:

* Get [Dep](https://github.com/golang/dep): `go get -u github.com/golang/dep/cmd/dep`
* Edit the `Gopkg.toml` file to use a more recent version of the dependency
* Run `dep ensure -update`
  * The `Gopkg.lock` file and the `vendor` dir should have changed
* Review the changes, run the tests, and commit

