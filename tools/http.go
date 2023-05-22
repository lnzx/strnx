package tools

import (
	"crypto/tls"
	"net/http"
	"time"
)

var client *http.Client

func init() {
	transport := &http.Transport{
		TLSClientConfig:       &tls.Config{InsecureSkipVerify: true},
		MaxIdleConns:          100,
		MaxIdleConnsPerHost:   100,
		IdleConnTimeout:       90 * time.Second,
		TLSHandshakeTimeout:   10 * time.Second,
		ExpectContinueTimeout: 10 * time.Second,
	}
	client = &http.Client{Transport: transport, Timeout: 1 * time.Minute}
}

func Get(url string) (resp *http.Response, err error) {
	return client.Get(url)
}

func Do(req *http.Request) (resp *http.Response, err error) {
	return client.Do(req)
}
