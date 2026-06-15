package permission_test

import (
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"

	"github.com/Hari-Krishna-Moorthy/multi-tenant-IAM/domain/permission"
)

var _ = Describe("Permission", func() {
	Context("When creating a new permission", func() {
		It("should correctly format the permission string", func() {
			p := permission.New("global", "billing", "read")
			Expect(p.String()).To(Equal("global:billing:read"))
		})
	})

	Context("When parsing a valid permission string", func() {
		It("should return a Permission object", func() {
			p, err := permission.Parse("tenant:inventory:write")
			Expect(err).NotTo(HaveOccurred())
			Expect(p.Scope).To(Equal("tenant"))
			Expect(p.ServiceName).To(Equal("inventory"))
			Expect(p.Action).To(Equal("write"))
		})
	})

	Context("When parsing an invalid permission string", func() {
		It("should return an error", func() {
			_, err := permission.Parse("invalid:format")
			Expect(err).To(HaveOccurred())
			Expect(err.Error()).To(ContainSubstring("invalid permission format"))
		})
	})
})
