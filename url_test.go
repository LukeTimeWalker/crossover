package main_test

import (
	main "github.com/LukeTimeWalker/crossover"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

var _ = Describe("LoadUrlParams", func() {
	It("should parse URL parameters correctly", func() {
		u := "https://example.com/?key1=value1&key2=value2&key3=value3"
		expectedParams := map[string]interface{}{
			"key1": "value1",
			"key2": "value2",
			"key3": "value3",
		}

		params, err := main.LoadUrlParams(u)
		Expect(err).To(BeNil())
		Expect(params).To(Equal(expectedParams))
	})

	It("should handle URL params parsing errors", func() {
		u := "https://example.com/?key1=value1&key2=value2&key3=value3%ZZ"

		params, err := main.LoadUrlParams(u)
		Expect(err).NotTo(BeNil())
		Expect(params).To(BeNil())
	})
})
