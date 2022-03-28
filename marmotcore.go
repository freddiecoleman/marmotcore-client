package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
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
	Do(req *http.Request) (*http.Response, error)
	Get(url string) (resp *http.Response, err error)
	Post(url string, contentType string, body io.Reader) (resp *http.Response, err error)
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

func (mc MarmotcoreClient) postRequest(path string, body io.Reader) (resp *http.Response, err error) {
	return Client.Post(mc.url()+path, "application/json", body)
}

func (mc MarmotcoreClient) deleteRequest(path string) (resp *http.Response, err error) {
	req, err := http.NewRequest("DELETE", mc.url()+path, nil)
	if err != nil {
		fmt.Printf("error %s", err)
		return
	}
	return Client.Do(req)
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

type CreateNode struct {
	Region       string `json:"region"`
	InstanceType string `json:"instance_type"`
	ChiaVersion  string `json:"chia_version"`
	Network      string `json:"network"`
}

type NodesResponse struct {
	Nodes []Node `json:"nodes"`
}

type NodeResponse struct {
	Node Node `json:"node"`
}

type CreateNodeResponse struct {
	NodeId string `json:"node_id"`
}

type DeleteNodeResponse struct {
	Deleted bool `json:"deleted"`
}

func (mc MarmotcoreClient) GetNodes() (NodesResponse, error) {
	var nodes NodesResponse

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

func (mc MarmotcoreClient) CreateNode(createNode *CreateNode) (CreateNodeResponse, error) {
	var createNodeResponse CreateNodeResponse

	createNodeBytes, err := json.Marshal(createNode)

	resp, err := mc.postRequest("/nodes", bytes.NewBuffer(createNodeBytes))

	if err != nil {
		fmt.Printf("Error %s", err)
		return createNodeResponse, err
	}

	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)

	json.Unmarshal(body, &createNodeResponse)

	return createNodeResponse, nil
}

func (mc MarmotcoreClient) GetNode(nodeId string) (NodeResponse, error) {
	var node NodeResponse

	resp, err := mc.getRequest("/nodes/" + nodeId)

	if err != nil {
		fmt.Printf("Error %s", err)
		return node, err
	}

	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)

	json.Unmarshal(body, &node)

	return node, nil
}

func (mc MarmotcoreClient) DeleteNode(nodeId string) (DeleteNodeResponse, error) {
	var deleteNode DeleteNodeResponse

	resp, err := mc.deleteRequest("/nodes/" + nodeId)

	if err != nil {
		fmt.Printf("Error %s", err)
		return deleteNode, err
	}

	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)

	json.Unmarshal(body, &deleteNode)

	return deleteNode, nil
}

func main() {
	mc := &MarmotcoreClient{
		protocol:   "http",
		host:       "localhost",
		port:       "3000",
		apiVersion: "v1",
	}
	key, err := mc.GetKey("chia-1.3.*-testnet-testUserId-rocket-silence-shoot-313")
	if err != nil {
		fmt.Printf("Error %s", err)
		return
	}
	fmt.Print(key)
}

type Key struct {
	UserId string `json:"user_id"`
	NodeId string `json:"node_id"`
	Key    string `json:"key"`
	Cert   string `json:"cert"`
}

type KeyResponse struct {
	Key Key `json:"key"`
}

func (mc MarmotcoreClient) GetKey(nodeId string) (KeyResponse, error) {
	var key KeyResponse

	resp, err := mc.getRequest("/keys/" + nodeId)

	if err != nil {
		fmt.Printf("Error %s", err)
		return key, err
	}

	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)

	json.Unmarshal(body, &key)

	return key, nil
}
