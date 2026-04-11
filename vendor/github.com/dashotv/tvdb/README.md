# dashotv/tvdb

Golang TVDB Client (Alpha)

[![Build Status](https://drone.dasho.net/api/badges/dashotv/tvdb/status.svg?ref=refs/heads/main)](https://drone.dasho.net/dashotv/tvdb)
[![GoDoc](https://godoc.org/github.com/dashotv/tvdb?status.svg)](https://godoc.org/github.com/dashotv/tvdb)
[![Go Report Card](https://goreportcard.com/badge/github.com/dashotv/tvdb)](https://goreportcard.com/report/github.com/dashotv/tvdb)
![License](https://img.shields.io/badge/license-MIT-blue.svg)

## Generated Code

This package is generated from the [TVDB OpenAPI](https://api.thetvdb.com/swagger) specification using the [Speakeasy](https://speakeasyapi.dev) code generator.

The generated code is in the `openapi` directory, and I've written scripts to wrap
it with a more convenient interface.

See the `Generation` section below for more information.

> [!NOTE]
> I have needed to make a few tweaks to the openapi spec to get things to work correctly in go.

## Additional Documentation

The Speakeasy generator also generates documentation for the API. This is available
in the openapi directory's [README.md](openapi/README.md) and the [docs](openapi/docs) directory.

## Status

The code is currently in alpha, and is not ready for production use. It is being used
by the DashoTV project, but is not yet stable.

## Usage

Install the package with:

```bash
go get github.com/dashotv/tvdb
```

Import the package with:

```go
import "github.com/dashotv/tvdb"
```

Create a new client with:

```go
client := tvdb.New(apikey, token)
```

If you don't already have a token, you can authenticate with your `apikey` using the `Login` method:

```go
// Authenticate with your API key. This will return a client with the
// token already configured.
client, err := tvdb.Login(apikey)
if err != nil {
    // handle error
}
client.Token // the token
```

> [!IMPORTANT]
> You should store the token somewhere, by default the token is viable for 30 days. TVDB doesn't appear to care if you auth every call, but it adds a lot of overhead.

If you already have the token, you can create a client with your `apikey` and `token`:

```go
client := tvdb.New(apikey, token)
```

## Development

Create a local `.env` file with the following content:

```
# .env
TVDB_API_URL=https://api4.thetvdb.com/v4
TVDB_API_KEY=your_api_key
TVDB_API_TOKEN=your_api_token
```

Run the following to get the make targets:

```
> make help

Usage:
  make <target>

Targets:
  General:
    generate            Generate code from openapi.yml spec
    clean               Remove build related file
  Test:
    test                Run the tests of the project
    coverage            Run the tests of the project and export the coverage
  Lint:
    lint                Run all available linters
    lint-go             Use golintci-lint on your project
  Help:
    help                Show this help.

```

### Generation

To update the generated code, ensure you have the dependencies installed:

- Speakeasy CLI - See [Speakeasy](https://speakeasyapi.dev) for more information.
- Ruby - See [Ruby](https://www.ruby-lang.org/en/documentation/installation/) for more information.

Then run:

```bash
make generate
```

This will run the `scripts/generate.sh` script, which will:

- Run the speakeasy cli generator
- Rearrange the generated code
- Run a ruby script to build the wrapper client code

### Notes

There are several operations on the api that are disabled (this is handled in the
ruby script). These are operations that require additional priveleges (actions on
behalf of the user) or that don't work correctly. These are configured as an array
in the ruby script.
