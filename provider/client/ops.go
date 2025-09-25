package client 

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
)

// GetOps - Returns list of ops groups
func (c *Client) GetOps() ([]Ops, error) {
	req, err := http.NewRequest("GET", fmt.Sprintf("%s/op", c.endpoint), nil)
	if err != nil {
		return nil, err
	}

	body, err := c.doRequest(req)
	if err != nil {
		return nil, err
	}

	var ops []Ops
	if err := json.Unmarshal(body, &ops); err != nil {
		return nil, err
	}
	return ops, nil
}

// GetOpsByID - Returns a ops group by ID
func (c *Client) GetOpsByID(opsID string) (*Ops, error) {
	req, err := http.NewRequest("GET", fmt.Sprintf("%s/op/id/%s", c.endpoint, opsID), nil)
	if err != nil {
		return nil, err
	}

	body, err := c.doRequest(req)
	if err != nil {
		return nil, err
	}

	var ops Ops
	if err := json.Unmarshal(body, &ops); err != nil {
		return nil, err
	}
	return &ops, nil
}

func (c *Client) CreateOps(ops Ops) (*Ops, error) {
	b, err := json.Marshal(ops)
	if err != nil {
		return nil, err
	}

	req, err := http.NewRequest("POST", fmt.Sprintf("%s/op", c.endpoint), strings.NewReader(string(b)))
	if err != nil {
		return nil, err
	}
	req.Header.Set("Content-Type", "application/json")

	body, err := c.doRequest(req)
	if err != nil {
		return nil, err
	}

	var created Ops
	if err := json.Unmarshal(body, &created); err != nil {
		return nil, err
	}
	return &created, nil
}

func (c *Client) UpdateOps(opsID string, ops Ops) (*Ops, error) {
	rb, err := json.Marshal(ops)
	if err != nil {
		return nil, err
	}

	req, err := http.NewRequest("PUT", fmt.Sprintf("%s/op/%s", c.endpoint, opsID), strings.NewReader(string(rb)))
	if err != nil {
		return nil, err
	}

	body, err := c.doRequest(req)
	if err != nil {
		return nil, err
	}

	updatedOps := Ops{}
	err = json.Unmarshal(body, &updatedOps)
	if err != nil {
		return nil, err
	}

	return &updatedOps, nil
}

func (c *Client) DeleteOps(opsID string) error {
	req, err := http.NewRequest("DELETE", fmt.Sprintf("%s/op/%s", c.endpoint, opsID), nil)
	if err != nil {
		return err
	}

	_, err = c.doRequest(req)
	return err
}