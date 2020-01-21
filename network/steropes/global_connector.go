package steropes

import (
	"net/http"
	"time"
)

func Header(header http.Header) *Connector {
	return globalConnector.Header(header)
}

func Put(req, resp interface{}) error {
	return globalConnector.Put(req, resp)
}

func Get(req, resp interface{}) error {
	return globalConnector.Get(req, resp)
}

func Post(req, resp interface{}) error {
	return globalConnector.Post(req, resp)
}

func Delete(req, resp interface{}) error {
	return globalConnector.Delete(req, resp)
}

var globalConnector = NewConnector(&http.Client{
	Transport: &http.Transport{
		MaxIdleConns:        20,
		MaxIdleConnsPerHost: 2,
		MaxConnsPerHost:     10,
		IdleConnTimeout:     time.Minute,
	},
})
