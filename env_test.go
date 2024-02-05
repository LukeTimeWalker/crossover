package main_test

import (
	"os"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"

	main "github.com/LukeTimeWalker/crossover"
)

var _ = Describe("Env", func() {
	var environmenVariables []string

	Describe("Reading", func() {
		os.Setenv("TEST_ENV_2", "")
		os.Setenv("TEST_ENV_1", "1")

		Context("one env that is set", func() {
			environmenVariables = []string{"TEST_ENV_1"}
			envs := main.LoadEnvs(environmenVariables)

			It("should return a map with the length 1", func() {
				Expect(len(envs)).To(Equal(1))
			})
			It("should return a map with correspondig key value pair", func() {
				Expect(envs["TEST_ENV_1"]).To(Equal("1"))
			})
		})
		Context("one env that is set and one env that is empty", func() {
			environmenVariables = []string{"TEST_ENV_1", "TEST_ENV_2"}
			envs := main.LoadEnvs(environmenVariables)

			It("should return a map with the length 2", func() {
				Expect(len(envs)).To(Equal(2))
			})
			It("should return a map with one value that is empty", func() {
				Expect(envs["TEST_ENV_2"]).To(Equal(""))
			})
			It("should return a map with correspondig key value pair", func() {
				Expect(envs["TEST_ENV_1"]).To(Equal("1"))
			})
		})
		Context("that is not set", func() {
			environmenVariables = []string{"TEST_ENV_3"}
			envs := main.LoadEnvs(environmenVariables)

			It("should return a map without the corresponding key", func() {
				Expect(envs["TEST_ENV_3"]).To(BeNil())
			})
		})

	})

})
