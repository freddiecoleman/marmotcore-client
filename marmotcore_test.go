package main

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGetNodes(t *testing.T) {
	mc := &MarmotcoreClient{
		protocol:   "http",
		host:       "localhost",
		port:       "3000",
		apiVersion: "1.3.*",
	}
	Client = &MockClient{}

	nodes, err := mc.GetNodes()

	if err != nil {
		t.Errorf("Error: %d", err)
	}

	assert.EqualValues(t, nodes.Nodes[0].ChiaVersion, "111")

}
