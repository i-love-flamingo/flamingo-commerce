//go:build integration
// +build integration

package graphql_test

import (
	"net/http"
	"testing"

	"flamingo.me/flamingo-commerce/v3/test/integrationtest"
	"flamingo.me/flamingo-commerce/v3/test/integrationtest/projecttest/helper"
)

func Test_CommerceProductGet(t *testing.T) {
	baseURL := "http://" + FlamingoURL
	e := integrationtest.NewHTTPExpect(t, baseURL)
	response := helper.GraphQlRequest(t, e, loadGraphQL(t, "product_get", nil)).Expect().Status(http.StatusOK)

	expectedState := map[string]interface{}{
		"simple": map[string]interface{}{
			"title":           "TypeSimple product",
			"marketPlaceCode": "fake_simple",
			"meta":            map[string]interface{}{"keywords": []interface{}{"keywords"}},
			"price":           map[string]interface{}{"activeBase": map[string]interface{}{"amount": 1.0, "currency": "€"}},
		},
		"configurable": map[string]interface{}{
			"title":           "TypeConfigurable product",
			"marketPlaceCode": "fake_configurable",
			"meta":            map[string]interface{}{"keywords": []interface{}{"keywords"}},
			"variantSelection": map[string]interface{}{
				"variants": []interface{}{
					map[string]interface{}{
						"attributes": []interface{}{
							map[string]interface{}{
								"key":   "color",
								"value": "White",
							},
							map[string]interface{}{
								"key":   "size",
								"value": "M",
							},
						},
						"variant": map[string]interface{}{
							"marketplaceCode": "shirt-white-m",
						},
					},
					map[string]interface{}{
						"attributes": []interface{}{
							map[string]interface{}{
								"key":   "color",
								"value": "Black",
							},
							map[string]interface{}{
								"key":   "size",
								"value": "L",
							},
						},
						"variant": map[string]interface{}{
							"marketplaceCode": "shirt-black-l",
						},
					},
					map[string]interface{}{
						"attributes": []interface{}{
							map[string]interface{}{
								"key":   "color",
								"value": "Red",
							},
							map[string]interface{}{
								"key":   "size",
								"value": "L",
							},
						},
						"variant": map[string]interface{}{
							"marketplaceCode": "shirt-red-l",
						},
					},
					map[string]interface{}{
						"attributes": []interface{}{
							map[string]interface{}{
								"key":   "color",
								"value": "Red",
							},
							map[string]interface{}{
								"key":   "size",
								"value": "M",
							},
						},
						"variant": map[string]interface{}{
							"marketplaceCode": "shirt-red-m",
						},
					},
				},
				"attributes": []interface{}{
					map[string]interface{}{
						"label": "Color",
						"code":  "color",
						"options": []interface{}{
							map[string]interface{}{
								"label":    "White",
								"unitCode": "",
								"otherAttributesRestrictions": []interface{}{
									map[string]interface{}{
										"code":             "size",
										"availableOptions": []interface{}{"M"},
									},
								},
							},
							map[string]interface{}{
								"label":    "Black",
								"unitCode": "",
								"otherAttributesRestrictions": []interface{}{
									map[string]interface{}{
										"code":             "size",
										"availableOptions": []interface{}{"L"},
									},
								},
							},
							map[string]interface{}{
								"label":    "Red",
								"unitCode": "",
								"otherAttributesRestrictions": []interface{}{
									map[string]interface{}{
										"code":             "size",
										"availableOptions": []interface{}{"L", "M"},
									},
								},
							},
						},
					},
					map[string]interface{}{
						"label": "Size",
						"code":  "size",
						"options": []interface{}{
							map[string]interface{}{
								"label":    "M",
								"unitCode": "",
								"otherAttributesRestrictions": []interface{}{
									map[string]interface{}{
										"code":             "color",
										"availableOptions": []interface{}{"White", "Red"},
									},
								},
							},
							map[string]interface{}{
								"label":    "L",
								"unitCode": "",
								"otherAttributesRestrictions": []interface{}{
									map[string]interface{}{
										"code":             "color",
										"availableOptions": []interface{}{"Black", "Red"},
									},
								},
							},
						},
					},
				},
			},
		},
		"active_variant": map[string]interface{}{
			"type":  "configurable_with_activevariant",
			"title": "Shirt Black L",
			"meta":  map[string]interface{}{"keywords": []interface{}{"keywords"}},
			"variationSelections": []interface{}{
				map[string]interface{}{"code": "color", "label": "Color"},
				map[string]interface{}{"code": "size", "label": "Size"},
			},
			"activeVariationSelections": []interface{}{
				map[string]interface{}{"code": "color", "label": "Color", "value": "Black"},
				map[string]interface{}{"code": "size", "label": "Size", "value": "L"},
			},
		},
	}

	assertResponseForExpectedState(t, response, expectedState)
}
