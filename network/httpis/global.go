package httpis

import (
	"net/http"
	"time"
)

func URL(url string) *Connector {
	return global.URL(url)
}

func Header(header http.Header) *Connector {
	return global.Header(header)
}

func Put(req, resp interface{}) error {
	return global.Put(req, resp)
}

func Get(req, resp interface{}) error {
	return global.Get(req, resp)
}

func Post(req, resp interface{}) error {
	return global.Post(req, resp)
}

func Delete(req, resp interface{}) error {
	return global.Delete(req, resp)
}

var global = NewConnector(&http.Client{
	Transport: &http.Transport{
		MaxIdleConns:        20,
		MaxIdleConnsPerHost: 2,
		MaxConnsPerHost:     10,
		IdleConnTimeout:     time.Minute,
	},
})
