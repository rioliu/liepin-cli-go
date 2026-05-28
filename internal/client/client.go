package client

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strings"
	"time"
)

type Config struct {
	Token   string
	BaseURL string
	Timeout time.Duration
}

type Client struct {
	config     Config
	httpClient *http.Client
}

func New(config Config) *Client {
	return &Client{
		config: config,
		httpClient: &http.Client{
			Timeout: config.Timeout,
		},
	}
}

type RequestError struct {
	StatusCode int
	Body       string
}

func (e *RequestError) Error() string {
	parts := []string{fmt.Sprintf("%d", e.StatusCode)}
	if e.Body != "" {
		parts = append(parts, e.Body)
	}
	return strings.Join(parts, " ")
}

type AuthorizationError struct {
	StatusCode int
	Body       string
}

func (e *AuthorizationError) Error() string {
	parts := []string{fmt.Sprintf("%d", e.StatusCode)}
	if e.Body != "" {
		parts = append(parts, e.Body)
	}
	return strings.Join(parts, " ")
}

func (c *Client) buildURL(path string) (string, error) {
	base := strings.TrimRight(c.config.BaseURL, "/") + "/"
	rel := strings.TrimLeft(path, "/")
	return url.JoinPath(base, rel)
}

func (c *Client) Get(path string) (any, error) {
	reqURL, err := c.buildURL(path)
	if err != nil {
		return nil, fmt.Errorf("failed to build URL: %w", err)
	}

	req, err := http.NewRequest(http.MethodGet, reqURL, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}
	req.Header.Set("x-user-token", c.config.Token)

	return c.doRequest(req)
}

func (c *Client) Post(path string, payload any) (any, error) {
	reqURL, err := c.buildURL(path)
	if err != nil {
		return nil, fmt.Errorf("failed to build URL: %w", err)
	}

	body, err := json.Marshal(payload)
	if err != nil {
		return nil, fmt.Errorf("failed to serialize request body: %w", err)
	}

	req, err := http.NewRequest(http.MethodPost, reqURL, bytes.NewReader(body))
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}
	req.Header.Set("x-user-token", c.config.Token)
	req.Header.Set("Content-Type", "application/json")

	return c.doRequest(req)
}

func (c *Client) doRequest(req *http.Request) (any, error) {
	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("request failed: %w", err)
	}
	defer resp.Body.Close()

	bodyBytes, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response: %w", err)
	}
	bodyStr := strings.TrimSpace(string(bodyBytes))

	if resp.StatusCode == http.StatusUnauthorized || resp.StatusCode == http.StatusForbidden {
		return nil, &AuthorizationError{StatusCode: resp.StatusCode, Body: bodyStr}
	}
	if resp.StatusCode >= 400 {
		return nil, &RequestError{StatusCode: resp.StatusCode, Body: bodyStr}
	}

	contentType := resp.Header.Get("Content-Type")
	if strings.Contains(strings.ToLower(contentType), "application/json") {
		var result any
		if err := json.Unmarshal(bodyBytes, &result); err != nil {
			return nil, fmt.Errorf("invalid JSON response: %w", err)
		}
		return result, nil
	}

	if bodyStr == "" {
		return nil, nil
	}
	return bodyStr, nil
}
