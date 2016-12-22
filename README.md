usps
-------------
Zip Code validation using the USPS API in GO



Examples
-------------

### Authentication

If you already have the USERID, creating the client is simple:

````go
api := usps.NewUSPSApi("USERID")
````

### Full Example

````go
package main

import (
	"encoding/json"
	"fmt"
	"usps"
)

func main(){
	api := usps.NewUSPSApi("USERID")

	results, err := api.ValidateZip("91362")
	if err != nil {
		fmt.Println(err.Error())
		return
	}

	outgoingJSON, err := json.Marshal(results)
	if err != nil {
		fmt.Println(err.Error())
		return
	}

	fmt.Println(string(outgoingJSON))
}
````



Error Handling
---------------------------------

Errors are returned as the Error Description from the USPS API response


Google App Engine settings
---------------------------------

````go
	api := usps.NewUSPSApi("USERID")
	c := appengine.NewContext(r)
	api.HTTPClient.Transport = &urlfetch.Transport{Context: c}
````
