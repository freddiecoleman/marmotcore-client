package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
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

type Node struct {
	UserId       string `json:"user_id"`
	CreatedTime  int64  `json:"created_time"`
	NodeId       string `json:"node_id"`
	PublicIp     string `json:"public_ip"`
	Region       string `json:"region"`
	InstanceType string `json:"instance_type"`
	ChiaVersion  string `json:"chia_version"`
	Network      string `json:"network"`
	State        string `json:"state"`
	Deleted      bool   `json:"deleted"`
	DeletedTime  int    `json:"deleted_time,omitempty"`
}

type Nodes struct {
	Nodes []Node `json:"nodes"`
}

func (mc MarmotcoreClient) GetNodes() (Nodes, error) {
	var nodes Nodes

	resp, err := mc.getRequest("/nodes")

	if err != nil {
		fmt.Printf("Error %s", err)
		return nodes, err
	}

	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)

	json.Unmarshal(body, &nodes)

	return nodes, nil
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
