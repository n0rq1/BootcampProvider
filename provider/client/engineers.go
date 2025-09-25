package client

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
)

// GetEngineers - Returns list of engineers (no auth required)
func (c *Client) GetEngineers() ([]Engineer, error) {
	req, err := http.NewRequest("GET", fmt.Sprintf("%s/engineers", c.endpoint), nil)
	if err != nil {
		return nil, err
	}

	body, err := c.doRequest(req)
	if err != nil {
		return nil, err
	}

	engineers := []Engineer{}
	err = json.Unmarshal(body, &engineers)
	if err != nil {
		return nil, err
	}

	return engineers, nil
}

func (c *Client) CreateEngineer(engineer Engineer) (*Engineer, error) {
	rb, err := json.Marshal(engineer)
	if err != nil {
		return nil, err
	}

	req, err := http.NewRequest("POST", fmt.Sprintf("%s/engineers", c.endpoint), strings.NewReader(string(rb)))
	if err != nil {
		return nil, err
	}

	// Set headers
	req.Header.Set("Content-Type", "application/json")

	body, err := c.doRequest(req)
	if err != nil {
		return nil, err
	}

	newEngineer := Engineer{}
	err = json.Unmarshal(body, &newEngineer)
	if err != nil {
		return nil, err
	}

	return &newEngineer, nil
}

func (c *Client) GetEngineer(engineerID string) (*Engineer, error) {
	req, err := http.NewRequest("GET", fmt.Sprintf("%s/engineers/id/%s", c.endpoint, engineerID), nil)
	if err != nil {
		return nil, err
	}

	body, err := c.doRequest(req)
	if err != nil {
		return nil, err
	}

	engineer := Engineer{}
	if err := json.Unmarshal(body, &engineer); err != nil {
		return nil, err
	}

	return &engineer, nil
}

func (c *Client) UpdateEngineer(engineerID string, engineer Engineer) (*Engineer, error) {
	rb, err := json.Marshal(engineer)
	if err != nil {
		return nil, err
	}

	req, err := http.NewRequest("PUT", fmt.Sprintf("%s/engineers/%s", c.endpoint, engineerID), strings.NewReader(string(rb)))
	if err != nil {
		return nil, err
	}

	body, err := c.doRequest(req)
	if err != nil {
		return nil, err
	}

	updatedEngineer := Engineer{}
	err = json.Unmarshal(body, &updatedEngineer)
	if err != nil {
		return nil, err
	}

	return &updatedEngineer, nil
}

func (c *Client) DeleteEngineer(engineerID string) error {
	req, err := http.NewRequest("DELETE", fmt.Sprintf("%s/engineers/%s", c.endpoint, engineerID), nil)
	if err != nil {
		return err
	}

	_, err = c.doRequest(req)
	return err
}