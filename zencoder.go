package zencoder

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net"
	"net/http"
	"time"
)

const (
	DEFAULT_ZENCODER_API_ENDPOINT = "https://app.zencoder.com/api/v2/jobs"
	DEFAULT_RESPONSE_TYPE         = "application/json"
)

type Client struct {
	apiKey       string
	apiEndpoint  string
	responseType string
	timeout      int
}

type Options struct {
	ApiKey       string
	ApiEndpoint  string
	ResponseType string
	Timeout      int
}

func NewClient(options *Options) (*Client, error) {
	if options == nil {
		err := fmt.Errorf("error: cannot init Zencoder client without Options")
		return nil, err
	}
	if len(options.ApiKey) == 0 {
		err := fmt.Errorf("error: must supply ApiKey option to init")
		return nil, err
	}
	responseType := DEFAULT_RESPONSE_TYPE
	if len(options.ResponseType) > 0 {
		if options.ResponseType == "application/xml" {
			responseType = "application/xml"
		} else {
			err := fmt.Errorf("error: unsupported response type. response type may be application/json (default) or application/xml")
			return nil, err
		}
	}
	timeout := 30
	if options.Timeout > 0 {
		timeout = options.Timeout
	}
	apiEndpoint := DEFAULT_ZENCODER_API_ENDPOINT
	if len(options.ApiEndpoint) > 0 {
		apiEndpoint = options.ApiEndpoint
	}

	return &Client{
		apiKey:       options.ApiKey,
		apiEndpoint:  apiEndpoint,
		responseType: responseType,
		timeout:      timeout,
	}, nil
}

func (c *Client) Zencode(input string, outputs []map[string]interface{}, notifications []string) (map[string]interface{}, error) {
	outputsStr, err := json.Marshal(outputs)
	if err != nil {
		return nil, err
	}
	notificationsStr, err := json.Marshal(notifications)
	if err != nil {
		return nil, err
	}
	reqStr := ""
	if notifications != nil && len(notifications) > 0 {
		reqStr = fmt.Sprintf("{\"input\":\"%s\",\"output\":%s, \"notifications\":%s}", input, outputsStr, notificationsStr)
	} else {
		reqStr = fmt.Sprintf("{\"input\":\"%s\",\"output\":%s\"}", input, outputsStr)
	}
	fmt.Printf("reqStr: %s", reqStr)
	if req, err := http.NewRequest("POST", c.apiEndpoint,
		bytes.NewBuffer([]byte(reqStr))); err != nil {
		return nil, err
	} else {
		req.Header.Add("Content-Type", c.responseType)
		req.Header.Add("Zencoder-Api-Key", c.apiKey)

		tr := http.Transport{
			Dial: func(network, addr string) (net.Conn, error) {
				return net.DialTimeout(network, addr, time.Duration(c.timeout)*time.Second)
			},
		}

		if res, err := tr.RoundTrip(req); err != nil {
			return nil, err
		} else {
			defer res.Body.Close()

			strResp, _ := ioutil.ReadAll(res.Body)
			if res.StatusCode >= 400 {
				return nil, fmt.Errorf("error: %s", string(strResp))
			}

			var response map[string]interface{}
			json.Unmarshal(strResp, &response)

			return response, nil
		}
	}
}
