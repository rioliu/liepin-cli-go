// Package client implements an HTTP client for the Liepin open-agent API.
// It centralizes token handling, base-URL joining, TLS configuration, and
// the typed error returned for authorization, request, and TLS failures.
package client

import (
	"bytes"
	"crypto/tls"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strings"
	"time"
)

// Config holds the runtime parameters used to construct a Client.
type Config struct {
	Token              string
	BaseURL            string
	Timeout            time.Duration
	InsecureSkipVerify bool
}

// Client is an HTTP client that talks to the Liepin open-agent API. It
// automatically attaches the configured x-user-token to outgoing requests
// and decodes JSON responses.
type Client struct {
	config     Config
	httpClient *http.Client
}

// HTTP header names and media types used by the Liepin client.
const (
	HeaderXUserToken  = "x-user-token"
	HeaderContentType = "Content-Type"
	MediaTypeJSON     = "application/json"
)

// DefaultRateLimitWait is the fallback duration to wait before retrying when
// the server returns 429 but no Retry-After or X-RateLimit-Reset header.
const DefaultRateLimitWait = 5 * time.Second

// TLSError wraps a transport-layer TLS or x509 failure so callers can
// distinguish certificate issues from generic request errors and surface a
// helpful hint (such as the --insecure flag).
type TLSError struct {
	Err error
}

func (e *TLSError) Error() string {
	return e.Err.Error()
}

// New creates a Client using the supplied Config. When InsecureSkipVerify is
// true, the underlying HTTP transport is cloned with TLS verification
// disabled (intended for development against self-signed servers only).
func New(config Config) *Client {
	transport := http.DefaultTransport.(*http.Transport).Clone()
	if config.InsecureSkipVerify {
		transport.TLSClientConfig = &tls.Config{
			InsecureSkipVerify: true,
		}
	}
	return &Client{
		config: config,
		httpClient: &http.Client{
			Timeout:   config.Timeout,
			Transport: transport,
		},
	}
}

// RequestError represents a non-2xx HTTP response from the API that is not
// an authentication failure. Body contains the trimmed response body, if any.
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

// AuthorizationError is returned when the API responds with 401 Unauthorized
// or 403 Forbidden, indicating that the configured token is missing,
// expired, or lacks the required permissions.
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

// RateLimitError is returned when the API responds with 429 Too Many
// Requests. RetryAfter is parsed from the Retry-After header (seconds) or
// X-RateLimit-Reset header (Unix timestamp). Callers should sleep for the
// indicated duration before retrying.
type RateLimitError struct {
	StatusCode int
	Body       string
	RetryAfter time.Duration
}

func (e *RateLimitError) Error() string {
	parts := []string{fmt.Sprintf("%d", e.StatusCode)}
	if e.Body != "" {
		parts = append(parts, e.Body)
	}
	if e.RetryAfter > 0 {
		parts = append(parts, fmt.Sprintf("retry after %s", e.RetryAfter))
	}
	return strings.Join(parts, " ")
}

func (c *Client) buildURL(path string) (string, error) {
	base := strings.TrimRight(c.config.BaseURL, "/") + "/"
	rel := strings.TrimLeft(path, "/")
	return url.JoinPath(base, rel)
}

// Get performs an HTTP GET against the given API path (joined to BaseURL),
// attaches the configured token header, and returns the decoded response.
// Errors of type *AuthorizationError, *RequestError, or *TLSError signal
// specific failure modes; all other errors are wrapped with %w.
func (c *Client) Get(path string) (any, error) {
	reqURL, err := c.buildURL(path)
	if err != nil {
		return nil, fmt.Errorf("failed to build URL: %w", err)
	}

	req, err := http.NewRequest(http.MethodGet, reqURL, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}
	req.Header.Set(HeaderXUserToken, c.config.Token)

	return c.doRequest(req)
}

// Post performs an HTTP POST against the given API path (joined to BaseURL)
// with the supplied payload marshalled as JSON. It returns the decoded
// response body, or one of *AuthorizationError, *RequestError, *TLSError
// for the corresponding failure modes.
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
	req.Header.Set(HeaderXUserToken, c.config.Token)
	req.Header.Set(HeaderContentType, MediaTypeJSON)

	return c.doRequest(req)
}

func (c *Client) doRequest(req *http.Request) (any, error) {
	resp, err := c.httpClient.Do(req)
	if err != nil {
		if isTLSError(err) {
			return nil, &TLSError{Err: err}
		}
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
	if isRateLimit(resp.StatusCode, bodyBytes) {
		return nil, newRateLimitError(resp.StatusCode, bodyStr, resp.Header)
	}
	if resp.StatusCode >= 400 {
		return nil, &RequestError{StatusCode: resp.StatusCode, Body: bodyStr}
	}

	contentType := resp.Header.Get(HeaderContentType)
	if strings.Contains(strings.ToLower(contentType), MediaTypeJSON) {
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

// parseRetryAfter extracts the wait duration from rate-limit response headers.
// It checks, in order:
//  1. Retry-After header (seconds as integer, or HTTP-date)
//  2. X-RateLimit-Reset header (Unix timestamp in seconds)
//
// Returns 0 if no valid header is found.
func parseRetryAfter(h http.Header) time.Duration {
	// Standard Retry-After: seconds or HTTP-date
	if ra := h.Get("Retry-After"); ra != "" {
		if secs, err := time.ParseDuration(ra + "s"); err == nil && secs > 0 {
			return secs
		}
		if t, err := time.Parse(time.RFC1123, ra); err == nil {
			d := time.Until(t)
			if d > 0 {
				return d
			}
		}
	}
	// X-RateLimit-Reset: Unix timestamp
	if reset := h.Get("X-RateLimit-Reset"); reset != "" {
		if ts, err := time.Parse(time.UnixDate, reset); err == nil {
			d := time.Until(ts)
			if d > 0 {
				return d
			}
		}
		// Try as Unix timestamp (seconds since epoch)
		var unixSec int64
		if _, err := fmt.Sscanf(reset, "%d", &unixSec); err == nil && unixSec > 0 {
			d := time.Until(time.Unix(unixSec, 0))
			if d > 0 {
				return d
			}
		}
	}
	return 0
}

func newRateLimitError(statusCode int, body string, headers http.Header) *RateLimitError {
	retryAfter := parseRetryAfter(headers)
	if retryAfter == 0 {
		retryAfter = DefaultRateLimitWait
	}
	return &RateLimitError{
		StatusCode: statusCode,
		Body:       body,
		RetryAfter: retryAfter,
	}
}

// isRateLimit detects both HTTP-level and application-level rate limits:
//   - HTTP 429 status code
//   - HTTP 200 with code 429001 in the JSON body
func isRateLimit(statusCode int, body []byte) bool {
	if statusCode == http.StatusTooManyRequests {
		return true
	}
	var envelope struct {
		Code int `json:"code"`
	}
	if err := json.Unmarshal(body, &envelope); err != nil {
		return false
	}
	return envelope.Code == 429001
}

func isTLSError(err error) bool {
	if err == nil {
		return false
	}
	errStr := err.Error()
	return strings.Contains(errStr, "tls:") ||
		strings.Contains(errStr, "x509:") ||
		strings.Contains(errStr, "certificate")
}
