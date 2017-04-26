package domain

import (
	"testing"

	"fmt"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("Domain Test", func() {
	var cart Cart

	Context("Simple Cart Tests", func() {
		BeforeEach(func() {
			//TODO!
			cart := new(Cart)
			fmt.Println(cart)
		})
		It("Can add Items", func() {
			cartItem := Cartitem{"code1", 5, 2.5}
			cart.Add(cartItem)

			Expect(cart.GetLine(1)).To(Equal(cartItem))
		})

		It("Can add and update Items by code", func() {
			cart.AddOrUpdateByCode("code2", 2, 23)
			fmt.Println(cart)
			Expect(cart.GetLine(2).Qty).To(Equal(2))

			cart.AddOrUpdateByCode("code2", 2, 23)
			Expect(cart.GetLine(2).Qty).To(Equal(4))
		})
	})
})

func TestDomain(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Cart Suite")
}
