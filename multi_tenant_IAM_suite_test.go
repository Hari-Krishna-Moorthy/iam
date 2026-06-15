package multi_tenant_IAM_test

import (
	"testing"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

func TestMultiTenantIAM(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "MultiTenantIAM Suite")
}
