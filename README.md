# Go WhatsApp Multi-Device Implementation in REST API

This repository contains example of implementation [go.mau.fi/whatsmeow](https://go.mau.fi/whatsmeow/) package. This example is using a   [labstack/echo](https://github.com/labstack/echo) version 4.x.

## Getting Started

These instructions will get you a copy of the project up and running on your local machine for development and testing purposes.
See deployment for notes on how to deploy the project on a live system.

### Prerequisites

Prequisites package:
* Go (Go Programming Language)
* Make (Automated Execution using Makefile)

Optional package:
* GoReleaser (Go Automated Binaries Build)
* Docker (Application Containerization)

### Installing

Below is the instructions to make this codebase running:
* Create a Go Workspace directory and export it as the extended GOPATH directory
```
cd <your_go_workspace_directory>
export GOPATH=$GOPATH:"`pwd`"
```
* Under the Go Workspace directory create a source directory
```
mkdir -p src/github.com/dimaskiddo/go-whatsapp-multidevice-rest
```
* Move to the created directory and pull codebase
```
cd src/github.com/dimaskiddo/go-whatsapp-multidevice-rest
git clone -b master https://github.com/dimaskiddo/go-whatsapp-multidevice-rest.git .
```
* Run following command to pull dependecies package
```
make vendor
```
* Until this step you already can run this code by using this command
```
ln -sf .env.development .env
make run
```

## Running The Tests

Currently the test is not ready yet :)

## Deployment

To build this code to binaries for distribution purposes you can run following command:
```
make release
```
The build result will shown in build directory

## API Access

You can access any endpoint under **BASE_URL** environment variable which by default located at *.env* file.

Integrated API Documentation can be accessed in **BASE_URL**/docs/index.html or by default it's in `127.0.0.1:3000/api/v1/whatsapp/docs/index.html`

## Built With

* [Go](https://golang.org/) - Go Programming Languange
* [GoReleaser](https://github.com/goreleaser/goreleaser) - Go Automated Binaries Build
* [Make](https://www.gnu.org/software/make/) - GNU Make Automated Execution
* [Docker](https://www.docker.com/) - Application Containerization

## Authors

* **Dimas Restu Hidayanto** - *Initial Work* - [DimasKiddo](https://github.com/dimaskiddo)

See also the list of [contributors](https://github.com/dimaskiddo/go-whatsapp-multidevice-rest/contributors) who participated in this project

## Annotation

You can seek more information for the make command parameters in the [Makefile](https://github.com/dimaskiddo/go-whatsapp-multidevice-rest/-/raw/master/Makefile)