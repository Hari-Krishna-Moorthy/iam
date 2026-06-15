package permission_test

import (
	"testing"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

func TestPermission(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Permission Domain Suite")
}
