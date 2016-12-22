package usps

import (
	"encoding/xml"
	"errors"
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
func (a API) ValidateZip(zipCode string) (*Response, error) {
	client := a.HTTPClient
	r := &Response{}
	e := &Error{}

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
		if xmlerr := xml.Unmarshal(body, &r.CityStateLookupResponse); xmlerr != nil {
			return r, xmlerr
		}
	} else {
		return r, readerr
	}

	//Handle USPS error messages
	if r.CityStateLookupResponse.ZipCode == nil {
		xml.Unmarshal(body, &e)
		err := errors.New(e.Description.Text)
		return r, err
	}
	if r.CityStateLookupResponse.ZipCode.Error != nil {
		err := errors.New(r.CityStateLookupResponse.ZipCode.Error.Description.Text)
		return r, err
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

//Response for full response
type Response struct {
	CityStateLookupResponse *CityStateLookupResponse `xml:"CityStateLookupResponse,omitempty" json:"response,omitempty"`
}

//CityStateLookupResponse for success/error response
type CityStateLookupResponse struct {
	ZipCode *ZipCode `xml:"ZipCode,omitempty" json:"zip_code,omitempty"`
}

//ZipCode for success/error response
type ZipCode struct {
	ID    string `xml:"ID,attr"  json:"id,"`
	City  *City  `xml:"City,omitempty" json:"city,omitempty"`
	State *State `xml:"State,omitempty" json:"state,omitempty"`
	Zip5  *Zip5  `xml:"Zip5,omitempty" json:"zip5,omitempty"`
	Error *Error `xml:"Error,omitempty" json:"error,omitempty"`
}

//Zip5 for successful response
type Zip5 struct {
	Text string `xml:",chardata" json:",omitempty"`
}

//City for successful response
type City struct {
	Text string `xml:",chardata" json:",omitempty"`
}

//State for successful response
type State struct {
	Text string `xml:",chardata" json:",omitempty"`
}

//Error for error response
type Error struct {
	Description *Description `xml:" Description,omitempty" json:"description,omitempty"`
	HelpContext *HelpContext `xml:" HelpContext,omitempty" json:"help_context,omitempty"`
	HelpFile    *HelpFile    `xml:" HelpFile,omitempty" json:"help_file,omitempty"`
	Number      *Number      `xml:" Number,omitempty" json:"number,omitempty"`
	Source      *Source      `xml:" Source,omitempty" json:"source,omitempty"`
}

//Description for for error response
type Description struct {
	Text string `xml:",chardata" json:",omitempty"`
}

//HelpContext for error response
type HelpContext struct {
}

//HelpFile for error response
type HelpFile struct {
}

//Number for error response
type Number struct {
	Text string `xml:",chardata" json:",omitempty"`
}

//Source for error response
type Source struct {
	Text string `xml:",chardata" json:",omitempty"`
}
