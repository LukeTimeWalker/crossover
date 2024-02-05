package main_test

import (
	main "github.com/LukeTimeWalker/crossover"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

var _ = Describe("Template", func() {

	Describe("Creating a template", func() {
		Context("that is valid", func() {
			template, err := main.CreateTemplate("testTemplate", "'Bearer {{ .env.TEST_TOKEN }}'")
			It("should return no error", func() {
				Expect(err).To(BeNil())
			})
			It("should return template", func() {
				Expect(template).NotTo(BeNil())
			})
		})
		Context("that is not valid", func() {
			template, err := main.CreateTemplate("testTemplate", "'Bearer {{ .env.TEST_TOKEN }'")
			It("should return an error", func() {
				Expect(err).NotTo(BeNil())
			})
			It("should return template", func() {
				Expect(template).To(BeNil())
			})
		})
	})

	Describe("Executing a template", func() {

	})
})
