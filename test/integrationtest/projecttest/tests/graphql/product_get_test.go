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
			"meta":            map[string]interface{}{"keywords": []interface{}{string("keywords")}},
			"price":           map[string]interface{}{"activeBase": map[string]interface{}{"amount": 1.0, "currency": "€"}},
		},
		"configurable": map[string]interface{}{
			"title":           "TypeConfigurable product",
			"marketPlaceCode": "fake_configurable",
			"meta":            map[string]interface{}{"keywords": []interface{}{string("keywords")}},
			"variantSelection": map[string]interface{}{
				"variants": []interface{}{
					map[string]interface{}{
						"attributes": []interface{}{
							map[string]interface{}{
								"key":   string("color"),
								"value": string("White"),
							},
							map[string]interface{}{
								"key":   string("size"),
								"value": string("M"),
							},
						},
						"variant": map[string]interface{}{
							"marketplaceCode": "shirt-white-m",
						},
					},
					map[string]interface{}{
						"attributes": []interface{}{
							map[string]interface{}{
								"key":   string("color"),
								"value": string("Black"),
							},
							map[string]interface{}{
								"key":   string("size"),
								"value": string("L"),
							},
						},
						"variant": map[string]interface{}{
							"marketplaceCode": "shirt-black-l",
						},
					},
					map[string]interface{}{
						"attributes": []interface{}{
							map[string]interface{}{
								"key":   string("color"),
								"value": string("Red"),
							},
							map[string]interface{}{
								"key":   string("size"),
								"value": string("L"),
							},
						},
						"variant": map[string]interface{}{
							"marketplaceCode": "shirt-red-l",
						},
					},
					map[string]interface{}{
						"attributes": []interface{}{
							map[string]interface{}{
								"key":   string("color"),
								"value": string("Red"),
							},
							map[string]interface{}{
								"key":   string("size"),
								"value": string("M"),
							},
						},
						"variant": map[string]interface{}{
							"marketplaceCode": "shirt-red-m",
						},
					},
				},
				"attributes": []interface{}{
					map[string]interface{}{
						"label": string("Color"),
						"code":  string("color"),
						"options": []interface{}{
							map[string]interface{}{
								"label":    string("White"),
								"unitCode": string(""),
								"otherAttributesRestrictions": []interface{}{
									map[string]interface{}{
										"code":             string("size"),
										"availableOptions": []interface{}{"M"},
									},
								},
							},
							map[string]interface{}{
								"label":    string("Black"),
								"unitCode": string(""),
								"otherAttributesRestrictions": []interface{}{
									map[string]interface{}{
										"code":             string("size"),
										"availableOptions": []interface{}{"L"},
									},
								},
							},
							map[string]interface{}{
								"label":    string("Red"),
								"unitCode": string(""),
								"otherAttributesRestrictions": []interface{}{
									map[string]interface{}{
										"code":             string("size"),
										"availableOptions": []interface{}{"L", "M"},
									},
								},
							},
						},
					},
					map[string]interface{}{
						"label": string("Size"),
						"code":  string("size"),
						"options": []interface{}{
							map[string]interface{}{
								"label":    string("M"),
								"unitCode": string(""),
								"otherAttributesRestrictions": []interface{}{
									map[string]interface{}{
										"code":             string("color"),
										"availableOptions": []interface{}{"White", "Red"},
									},
								},
							},
							map[string]interface{}{
								"label":    string("L"),
								"unitCode": string(""),
								"otherAttributesRestrictions": []interface{}{
									map[string]interface{}{
										"code":             string("color"),
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
			"meta":  map[string]interface{}{"keywords": []interface{}{string("keywords")}},
			"variationSelections": []interface{}{
				map[string]interface{}{"code": string("color"), "label": string("Color")},
				map[string]interface{}{"code": string("size"), "label": string("Size")},
			},
			"activeVariationSelections": []interface{}{
				map[string]interface{}{"code": string("color"), "label": string("Color"), "value": string("Black")},
				map[string]interface{}{"code": string("size"), "label": string("Size"), "value": string("L")},
			},
		},
	}

	assertResponseForExpectedState(t, response, expectedState)
}
