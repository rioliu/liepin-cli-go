package authstore_test

import (
	"os"
	"path/filepath"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"

	"github.com/rioliu/liepin-cli-go/internal/authstore"
)

var _ = Describe("Authstore", func() {
	var dir string

	BeforeEach(func() {
		dir, _ = os.MkdirTemp("", "auth-test")
		DeferCleanup(func() { os.RemoveAll(dir) })
		os.Setenv("XDG_CONFIG_HOME", dir)
		DeferCleanup(func() { os.Unsetenv("XDG_CONFIG_HOME") })
	})

	Describe("SaveToken", func() {
		It("saves and loads token", func() {
			Expect(authstore.SaveToken("x-user-token-abc")).To(Succeed())

			token, err := authstore.LoadToken()
			Expect(err).NotTo(HaveOccurred())
			Expect(token).To(Equal("x-user-token-abc"))
		})

		It("updates existing config", func() {
			Expect(authstore.SaveToken("first")).To(Succeed())
			Expect(authstore.SaveToken("second")).To(Succeed())

			token, err := authstore.LoadToken()
			Expect(err).NotTo(HaveOccurred())
			Expect(token).To(Equal("second"))
		})
	})

	Describe("LoadToken", func() {
		It("returns empty string for empty config", func() {
			token, err := authstore.LoadToken()
			Expect(err).NotTo(HaveOccurred())
			Expect(token).To(BeEmpty())
		})
	})

	Describe("ClearToken", func() {
		It("clears token when present", func() {
			Expect(authstore.SaveToken("token-123")).To(Succeed())

			cleared, err := authstore.ClearToken()
			Expect(err).NotTo(HaveOccurred())
			Expect(cleared).To(BeTrue())

			token, err := authstore.LoadToken()
			Expect(err).NotTo(HaveOccurred())
			Expect(token).To(BeEmpty())
		})

		It("returns false when no config exists", func() {
			cleared, err := authstore.ClearToken()
			Expect(err).NotTo(HaveOccurred())
			Expect(cleared).To(BeFalse())
		})

		It("returns false when token field is missing", func() {
			cfgDir := filepath.Join(dir, "liepin-cli")
			os.MkdirAll(cfgDir, 0755)
			os.WriteFile(filepath.Join(cfgDir, "config.json"), []byte(`{}`), 0644)

			cleared, err := authstore.ClearToken()
			Expect(err).NotTo(HaveOccurred())
			Expect(cleared).To(BeFalse())
		})
	})
})
