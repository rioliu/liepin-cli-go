package client_test

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"time"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"

	"github.com/rioliu/liepin-cli-go/internal/client"
)

func newClient(srv *httptest.Server) *client.Client {
	return client.New(client.Config{Token: "token-123", BaseURL: srv.URL, Timeout: 0})
}

func jsonOK(w http.ResponseWriter) {
	w.Header().Set(client.HeaderContentType, client.MediaTypeJSON)
	json.NewEncoder(w).Encode(map[string]bool{"ok": true})
}

var _ = Describe("Client GET", func() {
	It("sends auth header", func() {
		srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			defer GinkgoRecover()
			Expect(r.Header.Get(client.HeaderXUserToken)).To(Equal("token-123"))
			Expect(r.Method).To(Equal(http.MethodGet))
			jsonOK(w)
		}))
		defer srv.Close()

		c := newClient(srv)
		result, err := c.Get("/mcp/get-resume")
		Expect(err).NotTo(HaveOccurred())
		Expect(result.(map[string]any)["ok"]).To(BeTrue())
	})

	It("preserves base URL path prefix", func() {
		srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			defer GinkgoRecover()
			Expect(r.URL.Path).To(Equal("/proxy/api/mcp/get-resume"))
			jsonOK(w)
		}))
		defer srv.Close()

		c := client.New(client.Config{Token: "token-123", BaseURL: srv.URL + "/proxy/api", Timeout: 0})
		_, err := c.Get("/mcp/get-resume")
		Expect(err).NotTo(HaveOccurred())
	})

	It("returns AuthorizationError for 401", func() {
		srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
			w.WriteHeader(http.StatusUnauthorized)
		}))
		defer srv.Close()

		c := newClient(srv)
		_, err := c.Get("/any")
		Expect(err).To(BeAssignableToTypeOf(&client.AuthorizationError{}))
		authErr := err.(*client.AuthorizationError)
		Expect(authErr.StatusCode).To(Equal(401))
	})

	It("returns RequestError for 500", func() {
		srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
			w.WriteHeader(http.StatusInternalServerError)
		}))
		defer srv.Close()

		c := newClient(srv)
		_, err := c.Get("/any")
		Expect(err).To(BeAssignableToTypeOf(&client.RequestError{}))
		reqErr := err.(*client.RequestError)
		Expect(reqErr.StatusCode).To(Equal(500))
	})

	It("returns error for invalid JSON response", func() {
		srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
			w.Header().Set(client.HeaderContentType, client.MediaTypeJSON)
			w.Write([]byte("not-json"))
		}))
		defer srv.Close()

		c := newClient(srv)
		_, err := c.Get("/bad-json")
		Expect(err).To(HaveOccurred())
	})

	It("returns plain text as string", func() {
		srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
			w.Header().Set(client.HeaderContentType, "text/plain")
			w.Write([]byte("hello"))
		}))
		defer srv.Close()

		c := newClient(srv)
		result, err := c.Get("/text")
		Expect(err).NotTo(HaveOccurred())
		Expect(result).To(Equal("hello"))
	})

	It("returns nil for empty response", func() {
		srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
			w.WriteHeader(http.StatusOK)
		}))
		defer srv.Close()

		c := newClient(srv)
		result, err := c.Get("/empty")
		Expect(err).NotTo(HaveOccurred())
		Expect(result).To(BeNil())
	})

	It("returns error for bad URL", func() {
		c := client.New(client.Config{Token: "t", BaseURL: "http://127.0.0.1:0", Timeout: 0})
		_, err := c.Get("/x")
		Expect(err).To(HaveOccurred())
	})

	It("returns RateLimitError for 429 with Retry-After header", func() {
		srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
			w.Header().Set("Retry-After", "30")
			w.WriteHeader(http.StatusTooManyRequests)
			w.Write([]byte(`{"msg":"too many requests"}`))
		}))
		defer srv.Close()

		c := newClient(srv)
		_, err := c.Get("/any")
		Expect(err).To(BeAssignableToTypeOf(&client.RateLimitError{}))
		rlErr := err.(*client.RateLimitError)
		Expect(rlErr.StatusCode).To(Equal(429))
		Expect(rlErr.RetryAfter).To(Equal(30 * time.Second))
		Expect(rlErr.Body).To(ContainSubstring("too many requests"))
	})

	It("returns RateLimitError for 429 with X-RateLimit-Reset header", func() {
		resetTime := time.Now().Add(60 * time.Second).Unix()
		srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
			w.Header().Set("X-RateLimit-Reset", fmt.Sprintf("%d", resetTime))
			w.WriteHeader(http.StatusTooManyRequests)
		}))
		defer srv.Close()

		c := newClient(srv)
		_, err := c.Get("/any")
		Expect(err).To(BeAssignableToTypeOf(&client.RateLimitError{}))
		rlErr := err.(*client.RateLimitError)
		Expect(rlErr.RetryAfter).To(BeNumerically(">=", 59*time.Second))
		Expect(rlErr.RetryAfter).To(BeNumerically("<=", 61*time.Second))
	})

	It("returns RateLimitError for 429 without headers uses default fallback", func() {
		srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
			w.WriteHeader(http.StatusTooManyRequests)
		}))
		defer srv.Close()

		c := newClient(srv)
		_, err := c.Get("/any")
		Expect(err).To(BeAssignableToTypeOf(&client.RateLimitError{}))
		rlErr := err.(*client.RateLimitError)
		Expect(rlErr.RetryAfter).To(Equal(client.DefaultRateLimitWait))
	})

	It("returns RateLimitError for 429 with negative Retry-After uses fallback", func() {
		srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
			w.Header().Set("Retry-After", "-5")
			w.WriteHeader(http.StatusTooManyRequests)
		}))
		defer srv.Close()

		c := newClient(srv)
		_, err := c.Get("/any")
		Expect(err).To(BeAssignableToTypeOf(&client.RateLimitError{}))
		rlErr := err.(*client.RateLimitError)
		Expect(rlErr.RetryAfter).To(Equal(client.DefaultRateLimitWait))
	})
})

var _ = Describe("Client POST", func() {
	It("sends JSON payload with auth header", func() {
		srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			defer GinkgoRecover()
			Expect(r.Header.Get(client.HeaderXUserToken)).To(Equal("token-123"))
			Expect(r.Method).To(Equal(http.MethodPost))
			Expect(r.Header.Get(client.HeaderContentType)).To(Equal(client.MediaTypeJSON))
			jsonOK(w)
		}))
		defer srv.Close()

		c := newClient(srv)
		payload := map[string]any{"jobId": 1, "jobKind": "2"}
		result, err := c.Post("/mcp/apply-job", payload)
		Expect(err).NotTo(HaveOccurred())
		Expect(result.(map[string]any)["ok"]).To(BeTrue())
	})

	It("returns AuthorizationError for 403", func() {
		srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
			w.WriteHeader(http.StatusForbidden)
		}))
		defer srv.Close()

		c := newClient(srv)
		_, err := c.Post("/any", map[string]any{})
		Expect(err).To(BeAssignableToTypeOf(&client.AuthorizationError{}))
		authErr := err.(*client.AuthorizationError)
		Expect(authErr.StatusCode).To(Equal(403))
	})
})
