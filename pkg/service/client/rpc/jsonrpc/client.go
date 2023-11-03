package jsonrpc

import (
	"bytes"
	"encoding/json"
	"net/http"
)

type RPCClient interface {
	Post(request interface{}) (*http.Response, error)
}

type RPCClientImpl struct {
	client *http.Client
	url    string
}

func NewRPCClient(url string) RPCClient {
	return &RPCClientImpl{
		client: &http.Client{
			Transport: &http.Transport{
				// Set to true to prevent GZIP-bomb DoS attacks
				DisableCompression: true,
			},
		},
		url: url,
	}
}

func (rc *RPCClientImpl) Post(request interface{}) (*http.Response, error) {
	reqBytes, err := json.Marshal(request)
	if err != nil {
		return nil, err
	}
	reqBuf := bytes.NewBuffer(reqBytes)
	httpReq, err := http.NewRequest(http.MethodPost, rc.url, reqBuf)
	if err != nil {
		return nil, err
	}

	httpReq.Header.Set("Content-type", "application/json")
	return rc.client.Do(httpReq)
}
