package usps

import (
	"encoding/xml"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
)

//USPS API Signature
const (
	Scheme = "http://"
	Host   = "production.shippingapis.com"
	Path   = "/ShippingAPI.dll"
	Type   = "CityStateLookup"
)

//ValidateZip returns non empty Response if successful
func (a API) ValidateZip(zipCode string) (*CityStateLookupResponse, error) {
	client := a.HTTPClient
	r := &CityStateLookupResponse{}
	e := &apiError{}

	xmlVals := "<CityStateLookupRequest USERID='" + a.Credentials + "'><ZipCode ID='0'> <Zip5>" + zipCode + "</Zip5></ZipCode></CityStateLookupRequest>"

	//Build out URL
	u, err := url.Parse(Scheme + Host)
	if err != nil {
		return r, err
	}

	u.Path = Path
	q := u.Query()
	q.Set("API", Type)
	q.Set("XML", xmlVals)
	u.RawQuery = q.Encode()

	//Get Request
	resp, err := client.Get(u.String())
	if err != nil {
		return r, err
	}

	//Parse XML
	body, readerr := ioutil.ReadAll(resp.Body)
	if readerr == nil {
		//load xml object
		if xmlerr := xml.Unmarshal(body, &r); xmlerr != nil {
			return r, xmlerr
		}
	} else {
		return r, readerr
	}

	//Handle USPS error messages
	if r.ZipCode.ID == "" {
		if err := xml.Unmarshal(body, &e); err != nil {
			return r, nil
		}
		return r, e
	}

	return r, nil
}

//NewUSPSApi returns an API struct
func NewUSPSApi(username string) *API {
	c := &API{
		Credentials: username,
		HTTPClient:  http.DefaultClient,
	}
	return c
}

// *****************************************************************************
// Structs for the USPS API
// *****************************************************************************

//API struct for USPS API settings
type API struct {
	Credentials string
	HTTPClient  *http.Client
}

// CityStateLookupResponse is the XML response for the CityStateLookupRequest
// type.
type CityStateLookupResponse struct {
	XMLName xml.Name `xml:"CityStateLookupResponse" json:"-"`
	ZipCode struct {
		ID    string `xml:"ID,attr,omitempty" json:"id,omitempty"`
		Zip5  string `xml:"Zip5,omitempty" json:"zip5,omitempty"`
		City  string `xml:"City,omitempty" json:"city,omitempty"`
		State string `xml:"State,omitempty" json:"state,omitempty"`
	} `xml:"ZipCode,omitempty" json:"zipcode,omitempty"`
}

// apiError is the XML structure for errors returned by the API.
type apiError struct {
	XMLName     xml.Name `xml:"Error" json:"-"`
	Number      string   `xml:"Number,omitempty" json:"number"`
	Description string   `xml:"Description,omitempty" json:"description"`
	Source      string   `xml:"Source,omitempty" json:"source"`
}

// Implement the error interface.
func (e *apiError) Error() string {
	return fmt.Sprintf("number: %s, source: %s; %s", e.Number, e.Source, e.Description)
}
