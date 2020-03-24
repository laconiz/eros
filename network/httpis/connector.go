package httpis

import (
	"bytes"
	"github.com/laconiz/eros/utils/json"
	"io"
	"io/ioutil"
	"net/http"
	"strings"
)

// ---------------------------------------------------------------------------------------------------------------------

func NewConnector(client *http.Client) *Connector {

	header := http.Header{
		"Accept":       []string{"application/json"},
		"content-type": []string{"application/json;charset=utf-8"},
	}

	return &Connector{client: client, header: header}
}

// ---------------------------------------------------------------------------------------------------------------------

type Connector struct {
	client *http.Client // 客户端
	url    string       // 地址
	method string       // 方法
	header http.Header  // 请求头
}

// ---------------------------------------------------------------------------------------------------------------------

func (connector *Connector) clone() *Connector {

	return &Connector{
		client: connector.client,
		url:    connector.url,
		method: connector.method,
		header: connector.header.Clone(),
	}
}

func (connector *Connector) URL(url string) *Connector {

	const (
		httpPrefix  = "http://"
		httpsPrefix = "https://"
	)
	if !strings.HasPrefix(url, httpPrefix) && !strings.HasPrefix(url, httpsPrefix) {
		url = httpPrefix + url
	}

	n := connector.clone()
	n.url = url
	return n
}

func (connector *Connector) Method(method string) *Connector {
	n := connector.clone()
	n.method = method
	return n
}

func (connector *Connector) Header(header http.Header) *Connector {
	n := connector.clone()
	n.header = header
	return n
}

// ---------------------------------------------------------------------------------------------------------------------

func (connector *Connector) Do(req, resp interface{}) error {

	var reader io.Reader
	if req != nil {
		raw, err := json.Marshal(req)
		if err != nil {
			return err
		}
		reader = bytes.NewReader(raw)
	}

	request, err := http.NewRequest(connector.method, connector.url, reader)
	if err != nil {
		return err
	}
	request.Header = connector.header

	response, err := connector.client.Do(request)
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

// ---------------------------------------------------------------------------------------------------------------------

func (connector *Connector) Put(req, resp interface{}) error {
	return connector.Method(http.MethodPut).Do(req, resp)
}

func (connector *Connector) Get(req, resp interface{}) error {
	return connector.Method(http.MethodGet).Do(req, resp)
}

func (connector *Connector) Post(req, resp interface{}) error {
	return connector.Method(http.MethodPost).Do(req, resp)
}

func (connector *Connector) Delete(req, resp interface{}) error {
	return connector.Method(http.MethodDelete).Do(req, resp)
}
