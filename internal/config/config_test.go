package config_test

import (
	"os"
	"path/filepath"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"

	"github.com/rioliu/liepin-cli-go/internal/authstore"
	"github.com/rioliu/liepin-cli-go/internal/config"
)

var _ = Describe("ResolveConfig", func() {
	Describe("token resolution", func() {
		It("prefers flag token over env", func() {
			os.Setenv("LIEPIN_USER_TOKEN", "env-token")
			defer os.Unsetenv("LIEPIN_USER_TOKEN")

			cfg, err := config.ResolveConfig("flag-token", "", "pretty", 30.0)
			Expect(err).NotTo(HaveOccurred())
			Expect(cfg.Token).To(Equal("flag-token"))
			Expect(cfg.TokenSource).To(Equal("cli"))
		})

		It("uses env token when no flag", func() {
			os.Setenv("LIEPIN_USER_TOKEN", "env-token")
			defer os.Unsetenv("LIEPIN_USER_TOKEN")

			cfg, err := config.ResolveConfig("", "", "pretty", 30.0)
			Expect(err).NotTo(HaveOccurred())
			Expect(cfg.Token).To(Equal("env-token"))
			Expect(cfg.TokenSource).To(Equal("env"))
		})

		It("uses config file when no flag or env", func() {
			dir, _ := os.MkdirTemp("", "cfg-test")
			defer os.RemoveAll(dir)
			os.Setenv("XDG_CONFIG_HOME", dir)
			defer os.Unsetenv("XDG_CONFIG_HOME")
			os.Unsetenv("LIEPIN_USER_TOKEN")

			Expect(authstore.SaveToken("config-token")).To(Succeed())

			cfg, err := config.ResolveConfig("", "", "pretty", 30.0)
			Expect(err).NotTo(HaveOccurred())
			Expect(cfg.Token).To(Equal("config-token"))
			Expect(cfg.TokenSource).To(Equal("config"))
		})

		It("returns MissingTokenError when no token anywhere", func() {
			dir, _ := os.MkdirTemp("", "cfg-test")
			defer os.RemoveAll(dir)
			os.Setenv("XDG_CONFIG_HOME", dir)
			defer os.Unsetenv("XDG_CONFIG_HOME")
			os.Unsetenv("LIEPIN_USER_TOKEN")

			_, err := config.ResolveConfig("", "", "pretty", 30.0)
			Expect(err).To(BeAssignableToTypeOf(&config.MissingTokenError{}))
		})
	})

	Describe("base URL", func() {
		BeforeEach(func() {
			os.Setenv("LIEPIN_USER_TOKEN", "t")
			DeferCleanup(func() { os.Unsetenv("LIEPIN_USER_TOKEN") })
		})

		It("uses default when not specified", func() {
			cfg, err := config.ResolveConfig("", "", "pretty", 30.0)
			Expect(err).NotTo(HaveOccurred())
			Expect(cfg.BaseURL).To(Equal(config.DefaultBaseURL))
		})

		It("uses custom base URL when specified", func() {
			cfg, err := config.ResolveConfig("", "https://custom.example.com", "pretty", 30.0)
			Expect(err).NotTo(HaveOccurred())
			Expect(cfg.BaseURL).To(Equal("https://custom.example.com"))
		})
	})

	Describe("output format", func() {
		BeforeEach(func() {
			os.Setenv("LIEPIN_USER_TOKEN", "t")
			DeferCleanup(func() { os.Unsetenv("LIEPIN_USER_TOKEN") })
		})

		It("defaults to pretty", func() {
			cfg, err := config.ResolveConfig("", "", "", 30.0)
			Expect(err).NotTo(HaveOccurred())
			Expect(cfg.Output).To(Equal("pretty"))
		})

		It("rejects invalid output", func() {
			_, err := config.ResolveConfig("", "", "yaml", 30.0)
			Expect(err).To(HaveOccurred())
		})
	})

	Describe("timeout", func() {
		BeforeEach(func() {
			os.Setenv("LIEPIN_USER_TOKEN", "t")
			DeferCleanup(func() { os.Unsetenv("LIEPIN_USER_TOKEN") })
		})

		It("defaults to 30s", func() {
			cfg, err := config.ResolveConfig("", "", "pretty", 0)
			Expect(err).NotTo(HaveOccurred())
			Expect(cfg.Timeout.Seconds()).To(Equal(30.0))
		})
	})
})

var _ = Describe("Config file path", func() {
	It("writes to XDG_CONFIG_HOME/liepin-cli/config.json", func() {
		dir, _ := os.MkdirTemp("", "cfg-path")
		defer os.RemoveAll(dir)
		os.Setenv("XDG_CONFIG_HOME", dir)
		defer os.Unsetenv("XDG_CONFIG_HOME")

		Expect(authstore.SaveToken("my-token")).To(Succeed())

		configFile := filepath.Join(dir, "liepin-cli", "config.json")
		data, err := os.ReadFile(configFile)
		Expect(err).NotTo(HaveOccurred())
		Expect(string(data)).NotTo(BeEmpty())
	})
})
