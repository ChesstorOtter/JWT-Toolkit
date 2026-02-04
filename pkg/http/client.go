package httpclient

import (
	"crypto/tls"
	"fmt"
	"net/http"
	"net/url"
	"time"
)

type Client struct {
	baseUrl  string
	endpoint string
	proxy    string
	header   map[string]string
	timeout  time.Duration
	insecure bool
}

func NewHTTPClient(baseUrl string, endpoint string, insecure bool, timeout time.Duration) *Client {
	return &Client{
		baseUrl:  baseUrl,
		endpoint: endpoint,
		insecure: insecure,
		timeout:  timeout,
		header:   make(map[string]string),
	}
}

func (c *Client) SetHeader(key, value string) {
	if c.header == nil {
		c.header = make(map[string]string)
	}
	c.header[key] = value
}

func (c *Client) createTransport() *http.Transport {
	tr := &http.Transport{
		TLSClientConfig: &tls.Config{
			InsecureSkipVerify: c.insecure,
		},
	}

	if c.proxy != "" {
		if u, err := url.Parse(c.proxy); err == nil {
			tr.Proxy = http.ProxyURL(u)
		} else {
			fmt.Printf("invalid proxy URL %q: %v\n", c.proxy, err)
		}
	}
	return tr
}

func (c *Client) httpClient() *http.Client {
	return &http.Client{
		Timeout:   c.timeout,
		Transport: c.createTransport(),
	}
}

func (c *Client) DoRequest() (*http.Response, error) {
	client := c.httpClient()

	req, err := http.NewRequest(http.MethodGet, c.baseUrl+c.endpoint, nil)
	if err != nil {
		return nil, err
	}
	for key, value := range c.header {
		req.Header.Set(key, value)
	}
	return client.Do(req)
}

func (c *Client) SetProxy(proxy string) {
	c.proxy = proxy
	if proxy != "" {
		fmt.Println("Proxy set to:", proxy)
	}
}

func (c *Client) SetTimeout(timeout time.Duration) {
	c.timeout = timeout
	if timeout == 0 {
		fmt.Println("Timeout set to: no timeout")
	}
}

func (c *Client) SetInsecure(insecure bool) {
	c.insecure = insecure
	if insecure {
		fmt.Println("Insecure TLS verification enabled")
	}
}

func (c *Client) GetFullUrl() string {
	return c.baseUrl + c.endpoint
}

func (c *Client) GetFullUrlWithQuery(queryParams map[string]string) string {
	fullUrl := c.baseUrl + c.endpoint
	if len(queryParams) > 0 {
		v := url.Values{}
		for key, value := range queryParams {
			v.Set(key, value)
		}
		fullUrl += "?" + v.Encode()
	}
	return fullUrl
}

func (c *Client) TestToken(token string) (*http.Response, error) {
	if token == "" {
		return nil, fmt.Errorf("token is empty")
	}
	client := c.httpClient()

	req, err := http.NewRequest(http.MethodGet, c.baseUrl+c.endpoint, nil)
	if err != nil {
		return nil, err
	}
	for key, value := range c.header {
		req.Header.Set(key, value)
	}
	req.Header.Set("Authorization", "Bearer "+token)

	return client.Do(req)
}

func (c *Client) Get(url string) (*http.Response, error) {
	client := c.httpClient()

	req, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		return nil, err
	}
	for key, value := range c.header {
		req.Header.Set(key, value)
	}

	return client.Do(req)
}

func (c *Client) Post(url string) (*http.Response, error) {
	client := c.httpClient()

	req, err := http.NewRequest(http.MethodPost, url, nil)
	if err != nil {
		return nil, err
	}
	for key, value := range c.header {
		req.Header.Set(key, value)
	}

	return client.Do(req)
}

func (c *Client) Put(url string) (*http.Response, error) {
	client := c.httpClient()

	req, err := http.NewRequest(http.MethodPut, url, nil)
	if err != nil {
		return nil, err
	}
	for key, value := range c.header {
		req.Header.Set(key, value)
	}

	return client.Do(req)
}

func (c *Client) Delete(url string) (*http.Response, error) {
	client := c.httpClient()

	req, err := http.NewRequest(http.MethodDelete, url, nil)
	if err != nil {
		return nil, err
	}
	for key, value := range c.header {
		req.Header.Set(key, value)
	}

	return client.Do(req)
}
