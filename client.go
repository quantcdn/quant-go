package quant

import (
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"strings"
)

const apiHost = "https://api.quantcdn.io"
const apiBase = "v1"

type Client struct {
	HttpClient *http.Client
	ApiToken   string
	ApiClient  string
	ApiProject string
	Host       string
	Base       string
}

func NewClient(token string, client string, project string) *Client {
	return &Client{
		HttpClient: http.DefaultClient,
		ApiToken:   token,
		ApiClient:  client,
		ApiProject: project,
		Host:       apiHost,
		Base:       apiBase,
	}
}

func (c *Client) NewRequest(path string, method string, buffer io.Reader) (*http.Request, error) {
	url := fmt.Sprintf("%s/%s", c.Host, c.Base)
	if !strings.HasPrefix(path, "/") {
		url = url + "/"
	}
	req, err := http.NewRequest(method, url+path, buffer)

	if err != nil {
		return nil, err
	}
	return req, nil
}

func (c *Client) doRequest(req *http.Request) ([]byte, error) {
	req.Header.Set("User-Agent", "Quant (+http://api.quantcdn.io/tf)")
	req.Header.Set("Quant-Token", c.ApiToken)
	req.Header.Set("Quant-Customer", c.ApiClient)
	req.Header.Set("Quant-Project", c.ApiProject)
	req.Header.Set("Content-Type", "application/json")

	res, err := c.HttpClient.Do(req)

	if err != nil {
		return nil, err
	}

	body, err := ioutil.ReadAll(res.Body)
	defer res.Body.Close()
	if err != nil {
		return nil, err
	}
	return body, nil
}
