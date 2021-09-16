// Code generated by github.com/99designs/gqlgen, DO NOT EDIT.

//+build !graphql

package graphql

import (
	"context"

	"flamingo.me/flamingo-commerce/v3/cart/domain/cart"
	"flamingo.me/flamingo-commerce/v3/cart/domain/validation"
	"flamingo.me/flamingo-commerce/v3/cart/interfaces/controller/forms"
	graphql1 "flamingo.me/flamingo-commerce/v3/cart/interfaces/graphql"
	"flamingo.me/flamingo-commerce/v3/cart/interfaces/graphql/dto"
	domain3 "flamingo.me/flamingo-commerce/v3/category/domain"
	graphql7 "flamingo.me/flamingo-commerce/v3/category/interfaces/graphql"
	"flamingo.me/flamingo-commerce/v3/category/interfaces/graphql/categorydto"
	graphql5 "flamingo.me/flamingo-commerce/v3/checkout/interfaces/graphql"
	dto1 "flamingo.me/flamingo-commerce/v3/checkout/interfaces/graphql/dto"
	graphql6 "flamingo.me/flamingo-commerce/v3/customer/interfaces/graphql"
	"flamingo.me/flamingo-commerce/v3/customer/interfaces/graphql/dtocustomer"
	domain1 "flamingo.me/flamingo-commerce/v3/price/domain"
	"flamingo.me/flamingo-commerce/v3/product/domain"
	graphql2 "flamingo.me/flamingo-commerce/v3/product/interfaces/graphql"
	graphqlproductdto "flamingo.me/flamingo-commerce/v3/product/interfaces/graphql/product/dto"
	domain2 "flamingo.me/flamingo-commerce/v3/search/domain"
	graphql3 "flamingo.me/flamingo-commerce/v3/search/interfaces/graphql"
	"flamingo.me/flamingo-commerce/v3/search/interfaces/graphql/searchdto"
	graphql4 "flamingo.me/graphql"
)

var _ ResolverRoot = new(rootResolver)

type rootResolver struct {
	rootResolverCommerce_Cart_AdditionalData          *rootResolverCommerce_Cart_AdditionalData
	rootResolverCommerce_Cart_Cart                    *rootResolverCommerce_Cart_Cart
	rootResolverCommerce_Cart_DefaultPaymentSelection *rootResolverCommerce_Cart_DefaultPaymentSelection
	rootResolverCommerce_Cart_DeliveryInfo            *rootResolverCommerce_Cart_DeliveryInfo
	rootResolverCommerce_Cart_Item                    *rootResolverCommerce_Cart_Item
	rootResolverCommerce_Cart_ShippingItem            *rootResolverCommerce_Cart_ShippingItem
	rootResolverCommerce_Product_PriceInfo            *rootResolverCommerce_Product_PriceInfo
	rootResolverCommerce_Search_Meta                  *rootResolverCommerce_Search_Meta
	rootResolverMutation                              *rootResolverMutation
	rootResolverQuery                                 *rootResolverQuery
}

func (r *rootResolver) Inject(
	rootResolverCommerce_Cart_AdditionalData *rootResolverCommerce_Cart_AdditionalData,
	rootResolverCommerce_Cart_Cart *rootResolverCommerce_Cart_Cart,
	rootResolverCommerce_Cart_DefaultPaymentSelection *rootResolverCommerce_Cart_DefaultPaymentSelection,
	rootResolverCommerce_Cart_DeliveryInfo *rootResolverCommerce_Cart_DeliveryInfo,
	rootResolverCommerce_Cart_Item *rootResolverCommerce_Cart_Item,
	rootResolverCommerce_Cart_ShippingItem *rootResolverCommerce_Cart_ShippingItem,
	rootResolverCommerce_Product_PriceInfo *rootResolverCommerce_Product_PriceInfo,
	rootResolverCommerce_Search_Meta *rootResolverCommerce_Search_Meta,
	rootResolverMutation *rootResolverMutation,
	rootResolverQuery *rootResolverQuery,
) {
	r.rootResolverCommerce_Cart_AdditionalData = rootResolverCommerce_Cart_AdditionalData
	r.rootResolverCommerce_Cart_Cart = rootResolverCommerce_Cart_Cart
	r.rootResolverCommerce_Cart_DefaultPaymentSelection = rootResolverCommerce_Cart_DefaultPaymentSelection
	r.rootResolverCommerce_Cart_DeliveryInfo = rootResolverCommerce_Cart_DeliveryInfo
	r.rootResolverCommerce_Cart_Item = rootResolverCommerce_Cart_Item
	r.rootResolverCommerce_Cart_ShippingItem = rootResolverCommerce_Cart_ShippingItem
	r.rootResolverCommerce_Product_PriceInfo = rootResolverCommerce_Product_PriceInfo
	r.rootResolverCommerce_Search_Meta = rootResolverCommerce_Search_Meta
	r.rootResolverMutation = rootResolverMutation
	r.rootResolverQuery = rootResolverQuery
}

func (r *rootResolver) Commerce_Cart_AdditionalData() Commerce_Cart_AdditionalDataResolver {
	return r.rootResolverCommerce_Cart_AdditionalData
}
func (r *rootResolver) Commerce_Cart_Cart() Commerce_Cart_CartResolver {
	return r.rootResolverCommerce_Cart_Cart
}
func (r *rootResolver) Commerce_Cart_DefaultPaymentSelection() Commerce_Cart_DefaultPaymentSelectionResolver {
	return r.rootResolverCommerce_Cart_DefaultPaymentSelection
}
func (r *rootResolver) Commerce_Cart_DeliveryInfo() Commerce_Cart_DeliveryInfoResolver {
	return r.rootResolverCommerce_Cart_DeliveryInfo
}
func (r *rootResolver) Commerce_Cart_Item() Commerce_Cart_ItemResolver {
	return r.rootResolverCommerce_Cart_Item
}
func (r *rootResolver) Commerce_Cart_ShippingItem() Commerce_Cart_ShippingItemResolver {
	return r.rootResolverCommerce_Cart_ShippingItem
}
func (r *rootResolver) Commerce_Product_PriceInfo() Commerce_Product_PriceInfoResolver {
	return r.rootResolverCommerce_Product_PriceInfo
}
func (r *rootResolver) Commerce_Search_Meta() Commerce_Search_MetaResolver {
	return r.rootResolverCommerce_Search_Meta
}
func (r *rootResolver) Mutation() MutationResolver {
	return r.rootResolverMutation
}
func (r *rootResolver) Query() QueryResolver {
	return r.rootResolverQuery
}

type rootResolverCommerce_Cart_AdditionalData struct {
	resolveCustomAttributes func(ctx context.Context, obj *cart.AdditionalData) (*dto.CustomAttributes, error)
}

func (r *rootResolverCommerce_Cart_AdditionalData) Inject(
	commerce_Cart_AdditionalDataCustomAttributes *graphql1.CommerceCartAdditionalDataResolver,
) {
	r.resolveCustomAttributes = commerce_Cart_AdditionalDataCustomAttributes.CustomAttributes
}

func (r *rootResolverCommerce_Cart_AdditionalData) CustomAttributes(ctx context.Context, obj *cart.AdditionalData) (*dto.CustomAttributes, error) {
	return r.resolveCustomAttributes(ctx, obj)
}

type rootResolverCommerce_Cart_Cart struct {
	resolveGetDeliveryByCode func(ctx context.Context, obj *cart.Cart, deliveryCode string) (*cart.Delivery, error)
}

func (r *rootResolverCommerce_Cart_Cart) Inject(
	commerce_Cart_CartGetDeliveryByCode *graphql1.Resolver,
) {
	r.resolveGetDeliveryByCode = commerce_Cart_CartGetDeliveryByCode.GetDeliveryByCodeWithoutBool
}

func (r *rootResolverCommerce_Cart_Cart) GetDeliveryByCode(ctx context.Context, obj *cart.Cart, deliveryCode string) (*cart.Delivery, error) {
	return r.resolveGetDeliveryByCode(ctx, obj, deliveryCode)
}

type rootResolverCommerce_Cart_DefaultPaymentSelection struct {
	resolveCartSplit func(ctx context.Context, obj *cart.DefaultPaymentSelection) ([]*dto.PaymentSelectionSplit, error)
}

func (r *rootResolverCommerce_Cart_DefaultPaymentSelection) Inject(
	commerce_Cart_DefaultPaymentSelectionCartSplit *graphql1.CommerceCartQueryResolver,
) {
	r.resolveCartSplit = commerce_Cart_DefaultPaymentSelectionCartSplit.CartSplit
}

func (r *rootResolverCommerce_Cart_DefaultPaymentSelection) CartSplit(ctx context.Context, obj *cart.DefaultPaymentSelection) ([]*dto.PaymentSelectionSplit, error) {
	return r.resolveCartSplit(ctx, obj)
}

type rootResolverCommerce_Cart_DeliveryInfo struct {
	resolveAdditionalData func(ctx context.Context, obj *cart.DeliveryInfo) (*dto.CustomAttributes, error)
}

func (r *rootResolverCommerce_Cart_DeliveryInfo) Inject(
	commerce_Cart_DeliveryInfoAdditionalData *graphql1.CommerceCartDeliveryInfoResolver,
) {
	r.resolveAdditionalData = commerce_Cart_DeliveryInfoAdditionalData.AdditionalData
}

func (r *rootResolverCommerce_Cart_DeliveryInfo) AdditionalData(ctx context.Context, obj *cart.DeliveryInfo) (*dto.CustomAttributes, error) {
	return r.resolveAdditionalData(ctx, obj)
}

type rootResolverCommerce_Cart_Item struct {
	resolveAppliedDiscounts func(ctx context.Context, obj *cart.Item) (*dto.CartAppliedDiscounts, error)
}

func (r *rootResolverCommerce_Cart_Item) Inject(
	commerce_Cart_ItemAppliedDiscounts *dto.CartAppliedDiscountsResolver,
) {
	r.resolveAppliedDiscounts = commerce_Cart_ItemAppliedDiscounts.ForItem
}

func (r *rootResolverCommerce_Cart_Item) AppliedDiscounts(ctx context.Context, obj *cart.Item) (*dto.CartAppliedDiscounts, error) {
	return r.resolveAppliedDiscounts(ctx, obj)
}

type rootResolverCommerce_Cart_ShippingItem struct {
	resolveAppliedDiscounts func(ctx context.Context, obj *cart.ShippingItem) (*dto.CartAppliedDiscounts, error)
}

func (r *rootResolverCommerce_Cart_ShippingItem) Inject(
	commerce_Cart_ShippingItemAppliedDiscounts *dto.CartAppliedDiscountsResolver,
) {
	r.resolveAppliedDiscounts = commerce_Cart_ShippingItemAppliedDiscounts.ForShippingItem
}

func (r *rootResolverCommerce_Cart_ShippingItem) AppliedDiscounts(ctx context.Context, obj *cart.ShippingItem) (*dto.CartAppliedDiscounts, error) {
	return r.resolveAppliedDiscounts(ctx, obj)
}

type rootResolverCommerce_Product_PriceInfo struct {
	resolveActiveBase func(ctx context.Context, obj *domain.PriceInfo) (*domain1.Price, error)
}

func (r *rootResolverCommerce_Product_PriceInfo) Inject(
	commerce_Product_PriceInfoActiveBase *graphql2.CommerceProductQueryResolver,
) {
	r.resolveActiveBase = commerce_Product_PriceInfoActiveBase.ActiveBase
}

func (r *rootResolverCommerce_Product_PriceInfo) ActiveBase(ctx context.Context, obj *domain.PriceInfo) (*domain1.Price, error) {
	return r.resolveActiveBase(ctx, obj)
}

type rootResolverCommerce_Search_Meta struct {
	resolveSortOptions func(ctx context.Context, obj *domain2.SearchMeta) ([]*searchdto.CommerceSearchSortOption, error)
}

func (r *rootResolverCommerce_Search_Meta) Inject(
	commerce_Search_MetaSortOptions *graphql3.CommerceSearchQueryResolver,
) {
	r.resolveSortOptions = commerce_Search_MetaSortOptions.SortOptions
}

func (r *rootResolverCommerce_Search_Meta) SortOptions(ctx context.Context, obj *domain2.SearchMeta) ([]*searchdto.CommerceSearchSortOption, error) {
	return r.resolveSortOptions(ctx, obj)
}

type rootResolverMutation struct {
	resolveFlamingo                                   func(ctx context.Context) (*string, error)
	resolveCommerceCartAddToCart                      func(ctx context.Context, marketplaceCode string, qty int, deliveryCode string) (*dto.DecoratedCart, error)
	resolveCommerceCartDeleteCartDelivery             func(ctx context.Context, deliveryCode string) (*dto.DecoratedCart, error)
	resolveCommerceCartDeleteItem                     func(ctx context.Context, itemID string, deliveryCode string) (*dto.DecoratedCart, error)
	resolveCommerceCartUpdateItemQty                  func(ctx context.Context, itemID string, deliveryCode string, qty int) (*dto.DecoratedCart, error)
	resolveCommerceCartUpdateBillingAddress           func(ctx context.Context, addressForm *forms.AddressForm) (*dto.BillingAddressForm, error)
	resolveCommerceCartUpdateSelectedPayment          func(ctx context.Context, gateway string, method string) (*dto.SelectedPaymentResult, error)
	resolveCommerceCartApplyCouponCodeOrGiftCard      func(ctx context.Context, code string) (*dto.DecoratedCart, error)
	resolveCommerceCartRemoveGiftCard                 func(ctx context.Context, giftCardCode string) (*dto.DecoratedCart, error)
	resolveCommerceCartRemoveCouponCode               func(ctx context.Context, couponCode string) (*dto.DecoratedCart, error)
	resolveCommerceCartUpdateDeliveryAddresses        func(ctx context.Context, deliveryAdresses []*forms.DeliveryForm) ([]*dto.DeliveryAddressForm, error)
	resolveCommerceCartUpdateDeliveryShippingOptions  func(ctx context.Context, shippingOptions []*dto.DeliveryShippingOption) (*dto.UpdateShippingOptionsResult, error)
	resolveCommerceCartClean                          func(ctx context.Context) (bool, error)
	resolveCommerceCartUpdateAdditionalData           func(ctx context.Context, additionalData []*dto.KeyValue) (*dto.DecoratedCart, error)
	resolveCommerceCartUpdateDeliveriesAdditionalData func(ctx context.Context, data []*dto.DeliveryAdditionalData) (*dto.DecoratedCart, error)
	resolveCommerceCheckoutStartPlaceOrder            func(ctx context.Context, returnURL string) (*dto1.StartPlaceOrderResult, error)
	resolveCommerceCheckoutCancelPlaceOrder           func(ctx context.Context) (bool, error)
	resolveCommerceCheckoutClearPlaceOrder            func(ctx context.Context) (bool, error)
	resolveCommerceCheckoutRefreshPlaceOrder          func(ctx context.Context) (*dto1.PlaceOrderContext, error)
	resolveCommerceCheckoutRefreshPlaceOrderBlocking  func(ctx context.Context) (*dto1.PlaceOrderContext, error)
}

func (r *rootResolverMutation) Inject(
	mutationFlamingo *graphql4.FlamingoQueryResolver,
	mutationCommerceCartAddToCart *graphql1.CommerceCartMutationResolver,
	mutationCommerceCartDeleteCartDelivery *graphql1.CommerceCartMutationResolver,
	mutationCommerceCartDeleteItem *graphql1.CommerceCartMutationResolver,
	mutationCommerceCartUpdateItemQty *graphql1.CommerceCartMutationResolver,
	mutationCommerceCartUpdateBillingAddress *graphql1.CommerceCartMutationResolver,
	mutationCommerceCartUpdateSelectedPayment *graphql1.CommerceCartMutationResolver,
	mutationCommerceCartApplyCouponCodeOrGiftCard *graphql1.CommerceCartMutationResolver,
	mutationCommerceCartRemoveGiftCard *graphql1.CommerceCartMutationResolver,
	mutationCommerceCartRemoveCouponCode *graphql1.CommerceCartMutationResolver,
	mutationCommerceCartUpdateDeliveryAddresses *graphql1.CommerceCartMutationResolver,
	mutationCommerceCartUpdateDeliveryShippingOptions *graphql1.CommerceCartMutationResolver,
	mutationCommerceCartClean *graphql1.CommerceCartMutationResolver,
	mutationCommerceCartUpdateAdditionalData *graphql1.CommerceCartMutationResolver,
	mutationCommerceCartUpdateDeliveriesAdditionalData *graphql1.CommerceCartMutationResolver,
	mutationCommerceCheckoutStartPlaceOrder *graphql5.CommerceCheckoutMutationResolver,
	mutationCommerceCheckoutCancelPlaceOrder *graphql5.CommerceCheckoutMutationResolver,
	mutationCommerceCheckoutClearPlaceOrder *graphql5.CommerceCheckoutMutationResolver,
	mutationCommerceCheckoutRefreshPlaceOrder *graphql5.CommerceCheckoutMutationResolver,
	mutationCommerceCheckoutRefreshPlaceOrderBlocking *graphql5.CommerceCheckoutMutationResolver,
) {
	r.resolveFlamingo = mutationFlamingo.Flamingo
	r.resolveCommerceCartAddToCart = mutationCommerceCartAddToCart.CommerceAddToCart
	r.resolveCommerceCartDeleteCartDelivery = mutationCommerceCartDeleteCartDelivery.CommerceDeleteCartDelivery
	r.resolveCommerceCartDeleteItem = mutationCommerceCartDeleteItem.CommerceDeleteItem
	r.resolveCommerceCartUpdateItemQty = mutationCommerceCartUpdateItemQty.CommerceUpdateItemQty
	r.resolveCommerceCartUpdateBillingAddress = mutationCommerceCartUpdateBillingAddress.CommerceCartUpdateBillingAddress
	r.resolveCommerceCartUpdateSelectedPayment = mutationCommerceCartUpdateSelectedPayment.CommerceCartUpdateSelectedPayment
	r.resolveCommerceCartApplyCouponCodeOrGiftCard = mutationCommerceCartApplyCouponCodeOrGiftCard.CommerceCartApplyCouponCodeOrGiftCard
	r.resolveCommerceCartRemoveGiftCard = mutationCommerceCartRemoveGiftCard.CommerceCartRemoveGiftCard
	r.resolveCommerceCartRemoveCouponCode = mutationCommerceCartRemoveCouponCode.CommerceCartRemoveCouponCode
	r.resolveCommerceCartUpdateDeliveryAddresses = mutationCommerceCartUpdateDeliveryAddresses.CommerceCartUpdateDeliveryAddresses
	r.resolveCommerceCartUpdateDeliveryShippingOptions = mutationCommerceCartUpdateDeliveryShippingOptions.CommerceCartUpdateDeliveryShippingOptions
	r.resolveCommerceCartClean = mutationCommerceCartClean.CartClean
	r.resolveCommerceCartUpdateAdditionalData = mutationCommerceCartUpdateAdditionalData.UpdateAdditionalData
	r.resolveCommerceCartUpdateDeliveriesAdditionalData = mutationCommerceCartUpdateDeliveriesAdditionalData.UpdateDeliveriesAdditionalData
	r.resolveCommerceCheckoutStartPlaceOrder = mutationCommerceCheckoutStartPlaceOrder.CommerceCheckoutStartPlaceOrder
	r.resolveCommerceCheckoutCancelPlaceOrder = mutationCommerceCheckoutCancelPlaceOrder.CommerceCheckoutCancelPlaceOrder
	r.resolveCommerceCheckoutClearPlaceOrder = mutationCommerceCheckoutClearPlaceOrder.CommerceCheckoutClearPlaceOrder
	r.resolveCommerceCheckoutRefreshPlaceOrder = mutationCommerceCheckoutRefreshPlaceOrder.CommerceCheckoutRefreshPlaceOrder
	r.resolveCommerceCheckoutRefreshPlaceOrderBlocking = mutationCommerceCheckoutRefreshPlaceOrderBlocking.CommerceCheckoutRefreshPlaceOrderBlocking
}

func (r *rootResolverMutation) Flamingo(ctx context.Context) (*string, error) {
	return r.resolveFlamingo(ctx)
}
func (r *rootResolverMutation) CommerceCartAddToCart(ctx context.Context, marketplaceCode string, qty int, deliveryCode string) (*dto.DecoratedCart, error) {
	return r.resolveCommerceCartAddToCart(ctx, marketplaceCode, qty, deliveryCode)
}
func (r *rootResolverMutation) CommerceCartDeleteCartDelivery(ctx context.Context, deliveryCode string) (*dto.DecoratedCart, error) {
	return r.resolveCommerceCartDeleteCartDelivery(ctx, deliveryCode)
}
func (r *rootResolverMutation) CommerceCartDeleteItem(ctx context.Context, itemID string, deliveryCode string) (*dto.DecoratedCart, error) {
	return r.resolveCommerceCartDeleteItem(ctx, itemID, deliveryCode)
}
func (r *rootResolverMutation) CommerceCartUpdateItemQty(ctx context.Context, itemID string, deliveryCode string, qty int) (*dto.DecoratedCart, error) {
	return r.resolveCommerceCartUpdateItemQty(ctx, itemID, deliveryCode, qty)
}
func (r *rootResolverMutation) CommerceCartUpdateBillingAddress(ctx context.Context, addressForm *forms.AddressForm) (*dto.BillingAddressForm, error) {
	return r.resolveCommerceCartUpdateBillingAddress(ctx, addressForm)
}
func (r *rootResolverMutation) CommerceCartUpdateSelectedPayment(ctx context.Context, gateway string, method string) (*dto.SelectedPaymentResult, error) {
	return r.resolveCommerceCartUpdateSelectedPayment(ctx, gateway, method)
}
func (r *rootResolverMutation) CommerceCartApplyCouponCodeOrGiftCard(ctx context.Context, code string) (*dto.DecoratedCart, error) {
	return r.resolveCommerceCartApplyCouponCodeOrGiftCard(ctx, code)
}
func (r *rootResolverMutation) CommerceCartRemoveGiftCard(ctx context.Context, giftCardCode string) (*dto.DecoratedCart, error) {
	return r.resolveCommerceCartRemoveGiftCard(ctx, giftCardCode)
}
func (r *rootResolverMutation) CommerceCartRemoveCouponCode(ctx context.Context, couponCode string) (*dto.DecoratedCart, error) {
	return r.resolveCommerceCartRemoveCouponCode(ctx, couponCode)
}
func (r *rootResolverMutation) CommerceCartUpdateDeliveryAddresses(ctx context.Context, deliveryAdresses []*forms.DeliveryForm) ([]*dto.DeliveryAddressForm, error) {
	return r.resolveCommerceCartUpdateDeliveryAddresses(ctx, deliveryAdresses)
}
func (r *rootResolverMutation) CommerceCartUpdateDeliveryShippingOptions(ctx context.Context, shippingOptions []*dto.DeliveryShippingOption) (*dto.UpdateShippingOptionsResult, error) {
	return r.resolveCommerceCartUpdateDeliveryShippingOptions(ctx, shippingOptions)
}
func (r *rootResolverMutation) CommerceCartClean(ctx context.Context) (bool, error) {
	return r.resolveCommerceCartClean(ctx)
}
func (r *rootResolverMutation) CommerceCartUpdateAdditionalData(ctx context.Context, additionalData []*dto.KeyValue) (*dto.DecoratedCart, error) {
	return r.resolveCommerceCartUpdateAdditionalData(ctx, additionalData)
}
func (r *rootResolverMutation) CommerceCartUpdateDeliveriesAdditionalData(ctx context.Context, data []*dto.DeliveryAdditionalData) (*dto.DecoratedCart, error) {
	return r.resolveCommerceCartUpdateDeliveriesAdditionalData(ctx, data)
}
func (r *rootResolverMutation) CommerceCheckoutStartPlaceOrder(ctx context.Context, returnURL string) (*dto1.StartPlaceOrderResult, error) {
	return r.resolveCommerceCheckoutStartPlaceOrder(ctx, returnURL)
}
func (r *rootResolverMutation) CommerceCheckoutCancelPlaceOrder(ctx context.Context) (bool, error) {
	return r.resolveCommerceCheckoutCancelPlaceOrder(ctx)
}
func (r *rootResolverMutation) CommerceCheckoutClearPlaceOrder(ctx context.Context) (bool, error) {
	return r.resolveCommerceCheckoutClearPlaceOrder(ctx)
}
func (r *rootResolverMutation) CommerceCheckoutRefreshPlaceOrder(ctx context.Context) (*dto1.PlaceOrderContext, error) {
	return r.resolveCommerceCheckoutRefreshPlaceOrder(ctx)
}
func (r *rootResolverMutation) CommerceCheckoutRefreshPlaceOrderBlocking(ctx context.Context) (*dto1.PlaceOrderContext, error) {
	return r.resolveCommerceCheckoutRefreshPlaceOrderBlocking(ctx)
}

type rootResolverQuery struct {
	resolveFlamingo                         func(ctx context.Context) (*string, error)
	resolveCommerceProduct                  func(ctx context.Context, marketPlaceCode string, variantMarketPlaceCode *string) (graphqlproductdto.Product, error)
	resolveCommerceProductSearch            func(ctx context.Context, searchRequest searchdto.CommerceSearchRequest) (*graphql2.SearchResultDTO, error)
	resolveCommerceCustomerStatus           func(ctx context.Context) (*dtocustomer.CustomerStatusResult, error)
	resolveCommerceCustomer                 func(ctx context.Context) (*dtocustomer.CustomerResult, error)
	resolveCommerceCartDecoratedCart        func(ctx context.Context) (*dto.DecoratedCart, error)
	resolveCommerceCartValidator            func(ctx context.Context) (*validation.Result, error)
	resolveCommerceCartQtyRestriction       func(ctx context.Context, marketplaceCode string, variantCode *string, deliveryCode string) (*validation.RestrictionResult, error)
	resolveCommerceCheckoutActivePlaceOrder func(ctx context.Context) (bool, error)
	resolveCommerceCheckoutCurrentContext   func(ctx context.Context) (*dto1.PlaceOrderContext, error)
	resolveCommerceCategoryTree             func(ctx context.Context, activeCategoryCode string) (domain3.Tree, error)
	resolveCommerceCategory                 func(ctx context.Context, categoryCode string, categorySearchRequest *searchdto.CommerceSearchRequest) (*categorydto.CategorySearchResult, error)
}

func (r *rootResolverQuery) Inject(
	queryFlamingo *graphql4.FlamingoQueryResolver,
	queryCommerceProduct *graphql2.CommerceProductQueryResolver,
	queryCommerceProductSearch *graphql2.CommerceProductQueryResolver,
	queryCommerceCustomerStatus *graphql6.CustomerResolver,
	queryCommerceCustomer *graphql6.CustomerResolver,
	queryCommerceCartDecoratedCart *graphql1.CommerceCartQueryResolver,
	queryCommerceCartValidator *graphql1.CommerceCartQueryResolver,
	queryCommerceCartQtyRestriction *graphql1.CommerceCartQueryResolver,
	queryCommerceCheckoutActivePlaceOrder *graphql5.CommerceCheckoutQueryResolver,
	queryCommerceCheckoutCurrentContext *graphql5.CommerceCheckoutQueryResolver,
	queryCommerceCategoryTree *graphql7.CommerceCategoryQueryResolver,
	queryCommerceCategory *graphql7.CommerceCategoryQueryResolver,
) {
	r.resolveFlamingo = queryFlamingo.Flamingo
	r.resolveCommerceProduct = queryCommerceProduct.CommerceProduct
	r.resolveCommerceProductSearch = queryCommerceProductSearch.CommerceProductSearch
	r.resolveCommerceCustomerStatus = queryCommerceCustomerStatus.CommerceCustomerStatus
	r.resolveCommerceCustomer = queryCommerceCustomer.CommerceCustomer
	r.resolveCommerceCartDecoratedCart = queryCommerceCartDecoratedCart.CommerceCart
	r.resolveCommerceCartValidator = queryCommerceCartValidator.CommerceCartValidator
	r.resolveCommerceCartQtyRestriction = queryCommerceCartQtyRestriction.CommerceCartQtyRestriction
	r.resolveCommerceCheckoutActivePlaceOrder = queryCommerceCheckoutActivePlaceOrder.CommerceCheckoutActivePlaceOrder
	r.resolveCommerceCheckoutCurrentContext = queryCommerceCheckoutCurrentContext.CommerceCheckoutCurrentContext
	r.resolveCommerceCategoryTree = queryCommerceCategoryTree.CommerceCategoryTree
	r.resolveCommerceCategory = queryCommerceCategory.CommerceCategory
}

func (r *rootResolverQuery) Flamingo(ctx context.Context) (*string, error) {
	return r.resolveFlamingo(ctx)
}
func (r *rootResolverQuery) CommerceProduct(ctx context.Context, marketPlaceCode string, variantMarketPlaceCode *string) (graphqlproductdto.Product, error) {
	return r.resolveCommerceProduct(ctx, marketPlaceCode, variantMarketPlaceCode)
}
func (r *rootResolverQuery) CommerceProductSearch(ctx context.Context, searchRequest searchdto.CommerceSearchRequest) (*graphql2.SearchResultDTO, error) {
	return r.resolveCommerceProductSearch(ctx, searchRequest)
}
func (r *rootResolverQuery) CommerceCustomerStatus(ctx context.Context) (*dtocustomer.CustomerStatusResult, error) {
	return r.resolveCommerceCustomerStatus(ctx)
}
func (r *rootResolverQuery) CommerceCustomer(ctx context.Context) (*dtocustomer.CustomerResult, error) {
	return r.resolveCommerceCustomer(ctx)
}
func (r *rootResolverQuery) CommerceCartDecoratedCart(ctx context.Context) (*dto.DecoratedCart, error) {
	return r.resolveCommerceCartDecoratedCart(ctx)
}
func (r *rootResolverQuery) CommerceCartValidator(ctx context.Context) (*validation.Result, error) {
	return r.resolveCommerceCartValidator(ctx)
}
func (r *rootResolverQuery) CommerceCartQtyRestriction(ctx context.Context, marketplaceCode string, variantCode *string, deliveryCode string) (*validation.RestrictionResult, error) {
	return r.resolveCommerceCartQtyRestriction(ctx, marketplaceCode, variantCode, deliveryCode)
}
func (r *rootResolverQuery) CommerceCheckoutActivePlaceOrder(ctx context.Context) (bool, error) {
	return r.resolveCommerceCheckoutActivePlaceOrder(ctx)
}
func (r *rootResolverQuery) CommerceCheckoutCurrentContext(ctx context.Context) (*dto1.PlaceOrderContext, error) {
	return r.resolveCommerceCheckoutCurrentContext(ctx)
}
func (r *rootResolverQuery) CommerceCategoryTree(ctx context.Context, activeCategoryCode string) (domain3.Tree, error) {
	return r.resolveCommerceCategoryTree(ctx, activeCategoryCode)
}
func (r *rootResolverQuery) CommerceCategory(ctx context.Context, categoryCode string, categorySearchRequest *searchdto.CommerceSearchRequest) (*categorydto.CategorySearchResult, error) {
	return r.resolveCommerceCategory(ctx, categoryCode, categorySearchRequest)
}

func direct(root *rootResolver) map[string]interface{} {
	return map[string]interface{}{
		"Commerce_Cart_AdditionalData.CustomAttributes":       root.Commerce_Cart_AdditionalData().CustomAttributes,
		"Commerce_Cart_Cart.GetDeliveryByCode":                root.Commerce_Cart_Cart().GetDeliveryByCode,
		"Commerce_Cart_DefaultPaymentSelection.CartSplit":     root.Commerce_Cart_DefaultPaymentSelection().CartSplit,
		"Commerce_Cart_DeliveryInfo.AdditionalData":           root.Commerce_Cart_DeliveryInfo().AdditionalData,
		"Commerce_Cart_Item.AppliedDiscounts":                 root.Commerce_Cart_Item().AppliedDiscounts,
		"Commerce_Cart_ShippingItem.AppliedDiscounts":         root.Commerce_Cart_ShippingItem().AppliedDiscounts,
		"Commerce_Product_PriceInfo.ActiveBase":               root.Commerce_Product_PriceInfo().ActiveBase,
		"Commerce_Search_Meta.SortOptions":                    root.Commerce_Search_Meta().SortOptions,
		"Mutation.Flamingo":                                   root.Mutation().Flamingo,
		"Mutation.CommerceCartAddToCart":                      root.Mutation().CommerceCartAddToCart,
		"Mutation.CommerceCartDeleteCartDelivery":             root.Mutation().CommerceCartDeleteCartDelivery,
		"Mutation.CommerceCartDeleteItem":                     root.Mutation().CommerceCartDeleteItem,
		"Mutation.CommerceCartUpdateItemQty":                  root.Mutation().CommerceCartUpdateItemQty,
		"Mutation.CommerceCartUpdateBillingAddress":           root.Mutation().CommerceCartUpdateBillingAddress,
		"Mutation.CommerceCartUpdateSelectedPayment":          root.Mutation().CommerceCartUpdateSelectedPayment,
		"Mutation.CommerceCartApplyCouponCodeOrGiftCard":      root.Mutation().CommerceCartApplyCouponCodeOrGiftCard,
		"Mutation.CommerceCartRemoveGiftCard":                 root.Mutation().CommerceCartRemoveGiftCard,
		"Mutation.CommerceCartRemoveCouponCode":               root.Mutation().CommerceCartRemoveCouponCode,
		"Mutation.CommerceCartUpdateDeliveryAddresses":        root.Mutation().CommerceCartUpdateDeliveryAddresses,
		"Mutation.CommerceCartUpdateDeliveryShippingOptions":  root.Mutation().CommerceCartUpdateDeliveryShippingOptions,
		"Mutation.CommerceCartClean":                          root.Mutation().CommerceCartClean,
		"Mutation.CommerceCartUpdateAdditionalData":           root.Mutation().CommerceCartUpdateAdditionalData,
		"Mutation.CommerceCartUpdateDeliveriesAdditionalData": root.Mutation().CommerceCartUpdateDeliveriesAdditionalData,
		"Mutation.CommerceCheckoutStartPlaceOrder":            root.Mutation().CommerceCheckoutStartPlaceOrder,
		"Mutation.CommerceCheckoutCancelPlaceOrder":           root.Mutation().CommerceCheckoutCancelPlaceOrder,
		"Mutation.CommerceCheckoutClearPlaceOrder":            root.Mutation().CommerceCheckoutClearPlaceOrder,
		"Mutation.CommerceCheckoutRefreshPlaceOrder":          root.Mutation().CommerceCheckoutRefreshPlaceOrder,
		"Mutation.CommerceCheckoutRefreshPlaceOrderBlocking":  root.Mutation().CommerceCheckoutRefreshPlaceOrderBlocking,
		"Query.Flamingo":                         root.Query().Flamingo,
		"Query.CommerceProduct":                  root.Query().CommerceProduct,
		"Query.CommerceProductSearch":            root.Query().CommerceProductSearch,
		"Query.CommerceCustomerStatus":           root.Query().CommerceCustomerStatus,
		"Query.CommerceCustomer":                 root.Query().CommerceCustomer,
		"Query.CommerceCartDecoratedCart":        root.Query().CommerceCartDecoratedCart,
		"Query.CommerceCartValidator":            root.Query().CommerceCartValidator,
		"Query.CommerceCartQtyRestriction":       root.Query().CommerceCartQtyRestriction,
		"Query.CommerceCheckoutActivePlaceOrder": root.Query().CommerceCheckoutActivePlaceOrder,
		"Query.CommerceCheckoutCurrentContext":   root.Query().CommerceCheckoutCurrentContext,
		"Query.CommerceCategoryTree":             root.Query().CommerceCategoryTree,
		"Query.CommerceCategory":                 root.Query().CommerceCategory,
	}
}
