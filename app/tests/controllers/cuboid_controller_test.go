package controllers_test

import (
	. "cuboid-challenge/app/models"
	"cuboid-challenge/app/tests/testutils"
	"fmt"
	"net/http"
	"net/http/httptest"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("Cuboid Controller", func() {
	testutils.LoadEnv()
	testutils.ConnectDB()
	testutils.ClearDB()

	AfterEach(func() {
		testutils.ClearDB()
	})

	var recorder *httptest.ResponseRecorder
	var bag *Bag

	BeforeEach(func() {
		bag = &Bag{
			Title:  "A bag",
			Volume: 5,
			Cuboids: []Cuboid{
				{Width: 1, Height: 1, Depth: 1},
				{Width: 1, Height: 1, Depth: 2},
			},
		}
		testutils.AddRecords(bag)
	})

	Describe("List", func() {
		BeforeEach(func() {
			recorder = testutils.MockRequest(http.MethodGet, "/cuboids", nil)
		})

		It("Response HTTP status code 200", func() {
			Expect(recorder.Code).To(Equal(200))
		})

		It("Fetch all cuboids", func() {
			l, _ := testutils.DeserializeList(recorder.Body.String())
			Expect(len(l)).To(Equal(2))
			for _, dataMap := range l {
				Expect(dataMap["width"]).ToNot(BeNil())
				Expect(dataMap["height"]).ToNot(BeNil())
				Expect(dataMap["depth"]).ToNot(BeNil())
				Expect(dataMap["volume"]).ToNot(BeNil())
				Expect(dataMap["bagId"]).To(BeEquivalentTo(bag.ID))
			}
		})
	})

	Describe("Get", func() {
		var cuboidID uint

		JustBeforeEach(func() {
			recorder = testutils.MockRequest(http.MethodGet, "/cuboids/"+fmt.Sprintf("%v", cuboidID), nil)
		})

		Context("When the cuboid is present", func() {
			BeforeEach(func() {
				cuboidID = bag.Cuboids[0].ID
			})

			It("Response HTTP status code 200", func() {
				Expect(recorder.Code).To(Equal(200))
			})

			It("Get the cuboid", func() {
				dataMap, _ := testutils.Deserialize(recorder.Body.String())
				Expect(dataMap["id"]).To(BeEquivalentTo(bag.Cuboids[0].ID))
				Expect(dataMap["width"]).ToNot(BeNil())
				Expect(dataMap["height"]).ToNot(BeNil())
				Expect(dataMap["depth"]).ToNot(BeNil())
				Expect(dataMap["volume"]).ToNot(BeNil())
				Expect(dataMap["bagId"]).ToNot(BeNil())
			})
		})

		Context("When cuboid is not present", func() {
			BeforeEach(func() {
				cuboidID = 9999
			})

			It("Response HTTP status code 404", func() {
				Expect(recorder.Code).To(Equal(404))
			})
		})
	})

	Describe("Create", func() {
		cuboidPayload := map[string]interface{}{}

		BeforeEach(func() {
			cuboidPayload = map[string]interface{}{
				"width":  1,
				"height": 1,
				"depth":  1,
				"bagId":  bag.ID,
			}
		})

		JustBeforeEach(func() {
			body, _ := testutils.SerializeToString(cuboidPayload)
			recorder = testutils.MockRequest(http.MethodPost, "/cuboids", &body)
		})

		It("Response HTTP status code 201", func() {
			Expect(recorder.Code).To(Equal(201))
		})

		It("Returns the created cuboid", func() {
			dataMap, _ := testutils.Deserialize(recorder.Body.String())
			Expect(dataMap["width"]).ToNot(BeNil())
			Expect(dataMap["height"]).ToNot(BeNil())
			Expect(dataMap["depth"]).ToNot(BeNil())
			Expect(dataMap["volume"]).ToNot(BeNil())
			Expect(dataMap["bagId"]).To(BeEquivalentTo(bag.ID))
		})

		Context("When cuboid does not fit into the bag", func() {
			BeforeEach(func() {
				cuboidPayload["width"] = 3
			})

			It("Does not create the Cuboid", func() {
				Expect(recorder.Code).To(Equal(400))
				dataMap, _ := testutils.Deserialize(recorder.Body.String())
				Expect(dataMap["error"]).To(Equal("Insufficient capacity in bag"))
			})
		})

		Context("When the bag is disabled", func() {
			BeforeEach(func() {
				bag.SetDisabled(true)
				bag.Cuboids = []Cuboid{}
				testutils.UpdateRecords(&bag)
			})

			It("Does not admit new cuboids", func() {
				Expect(recorder.Code).To(Equal(400))
				dataMap, _ := testutils.Deserialize(recorder.Body.String())
				Expect(dataMap["error"]).To(Equal("Bag is disabled"))
			})
		})
	})

	// DO NOT modify the tests ABOVE
	// IMPLEMENT the tests BELOW

	Describe("Update", func() {
		bag = &Bag{
			Title:  "A bag",
			Volume: 5,
			Cuboids: []Cuboid{
				{Width: 1, Height: 1, Depth: 1},
				{Width: 1, Height: 1, Depth: 2},
			},
		}
		testutils.AddRecords(bag)

		cuboidPayload := map[string]interface{}{}

		BeforeEach(func() {
			cuboidPayload = map[string]interface{}{
				"id":     bag.Cuboids[0].ID,
				"width":  1,
				"height": 1,
				"depth":  1,
				"bagId":  bag.ID,
			}
		})

		JustBeforeEach(func() {
			body, _ := testutils.SerializeToString(cuboidPayload)
			recorder = testutils.MockRequest(http.MethodPut, "/cuboids", &body)
		})

		It("Response HTTP status code 200", func() {
			Expect(recorder.Code).To(Equal(200))
		})

		It("Returns the updated cuboid", func() {
			dataMap, _ := testutils.Deserialize(recorder.Body.String())
			Expect(dataMap["width"]).ToNot(BeNil())
			Expect(dataMap["height"]).ToNot(BeNil())
			Expect(dataMap["depth"]).ToNot(BeNil())
			Expect(dataMap["volume"]).ToNot(BeNil())
			Expect(dataMap["bagId"]).To(BeEquivalentTo(bag.ID))
		})

		Context("When cuboid does not fit into the bag", func() {
			BeforeEach(func() {
				cuboidPayload["width"] = 3
			})

			It("Response HTTP status code 400", func() {
				Expect(recorder.Code).To(Equal(400))
			})

			It("Response a JSON with error message 'Insufficient capacity in bag'", func() {
				dataMap, _ := testutils.Deserialize(recorder.Body.String())
				Expect(dataMap["error"]).To(Equal("Insufficient capacity in bag"))
			})
		})

		Context("When cuboid is not present", func() {
			BeforeEach(func() {
				cuboidPayload["id"] = 9999
			})
			It("Response HTTP status code 404", func() {
				Expect(recorder.Code).To(Equal(404))
			})
		})
	})

	Describe("Delete", func() {
		var cuboidID uint

		JustBeforeEach(func() {
			recorder = testutils.MockRequest(http.MethodDelete, "/cuboids/"+fmt.Sprintf("%v", cuboidID), nil)
		})

		Context("When the cuboid is present", func() {
			BeforeEach(func() {
				cuboidID = bag.Cuboids[0].ID
			})

			It("Response HTTP status code 200", func() {
				Expect(recorder.Code).To(Equal(200))
			})

			It("Remove the cuboid", func() {
				dataMap, _ := testutils.Deserialize(recorder.Body.String())
				Expect(dataMap["id"]).To(BeEquivalentTo(bag.Cuboids[0].ID))
			})
		})

		Context("When cuboid is not present", func() {
			BeforeEach(func() {
				cuboidID = 9999
			})

			It("Response HTTP status code 404", func() {
				Expect(recorder.Code).To(Equal(404))
			})
		})
	})
})
