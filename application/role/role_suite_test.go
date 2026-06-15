package role_test

import (
	"testing"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

func TestRoleApplication(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Role Application Suite")
}
