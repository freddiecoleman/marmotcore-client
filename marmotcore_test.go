package main

import (
	"bytes"
	"io"
	"io/ioutil"
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
)

type MockClient struct {
	GetFunc  func(url string) (resp *http.Response, err error)
	DoFunc   func(url string) (*http.Response, error)
	PostFunc func(url string, contentType string, body io.Reader) (resp *http.Response, err error)
}

var GetFunc func(url string) (*http.Response, error)
var DoFunc func(req *http.Request) (*http.Response, error)
var PostFunc func(url string, contentType string, body io.Reader) (resp *http.Response, err error)

func (m *MockClient) Get(url string) (*http.Response, error) {
	return GetFunc(url)
}

func (m *MockClient) Do(req *http.Request) (*http.Response, error) {
	return DoFunc(req)
}

func (m *MockClient) Post(url string, contentType string, body io.Reader) (resp *http.Response, err error) {
	return PostFunc(url, contentType, body)
}

func newNode(userId string, createdTime int64, nodeId string, publicIp string, region string, instanceType string, chiaVersion string, network string, state string, deleted bool) *Node {
	return &Node{
		UserId:       userId,
		CreatedTime:  createdTime,
		NodeId:       nodeId,
		PublicIp:     publicIp,
		Region:       region,
		InstanceType: instanceType,
		ChiaVersion:  chiaVersion,
		Network:      network,
		State:        state,
		Deleted:      deleted,
	}
}

func TestGetNodes(t *testing.T) {
	json := `{"nodes":[{"user_id":"testUserId","deleted":false,"instance_type":"node.small","chia_version":"1.3.*","region":"us-west-2","public_ip":"54.71.136.33","created_time":1648394251715,"node_id":"chia-1.3.*-testnet-testUserId-rest-equally-rabbit-1668","state":"R","network":"testnet"}]}`
	r := ioutil.NopCloser(bytes.NewReader([]byte(json)))

	GetFunc = func(url string) (*http.Response, error) {
		assert.EqualValues(t, "http://localhost:3000/v1/nodes", url)

		return &http.Response{
			StatusCode: 200,
			Body:       r,
		}, nil
	}
	mc := &MarmotcoreClient{
		protocol:   "http",
		host:       "localhost",
		port:       "3000",
		apiVersion: "v1",
	}
	Client = &MockClient{}

	nodes, err := mc.GetNodes()

	if err != nil {
		t.Errorf("Error: %d", err)
	}

	var createdTime int64
	createdTime = 1648394251715

	assert.EqualValues(t, NodesResponse{
		Nodes: []Node{
			*newNode("testUserId", createdTime, "chia-1.3.*-testnet-testUserId-rest-equally-rabbit-1668", "54.71.136.33", "us-west-2", "node.small", "1.3.*", "testnet", "R", false),
		},
	}, nodes)
}

func TestGetNode(t *testing.T) {
	json := `{"node":{"user_id":"testUserId","deleted":false,"instance_type":"node.small","chia_version":"1.3.*","region":"us-west-2","public_ip":"54.71.136.33","created_time":1648394251715,"node_id":"chia-1.3.*-testnet-testUserId-rest-equally-rabbit-1668","state":"R","network":"testnet"}}`
	r := ioutil.NopCloser(bytes.NewReader([]byte(json)))

	GetFunc = func(url string) (*http.Response, error) {
		assert.EqualValues(t, "http://localhost:3000/v1/nodes/chia-1.3.*-testnet-testUserId-rest-equally-rabbit-1668", url)

		return &http.Response{
			StatusCode: 200,
			Body:       r,
		}, nil
	}
	mc := &MarmotcoreClient{
		protocol:   "http",
		host:       "localhost",
		port:       "3000",
		apiVersion: "v1",
	}
	Client = &MockClient{}

	node, err := mc.GetNode("chia-1.3.*-testnet-testUserId-rest-equally-rabbit-1668")

	if err != nil {
		t.Errorf("Error: %d", err)
	}

	var createdTime int64
	createdTime = 1648394251715

	assert.EqualValues(t, NodeResponse{
		Node: *newNode("testUserId", createdTime, "chia-1.3.*-testnet-testUserId-rest-equally-rabbit-1668", "54.71.136.33", "us-west-2", "node.small", "1.3.*", "testnet", "R", false),
	}, node)
}

func TestDeleteNode(t *testing.T) {
	json := `{"deleted":true}`
	r := ioutil.NopCloser(bytes.NewReader([]byte(json)))

	DoFunc = func(req *http.Request) (*http.Response, error) {
		assert.EqualValues(t, "http://localhost:3000/v1/nodes/chia-1.3.*-testnet-testUserId-rest-equally-rabbit-1668", req.URL.String())

		return &http.Response{
			StatusCode: 200,
			Body:       r,
		}, nil
	}
	mc := &MarmotcoreClient{
		protocol:   "http",
		host:       "localhost",
		port:       "3000",
		apiVersion: "v1",
	}
	Client = &MockClient{}

	node, err := mc.DeleteNode("chia-1.3.*-testnet-testUserId-rest-equally-rabbit-1668")

	if err != nil {
		t.Errorf("Error: %d", err)
	}

	assert.EqualValues(t, DeleteNodeResponse{
		Deleted: true,
	}, node)
}

func TestCreateNode(t *testing.T) {
	json := `{"node_id": "chia-1.3.*-testnet-testUserId-rest-equally-rabbit-1668"}`
	r := ioutil.NopCloser(bytes.NewReader([]byte(json)))

	PostFunc = func(url string, contentType string, body io.Reader) (resp *http.Response, err error) {
		assert.EqualValues(t, "http://localhost:3000/v1/nodes", url)

		return &http.Response{
			StatusCode: 200,
			Body:       r,
		}, nil
	}
	mc := &MarmotcoreClient{
		protocol:   "http",
		host:       "localhost",
		port:       "3000",
		apiVersion: "v1",
	}
	Client = &MockClient{}

	createNode := &CreateNode{
		Region:       "us-west-2",
		InstanceType: "node.small",
		ChiaVersion:  "1.3.*",
		Network:      "testnet",
	}
	node, err := mc.CreateNode(createNode)

	if err != nil {
		t.Errorf("Error: %d", err)
	}

	assert.EqualValues(t, CreateNodeResponse{
		NodeId: "chia-1.3.*-testnet-testUserId-rest-equally-rabbit-1668",
	}, node)
}
