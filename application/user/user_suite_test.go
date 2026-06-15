package user_test

import (
	"testing"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

func TestUserApplication(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "User Application Suite")
}
