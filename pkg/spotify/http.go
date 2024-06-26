package spotify

import (
	"context"
	"fmt"
	"net/http"

	jsoniter "github.com/json-iterator/go"
	"github.com/rs/zerolog/log"
)

const BaseURL = "https://api.spotify.com/v1"

type Doer interface {
	Do(req *http.Request) (*http.Response, error)
}

type Client struct {
	Doer
}

func New(doer Doer) *Client {
	return &Client{Doer: doer}
}

type ErrorResponse struct {
	Status  int    `json:"status"`
	Message string `json:"message"`
}

func (e *ErrorResponse) Error() string {
	return fmt.Sprintf("unexpected status code: %d %s", e.Status, e.Message)
}

func (c *Client) GetJSON(ctx context.Context, method, path string, response interface{}) error {
	req, err := http.NewRequestWithContext(ctx, method, endpoint(path), nil)
	if err != nil {
		return err
	}

	resp, err := c.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		errorResponse := &ErrorResponse{
			Status:  resp.StatusCode,
			Message: "unexpected status code",
		}
		// Attempt to decode the error response, but don't worry if it fails.
		err := jsoniter.NewDecoder(resp.Body).Decode(errorResponse)
		if err != nil {
			log.Err(err).Msg("failed to decode Spotify error response")
		}
		return errorResponse
	}

	return jsoniter.NewDecoder(resp.Body).Decode(response)
}

func endpoint(path string) string {
	return BaseURL + path
}
