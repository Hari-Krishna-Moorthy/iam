package auth_test

import (
	"testing"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

func TestAuthApplication(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Auth Application Suite")
}
