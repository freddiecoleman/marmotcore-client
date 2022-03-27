package main

import (
	"bytes"
	"io/ioutil"
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
)

type MockClient struct {
	GetFunc func(url string) (resp *http.Response, err error)
}

var GetGetFunc func(url string) (*http.Response, error)

func (m *MockClient) Get(url string) (*http.Response, error) {
	return GetGetFunc(url)
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

	GetGetFunc = func(url string) (*http.Response, error) {
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

	assert.EqualValues(t, Nodes{
		Nodes: []Node{
			*newNode("testUserId", createdTime, "chia-1.3.*-testnet-testUserId-rest-equally-rabbit-1668", "54.71.136.33", "us-west-2", "node.small", "1.3.*", "testnet", "R", false),
		},
	}, nodes)
}
