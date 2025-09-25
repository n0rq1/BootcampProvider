package client

import (
	"fmt"
	"io"
	"net/http"
	"time"
)

type Client struct {
	endpoint    string
	HTTPClient *http.Client
}

// NewClient -
func NewClient(endpoint *string) (*Client, error) {
    c := Client{
        HTTPClient: &http.Client{Timeout: 10 * time.Second},
        endpoint:   "",
    }

    if endpoint != nil {
        c.endpoint = *endpoint
    }

    return &c, nil
}

func (c *Client) doRequest(req *http.Request) ([]byte, error) {
    // Basic request logging to aid debugging
    fmt.Printf("[client] HTTP %s %s\n", req.Method, req.URL.String())
    res, err := c.HTTPClient.Do(req)
    if err != nil {
        return nil, err
    }
    defer res.Body.Close()

	body, err := io.ReadAll(res.Body)
	if err != nil {
		return nil, err
	}

    if !(res.StatusCode/100 == 2){
        // Log non-2xx to aid debugging
        fmt.Printf("[client] HTTP %s %s -> %d, body: %s\n", req.Method, req.URL.String(), res.StatusCode, string(body))
        return nil, fmt.Errorf("status: %d, body: %s", res.StatusCode, body)
    }

    return body, err
}