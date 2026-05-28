package authstore_test

import (
	"testing"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

func TestAuthstore(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Authstore Suite")
}
