package usps

import (
	"encoding/xml"
	"fmt"
	"io"
	"net/http"
	"net/url"
)

// Client is a USPS API client.
type Client struct {
	userID   string
	endpoint string
	client   *http.Client
}

// An Option sets an option on a Client. It has private methods to prevent its
// use outside of this package.
type Option interface {
	set(*Client)
}

// A function adapter that implements the Option interface.
type optionFunc func(*Client)

func (fn optionFunc) set(r *Client) { fn(r) }

// Configure a Client.
func (c *Client) setOption(options ...Option) {
	for _, opt := range options {
		opt.set(c)
	}
}

// WithEndpoint sets the endpoint of the USPS API.
func WithEndpoint(endpoint string) Option {
	return optionFunc(func(c *Client) {
		c.endpoint = endpoint
	})
}

// WithHTTPClient sets the http.Client used to communicate with the USPS API.
func WithHTTPClient(client *http.Client) Option {
	return optionFunc(func(c *Client) {
		c.client = client
	})
}

// NewClient returns a USPS API client.
func NewClient(userID string, options ...Option) *Client {
	c := &Client{
		userID:   userID,
		endpoint: "http://production.shippingapis.com/ShippingAPI.dll",
		client:   http.DefaultClient,
	}
	c.setOption(options...)
	return c
}

//ValidateZip returns non empty Response if successful
func (c *Client) ValidateZip(zipCode string) (*CityStateLookupResponse, error) {
	req, err := http.NewRequest("GET", c.endpoint, nil)
	if err != nil {
		return nil, err
	}

	// Construct the URL encoded query
	query := `<CityStateLookupRequest USERID=%q><ZipCode ID="0"><Zip5>%s</Zip5></ZipCode></CityStateLookupRequest>`
	req.URL.RawQuery = fmt.Sprintf("API=CityStateLookup&XML=%s", url.QueryEscape(fmt.Sprintf(query, c.userID, zipCode)))

	// Get the request
	resp, err := c.client.Do(req)
	if err != nil {
		return nil, err
	}

	var (
		decoder = xml.NewDecoder(resp.Body)
		zipResp *CityStateLookupResponse
		apiErr  *APIError
	)

	for {
		// Read tokens in the XML document from the stream (resp.Body)
		t, err := decoder.Token()
		if err == io.EOF || t == nil {
			break // end of stream
		}
		if err != nil {
			return nil, err
		}
		switch se := t.(type) {
		case xml.StartElement:
			switch se.Name.Local {
			case "CityStateLookupResponse":
				if err = decoder.DecodeElement(&zipResp, &se); err != nil {
					return nil, err
				}
			case "Error":
				if err = decoder.DecodeElement(&apiErr, &se); err != nil {
					return nil, err
				}
			default:
				return nil, fmt.Errorf("unknown element: %q", se.Name.Local)
			}
		}
	}

	if apiErr != nil {
		return nil, apiErr
	}

	if zipResp.ZipCode.Error != nil {
		return nil, zipResp.ZipCode.Error
	}

	return zipResp, nil
}

// CityStateLookupResponse is the XML response for the CityStateLookupRequest
// type.
type CityStateLookupResponse struct {
	XMLName xml.Name `xml:"CityStateLookupResponse" json:"-"`
	ZipCode struct {
		ID    string    `xml:"ID,attr,omitempty" json:"id,omitempty"`
		Zip5  string    `xml:"Zip5,omitempty" json:"zip5,omitempty"`
		City  string    `xml:"City,omitempty" json:"city,omitempty"`
		State string    `xml:"State,omitempty" json:"state,omitempty"`
		Error *APIError `xml:"Error,omitempty" json:"error,omitempty"`
	} `xml:"ZipCode,omitempty" json:"zipcode,omitempty"`
}

// APIError is the XML structure for errors returned by the API.
type APIError struct {
	XMLName     xml.Name `xml:"Error" json:"-"`
	Number      string   `xml:"Number,omitempty" json:"number"`
	Description string   `xml:"Description,omitempty" json:"description"`
	Source      string   `xml:"Source,omitempty" json:"source"`
}

// Implement the error interface.
func (e *APIError) Error() string {
	return e.Description
}
