# usps
[![CircleCI](https://circleci.com/gh/dlutton/usps.svg?style=svg)](https://circleci.com/gh/dlutton/usps)
[![GoDoc](https://godoc.org/github.com/gopher/gopher?status.svg)](https://godoc.org/github.com/dlutton/usps)
[![Go Report Card](https://goreportcard.com/badge/github.com/dlutton/usps)](https://goreportcard.com/report/github.com/dlutton/usps)

Zip Code validation using the USPS API in Go.

## Examples

### Authentication

If you already have the USERID, creating the client is simple:

```go
client := usps.New("USERID")
```

### Google App Engine settings

```go
client := usps.NewClient("USERID", usps.WithHTTPClient(&http.Client{
	Transport: &urlfetch.Transport{Context: appengine.NewContext(r)},
}))
```

### Full Example

```go
package main

import (
	"encoding/json"
	"log"

	"github.com/dlutton/usps"
)

func main() {
	client := usps.NewClient("USERID")

	results, err := client.ValidateZip("91362")
	if err != nil {
		if apiErr, ok := err.(*usps.APIError); ok {
			log.Fatalf("number: %s, source: %s; %s", apiErr.Number, apiErr.Source, apiErr.Description)
		}
		log.Fatal(err)
	}

	output, err := json.Marshal(results)
	if err != nil {
		log.Fatal(err)
	}

	log.Print(string(output))
}
```

### Development

```
go get -v -u github.com/dlutton/usps
```

### License

This software is created under Apache v2 License. For the full license text, please see License.md
