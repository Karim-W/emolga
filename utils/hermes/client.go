package hermes

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"time"

	"go.uber.org/fx"
)

type HttpClient struct {
	client       *http.Client
	responseData []byte
	Error        error
	responseCode int
}

func (c HttpClient) Get(endpoint string, headers map[string]string) HttpClient {
	req, err := http.NewRequest("GET", endpoint, nil)

	if err != nil {
		c.Error = err
		return c
	}

	c.constructHeaders(req, headers)

	response, err := c.client.Do(req)
	if err != nil {
		c.Error = err
		return c
	}
	c.responseCode = response.StatusCode

	defer response.Body.Close()
	body, err := ioutil.ReadAll(response.Body)

	if err != nil {
		c.Error = err
		return c
	}
	c.responseData = body

	return c
}

func (c HttpClient) Post(endpoint string, headers map[string]string, body interface{}) HttpClient {
	strBody, err := json.Marshal(body)

	if err != nil {
		c.Error = err
		return c
	}

	req, err := http.NewRequest("POST", endpoint, bytes.NewBuffer(strBody))

	if err != nil {
		c.Error = err
		return c
	}

	if headers == nil {
		headers = make(map[string]string, 1)
	}

	if _, ok := headers["Content-Type"]; !ok {
		headers["Content-Type"] = "application/json"
	}

	c.constructHeaders(req, headers)

	response, err := c.client.Do(req)
	if err != nil {
		c.Error = err
		return c
	}
	c.responseCode = response.StatusCode

	defer response.Body.Close()
	resBody, err := ioutil.ReadAll(response.Body)

	if err != nil {
		c.Error = err
		return c
	}
	c.responseData = resBody
	return c
}

func (c HttpClient) Put(endpoint string, headers map[string]string, body interface{}) HttpClient {
	strBody, err := json.Marshal(body)

	if err != nil {
		c.Error = err
		return c
	}

	req, err := http.NewRequest("PUT", endpoint, bytes.NewBuffer(strBody))

	if err != nil {
		c.Error = err
		return c
	}

	if _, ok := headers["Content-Type"]; !ok {
		headers["Content-Type"] = "application/json"
	}

	c.constructHeaders(req, headers)

	response, err := c.client.Do(req)

	if err != nil {
		c.Error = err
		return c
	}
	c.responseCode = response.StatusCode

	defer response.Body.Close()
	resBody, err := ioutil.ReadAll(response.Body)

	if err != nil {
		c.Error = err
		return c
	}
	c.responseData = resBody
	return c
}

func (c HttpClient) Patch(endpoint string, headers map[string]string, body interface{}) HttpClient {
	strBody, err := json.Marshal(body)

	if err != nil {
		c.Error = err
		return c
	}

	req, err := http.NewRequest("PATCH", endpoint, bytes.NewBuffer(strBody))

	if err != nil {
		c.Error = err
		return c
	}

	if _, ok := headers["Content-Type"]; !ok {
		headers["Content-Type"] = "application/json"
	}

	c.constructHeaders(req, headers)

	response, err := c.client.Do(req)

	if err != nil {
		c.Error = err
		return c
	}
	c.responseCode = response.StatusCode

	defer response.Body.Close()
	resBody, err := ioutil.ReadAll(response.Body)

	if err != nil {
		c.Error = err
		return c
	}

	c.responseData = resBody
	return c
}

func (c HttpClient) Del(endpoint string, headers map[string]string) HttpClient {

	req, err := http.NewRequest("DELETE", endpoint, nil)

	if err != nil {
		c.Error = err
		return c
	}

	if _, ok := headers["Content-Type"]; !ok {
		headers["Content-Type"] = "application/json"
	}

	c.constructHeaders(req, headers)

	response, err := c.client.Do(req)

	if err != nil {
		c.Error = err
		return c
	}
	c.responseCode = response.StatusCode

	defer response.Body.Close()
	resBody, err := ioutil.ReadAll(response.Body)

	if err != nil {
		c.Error = err
		return c
	}

	c.responseData = resBody
	return c
}

func (c HttpClient) Result(res ...interface{}) (int, bool, error) {
	if c.Error == nil {
		var ok = true
		if len(c.responseData) > 0 {
			for i := range res {
				err := json.Unmarshal(c.responseData, &res[i])
				ok = ok || (err == nil)
				c.Error = err
			}
			return c.responseCode, ok, c.Error
		} else if c.responseCode == 204 {
			return c.responseCode, ok, nil
		} else {
			return c.responseCode, ok, fmt.Errorf("no response data")
		}
	} else {
		return c.responseCode, false, c.Error
	}
}

func (c HttpClient) constructHeaders(req *http.Request, headers map[string]string) {
	for key, val := range headers {
		req.Header.Add(key, val)
	}
}

func NewHttpClient() *HttpClient {
	client := HttpClient{}

	client.client = &http.Client{
		Timeout: time.Duration(600) * time.Second,
	}

	return &client
}

var Module = fx.Option(fx.Provide(NewHttpClient))
