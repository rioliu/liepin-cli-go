//go:build e2e

package e2e

import (
	"os"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"

	"github.com/rioliu/liepin-cli-go/internal/client"
	"github.com/rioliu/liepin-cli-go/internal/config"
)

var _ = Describe("E2E: read-only endpoints", func() {
	var c *client.Client

	BeforeEach(func() {
		token := os.Getenv("LIEPIN_USER_TOKEN")
		if token == "" {
			Skip("LIEPIN_USER_TOKEN not set")
		}

		baseURL := os.Getenv("LIEPIN_BASE_URL")

		cfg, err := config.ResolveConfig(token, baseURL, "pretty", 30.0)
		Expect(err).NotTo(HaveOccurred())

		c = client.New(client.Config{
			Token:   cfg.Token,
			BaseURL: cfg.BaseURL,
			Timeout: cfg.Timeout,
		})
	})

	It("GET /mcp/get-resume returns resume data", func() {
		data, err := c.Get("/mcp/get-resume")
		Expect(err).NotTo(HaveOccurred())
		Expect(data).NotTo(BeNil())

		m, ok := data.(map[string]any)
		Expect(ok).To(BeTrue(), "response should be a JSON object")
		Expect(m).To(HaveKey("errCode"))
		Expect(m).To(HaveKey("data"))
	})

	It("POST /mcp/search-job returns search results", func() {
		payload := map[string]any{
			"jobName": "开发",
			"address": "北京",
			"page":    0,
		}
		data, err := c.Post("/mcp/search-job", payload)
		Expect(err).NotTo(HaveOccurred())
		Expect(data).NotTo(BeNil())

		m, ok := data.(map[string]any)
		Expect(ok).To(BeTrue(), "response should be a JSON object")
		Expect(m).NotTo(BeEmpty())
	})
})

var _ = Describe("E2E: write endpoints", func() {
	var c *client.Client
	var isProd bool

	BeforeEach(func() {
		token := os.Getenv("LIEPIN_USER_TOKEN")
		if token == "" {
			Skip("LIEPIN_USER_TOKEN not set")
		}

		baseURL := os.Getenv("LIEPIN_BASE_URL")

		cfg, err := config.ResolveConfig(token, baseURL, "pretty", 30.0)
		Expect(err).NotTo(HaveOccurred())

		isProd = config.IsProduction(cfg.BaseURL)

		c = client.New(client.Config{
			Token:   cfg.Token,
			BaseURL: cfg.BaseURL,
			Timeout: cfg.Timeout,
		})
	})

	expectJSONObject := func(data any) map[string]any {
		Expect(data).NotTo(BeNil())
		m, ok := data.(map[string]any)
		Expect(ok).To(BeTrue(), "response should be a JSON object")
		return m
	}

	It("POST /mcp/apply-job", func() {
		if isProd {
			Skip("write test skipped against production — set LIEPIN_BASE_URL to a test environment")
		}
		payload := map[string]any{"jobId": 0, "jobKind": "0"}
		data, err := c.Post("/mcp/apply-job", payload)
		Expect(err).NotTo(HaveOccurred())
		expectJSONObject(data)
	})

	It("POST /mcp/update-base-info", func() {
		if isProd {
			Skip("write test skipped against production — set LIEPIN_BASE_URL to a test environment")
		}
		payload := map[string]any{"realName": "e2e-test"}
		data, err := c.Post("/mcp/update-base-info", payload)
		Expect(err).NotTo(HaveOccurred())
		expectJSONObject(data)
	})

	It("POST /mcp/update-self-assess", func() {
		if isProd {
			Skip("write test skipped against production — set LIEPIN_BASE_URL to a test environment")
		}
		payload := map[string]any{"selfAssess": "e2e test assessment"}
		data, err := c.Post("/mcp/update-self-assess", payload)
		Expect(err).NotTo(HaveOccurred())
		expectJSONObject(data)
	})

	It("POST /mcp/add-edu-exp", func() {
		if isProd {
			Skip("write test skipped against production — set LIEPIN_BASE_URL to a test environment")
		}
		payload := map[string]any{
			"school": "e2e-test-school",
			"start":  "202001",
			"end":    "202406",
			"degree": "5",
		}
		data, err := c.Post("/mcp/add-edu-exp", payload)
		Expect(err).NotTo(HaveOccurred())
		expectJSONObject(data)
	})

	It("POST /mcp/update-edu-exp", func() {
		if isProd {
			Skip("write test skipped against production — set LIEPIN_BASE_URL to a test environment")
		}
		payload := map[string]any{"eduId": 0, "school": "e2e-test-updated"}
		data, err := c.Post("/mcp/update-edu-exp", payload)
		Expect(err).NotTo(HaveOccurred())
		expectJSONObject(data)
	})

	It("POST /mcp/add-work-exp", func() {
		if isProd {
			Skip("write test skipped against production — set LIEPIN_BASE_URL to a test environment")
		}
		payload := map[string]any{
			"compName":  "e2e-test-company",
			"workStart": "202001",
			"workEnd":   "202406",
		}
		data, err := c.Post("/mcp/add-work-exp", payload)
		Expect(err).NotTo(HaveOccurred())
		expectJSONObject(data)
	})

	It("POST /mcp/update-work-exp", func() {
		if isProd {
			Skip("write test skipped against production — set LIEPIN_BASE_URL to a test environment")
		}
		payload := map[string]any{"workId": 0, "compName": "e2e-test-updated"}
		data, err := c.Post("/mcp/update-work-exp", payload)
		Expect(err).NotTo(HaveOccurred())
		expectJSONObject(data)
	})

	It("POST /mcp/add-project-exp", func() {
		if isProd {
			Skip("write test skipped against production — set LIEPIN_BASE_URL to a test environment")
		}
		payload := map[string]any{
			"name":  "e2e-test-project",
			"start": "202301",
			"end":   "202312",
		}
		data, err := c.Post("/mcp/add-project-exp", payload)
		Expect(err).NotTo(HaveOccurred())
		expectJSONObject(data)
	})

	It("POST /mcp/update-project-exp", func() {
		if isProd {
			Skip("write test skipped against production — set LIEPIN_BASE_URL to a test environment")
		}
		payload := map[string]any{"id": 0, "name": "e2e-test-updated"}
		data, err := c.Post("/mcp/update-project-exp", payload)
		Expect(err).NotTo(HaveOccurred())
		expectJSONObject(data)
	})

	It("POST /mcp/add-job-want", func() {
		if isProd {
			Skip("write test skipped against production — set LIEPIN_BASE_URL to a test environment")
		}
		payload := map[string]any{
			"jobtitle": "e2e-test-title",
			"dq":       "北京",
		}
		data, err := c.Post("/mcp/add-job-want", payload)
		Expect(err).NotTo(HaveOccurred())
		expectJSONObject(data)
	})

	It("POST /mcp/update-job-want", func() {
		if isProd {
			Skip("write test skipped against production — set LIEPIN_BASE_URL to a test environment")
		}
		payload := map[string]any{"id": 0, "jobtitle": "e2e-test-updated"}
		data, err := c.Post("/mcp/update-job-want", payload)
		Expect(err).NotTo(HaveOccurred())
		expectJSONObject(data)
	})
})
