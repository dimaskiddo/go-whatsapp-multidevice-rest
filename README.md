# Go WhatsApp Multi-Device Implementation in REST API

This repository contains example of implementation [go.mau.fi/whatsmeow](https://go.mau.fi/whatsmeow/) package with Multi-Session/Account Support. This example is using a [labstack/echo](https://github.com/labstack/echo) version 4.x.

## Features

- Multi-Session/Account Support
- Multi-Device Support
- WhatsApp Authentication (QR Code and Logout)
- WhatsApp Messaging Send Text
- WhatsApp Messaging Send Media (Document, Image, Audio, Video, Sticker)
- WhatsApp Messaging Send Location
- WhatsApp Messaging Send Contact
- WhatsApp Messaging Send Link
- And Much More ...

## Getting Started

These instructions will get you a copy of the project up and running on your local machine for development and testing purposes.
See deployment section for notes on how to deploy the project on a live system.

### Prerequisites

Prequisites packages:
* Go (Go Programming Language)
* Swag (Go Annotations Converter to Swagger Documentation)
* GoReleaser (Go Automated Binaries Build)
* Make (Automated Execution using Makefile)

Optional packages:
* Docker (Application Containerization)

### Deployment

#### **Using Container**

1) Install Docker CE based on the [manual documentation](https://docs.docker.com/desktop/)

2) Run the following command on your Terminal or PowerShell
```sh
docker run -d \
  -p 3000:3000 \
  --name go-whatsapp-multidevice \
  --rm dimaskiddo/go-whatsapp-multidevice-rest:latest
```

3) Now it should be accessible in your machine by accessing `localhost:3000/api/v1/whatsapp` or `127.0.0.1:3000/api/v1/whatsapp`

4) Try to use integrated API docs that accesible in `localhost:3000/api/v1/whatsapp/docs/` or `127.0.0.1:3000/api/v1/whatsapp/docs/`

#### **Using Pre-Build Binaries**

1) Download Pre-Build Binaries from the [release page](https://github.com/dimaskiddo/go-whatsapp-multidevice-rest/releases)

2) Extract the zipped file

3) Copy the `.env.default` file as `.env` file

4) Run the pre-build binary
```sh
# MacOS / Linux
chmod 755 go-whatsapp-multidevice-rest
./go-whatsapp-multidevice-rest

# Windows
# You can double click it or using PowerShell
.\go-whatsapp-multidevice-rest.exe
```

5) Now it should be accessible in your machine by accessing `localhost:3000/api/v1/whatsapp` or `127.0.0.1:3000/api/v1/whatsapp`

6) Try to use integrated API docs that accesible in `localhost:3000/api/v1/whatsapp/docs/` or `127.0.0.1:3000/api/v1/whatsapp/docs/`

#### **Build From Source**

Below is the instructions to make this source code running:

1) Create a Go Workspace directory and export it as the extended GOPATH directory
```sh
cd <your_go_workspace_directory>
export GOPATH=$GOPATH:"`pwd`"
```

2) Under the Go Workspace directory create a source directory
```sh
mkdir -p src/github.com/dimaskiddo/go-whatsapp-multidevice-rest
```

3) Move to the created directory and pull codebase
```sh
cd src/github.com/dimaskiddo/go-whatsapp-multidevice-rest
git clone -b master https://github.com/dimaskiddo/go-whatsapp-multidevice-rest.git .
```

4) Run following command to pull vendor packages
```sh
make vendor
```

5) Link or copy environment variables file
```sh
ln -sf .env.development .env
# - OR -
cp .env.development .env
```

6) Until this step you already can run this code by using this command
```sh
make run
```

7) *(Optional)* Use following command to build this code into binary spesific platform
```sh
make build
```

8) *(Optional)* To make mass binaries distribution you can use following command
```sh
make release
```

9) Now it should be accessible in your machine by accessing `localhost:3000/api/v1/whatsapp` or `127.0.0.1:3000/api/v1/whatsapp`

10) Try to use integrated API docs that accesible in `localhost:3000/api/v1/whatsapp/docs/` or `127.0.0.1:3000/api/v1/whatsapp/docs/`

## API Access

You can access any endpoint under **HTTP_BASE_URL** environment variable which by default located at `.env` file.

Integrated API Documentation can be accessed in `<HTTP_BASE_URL>/docs/` or by default it's in `localhost:3000/api/v1/whatsapp/docs/` or `127.0.0.1:3000/api/v1/whatsapp/docs/`

## Running The Tests

Currently the test is not ready yet :)

## Built With

* [Go](https://golang.org/) - Go Programming Languange
* [Swag](https://github.com/swaggo/swag) - Go Annotations Converter to Swagger Documentation
* [GoReleaser](https://github.com/goreleaser/goreleaser) - Go Automated Binaries Build
* [Make](https://www.gnu.org/software/make/) - GNU Make Automated Execution
* [Docker](https://www.docker.com/) - Application Containerization

## Authors

* **Dimas Restu Hidayanto** - *Initial Work* - [DimasKiddo](https://github.com/dimaskiddo)

See also the list of [contributors](https://github.com/dimaskiddo/go-whatsapp-multidevice-rest/contributors) who participated in this project

## Annotation

You can seek more information for the make command parameters in the [Makefile](https://github.com/dimaskiddo/go-whatsapp-multidevice-rest/-/raw/master/Makefile)

## License

Copyright (C) 2022 Dimas Restu Hidayanto

Permission is hereby granted, free of charge, to any person obtaining a copy of this software and associated documentation files (the "Software"), to deal in the Software without restriction, including without limitation the rights to use, copy, modify, merge, publish, distribute, sublicense, and/or sell copies of the Software, and to permit persons to whom the Software is furnished to do so, subject to the following conditions:

The above copyright notice and this permission notice shall be included in all copies or substantial portions of the Software.

THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY, FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM, OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE SOFTWARE.
