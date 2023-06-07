package midjourney

import (
	"fmt"
	"io/ioutil"
	"net/http"
)

type Client struct {
	*http.Client

	Config *Config
}

func NewClient(cfg *Config) *Client {
	return &Client{http.DefaultClient, cfg}
}

func checkResponse(resp *http.Response) error {
	if resp.StatusCode >= 400 {
		body, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			return fmt.Errorf("Call ioutil.ReadAll failed, err: %w", err)
		}

		return fmt.Errorf("resp.StatusCode: %d, body: %s", resp.StatusCode, string(body))
	}

	return nil
}
