package marmotcoreclient

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
		Protocol:   "http",
		Host:       "localhost",
		Port:       "3000",
		ApiVersion: "v1",
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
		Protocol:   "http",
		Host:       "localhost",
		Port:       "3000",
		ApiVersion: "v1",
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
		Protocol:   "http",
		Host:       "localhost",
		Port:       "3000",
		ApiVersion: "v1",
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
		Protocol:   "http",
		Host:       "localhost",
		Port:       "3000",
		ApiVersion: "v1",
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

func newKey(userId string, nodeId string, key string, cert string) *Key {
	return &Key{
		UserId: userId,
		NodeId: nodeId,
		Key:    key,
		Cert:   cert,
	}
}

func TestGetKey(t *testing.T) {
	json := `{"key":{"cert":"LS0tLS1CRUdJTiBDRVJUSUZJQ0FURS0tLS0tCk1JSURMRENDQWhTZ0F3SUJBZ0lVTC9mNlU2eEhCeXZBUlNHcjJWTWNlaEZZSmVrd0RRWUpLb1pJaHZjTkFRRUwKQlFBd1JERU5NQXNHQTFVRUNnd0VRMmhwWVRFUU1BNEdBMVVFQXd3SFEyaHBZU0JEUVRFaE1COEdBMVVFQ3d3WQpUM0puWVc1cFl5QkdZWEp0YVc1bklFUnBkbWx6YVc5dU1DQVhEVEl5TURNeU5qRTJNRGMwTmxvWUR6SXhNREF3Ck9EQXlNREF3TURBd1dqQkJNUTB3Q3dZRFZRUUREQVJEYUdsaE1RMHdDd1lEVlFRS0RBUkRhR2xoTVNFd0h3WUQKVlFRTERCaFBjbWRoYm1saklFWmhjbTFwYm1jZ1JHbDJhWE5wYjI0d2dnRWlNQTBHQ1NxR1NJYjNEUUVCQVFVQQpBNElCRHdBd2dnRUtBb0lCQVFDK1VsbExUcDljOFJKUnA1bktzNDhUb0lDS0MwaTBnNC9PSUV0bVFMNWN2Q3hzCk5tMk16QXkwUnVOM0tiUHdyd2wxM0p6N3hxVXl2bGFEYjIrZzlxTGhYeW9BL0hlVFZxS2V5elVtVUovTEJBMGMKdFFNQjNOcWpkczNOQWs2UTM4VmhLS3VVQytOMU1GTmVYNFo5bXF6WWhvMjFFM0V3Nm14enlTa2pERE93a2t6bgpBWUFsK3ltc24rTi93Wmk4VUswb1EvNTlUYVFtR1p0VmVYalNkMlNYcTljbU1wZitnTVJqSzdxRzE0T1RsSG1pCi9vU0EzR29VdU56V254ZTdzMEhPWnRSSHd1ZHpuZ1BxOHZCejVLbDVUaUVVZ3E4dWYzSTM2WDlNNWFkWndaT2UKb1JwWUhhbFRjYUZIb09WdXEvVE0zSmsva3Ywd3RCQTJrM2ZyVDhrN0FnTUJBQUdqRnpBVk1CTUdBMVVkRVFRTQpNQXFDQ0dOb2FXRXVibVYwTUEwR0NTcUdTSWIzRFFFQkN3VUFBNElCQVFBOUF6R1JSN0M3ZzQ0YTRuclNxcHZUCi9lMTR1cXVTMmxjd1RIWFlaRG0zdzM0Q2daM2FIWlVIU1FFSDdjYlN0MFlIWDJYOVF4VFZaZHROa1VKc1IxKzcKUURQWUt2TWJnNjVRc2R0aVQrbURXYUFGbFc1UDhVazRRNzZNQ3JFY2RnK0lqQVRHNS9IZ3NHckJJdHdGU2x6UgpTTHRuVHNodm1jTFVrY01oZHJ4UFdXa21IelRqdE5KVlpBeStSV1dDNG1vV1ZNdTFwSGg5TmtEM1dZS0gwQWEvCnNsZ2hZMkowVTVCQnQwcUM5YzJ4NlNqTWNSeWE1RkVGZisxMTRMZ3ZRWHYwSys2Y3NrbDlZaEJqQ3FWR2FuaEsKSjk3RUFGMmNFNnFKRDdIaUZodVM3enlSZ3pyeXVYemRKU1B6YTFJbS9idDJxc0t6dHV4VFFCaG00RE52UzlPagotLS0tLUVORCBDRVJUSUZJQ0FURS0tLS0tCg==","user_id":"testUserId","key":"LS0tLS1CRUdJTiBSU0EgUFJJVkFURSBLRVktLS0tLQpNSUlFb3dJQkFBS0NBUUVBdmxKWlMwNmZYUEVTVWFlWnlyT1BFNkNBaWd0SXRJT1B6aUJMWmtDK1hMd3NiRFp0CmpNd010RWJqZHltejhLOEpkZHljKzhhbE1yNVdnMjl2b1BhaTRWOHFBUHgzazFhaW5zczFKbENmeXdRTkhMVUQKQWR6YW8zYk56UUpPa04vRllTaXJsQXZqZFRCVFhsK0dmWnFzMklhTnRSTnhNT3BzYzhrcEl3d3pzSkpNNXdHQQpKZnNwckovamY4R1l2RkN0S0VQK2ZVMmtKaG1iVlhsNDBuZGtsNnZYSmpLWC9vREVZeXU2aHRlRGs1UjVvdjZFCmdOeHFGTGpjMXA4WHU3TkJ6bWJVUjhMbmM1NEQ2dkx3YytTcGVVNGhGSUt2TG45eU4rbC9UT1duV2NHVG5xRWEKV0IycFUzR2hSNkRsYnF2MHpOeVpQNUw5TUxRUU5wTjM2MC9KT3dJREFRQUJBb0lCQUVXSFJBUFU5emMzQXFBKwpBRnNpK2RRTHdLbXRzYVB3cENxRGRjZ25RdVFTQXBDTjJidWtGOGdNVmJFRTFTM0taRVJna2lFUTB1MFV4L3BBClRhR2FPS3JvM3lsNUVoNExlbUZVajFndUFCSmtxbjJnU3piMC9oTFZwaDJOQ0RLNDdSeXZoTzhTNE1mQkhkUE8KUmczQTRnQkFONmk5TDQ3Vk5VV2ZhWjNhS1ZBSm9ROVJzRzlmMkh1cUJWdGE4bnZCQWJBeFIvbEN6UVk3UGRRQwpjUmhjRHFNTmZkUzBoV2F5T2xJTHdoekdSbTFSQ1pZVU9kNVNtSE90Mk1lWHZiK1B1N2dRdWpHSUQ2SmlKcjIzCnpYTmRET3pkWVlVdG81T0RZMGxiUmExWlR0RVJwM1IrTVhCYjlKUXdacUF4a2JXRWk5RGkxelJLM2hUQ3NFelcKWHlNWUd3RUNnWUVBNWlodmdQOGhYVzdsR1VBM1JRUStVYUJHVk1tb2l3U3FEMmYxT3pPdHBleHVveGJjVE9TMgp0UkFjUHREUFYxOEtYSDNuOGV6d1cwaGJhY0xEci80RWRteVBtMzFMOFMyc2FXdDF3VCtrMEl4OWM5bUdPN3doCmNpYU5vTk9EVHpMODVGcFlQSDhsMGdlSmlBV2wzK1plc3JZRnVBYWlGa3hNaWEySjZJWHRwWHNDZ1lFQTA3RGcKT2NjeFRvdVdTOWthVmQ4VlczODc2Yk9IbXUwMTcwcGpoLzBnWFdYaTE1eEN0QndKNTVhU2R2cXg3M25sbHRzLwpheU5UTmZrRGhVZytxUDU4SkJ4YkpLb3hKaHRUWXBsV2NaSWFWY0FURFRqRlpkcmRqYVpvRENWYUZvUk5Jd1d4CnlENGxxTGZwYit6Q1dqU1pIakNaN25nUm5mT3Z1enVvd0IzbXYwRUNnWUI3YStPZmdURWJNWVNaYmQ0MW5IanUKdk12Nlc4bU9Bd3BQQ0tody90MTN4TG52cXlxbjhWNG82bUs3TEs2RFkrdmlmWUlNTWNzU3FGS2MvRnlEMit2NAp0VzZ0S3h1MlVZL0xXRnpsTElQdGNlazBYc21rN3RYZ3FOdjZDbkszM2RmUGZNTWtiZXFTSG9pWjhLMXF5OWFzCmJ2L1NGM3lFQ0paaW5qVCtCQlBVVVFLQmdHckZXS0xydkF2UXpkS2R3dkd5M2hQVEhjWG0vaXQvSDJmOURpeUYKMkhBak5vSG5WNkYrVHVTWEJuS3FTVnJ2RVlUQU9zRndCTVZCUHF4WDN2cmZ1SCtDS2RwWldRYk9XNFZzcjdRQwpxL082T0NIQUU1Z25CdjR1QTJhMDVEWmRMb2JPbFhmWkdLZDdjMnMzY0dPTkNmbTNLN3lpcE1nVkkvcDh3Y1Y4CmoyakJBb0dCQU1zR0U5TisrOHIreHVHMjZtV0w2Vmc4SUEzVm9VWXFXY0xTYnhxSG5CbThjdjIya0ZQMVhYNi8KN1hsT0h1WE1ERjV6U2RWTFFINUU5MjUxbzg2OTQ0MUx2M3RPTmpSSVdRUnRPOWpDdnNJTERNaHl0QitibndvZQp3Q1NBOTJYYytOZVB6L0tVVHdNNVMxcklwbEg0S1N0UTZSUTVIVjh2UG5iRkQzQ2VNWUlKCi0tLS0tRU5EIFJTQSBQUklWQVRFIEtFWS0tLS0tCg==","node_id":"chia-1.3.*-testnet-testUserId-child-attention-actual-1049"}}`
	r := ioutil.NopCloser(bytes.NewReader([]byte(json)))

	GetFunc = func(url string) (*http.Response, error) {
		assert.EqualValues(t, "http://localhost:3000/v1/keys/chia-1.3.*-testnet-testUserId-child-attention-actual-1049", url)

		return &http.Response{
			StatusCode: 200,
			Body:       r,
		}, nil
	}
	mc := &MarmotcoreClient{
		Protocol:   "http",
		Host:       "localhost",
		Port:       "3000",
		ApiVersion: "v1",
	}
	Client = &MockClient{}

	key, err := mc.GetKey("chia-1.3.*-testnet-testUserId-child-attention-actual-1049")

	if err != nil {
		t.Errorf("Error: %d", err)
	}

	assert.EqualValues(t, KeyResponse{
		Key: *newKey("testUserId", "chia-1.3.*-testnet-testUserId-child-attention-actual-1049", "LS0tLS1CRUdJTiBSU0EgUFJJVkFURSBLRVktLS0tLQpNSUlFb3dJQkFBS0NBUUVBdmxKWlMwNmZYUEVTVWFlWnlyT1BFNkNBaWd0SXRJT1B6aUJMWmtDK1hMd3NiRFp0CmpNd010RWJqZHltejhLOEpkZHljKzhhbE1yNVdnMjl2b1BhaTRWOHFBUHgzazFhaW5zczFKbENmeXdRTkhMVUQKQWR6YW8zYk56UUpPa04vRllTaXJsQXZqZFRCVFhsK0dmWnFzMklhTnRSTnhNT3BzYzhrcEl3d3pzSkpNNXdHQQpKZnNwckovamY4R1l2RkN0S0VQK2ZVMmtKaG1iVlhsNDBuZGtsNnZYSmpLWC9vREVZeXU2aHRlRGs1UjVvdjZFCmdOeHFGTGpjMXA4WHU3TkJ6bWJVUjhMbmM1NEQ2dkx3YytTcGVVNGhGSUt2TG45eU4rbC9UT1duV2NHVG5xRWEKV0IycFUzR2hSNkRsYnF2MHpOeVpQNUw5TUxRUU5wTjM2MC9KT3dJREFRQUJBb0lCQUVXSFJBUFU5emMzQXFBKwpBRnNpK2RRTHdLbXRzYVB3cENxRGRjZ25RdVFTQXBDTjJidWtGOGdNVmJFRTFTM0taRVJna2lFUTB1MFV4L3BBClRhR2FPS3JvM3lsNUVoNExlbUZVajFndUFCSmtxbjJnU3piMC9oTFZwaDJOQ0RLNDdSeXZoTzhTNE1mQkhkUE8KUmczQTRnQkFONmk5TDQ3Vk5VV2ZhWjNhS1ZBSm9ROVJzRzlmMkh1cUJWdGE4bnZCQWJBeFIvbEN6UVk3UGRRQwpjUmhjRHFNTmZkUzBoV2F5T2xJTHdoekdSbTFSQ1pZVU9kNVNtSE90Mk1lWHZiK1B1N2dRdWpHSUQ2SmlKcjIzCnpYTmRET3pkWVlVdG81T0RZMGxiUmExWlR0RVJwM1IrTVhCYjlKUXdacUF4a2JXRWk5RGkxelJLM2hUQ3NFelcKWHlNWUd3RUNnWUVBNWlodmdQOGhYVzdsR1VBM1JRUStVYUJHVk1tb2l3U3FEMmYxT3pPdHBleHVveGJjVE9TMgp0UkFjUHREUFYxOEtYSDNuOGV6d1cwaGJhY0xEci80RWRteVBtMzFMOFMyc2FXdDF3VCtrMEl4OWM5bUdPN3doCmNpYU5vTk9EVHpMODVGcFlQSDhsMGdlSmlBV2wzK1plc3JZRnVBYWlGa3hNaWEySjZJWHRwWHNDZ1lFQTA3RGcKT2NjeFRvdVdTOWthVmQ4VlczODc2Yk9IbXUwMTcwcGpoLzBnWFdYaTE1eEN0QndKNTVhU2R2cXg3M25sbHRzLwpheU5UTmZrRGhVZytxUDU4SkJ4YkpLb3hKaHRUWXBsV2NaSWFWY0FURFRqRlpkcmRqYVpvRENWYUZvUk5Jd1d4CnlENGxxTGZwYit6Q1dqU1pIakNaN25nUm5mT3Z1enVvd0IzbXYwRUNnWUI3YStPZmdURWJNWVNaYmQ0MW5IanUKdk12Nlc4bU9Bd3BQQ0tody90MTN4TG52cXlxbjhWNG82bUs3TEs2RFkrdmlmWUlNTWNzU3FGS2MvRnlEMit2NAp0VzZ0S3h1MlVZL0xXRnpsTElQdGNlazBYc21rN3RYZ3FOdjZDbkszM2RmUGZNTWtiZXFTSG9pWjhLMXF5OWFzCmJ2L1NGM3lFQ0paaW5qVCtCQlBVVVFLQmdHckZXS0xydkF2UXpkS2R3dkd5M2hQVEhjWG0vaXQvSDJmOURpeUYKMkhBak5vSG5WNkYrVHVTWEJuS3FTVnJ2RVlUQU9zRndCTVZCUHF4WDN2cmZ1SCtDS2RwWldRYk9XNFZzcjdRQwpxL082T0NIQUU1Z25CdjR1QTJhMDVEWmRMb2JPbFhmWkdLZDdjMnMzY0dPTkNmbTNLN3lpcE1nVkkvcDh3Y1Y4CmoyakJBb0dCQU1zR0U5TisrOHIreHVHMjZtV0w2Vmc4SUEzVm9VWXFXY0xTYnhxSG5CbThjdjIya0ZQMVhYNi8KN1hsT0h1WE1ERjV6U2RWTFFINUU5MjUxbzg2OTQ0MUx2M3RPTmpSSVdRUnRPOWpDdnNJTERNaHl0QitibndvZQp3Q1NBOTJYYytOZVB6L0tVVHdNNVMxcklwbEg0S1N0UTZSUTVIVjh2UG5iRkQzQ2VNWUlKCi0tLS0tRU5EIFJTQSBQUklWQVRFIEtFWS0tLS0tCg==", "LS0tLS1CRUdJTiBDRVJUSUZJQ0FURS0tLS0tCk1JSURMRENDQWhTZ0F3SUJBZ0lVTC9mNlU2eEhCeXZBUlNHcjJWTWNlaEZZSmVrd0RRWUpLb1pJaHZjTkFRRUwKQlFBd1JERU5NQXNHQTFVRUNnd0VRMmhwWVRFUU1BNEdBMVVFQXd3SFEyaHBZU0JEUVRFaE1COEdBMVVFQ3d3WQpUM0puWVc1cFl5QkdZWEp0YVc1bklFUnBkbWx6YVc5dU1DQVhEVEl5TURNeU5qRTJNRGMwTmxvWUR6SXhNREF3Ck9EQXlNREF3TURBd1dqQkJNUTB3Q3dZRFZRUUREQVJEYUdsaE1RMHdDd1lEVlFRS0RBUkRhR2xoTVNFd0h3WUQKVlFRTERCaFBjbWRoYm1saklFWmhjbTFwYm1jZ1JHbDJhWE5wYjI0d2dnRWlNQTBHQ1NxR1NJYjNEUUVCQVFVQQpBNElCRHdBd2dnRUtBb0lCQVFDK1VsbExUcDljOFJKUnA1bktzNDhUb0lDS0MwaTBnNC9PSUV0bVFMNWN2Q3hzCk5tMk16QXkwUnVOM0tiUHdyd2wxM0p6N3hxVXl2bGFEYjIrZzlxTGhYeW9BL0hlVFZxS2V5elVtVUovTEJBMGMKdFFNQjNOcWpkczNOQWs2UTM4VmhLS3VVQytOMU1GTmVYNFo5bXF6WWhvMjFFM0V3Nm14enlTa2pERE93a2t6bgpBWUFsK3ltc24rTi93Wmk4VUswb1EvNTlUYVFtR1p0VmVYalNkMlNYcTljbU1wZitnTVJqSzdxRzE0T1RsSG1pCi9vU0EzR29VdU56V254ZTdzMEhPWnRSSHd1ZHpuZ1BxOHZCejVLbDVUaUVVZ3E4dWYzSTM2WDlNNWFkWndaT2UKb1JwWUhhbFRjYUZIb09WdXEvVE0zSmsva3Ywd3RCQTJrM2ZyVDhrN0FnTUJBQUdqRnpBVk1CTUdBMVVkRVFRTQpNQXFDQ0dOb2FXRXVibVYwTUEwR0NTcUdTSWIzRFFFQkN3VUFBNElCQVFBOUF6R1JSN0M3ZzQ0YTRuclNxcHZUCi9lMTR1cXVTMmxjd1RIWFlaRG0zdzM0Q2daM2FIWlVIU1FFSDdjYlN0MFlIWDJYOVF4VFZaZHROa1VKc1IxKzcKUURQWUt2TWJnNjVRc2R0aVQrbURXYUFGbFc1UDhVazRRNzZNQ3JFY2RnK0lqQVRHNS9IZ3NHckJJdHdGU2x6UgpTTHRuVHNodm1jTFVrY01oZHJ4UFdXa21IelRqdE5KVlpBeStSV1dDNG1vV1ZNdTFwSGg5TmtEM1dZS0gwQWEvCnNsZ2hZMkowVTVCQnQwcUM5YzJ4NlNqTWNSeWE1RkVGZisxMTRMZ3ZRWHYwSys2Y3NrbDlZaEJqQ3FWR2FuaEsKSjk3RUFGMmNFNnFKRDdIaUZodVM3enlSZ3pyeXVYemRKU1B6YTFJbS9idDJxc0t6dHV4VFFCaG00RE52UzlPagotLS0tLUVORCBDRVJUSUZJQ0FURS0tLS0tCg=="),
	}, key)
}

func TestGetKeys(t *testing.T) {
	json := `{"keys":[{"cert":"LS0tLS1CRUdJTiBDRVJUSUZJQ0FURS0tLS0tCk1JSURMRENDQWhTZ0F3SUJBZ0lVTC9mNlU2eEhCeXZBUlNHcjJWTWNlaEZZSmVrd0RRWUpLb1pJaHZjTkFRRUwKQlFBd1JERU5NQXNHQTFVRUNnd0VRMmhwWVRFUU1BNEdBMVVFQXd3SFEyaHBZU0JEUVRFaE1COEdBMVVFQ3d3WQpUM0puWVc1cFl5QkdZWEp0YVc1bklFUnBkbWx6YVc5dU1DQVhEVEl5TURNeU5qRTJNRGMwTmxvWUR6SXhNREF3Ck9EQXlNREF3TURBd1dqQkJNUTB3Q3dZRFZRUUREQVJEYUdsaE1RMHdDd1lEVlFRS0RBUkRhR2xoTVNFd0h3WUQKVlFRTERCaFBjbWRoYm1saklFWmhjbTFwYm1jZ1JHbDJhWE5wYjI0d2dnRWlNQTBHQ1NxR1NJYjNEUUVCQVFVQQpBNElCRHdBd2dnRUtBb0lCQVFDK1VsbExUcDljOFJKUnA1bktzNDhUb0lDS0MwaTBnNC9PSUV0bVFMNWN2Q3hzCk5tMk16QXkwUnVOM0tiUHdyd2wxM0p6N3hxVXl2bGFEYjIrZzlxTGhYeW9BL0hlVFZxS2V5elVtVUovTEJBMGMKdFFNQjNOcWpkczNOQWs2UTM4VmhLS3VVQytOMU1GTmVYNFo5bXF6WWhvMjFFM0V3Nm14enlTa2pERE93a2t6bgpBWUFsK3ltc24rTi93Wmk4VUswb1EvNTlUYVFtR1p0VmVYalNkMlNYcTljbU1wZitnTVJqSzdxRzE0T1RsSG1pCi9vU0EzR29VdU56V254ZTdzMEhPWnRSSHd1ZHpuZ1BxOHZCejVLbDVUaUVVZ3E4dWYzSTM2WDlNNWFkWndaT2UKb1JwWUhhbFRjYUZIb09WdXEvVE0zSmsva3Ywd3RCQTJrM2ZyVDhrN0FnTUJBQUdqRnpBVk1CTUdBMVVkRVFRTQpNQXFDQ0dOb2FXRXVibVYwTUEwR0NTcUdTSWIzRFFFQkN3VUFBNElCQVFBOUF6R1JSN0M3ZzQ0YTRuclNxcHZUCi9lMTR1cXVTMmxjd1RIWFlaRG0zdzM0Q2daM2FIWlVIU1FFSDdjYlN0MFlIWDJYOVF4VFZaZHROa1VKc1IxKzcKUURQWUt2TWJnNjVRc2R0aVQrbURXYUFGbFc1UDhVazRRNzZNQ3JFY2RnK0lqQVRHNS9IZ3NHckJJdHdGU2x6UgpTTHRuVHNodm1jTFVrY01oZHJ4UFdXa21IelRqdE5KVlpBeStSV1dDNG1vV1ZNdTFwSGg5TmtEM1dZS0gwQWEvCnNsZ2hZMkowVTVCQnQwcUM5YzJ4NlNqTWNSeWE1RkVGZisxMTRMZ3ZRWHYwSys2Y3NrbDlZaEJqQ3FWR2FuaEsKSjk3RUFGMmNFNnFKRDdIaUZodVM3enlSZ3pyeXVYemRKU1B6YTFJbS9idDJxc0t6dHV4VFFCaG00RE52UzlPagotLS0tLUVORCBDRVJUSUZJQ0FURS0tLS0tCg==","user_id":"testUserId","key":"LS0tLS1CRUdJTiBSU0EgUFJJVkFURSBLRVktLS0tLQpNSUlFb3dJQkFBS0NBUUVBdmxKWlMwNmZYUEVTVWFlWnlyT1BFNkNBaWd0SXRJT1B6aUJMWmtDK1hMd3NiRFp0CmpNd010RWJqZHltejhLOEpkZHljKzhhbE1yNVdnMjl2b1BhaTRWOHFBUHgzazFhaW5zczFKbENmeXdRTkhMVUQKQWR6YW8zYk56UUpPa04vRllTaXJsQXZqZFRCVFhsK0dmWnFzMklhTnRSTnhNT3BzYzhrcEl3d3pzSkpNNXdHQQpKZnNwckovamY4R1l2RkN0S0VQK2ZVMmtKaG1iVlhsNDBuZGtsNnZYSmpLWC9vREVZeXU2aHRlRGs1UjVvdjZFCmdOeHFGTGpjMXA4WHU3TkJ6bWJVUjhMbmM1NEQ2dkx3YytTcGVVNGhGSUt2TG45eU4rbC9UT1duV2NHVG5xRWEKV0IycFUzR2hSNkRsYnF2MHpOeVpQNUw5TUxRUU5wTjM2MC9KT3dJREFRQUJBb0lCQUVXSFJBUFU5emMzQXFBKwpBRnNpK2RRTHdLbXRzYVB3cENxRGRjZ25RdVFTQXBDTjJidWtGOGdNVmJFRTFTM0taRVJna2lFUTB1MFV4L3BBClRhR2FPS3JvM3lsNUVoNExlbUZVajFndUFCSmtxbjJnU3piMC9oTFZwaDJOQ0RLNDdSeXZoTzhTNE1mQkhkUE8KUmczQTRnQkFONmk5TDQ3Vk5VV2ZhWjNhS1ZBSm9ROVJzRzlmMkh1cUJWdGE4bnZCQWJBeFIvbEN6UVk3UGRRQwpjUmhjRHFNTmZkUzBoV2F5T2xJTHdoekdSbTFSQ1pZVU9kNVNtSE90Mk1lWHZiK1B1N2dRdWpHSUQ2SmlKcjIzCnpYTmRET3pkWVlVdG81T0RZMGxiUmExWlR0RVJwM1IrTVhCYjlKUXdacUF4a2JXRWk5RGkxelJLM2hUQ3NFelcKWHlNWUd3RUNnWUVBNWlodmdQOGhYVzdsR1VBM1JRUStVYUJHVk1tb2l3U3FEMmYxT3pPdHBleHVveGJjVE9TMgp0UkFjUHREUFYxOEtYSDNuOGV6d1cwaGJhY0xEci80RWRteVBtMzFMOFMyc2FXdDF3VCtrMEl4OWM5bUdPN3doCmNpYU5vTk9EVHpMODVGcFlQSDhsMGdlSmlBV2wzK1plc3JZRnVBYWlGa3hNaWEySjZJWHRwWHNDZ1lFQTA3RGcKT2NjeFRvdVdTOWthVmQ4VlczODc2Yk9IbXUwMTcwcGpoLzBnWFdYaTE1eEN0QndKNTVhU2R2cXg3M25sbHRzLwpheU5UTmZrRGhVZytxUDU4SkJ4YkpLb3hKaHRUWXBsV2NaSWFWY0FURFRqRlpkcmRqYVpvRENWYUZvUk5Jd1d4CnlENGxxTGZwYit6Q1dqU1pIakNaN25nUm5mT3Z1enVvd0IzbXYwRUNnWUI3YStPZmdURWJNWVNaYmQ0MW5IanUKdk12Nlc4bU9Bd3BQQ0tody90MTN4TG52cXlxbjhWNG82bUs3TEs2RFkrdmlmWUlNTWNzU3FGS2MvRnlEMit2NAp0VzZ0S3h1MlVZL0xXRnpsTElQdGNlazBYc21rN3RYZ3FOdjZDbkszM2RmUGZNTWtiZXFTSG9pWjhLMXF5OWFzCmJ2L1NGM3lFQ0paaW5qVCtCQlBVVVFLQmdHckZXS0xydkF2UXpkS2R3dkd5M2hQVEhjWG0vaXQvSDJmOURpeUYKMkhBak5vSG5WNkYrVHVTWEJuS3FTVnJ2RVlUQU9zRndCTVZCUHF4WDN2cmZ1SCtDS2RwWldRYk9XNFZzcjdRQwpxL082T0NIQUU1Z25CdjR1QTJhMDVEWmRMb2JPbFhmWkdLZDdjMnMzY0dPTkNmbTNLN3lpcE1nVkkvcDh3Y1Y4CmoyakJBb0dCQU1zR0U5TisrOHIreHVHMjZtV0w2Vmc4SUEzVm9VWXFXY0xTYnhxSG5CbThjdjIya0ZQMVhYNi8KN1hsT0h1WE1ERjV6U2RWTFFINUU5MjUxbzg2OTQ0MUx2M3RPTmpSSVdRUnRPOWpDdnNJTERNaHl0QitibndvZQp3Q1NBOTJYYytOZVB6L0tVVHdNNVMxcklwbEg0S1N0UTZSUTVIVjh2UG5iRkQzQ2VNWUlKCi0tLS0tRU5EIFJTQSBQUklWQVRFIEtFWS0tLS0tCg==","node_id":"chia-1.3.*-testnet-testUserId-child-attention-actual-1049"}]}`
	r := ioutil.NopCloser(bytes.NewReader([]byte(json)))

	GetFunc = func(url string) (*http.Response, error) {
		assert.EqualValues(t, "http://localhost:3000/v1/keys", url)

		return &http.Response{
			StatusCode: 200,
			Body:       r,
		}, nil
	}
	mc := &MarmotcoreClient{
		Protocol:   "http",
		Host:       "localhost",
		Port:       "3000",
		ApiVersion: "v1",
	}
	Client = &MockClient{}

	keys, err := mc.GetKeys()

	if err != nil {
		t.Errorf("Error: %d", err)
	}

	assert.EqualValues(t, KeysResponse{
		Keys: []Key{
			*newKey("testUserId", "chia-1.3.*-testnet-testUserId-child-attention-actual-1049", "LS0tLS1CRUdJTiBSU0EgUFJJVkFURSBLRVktLS0tLQpNSUlFb3dJQkFBS0NBUUVBdmxKWlMwNmZYUEVTVWFlWnlyT1BFNkNBaWd0SXRJT1B6aUJMWmtDK1hMd3NiRFp0CmpNd010RWJqZHltejhLOEpkZHljKzhhbE1yNVdnMjl2b1BhaTRWOHFBUHgzazFhaW5zczFKbENmeXdRTkhMVUQKQWR6YW8zYk56UUpPa04vRllTaXJsQXZqZFRCVFhsK0dmWnFzMklhTnRSTnhNT3BzYzhrcEl3d3pzSkpNNXdHQQpKZnNwckovamY4R1l2RkN0S0VQK2ZVMmtKaG1iVlhsNDBuZGtsNnZYSmpLWC9vREVZeXU2aHRlRGs1UjVvdjZFCmdOeHFGTGpjMXA4WHU3TkJ6bWJVUjhMbmM1NEQ2dkx3YytTcGVVNGhGSUt2TG45eU4rbC9UT1duV2NHVG5xRWEKV0IycFUzR2hSNkRsYnF2MHpOeVpQNUw5TUxRUU5wTjM2MC9KT3dJREFRQUJBb0lCQUVXSFJBUFU5emMzQXFBKwpBRnNpK2RRTHdLbXRzYVB3cENxRGRjZ25RdVFTQXBDTjJidWtGOGdNVmJFRTFTM0taRVJna2lFUTB1MFV4L3BBClRhR2FPS3JvM3lsNUVoNExlbUZVajFndUFCSmtxbjJnU3piMC9oTFZwaDJOQ0RLNDdSeXZoTzhTNE1mQkhkUE8KUmczQTRnQkFONmk5TDQ3Vk5VV2ZhWjNhS1ZBSm9ROVJzRzlmMkh1cUJWdGE4bnZCQWJBeFIvbEN6UVk3UGRRQwpjUmhjRHFNTmZkUzBoV2F5T2xJTHdoekdSbTFSQ1pZVU9kNVNtSE90Mk1lWHZiK1B1N2dRdWpHSUQ2SmlKcjIzCnpYTmRET3pkWVlVdG81T0RZMGxiUmExWlR0RVJwM1IrTVhCYjlKUXdacUF4a2JXRWk5RGkxelJLM2hUQ3NFelcKWHlNWUd3RUNnWUVBNWlodmdQOGhYVzdsR1VBM1JRUStVYUJHVk1tb2l3U3FEMmYxT3pPdHBleHVveGJjVE9TMgp0UkFjUHREUFYxOEtYSDNuOGV6d1cwaGJhY0xEci80RWRteVBtMzFMOFMyc2FXdDF3VCtrMEl4OWM5bUdPN3doCmNpYU5vTk9EVHpMODVGcFlQSDhsMGdlSmlBV2wzK1plc3JZRnVBYWlGa3hNaWEySjZJWHRwWHNDZ1lFQTA3RGcKT2NjeFRvdVdTOWthVmQ4VlczODc2Yk9IbXUwMTcwcGpoLzBnWFdYaTE1eEN0QndKNTVhU2R2cXg3M25sbHRzLwpheU5UTmZrRGhVZytxUDU4SkJ4YkpLb3hKaHRUWXBsV2NaSWFWY0FURFRqRlpkcmRqYVpvRENWYUZvUk5Jd1d4CnlENGxxTGZwYit6Q1dqU1pIakNaN25nUm5mT3Z1enVvd0IzbXYwRUNnWUI3YStPZmdURWJNWVNaYmQ0MW5IanUKdk12Nlc4bU9Bd3BQQ0tody90MTN4TG52cXlxbjhWNG82bUs3TEs2RFkrdmlmWUlNTWNzU3FGS2MvRnlEMit2NAp0VzZ0S3h1MlVZL0xXRnpsTElQdGNlazBYc21rN3RYZ3FOdjZDbkszM2RmUGZNTWtiZXFTSG9pWjhLMXF5OWFzCmJ2L1NGM3lFQ0paaW5qVCtCQlBVVVFLQmdHckZXS0xydkF2UXpkS2R3dkd5M2hQVEhjWG0vaXQvSDJmOURpeUYKMkhBak5vSG5WNkYrVHVTWEJuS3FTVnJ2RVlUQU9zRndCTVZCUHF4WDN2cmZ1SCtDS2RwWldRYk9XNFZzcjdRQwpxL082T0NIQUU1Z25CdjR1QTJhMDVEWmRMb2JPbFhmWkdLZDdjMnMzY0dPTkNmbTNLN3lpcE1nVkkvcDh3Y1Y4CmoyakJBb0dCQU1zR0U5TisrOHIreHVHMjZtV0w2Vmc4SUEzVm9VWXFXY0xTYnhxSG5CbThjdjIya0ZQMVhYNi8KN1hsT0h1WE1ERjV6U2RWTFFINUU5MjUxbzg2OTQ0MUx2M3RPTmpSSVdRUnRPOWpDdnNJTERNaHl0QitibndvZQp3Q1NBOTJYYytOZVB6L0tVVHdNNVMxcklwbEg0S1N0UTZSUTVIVjh2UG5iRkQzQ2VNWUlKCi0tLS0tRU5EIFJTQSBQUklWQVRFIEtFWS0tLS0tCg==", "LS0tLS1CRUdJTiBDRVJUSUZJQ0FURS0tLS0tCk1JSURMRENDQWhTZ0F3SUJBZ0lVTC9mNlU2eEhCeXZBUlNHcjJWTWNlaEZZSmVrd0RRWUpLb1pJaHZjTkFRRUwKQlFBd1JERU5NQXNHQTFVRUNnd0VRMmhwWVRFUU1BNEdBMVVFQXd3SFEyaHBZU0JEUVRFaE1COEdBMVVFQ3d3WQpUM0puWVc1cFl5QkdZWEp0YVc1bklFUnBkbWx6YVc5dU1DQVhEVEl5TURNeU5qRTJNRGMwTmxvWUR6SXhNREF3Ck9EQXlNREF3TURBd1dqQkJNUTB3Q3dZRFZRUUREQVJEYUdsaE1RMHdDd1lEVlFRS0RBUkRhR2xoTVNFd0h3WUQKVlFRTERCaFBjbWRoYm1saklFWmhjbTFwYm1jZ1JHbDJhWE5wYjI0d2dnRWlNQTBHQ1NxR1NJYjNEUUVCQVFVQQpBNElCRHdBd2dnRUtBb0lCQVFDK1VsbExUcDljOFJKUnA1bktzNDhUb0lDS0MwaTBnNC9PSUV0bVFMNWN2Q3hzCk5tMk16QXkwUnVOM0tiUHdyd2wxM0p6N3hxVXl2bGFEYjIrZzlxTGhYeW9BL0hlVFZxS2V5elVtVUovTEJBMGMKdFFNQjNOcWpkczNOQWs2UTM4VmhLS3VVQytOMU1GTmVYNFo5bXF6WWhvMjFFM0V3Nm14enlTa2pERE93a2t6bgpBWUFsK3ltc24rTi93Wmk4VUswb1EvNTlUYVFtR1p0VmVYalNkMlNYcTljbU1wZitnTVJqSzdxRzE0T1RsSG1pCi9vU0EzR29VdU56V254ZTdzMEhPWnRSSHd1ZHpuZ1BxOHZCejVLbDVUaUVVZ3E4dWYzSTM2WDlNNWFkWndaT2UKb1JwWUhhbFRjYUZIb09WdXEvVE0zSmsva3Ywd3RCQTJrM2ZyVDhrN0FnTUJBQUdqRnpBVk1CTUdBMVVkRVFRTQpNQXFDQ0dOb2FXRXVibVYwTUEwR0NTcUdTSWIzRFFFQkN3VUFBNElCQVFBOUF6R1JSN0M3ZzQ0YTRuclNxcHZUCi9lMTR1cXVTMmxjd1RIWFlaRG0zdzM0Q2daM2FIWlVIU1FFSDdjYlN0MFlIWDJYOVF4VFZaZHROa1VKc1IxKzcKUURQWUt2TWJnNjVRc2R0aVQrbURXYUFGbFc1UDhVazRRNzZNQ3JFY2RnK0lqQVRHNS9IZ3NHckJJdHdGU2x6UgpTTHRuVHNodm1jTFVrY01oZHJ4UFdXa21IelRqdE5KVlpBeStSV1dDNG1vV1ZNdTFwSGg5TmtEM1dZS0gwQWEvCnNsZ2hZMkowVTVCQnQwcUM5YzJ4NlNqTWNSeWE1RkVGZisxMTRMZ3ZRWHYwSys2Y3NrbDlZaEJqQ3FWR2FuaEsKSjk3RUFGMmNFNnFKRDdIaUZodVM3enlSZ3pyeXVYemRKU1B6YTFJbS9idDJxc0t6dHV4VFFCaG00RE52UzlPagotLS0tLUVORCBDRVJUSUZJQ0FURS0tLS0tCg=="),
		},
	}, keys)
}
