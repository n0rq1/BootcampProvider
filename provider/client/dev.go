package client

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
	
)

// GetDev - Returns list of dev groups
func (c *Client) GetDev() ([]Dev, error) {
	req, err := http.NewRequest("GET", fmt.Sprintf("%s/dev", c.endpoint), nil)
	if err != nil {
		return nil, err
	}

	body, err := c.doRequest(req)
	if err != nil {
		return nil, err
	}

	var devs []Dev
	if err := json.Unmarshal(body, &devs); err != nil {
		return nil, err
	}
	return devs, nil
}

// GetDevByID - Returns a dev group by ID
func (c *Client) GetDevByID(devID string) (*Dev, error) {
	req, err := http.NewRequest("GET", fmt.Sprintf("%s/dev/id/%s", c.endpoint, devID), nil)
	if err != nil {
		return nil, err
	}

	body, err := c.doRequest(req)
	if err != nil {
		return nil, err
	}

	var dev Dev
	if err := json.Unmarshal(body, &dev); err != nil {
		return nil, err
	}
	return &dev, nil
}

func (c *Client) CreateDev(dev Dev) (*Dev, error) {
	b, err := json.Marshal(dev)
	if err != nil {
		return nil, err
	}

	req, err := http.NewRequest("POST", fmt.Sprintf("%s/dev", c.endpoint), strings.NewReader(string(b)))
	if err != nil {
		return nil, err
	}
	req.Header.Set("Content-Type", "application/json")

	body, err := c.doRequest(req)
	if err != nil {
		return nil, err
	}

	var created Dev
	if err := json.Unmarshal(body, &created); err != nil {
		return nil, err
	}
	return &created, nil
}

func (c *Client) UpdateDev(devID string, dev Dev) (*Dev, error) {
	rb, err := json.Marshal(dev)
	if err != nil {
		return nil, err
	}

	req, err := http.NewRequest("PUT", fmt.Sprintf("%s/dev/%s", c.endpoint, devID), strings.NewReader(string(rb)))
	if err != nil {
		return nil, err
	}

	body, err := c.doRequest(req)
	if err != nil {
		return nil, err
	}

	updatedDev := Dev{}
	err = json.Unmarshal(body, &updatedDev)
	if err != nil {
		return nil, err
	}

	return &updatedDev, nil
}

func (c *Client) DeleteDev(devID string) error {
	req, err := http.NewRequest("DELETE", fmt.Sprintf("%s/dev/%s", c.endpoint, devID), nil)
	if err != nil {
		return err
	}

	_, err = c.doRequest(req)
	return err
}
