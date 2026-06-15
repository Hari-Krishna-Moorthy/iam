package repositories_test

import (
	"testing"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

func TestRedisRepositories(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Redis Repositories Suite")
}
