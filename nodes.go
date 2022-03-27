package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"time"
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

var client = http.Client{Timeout: time.Duration(10) * time.Second}

func GetNodes(host string) (error, *Nodes) {
	resp, err := client.Get("http://" + host + ":3000/v1/nodes")
	if err != nil {
		fmt.Printf("Error %s", err)
		return err, nil
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)

	var nodes Nodes

	json.Unmarshal(body, &nodes)

	return nil, &nodes
}

func main() {
	err, nodes := GetNodes("localhost")
	if err != nil {
		fmt.Printf("Error %s", err)
		return
	}
	fmt.Print(nodes.Nodes[0].ChiaVersion)
}
