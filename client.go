package keen

import (
	"encoding/json"
	"fmt"
	"gopkg.in/httgo/interfaces.v1/httpclient"
	"gopkg.in/nowk/jsonify.v1"
	"io"
	"net/http"
	"net/url"
)

const (
	SCHEME      = "https"
	HOST        = "api.keen.io"
	API_VERSION = "3.0"
)

type Client struct {
	Client    httpclient.Interface
	ReadKey   string
	WriteKey  string
	MasterKey string
	ProjectID string
}

func NewClient(id string, opts ...func(c *Client)) *Client {
	c := &Client{
		Client:    http.DefaultClient,
		ProjectID: id,
	}

	for _, v := range opts {
		v(c)
	}

	return c
}

func (c Client) Do(req *http.Request) (*http.Response, error) {
	return c.Client.Do(req)
}

func URL(path string) *url.URL {
	return &url.URL{
		Scheme: SCHEME,
		Host:   HOST,
		Path:   path,
	}
}

func NewRequest(reso Resource) (*http.Request, error) {
	if reso.Authorization() == "" {
		return nil, fmt.Errorf("key error: authorization key required")
	}

	r, err := jsonify.NewReader(reso.Data())
	if err != nil {
		return nil, err
	}

	u := URL(reso.Path())
	req, err := http.NewRequest("POST", u.String(), r)
	if err != nil {
		return nil, err
	}

	req.Header.Add("Authorization", reso.Authorization())
	req.Header.Add("Content-Type", "application/json")
	return req, nil
}

func decode(r io.ReadCloser, d interface{}) error {
	defer r.Close()
	return json.NewDecoder(r).Decode(d)
}

// Write writes an Event to keen returning an error if not created
func (c Client) Write(e Event) error {
	resp, err := c.WriteResp(e)
	if err != nil {
		return err
	}

	var d EventResponse
	if err := decode(resp.Body, &d); err != nil {
		return err
	}
	if !d.Created {
		return fmt.Errorf("event not created: [%s] %s", d.ErrorCode, d.Message)
	}

	return nil
}

// WriteResp returns the raw http response when writing an event to keen
func (c Client) WriteResp(evt Event) (*http.Response, error) {
	req, err := NewRequest(&EventResource{evt, c.WriteKey, c.ProjectID})
	if err != nil {
		return nil, err
	}

	return c.Do(req)
}

// Query makes a keen query and decods the results to a given object
func (c Client) Query(qry Query, d interface{}) error {
	resp, err := c.QueryResp(qry)
	if err != nil {
		return err
	}

	if err := decode(resp.Body, d); err != nil {
		return err
	}

	return nil
}

// QueryResp makes a keen query and returns the raw http response
func (c Client) QueryResp(qry Query) (*http.Response, error) {
	req, err := NewRequest(&QueryResource{qry, c.ReadKey, c.ProjectID})
	if err != nil {
		return nil, err
	}

	return c.Do(req)
}
