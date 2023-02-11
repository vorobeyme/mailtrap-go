# mailtrap-go

Unofficial mailtrap.io API client for Go.

## Install
```
go get github.com/vorobeyme/mailtrap-go
```

## Usage

```
import "github.com/vorobeyme/mailtrap-go"
```

Create a new Mailtrap client, then use the exposed services to access different parts of the Mailtrap API.

```
package main

import (
    "github.com/vorobeyme/mailtrap-go"
)

func main() {
    client := mailtrap.New("api-token")
    client.SendEmail.Send(&mailtrap.SendEmailRequest{})
}
```

## Examples

To send email:

```
request := &mailtrap.SendEmailRequest{
    From:    mailtrap.EmailAddress{Email: "jd@example.com", Name: "John Doe"},
    To:      []mailtrap.EmailAddress{{Email: "to@example.com"}},
    Subject: "Subject",
    Text:    "Hello, world!",
}

resp, _, err := mailtrap.SendEmail.Send(request)

if err != nil {
    fmt.Printf("Something went wrong: %s", err)
    return err
}
```

## License

[MIT License](./LICENSE)