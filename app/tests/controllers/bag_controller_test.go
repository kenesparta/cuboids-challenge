package controllers_test

import (
	"cuboid-challenge/app/tests/factories"
	"cuboid-challenge/app/tests/testutils"
	"fmt"
	"net/http"
	"net/http/httptest"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("Bag Controller", func() {
	testutils.LoadEnv()
	testutils.ConnectDB()
	testutils.ClearDB()

	AfterEach(func() {
		testutils.ClearDB()
	})

	var recorder *httptest.ResponseRecorder

	Describe("List", func() {
		BeforeEach(func() {
			for i := 0; i < 3; i++ {
				testutils.AddRecords(factories.Bag())
			}
			recorder = testutils.MockRequest(http.MethodGet, "/bags", nil)
		})

		It("Response HTTP status code 200", func() {
			Expect(recorder.Code).To(Equal(200))
		})

		It("Fetch all bags", func() {
			l, _ := testutils.DeserializeList(recorder.Body.String())
			Expect(len(l)).To(Equal(3))
			for _, m := range l {
				Expect(m["title"]).ToNot(BeNil())
				Expect(m["volume"]).ToNot(BeNil())
				Expect(m["cuboids"]).ToNot(BeNil())
			}
		})
	})

	Describe("Get", func() {
		var bagID uint
		bag := factories.Bag()

		BeforeEach(func() {
			testutils.AddRecords(&bag)
		})

		JustBeforeEach(func() {
			recorder = testutils.MockRequest(http.MethodGet, "/bags/"+fmt.Sprintf("%v", bagID), nil)
		})

		Context("When the bag is present", func() {
			BeforeEach(func() {
				bagID = bag.ID
			})

			It("Response HTTP status code 200", func() {
				Expect(recorder.Code).To(Equal(200))
			})

			It("Get the bag", func() {
				m, _ := testutils.Deserialize(recorder.Body.String())
				Expect(m["id"]).To(BeEquivalentTo(bag.ID))
				Expect(m["title"]).ToNot(BeNil())
				Expect(m["volume"]).ToNot(BeNil())
				Expect(m["cuboids"]).ToNot(BeNil())
			})
		})

		Context("When the bag is not present", func() {
			BeforeEach(func() {
				bagID = 99999
			})

			It("Response HTTP status code 404", func() {
				Expect(recorder.Code).To(Equal(404))
			})
		})
	})

	Describe("Create", func() {
		bagPayload := map[string]interface{}{}

		JustBeforeEach(func() {
			body, _ := testutils.SerializeToString(bagPayload)
			recorder = testutils.MockRequest(http.MethodPost, "/bags", &body)
		})

		Context("When the payload is valid", func() {
			BeforeEach(func() {
				bagPayload["title"] = "My title"
				bagPayload["volume"] = 10
			})

			It("Response HTTP status code 201", func() {
				Expect(recorder.Code).To(Equal(201))
			})

			It("Create the bag", func() {
				m, _ := testutils.Deserialize(recorder.Body.String())
				Expect(m["id"]).ToNot(BeNil())
				Expect(m["title"]).ToNot(BeNil())
				Expect(m["volume"]).ToNot(BeNil())
				Expect(m["cuboids"]).To(BeEmpty())
			})
		})

		Context("When the payload is invalid", func() {
			BeforeEach(func() {
				bagPayload["title"] = ""
				bagPayload["volume"] = 0
			})

			It("Response HTTP status code 400", func() {
				Expect(recorder.Code).To(Equal(400))
			})
		})
	})
})
