package payload_test

import (
	"os"
	"path/filepath"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"

	"github.com/rioliu/liepin-cli-go/internal/payload"
)

var _ = Describe("LoadPayloadFile", func() {
	It("reads and parses JSON", func() {
		dir, _ := os.MkdirTemp("", "payload-test")
		defer os.RemoveAll(dir)
		path := filepath.Join(dir, "payload.json")
		os.WriteFile(path, []byte(`{"jobId":12,"jobKind":"2"}`), 0644)

		data, err := payload.LoadPayloadFile(path)
		Expect(err).NotTo(HaveOccurred())
		Expect(data["jobId"]).To(Equal(float64(12)))
		Expect(data["jobKind"]).To(Equal("2"))
	})

	It("returns empty map for empty path", func() {
		data, err := payload.LoadPayloadFile("")
		Expect(err).NotTo(HaveOccurred())
		Expect(data).To(BeEmpty())
	})

	It("returns error for missing file", func() {
		_, err := payload.LoadPayloadFile("/nonexistent/path.json")
		Expect(err).To(HaveOccurred())
	})

	It("returns error for invalid JSON", func() {
		dir, _ := os.MkdirTemp("", "payload-test")
		defer os.RemoveAll(dir)
		path := filepath.Join(dir, "bad.json")
		os.WriteFile(path, []byte("{not json"), 0644)

		_, err := payload.LoadPayloadFile(path)
		Expect(err).To(HaveOccurred())
	})

	It("returns error for JSON array root", func() {
		dir, _ := os.MkdirTemp("", "payload-test")
		defer os.RemoveAll(dir)
		path := filepath.Join(dir, "array.json")
		os.WriteFile(path, []byte(`[1,2,3]`), 0644)

		_, err := payload.LoadPayloadFile(path)
		Expect(err).To(HaveOccurred())
	})
})

var _ = Describe("MergePayload", func() {
	It("prefers explicit overrides over base", func() {
		merged := payload.MergePayload(
			map[string]any{"jobId": float64(12), "jobKind": "2"},
			map[string]any{"jobKind": "1", "unused": nil},
		)
		Expect(merged["jobId"]).To(Equal(float64(12)))
		Expect(merged["jobKind"]).To(Equal("1"))
		Expect(merged).NotTo(HaveKey("unused"))
	})

	It("strips nil values", func() {
		merged := payload.MergePayload(
			map[string]any{"a": float64(1), "b": nil},
			map[string]any{"c": float64(2)},
		)
		Expect(merged["a"]).To(Equal(float64(1)))
		Expect(merged["c"]).To(Equal(float64(2)))
		Expect(merged).NotTo(HaveKey("b"))
	})

	It("handles empty base", func() {
		merged := payload.MergePayload(
			map[string]any{},
			map[string]any{"x": "y"},
		)
		Expect(merged["x"]).To(Equal("y"))
	})
})
