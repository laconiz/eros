package http

import (
	"bytes"
	"io/ioutil"
	"net/http"
	"regexp"
	"time"

	"github.com/laconiz/eros/json"
)

const defaultUrlPrefix = "http://"

var regUrlPrefix, _ = regexp.Compile(`^http`)

type Client struct {
	client *http.Client
	url    string
	method string
	header http.Header
}

func (c *Client) Clone() *Client {
	n := *c
	return &n
}

func (c *Client) Client(client *http.Client) *Client {
	n := c.Clone()
	n.client = client
	return n
}

func (c *Client) Url(url string) *Client {

	n := c.Clone()

	if !regUrlPrefix.MatchString(url) {
		url = defaultUrlPrefix + url
	}
	c.url = url

	return n
}

func (c *Client) Method(method string) *Client {
	n := c.Clone()
	n.method = method
	return n
}

func (c *Client) Header(header http.Header) *Client {
	n := c.Clone()
	n.header = header
	return n
}

func (c *Client) Do(req, resp interface{}) error {

	var stream []byte
	if req != nil {
		raw, err := json.Marshal(req)
		if err != nil {
			return err
		}
		stream = raw
	}

	request, err := http.NewRequest(c.method, c.url, bytes.NewReader(stream))
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

	stream, err = ioutil.ReadAll(response.Body)
	if err != nil {
		return err
	}

	return json.Unmarshal(stream, resp)
}

func NewClient() *Client {

	return &Client{
		client: &http.Client{
			Transport: &http.Transport{
				MaxIdleConns:        20,
				MaxIdleConnsPerHost: 2,
				MaxConnsPerHost:     10,
				IdleConnTimeout:     time.Minute,
			},
			Timeout: time.Second * 5,
		},
		url:    "",
		method: "",
		header: http.Header{
			"Accept":       []string{"application/json"},
			"content-type": []string{"application/json;charset=utf-8"},
		},
	}
}
