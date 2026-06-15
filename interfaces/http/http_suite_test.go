package http_test

import (
	"testing"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

func TestHttpInterface(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "HTTP Interfaces (E2E) Suite")
}
