package cart

import (
	"testing"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("Domain Test", func() {
	var cart *Cart

	Context("Simple Cart Tests", func() {
		BeforeEach(func() {
			cart = new(Cart)
		})
		It("Can add and get Items", func() {
			cartItem := Cartitem{MarketplaceCode: "code1", Qty: 5}
			cart.Cartitems = append(cart.Cartitems, cartItem)

			found, nr := cart.HasItem("code1", "")
			Expect(found).To(Equal(true))
			Expect(nr).To(Equal(1))
			Expect(cart.GetByLineNr(1)).To(Equal(&cartItem))
		})
	})
})

func TestDomain(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Cart Suite")
}
