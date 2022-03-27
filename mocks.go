package main

import "net/http"

type MockClient struct {
	GetFunc func(url string) (resp *http.Response, err error)
}

var GetGetFunc func(url string) (*http.Response, error)

func (m *MockClient) Get(url string) (*http.Response, error) {
	return GetGetFunc(url)
}
