package api_client

import (
	"encoding/json"
	"errors"
	"log"
	"net/http"
	"time"
)

type APIGroupMetrics struct {
	Hash string `json:"hash"`
}

type APIFileMetrics struct {
	Path string `json:"path"`
	Hash string `json:"hash"`
}

// APIMetrics contains all the configuration settings.  Both Groups and Files maps
// are keyed by a string.  In the case of Groups, it should match the name of
// the group to be processed (E.g. wheel).  For Files, the key is a freeform
// shortname that relates to the filename in Path (E.g. sshdcfg for
// /etc/ssh/sshd_config)
type APIMetrics struct {
	Groups map[string]APIGroupMetrics `json:"groups"`
	Files  map[string]APIFileMetrics  `json:"files"`
}

type Client struct {
	BaseURL    string
	HTTPClient *http.Client
}

type errorResponse struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
}

func NewClient(baseURLv1 string) *Client {
	return &Client{
		BaseURL: baseURLv1,
		HTTPClient: &http.Client{
			Timeout: time.Second * 10,
		},
	}
}

func (c *Client) GetMetrics() (*APIMetrics, error) {
	metrics := new(APIMetrics)
	req, err := http.NewRequest(http.MethodGet, c.BaseURL, nil)
	if err != nil {
		return metrics, err
	}
	req.Header.Set("Content-Type", "application/json; charset=utf-8")
	req.Header.Set("Accept", "application/json; charset=utf-8")
	res, err := c.HTTPClient.Do(req)
	if err != nil {
		return metrics, err
	}
	defer res.Body.Close()

	// Handle bad responses
	if res.StatusCode < http.StatusOK || res.StatusCode >= http.StatusBadRequest {
		errRes := new(errorResponse)
		if err = json.NewDecoder(res.Body).Decode(errRes); err == nil {
			return metrics, errors.New(errRes.Message)
		}
		log.Fatalf("unknown error, status code: %d", res.StatusCode)
	}

	if err = json.NewDecoder(res.Body).Decode(metrics); err != nil {
		return metrics, err
	}
	return metrics, nil
}
