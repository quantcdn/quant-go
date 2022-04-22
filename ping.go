package quant

import "encoding/json"

type Ping struct {
	Project string `json:"project"`
}

func (c *Client) Ping() (p Ping, err error) {
	req, err := c.NewRequest("ping", "GET")
	res, err := c.doRequest(req)
	json.Unmarshal(res, &p)
	return
}
