# unisrv (Unity Server for WebGL)

[![Go Reference](https://pkg.go.dev/badge/github.com/frozenbonito/unisrv.svg)](https://pkg.go.dev/github.com/frozenbonito/unisrv)
[![CI](https://github.com/frozenbonito/unisrv/actions/workflows/ci.yaml/badge.svg)](https://github.com/frozenbonito/unisrv/actions/workflows/ci.yaml)
[![License](https://img.shields.io/github/license/frozenbonito/unisrv)](https://github.com/frozenbonito/unisrv/blob/main/LICENSE)
[![Release](https://img.shields.io/github/v/release/frozenbonito/unisrv)](https://github.com/frozenbonito/unisrv/releases/latest)

`unisrv` is a preview server for Unity WebGL applications.

**Notice:** This project is focused on local preview only. Not recommended for production use.

### CLI

#### Installation

Download the binary from the [releases](https://github.com/frozenbonito/unisrv/releases).

#### Running

To start the server, execute the following in your WebGL build location:

```console
unisrv
```

Or specify the build location expressly:

```console
unisrv ./Build/
```

#### Configurations

The server is configurable via the following options or environment variables.

| Option              | Environment Variable      | Default Value | Description                                       |
| ------------------- | ------------------------- | ------------- | ------------------------------------------------- |
| `-base`             | `UNISRV_BASE`             |               | The base path for Unity application.              |
| `-disable-no-cache` | `UNISRV_DISABLE_NO_CACHE` | false         | Disable setting `Cache-Control: no-cache` header. |
| `-host`             | `UNISRV_HOST`             | `localhost`   | The hostname to listen on.                        |
| `-port`             | `UNISRV_PORT`             | 5000          | The port number to listen on.                     |
| `-read-timeout`     | `UNISRV_READ_TIMEOUT`     | 5             | The maximum duration for reading request.         |
| `-write-timeout`    | `UNISRV_WRITE_TIMEOUT`    | 5             | The maximum duration for writing response.        |

### Docker image

[Docker images](https://hub.docker.com/repository/docker/frozenbonito/unisrv) are also available.

```console
docker run --rm -v $(pwd):/app -p 5000:5000 frozenbonito/unisrv
```

Mount your Unity application to `/app` directory in the container.

### Library

It can also be used as a library for Go.

For example:

```go
package main

import (
	"net/http"

	"github.com/frozenbonito/unisrv"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
)

func main() {
	r := chi.NewRouter()
	r.Use(middleware.Logger)
	r.Mount("/", unisrv.NewHandler("/path/to/unity-build-location", nil))
	http.ListenAndServe(":8080", r)
}
```

See [go.dev](https://pkg.go.dev/github.com/frozenbonito/unisrv) for more details.
