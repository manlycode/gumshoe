package matchers_test

import (
	. "github.com/pivotal/gumshoe/repos/ginkgo"
	. "github.com/pivotal/gumshoe/repos/gomega"
	. "github.com/pivotal/gumshoe/repos/gomega/matchers"
)

var _ = Describe("BeTrue", func() {
	It("should handle true and false correctly", func() {
		Ω(true).Should(BeTrue())
		Ω(false).ShouldNot(BeTrue())
	})

	It("should only support booleans", func() {
		success, _, err := (&BeTrueMatcher{}).Match("foo")
		Ω(success).Should(BeFalse())
		Ω(err).Should(HaveOccured())
	})
})
