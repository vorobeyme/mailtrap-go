# mailtrap-go

[![Go](https://github.com/vorobeyme/mailtrap-go/actions/workflows/go.yml/badge.svg)](https://github.com/vorobeyme/mailtrap-go/actions/workflows/go.yml)
[![codecov](https://codecov.io/gh/vorobeyme/mailtrap-go/branch/main/graph/badge.svg?token=III91WIPLL)](https://codecov.io/gh/vorobeyme/mailtrap-go)
[![Go Report Card](https://goreportcard.com/badge/github.com/vorobeyme/mailtrap-go)](https://goreportcard.com/report/github.com/vorobeyme/mailtrap-go)


Unofficial Mailtrap API client for Go.

The public API documentation is available at [https://api-docs.mailtrap.io](https://api-docs.mailtrap.io/docs/mailtrap-api-docs).

**NOTE: THIS PACKAGE IS STILL UNDER DEVELOPMENT.**

## Installation
```
go get github.com/vorobeyme/mailtrap-go
```

## Usage

```go
import "github.com/vorobeyme/mailtrap-go"
```

Create a new Mailtrap client, then use the exposed services to access different parts of the Mailtrap API.

```go
package main

import (
    "log"

    "github.com/vorobeyme/mailtrap-go"
)

func main() {
    client, err := mailtrap.NewSendingClient("api-token")
    if err != nil {
        log.Fatal(err)
    }

    email := &mailtrap.SendEmailRequest{ ... }
    resp, _, err := client.SendEmail.Send(email)
}
```

## Examples

To find code examples that demonstrate how to call the Mailtrap API client for Go, see the [examples](/examples/) folder.


## License

[MIT License](./LICENSE)