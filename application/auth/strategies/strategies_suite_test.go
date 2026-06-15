package strategies_test

import (
	"testing"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

func TestAuthStrategies(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Auth Strategies Suite")
}
