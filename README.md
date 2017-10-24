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

$ ./currency-fetcher --config dev.config.yaml

## Service Configuration

The configuration file is specified on the command-line using the '--config' parameter as above.

The config file contains the following sections:

| Section | Parameter | Specific to Environment | Description | Dev Env values |
|:---:|:---:|:---:|---|---|
| **global** | | | | |
| | statsdhost | No | The location of the statsd client (If empty string, no stats are published) | "127.0.0.1:8125" |
| | refreshinterval | No | The interval between which to pull updates (in minutes) | 60 |
| **zookeeper** | | | | |
| | hosts | Yes | The list of zookeeper hosts to connect to | [ "zookeeper01.dev.uswest2.dmxleo.internal:2181", "zookeeper02.dev.uswest2.dmxleo.internal:2181", "zookeeper03.dev.uswest2.dmxleo.internal:2181" ] |
| | masterpath | No | The path in zookeeper used to determine which instance is master | "/currency-master-e82b2ac7-9c61-4523-b171-c472052044ff" |
| | timeoutms | No | Zookeeper timeout (in ms) | 10000 |
| **kafka** | | | | |
| | brokers | Yes | The list of Kafka brokers to connect to | [ "kafka01.dev.uswest2.dmxleo.internal:9092", "kafka02.dev.uswest2.dmxleo.internal:9092", "kafka03.dev.uswest2.dmxleo.internal:9092", "kafka04.dev.uswest2.dmxleo.internal:9092" ] |
| | clientid | Yes | The client ID to use when connecting to Kafka (this needs to be different on each instance) | cfkc01 |
| | enabled | No | Boolean indicating whether to publish messages on Kafka or not | true |
| | maxretries | No | The number of times to retry sending a message | 3 |
| | topic | No | The Kafka topic to publish to | CurrencyUpdate |
| **http** | | | | |
| | address | No | The HTTP address and port to bind to | ":8080" |
| | readtimeout | No | The HTTP read timeout (in seconds) | 5 |
| | writetimeout | No | The HTTP write timeout (in seconds) | 10 |
| **xe** | | | | |
| | endpoint | No | The XE.com API endpoint to hit | "https://xecdapi.xe.com/v1/convert_from?to=%s&from=%s&amount=1" |
| | username | No | The XE.com API username | |
| | password | No | The XE.com API password | |
| | fromcurrencies | No | The list of currencies that we fetch data for | ["AUD", "JPY", "EUR", "GBP", "CAD"] |
| | tocurrency | No | The currency we find rates for | "USD" |
| **log** | | | | |
| | level | Yes | The logging level (INFO, WARN, ERROR) | "INFO" |
| | filename | Yes | The file to write logs to | "/var/log/currency-fetcher/dev.log" |
| | json | Yes | Boolean, if true then logs are written in JSON format | false |
| | maxsizemb | No | The max. file size after which a new file is created (in MB) | 500 |
| | maxbackups | No | The max. number of backup files to keep | 7 |
| | maxagedays | No | The max. number of days to keep backup files for | 10 |
| | writestdout | No | Boolean, if true then logs are written to stdout instead of a file | false |


### Run on Docker

An alternative is to run the API in an isolated environment, using docker containers.

* Install [Docker](https://www.docker.com/) and [Docker Compose](https://docs.docker.com/compose/)
  * `brew cask install docker` on MacOS, and then launch the `Docker` application to start the docker daemon
* Clone this repository
  * `git clone https://github.com/dailymotion-leo/currency-fetcher.git`
* Start the application (and all its dependencies)
  * `make run-in-docker`
* The API is exposed on the port `8090`. And check the available HTTP endpoints.
* At the end, just hit `ctrl+c` in the docker output to stop the containers
  * If you want to clean at the end (except for the local database's data), run the following command: `docker-compose kill && docker-compose rm -f && docker system prune -f`

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

