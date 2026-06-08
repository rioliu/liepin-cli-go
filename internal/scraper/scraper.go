// Package scraper fetches job detail pages from liepin.com and extracts
// structured data from the embedded schema.org JobPosting JSON-LD.
package scraper

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"regexp"
	"strings"
	"time"
)

// JobDetail holds the structured data extracted from a liepin job page.
type JobDetail struct {
	Title                  string `json:"title"`
	Description            string `json:"description"`
	ExperienceRequirements string `json:"experienceRequirements"`
	EducationRequirements  string `json:"educationRequirements"`
	EmploymentType         string `json:"employmentType"`
	Location               string `json:"location"`
	Company                string `json:"company"`
	DatePosted             string `json:"datePosted"`
	ValidThrough           string `json:"validThrough"`
	URL                    string `json:"url"`
}

// ldJobPosting is the intermediate struct for parsing the schema.org
// JobPosting JSON-LD block embedded in the page.
type ldJobPosting struct {
	Context                string `json:"@context"`
	Type                   string `json:"@type"`
	Title                  string `json:"title"`
	Description            string `json:"description"`
	DatePosted             string `json:"datePosted"`
	ValidThrough           string `json:"validThrough"`
	EmploymentType         string `json:"employmentType"`
	ExperienceRequirements string `json:"experienceRequirements"`
	EducationRequirements  string `json:"educationRequirements"`
	URL                    string `json:"url"`
	HiringOrganization     struct {
		Name string `json:"name"`
	} `json:"hiringOrganization"`
	JobLocation struct {
		Address struct {
			StreetAddress string `json:"streetAddress"`
		} `json:"address"`
	} `json:"jobLocation"`
}

var (
	reJobPosting = regexp.MustCompile(`(?s)<script[^>]*type="application/ld\+json"[^>]*>(.*?)</script>`)
	reDescField  = regexp.MustCompile(`(?s)"description"\s*:\s*"((?:[^"\\]|\\.)*)"\s*[,}]`)
)

// FetchJobDetail fetches the job page at the given URL and extracts the
// structured JobPosting data. The URL is typically the jobDetailUrl from
// search results (e.g. https://www.liepin.com/a/12345.shtml for type 1
// or https://www.liepin.com/job/1912345.shtml for type 2).
func FetchJobDetail(url string) (*JobDetail, error) {
	client := &http.Client{Timeout: 15 * time.Second}
	req, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}
	req.Header.Set("User-Agent", "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/120.0.0.0 Safari/537.36")

	resp, err := client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("request failed: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("HTTP %d from %s", resp.StatusCode, url)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response: %w", err)
	}
	html := string(body)

	if strings.Contains(html, "此页面似乎不存在") {
		return nil, fmt.Errorf("job page not found: %s", url)
	}

	return parseJobPosting(html, url)
}

// parseJobPosting extracts the schema.org JobPosting JSON-LD from HTML.
func parseJobPosting(html, url string) (*JobDetail, error) {
	// Primary: find <script type="application/ld+json"> blocks containing JobPosting
	matches := reJobPosting.FindAllStringSubmatch(html, -1)
	for _, match := range matches {
		if len(match) < 2 || !strings.Contains(match[1], "JobPosting") {
			continue
		}
		ld, err := decodeJobPosting(match[1])
		if err == nil && ld.Title != "" {
			return ld.toJobDetail(url), nil
		}
	}

	// Fallback: scan all <script> blocks for JobPosting
	for _, block := range extractScriptBlocks(html) {
		if !strings.Contains(block, "JobPosting") {
			continue
		}
		ld, err := decodeJobPosting(block)
		if err == nil && ld.Title != "" {
			return ld.toJobDetail(url), nil
		}
		// If JSON decode fails, try regex for description
		if desc := reDescField.FindStringSubmatch(block); desc != nil {
			return &JobDetail{
				Description: unescapeJSON(desc[1]),
				URL:         url,
			}, nil
		}
	}

	return nil, fmt.Errorf("no JobPosting data found in page")
}

// decodeJobPosting parses JSON-LD content into ldJobPosting, tolerating
// literal newlines inside JSON strings (strict=false).
func decodeJobPosting(raw string) (*ldJobPosting, error) {
	var ld ldJobPosting
	dec := json.NewDecoder(strings.NewReader(raw))
	dec.UseNumber()
	if err := dec.Decode(&ld); err != nil {
		// Try with strict=false via Unmarshal
		if err2 := json.Unmarshal([]byte(raw), &ld); err2 != nil {
			// Last resort: fix control characters and retry
			fixed := fixControlChars(raw)
			if err3 := json.Unmarshal([]byte(fixed), &ld); err3 != nil {
				return nil, err3
			}
		}
	}
	return &ld, nil
}

// toJobDetail converts ldJobPosting to the public JobDetail struct.
func (ld *ldJobPosting) toJobDetail(url string) *JobDetail {
	loc := ld.JobLocation.Address.StreetAddress
	return &JobDetail{
		Title:                  ld.Title,
		Description:            ld.Description,
		ExperienceRequirements: ld.ExperienceRequirements,
		EducationRequirements:  ld.EducationRequirements,
		EmploymentType:         ld.EmploymentType,
		Location:               loc,
		Company:                ld.HiringOrganization.Name,
		DatePosted:             ld.DatePosted,
		ValidThrough:           ld.ValidThrough,
		URL:                    url,
	}
}

// extractScriptBlocks returns the content of all <script> tags.
func extractScriptBlocks(html string) []string {
	var blocks []string
	remaining := html
	for {
		tagStart := strings.Index(remaining, "<script")
		if tagStart < 0 {
			break
		}
		tagEnd := strings.Index(remaining[tagStart:], ">")
		if tagEnd < 0 {
			break
		}
		contentStart := tagStart + tagEnd + 1
		closeTag := strings.Index(remaining[contentStart:], "</script>")
		if closeTag < 0 {
			break
		}
		blocks = append(blocks, remaining[contentStart:contentStart+closeTag])
		remaining = remaining[contentStart+closeTag+9:]
	}
	return blocks
}

// fixControlChars replaces literal newlines/tabs inside JSON strings with
// escaped versions so that json.Unmarshal (strict mode) can parse them.
func fixControlChars(s string) string {
	var b strings.Builder
	inString := false
	escaped := false
	for i := 0; i < len(s); i++ {
		c := s[i]
		if escaped {
			b.WriteByte(c)
			escaped = false
			continue
		}
		if c == '\\' && inString {
			b.WriteByte(c)
			escaped = true
			continue
		}
		if c == '"' {
			inString = !inString
			b.WriteByte(c)
			continue
		}
		if inString {
			switch c {
			case '\n':
				b.WriteString("\\n")
			case '\r':
				b.WriteString("\\r")
			case '\t':
				b.WriteString("\\t")
			default:
				b.WriteByte(c)
			}
		} else {
			b.WriteByte(c)
		}
	}
	return b.String()
}

// unescapeJSON reverses common JSON escape sequences.
func unescapeJSON(s string) string {
	s = strings.ReplaceAll(s, "\\n", "\n")
	s = strings.ReplaceAll(s, "\\r", "\r")
	s = strings.ReplaceAll(s, "\\t", "\t")
	s = strings.ReplaceAll(s, "\\\"", "\"")
	s = strings.ReplaceAll(s, "\\\\", "\\")
	return s
}
