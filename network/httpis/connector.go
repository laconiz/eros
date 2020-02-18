package httpis

import (
	"bytes"
	"io"
	"io/ioutil"
	"net/http"
	"strings"

	"github.com/laconiz/eros/utils/json"
)

const (
	httpPrefix  = "http://"
	httpsPrefix = "https://"
)

type Connector struct {
	client *http.Client
	url    string
	method string
	header http.Header
}

func (c *Connector) clone() *Connector {
	return &Connector{client: c.client, url: c.url, method: c.method, header: c.header.Clone()}
}

func (c *Connector) URL(url string) *Connector {
	n := c.clone()
	if !strings.HasPrefix(url, httpPrefix) && !strings.HasPrefix(url, httpsPrefix) {
		url = httpPrefix + url
	}
	n.url = url
	return n
}

func (c *Connector) Method(method string) *Connector {
	n := c.clone()
	n.method = method
	return n
}

func (c *Connector) Header(header http.Header) *Connector {
	n := c.clone()
	n.header = header
	return n
}

func (c *Connector) Do(req, resp interface{}) error {

	var reader io.Reader
	if req != nil {
		raw, err := json.Marshal(req)
		if err != nil {
			return err
		}
		reader = bytes.NewReader(raw)
	}

	request, err := http.NewRequest(c.method, c.url, reader)
	if err != nil {
		return err
	}
	request.Header = c.header

	response, err := c.client.Do(request)
	if err != nil {
		return err
	}
	defer response.Body.Close()

	if resp == nil {
		return nil
	}

	stream, err := ioutil.ReadAll(response.Body)
	if err != nil {
		return err
	}

	return json.Unmarshal(stream, resp)
}

func (c *Connector) Put(req, resp interface{}) error {
	return c.clone().Method(http.MethodPut).Do(req, resp)
}

func (c *Connector) Get(req, resp interface{}) error {
	return c.clone().Method(http.MethodGet).Do(req, resp)
}

func (c *Connector) Post(req, resp interface{}) error {
	return c.clone().Method(http.MethodPost).Do(req, resp)
}

func (c *Connector) Delete(req, resp interface{}) error {
	return c.clone().Method(http.MethodDelete).Do(req, resp)
}

func NewConnector(client *http.Client) *Connector {
	return &Connector{
		client: client,
		header: http.Header{
			"Accept":       []string{"application/json"},
			"content-type": []string{"application/json;charset=utf-8"},
		},
	}
}
