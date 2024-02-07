package main_test

import (
	main "github.com/LukeTimeWalker/crossover"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

var _ = Describe("Template", func() {

	testTemplate := `
		{
			"title": "{{ .env.TITLE }}",
			"description": "Body: {{ .body.test }}  \n{{ len .body.array }}  \n{{ .url.url }}"
		}
	`

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
		template, _ := main.CreateTemplate("testTemplate", testTemplate)
		values := main.Variables{
			"env":  map[string]interface{}{"TITLE": "testTitle"},
			"body": map[string]interface{}{"test": "bodyContent", "array": []string{"1", "2", "3"}},
			"url":  map[string]interface{}{"url": "param"}}

		Context("with all values available", func() {
			result, err := main.ExecuteTemplate(template, values)
			It("should not throw an error", func() {
				Expect(err).To(BeNil())
			})
			It("should return the string with interpolated values", func() {
				Expect(result).To(Equal(`
		{
			"title": "testTitle",
			"description": "Body: bodyContent  \n3  \nparam"
		}
	`))
			})
		})
	})
})
