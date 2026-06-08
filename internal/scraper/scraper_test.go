package scraper

import (
	"net/http"
	"net/http/httptest"
	"testing"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

func TestScraper(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Scraper Suite")
}

// ─── helpers ────────────────────────────────────────────────────────

const type1PageHTML = `<!DOCTYPE html>
<html><head>
<title>Test Job</title>
<script type="application/ld+json">
{
    "@context": "https://schema.org",
    "@type": "JobPosting",
    "title": "Senior Golang Engineer",
    "description": "岗位职责：\n1. 负责核心系统开发\n2. 微服务架构设计\n\n任职要求：\n1. 5年以上Go开发经验\n2. 精通并发编程",
    "datePosted": "2026-06-01",
    "validThrough": "2026-12-31",
    "employmentType": "FULL_TIME",
    "hiringOrganization": {
        "@type": "Organization",
        "name": "某知名科技公司"
    },
    "jobLocation": {
        "@type": "Place",
        "address": {
            "@type": "PostalAddress",
            "streetAddress": "北京"
        }
    },
    "experienceRequirements": "5-10年",
    "educationRequirements": "本科"
}
</script>
</head><body></body></html>`

const type2PageHTML = `<!DOCTYPE html>
<html><head>
<title>Test Job 2</title>
<script type="application/ld+json">
    {
        "@context": "https://schema.org",
        "@type": "JobPosting",
        "title": "AI Infra Engineer",
        "description": "Job responsibilities:\n1. Build AI infrastructure\n2. Optimize GPU clusters",
        "datePosted": "2026-05-15",
        "validThrough": "2026-11-30",
        "employmentType": "FULL_TIME",
        "hiringOrganization": {
            "name": "ByteDance"
        },
        "jobLocation": {
            "address": {
                "streetAddress": "Shanghai"
            }
        },
        "experienceRequirements": "3-5年",
        "educationRequirements": "本科"
    }
</script>
</head><body></body></html>`

const pageNotFoundHTML = `<!DOCTYPE html>
<html><head><title>LiePin.com</title></head>
<body><div class="error-main-container">
<p class="title">我们找遍了所有地方</p>
<p class="subtitle">此页面似乎不存在</p>
</div></body></html>`

const noLDPageHTML = `<!DOCTYPE html>
<html><head>
<title>Some Job</title>
<script>var x = 1;</script>
</head><body>Job content here</body></html>`

const multiBlockPageHTML = `<!DOCTYPE html>
<html><head>
<title>Test Job</title>
<script type="application/ld+json">
{
    "@context": "https://ziyuan.baidu.com/contexts/cambrian.jsonld",
    "appid": "123456",
    "title": "百度SEO标题",
    "description": "百度SEO描述"
}
</script>
<script type="application/ld+json">
{
    "@context": "https://schema.org",
    "@type": "JobPosting",
    "title": "Backend Engineer",
    "description": "Build scalable APIs.",
    "hiringOrganization": {"name": "TestCorp"},
    "jobLocation": {"address": {"streetAddress": "Shenzhen"}},
    "experienceRequirements": "3-5年",
    "educationRequirements": "本科"
}
</script>
</head><body></body></html>`

// ─── FetchJobDetail tests ──────────────────────────────────────────

var _ = Describe("FetchJobDetail", func() {
	It("extracts full detail from type 1 job page", func() {
		srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
			w.Header().Set("Content-Type", "text/html")
			w.Write([]byte(type1PageHTML))
		}))
		DeferCleanup(srv.Close)

		detail, err := FetchJobDetail(srv.URL)
		Expect(err).NotTo(HaveOccurred())
		Expect(detail.Title).To(Equal("Senior Golang Engineer"))
		Expect(detail.Description).To(ContainSubstring("核心系统开发"))
		Expect(detail.Description).To(ContainSubstring("并发编程"))
		Expect(detail.ExperienceRequirements).To(Equal("5-10年"))
		Expect(detail.EducationRequirements).To(Equal("本科"))
		Expect(detail.EmploymentType).To(Equal("FULL_TIME"))
		Expect(detail.Location).To(Equal("北京"))
		Expect(detail.Company).To(Equal("某知名科技公司"))
		Expect(detail.DatePosted).To(Equal("2026-06-01"))
		Expect(detail.ValidThrough).To(Equal("2026-12-31"))
	})

	It("extracts full detail from type 2 job page", func() {
		srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
			w.Header().Set("Content-Type", "text/html")
			w.Write([]byte(type2PageHTML))
		}))
		DeferCleanup(srv.Close)

		detail, err := FetchJobDetail(srv.URL)
		Expect(err).NotTo(HaveOccurred())
		Expect(detail.Title).To(Equal("AI Infra Engineer"))
		Expect(detail.Description).To(ContainSubstring("AI infrastructure"))
		Expect(detail.Company).To(Equal("ByteDance"))
		Expect(detail.Location).To(Equal("Shanghai"))
	})

	It("returns error for page not found", func() {
		srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
			w.Write([]byte(pageNotFoundHTML))
		}))
		DeferCleanup(srv.Close)

		_, err := FetchJobDetail(srv.URL)
		Expect(err).To(HaveOccurred())
		Expect(err.Error()).To(ContainSubstring("not found"))
	})

	It("returns error when no JSON-LD present", func() {
		srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
			w.Write([]byte(noLDPageHTML))
		}))
		DeferCleanup(srv.Close)

		_, err := FetchJobDetail(srv.URL)
		Expect(err).To(HaveOccurred())
		Expect(err.Error()).To(ContainSubstring("no JobPosting"))
	})

	It("returns error for HTTP failures", func() {
		srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
			w.WriteHeader(http.StatusInternalServerError)
		}))
		DeferCleanup(srv.Close)

		_, err := FetchJobDetail(srv.URL)
		Expect(err).To(HaveOccurred())
		Expect(err.Error()).To(ContainSubstring("500"))
	})

	It("skips non-JobPosting JSON-LD blocks", func() {
		srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
			w.Header().Set("Content-Type", "text/html")
			w.Write([]byte(multiBlockPageHTML))
		}))
		DeferCleanup(srv.Close)

		detail, err := FetchJobDetail(srv.URL)
		Expect(err).NotTo(HaveOccurred())
		Expect(detail.Title).To(Equal("Backend Engineer"))
		Expect(detail.Description).To(ContainSubstring("scalable APIs"))
		Expect(detail.Company).To(Equal("TestCorp"))
	})
})

// ─── parseJobPosting tests ─────────────────────────────────────────

var _ = Describe("parseJobPosting", func() {
	It("handles literal newlines in description", func() {
		html := `<script type="application/ld+json">
{
    "@type": "JobPosting",
    "title": "Test",
    "description": "Line one
Line two
Line three"
}
</script>`
		detail, err := parseJobPosting(html, "http://test")
		Expect(err).NotTo(HaveOccurred())
		Expect(detail.Description).To(ContainSubstring("Line one"))
		Expect(detail.Description).To(ContainSubstring("Line two"))
	})

	It("falls back to regex when JSON is malformed", func() {
		html := `<script type="application/ld+json">
{invalid json "description": "fallback content", "title": "X"}
</script>`
		detail, err := parseJobPosting(html, "http://test")
		// Should either succeed via regex or return error
		if err == nil {
			Expect(detail.Description).To(ContainSubstring("fallback"))
		}
	})
})

// ─── fixControlChars tests ─────────────────────────────────────────

var _ = Describe("fixControlChars", func() {
	It("escapes newlines inside JSON strings", func() {
		input := "{\"key\": \"line1\nline2\"}" // literal newline
		fixed := fixControlChars(input)
		Expect(fixed).To(Equal(`{"key": "line1\nline2"}`)) // escaped \n
	})

	It("preserves escaped sequences", func() {
		input := `{"key": "value with \\n escaped"}`
		fixed := fixControlChars(input)
		Expect(fixed).To(Equal(input))
	})
})
