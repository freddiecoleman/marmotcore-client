package main

import (
	"fmt"
	"net/http"
	"time"
)

type MarmotcoreClient struct {
	protocol   string
	host       string
	port       string
	apiVersion string
}

type HTTPClient interface {
	Get(url string) (resp *http.Response, err error)
}

var Client HTTPClient

func init() {
	Client = &http.Client{Timeout: time.Duration(10) * time.Second}
}

func (mc MarmotcoreClient) url() string {
	return mc.protocol + "://" + mc.host + ":" + mc.port + "/" + mc.apiVersion
}

func (mc MarmotcoreClient) getRequest(path string) (resp *http.Response, err error) {
	return Client.Get(mc.url() + path)
}

func main() {
	mc := &MarmotcoreClient{
		protocol:   "http",
		host:       "localhost",
		port:       "3000",
		apiVersion: "v1",
	}
	nodes, err := mc.GetNodes()
	if err != nil {
		fmt.Printf("Error %s", err)
		return
	}
	fmt.Print(nodes)
}
