package draftout

import (
	"context"
	"encoding/json"
	"net/http"
	"net/url"
	"time"
)

type Client struct {
	client  *http.Client
	baseURL string
}

func New() *Client {
	client := &http.Client{
		Timeout: time.Second * 10,
	}
	return &Client{
		client:  client,
		baseURL: "https://draftoutmc.com",
	}
}

func (c *Client) get(ctx context.Context, path string, params url.Values, target any) error {
	u, err := url.Parse(c.baseURL + path)
	if err != nil {
		return err
	}
	if params != nil {
		u.RawQuery = params.Encode()
	}
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, u.String(), nil)
	req.Header.Set("Accept", "application/json")
	if err != nil {
		return err
	}
	res, err := c.client.Do(req)
	if err != nil {
		return err
	}
	defer res.Body.Close()

	err = json.NewDecoder(res.Body).Decode(target)
	if err != nil {
		return err
	}

	return nil
}
