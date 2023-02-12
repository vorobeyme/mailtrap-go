# mailtrap-go

[![Go](https://github.com/vorobeyme/mailtrap-go/actions/workflows/go.yml/badge.svg)](https://github.com/vorobeyme/mailtrap-go/actions/workflows/go.yml)

Unofficial Mailtrap API client for Go.

The public API documentation is available at [https://api-docs.mailtrap.io](https://api-docs.mailtrap.io/docs/mailtrap-api-docs).

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
    client, err := mailtrap.New("api-token")
    if err != nil {
        log.Fatal(err)
    }

    resp, _, err := client.SendEmail.Send(&mailtrap.SendEmailRequest{})
}
```

## Examples

To find code examples that demonstrate how to call the Mailtrap API client for Go, see the top-level [examples](/examples/) folder within this repository.


## License

[MIT License](./LICENSE)