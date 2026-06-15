package tenant_test

import (
	"testing"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

func TestTenantApplication(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Tenant Application Suite")
}
