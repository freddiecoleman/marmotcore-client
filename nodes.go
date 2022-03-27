package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
)

type Node struct {
	UserId       string `json:"user_id"`
	CreatedTime  int    `json:"created_time"`
	NodeId       string `json:"node_id"`
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
