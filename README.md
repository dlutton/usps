# usps

Zip Code validation using the USPS API in Go.

## Examples

### Authentication

If you already have the USERID, creating the client is simple:

```go
client := usps.New("USERID")
```

### Google App Engine settings

````go
client := usps.New("USERID", usps.HTTPClient(&http.Client{
	Transport: &urlfetch.Transport{Context: appengine.NewContext(r)},
}))
````

### Full Example

```go
package main

import (
	"encoding/json"
	"log"

	"github.com/dlutton/usps"
)

func main(){
	client := usps.New("USERID")

	results, err := client.ValidateZip("91362")
	if err != nil {
		log.Fatal(err)
	}

	output, err := json.Marshal(results)
	if err != nil {
		log.Fatal(err)
	}

	log.Print(string(output))
}
```

