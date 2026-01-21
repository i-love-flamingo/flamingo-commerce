package fake

import (
	// embed test data file directly
	_ "embed"
	"encoding/json"
	"errors"
	"os"
	"path/filepath"
	"strings"

	"flamingo.me/flamingo/v3/framework/flamingo"

	searchDomain "flamingo.me/flamingo-commerce/v3/search/domain"

	"flamingo.me/flamingo-commerce/v3/product/domain"
)

// convertSpecifications converts map[string]interface{} to domain.Specifications
// This is needed because json.Unmarshal stores interface{} fields as map[string]interface{}
// instead of the concrete domain.Specifications type.
// We use JSON marshal/unmarshal for idiomatic conversion that respects struct field names.
func convertSpecifications(value interface{}) domain.Specifications {
	// If already the correct type, return as-is
	if specs, ok := value.(domain.Specifications); ok {
		return specs
	}

	// Handle map[string]interface{} from JSON unmarshal by re-marshaling to JSON
	// and then unmarshaling into the concrete type
	jsonBytes, err := json.Marshal(value)
	if err != nil {
		return domain.Specifications{}
	}

	var specs domain.Specifications
	if err := json.Unmarshal(jsonBytes, &specs); err != nil {
		return domain.Specifications{}
	}

	return specs
}

// processSpecificationsAttribute processes the "specifications" attribute
// and converts it from map[string]interface{} to domain.Specifications
func processSpecificationsAttribute(attributes domain.Attributes) {
	if attributes == nil {
		return
	}

	attr, exists := attributes["specifications"]
	if !exists {
		return
	}

	converted := convertSpecifications(attr.RawValue)
	attributes["specifications"] = domain.Attribute{
		Code:      attr.Code,
		CodeLabel: attr.CodeLabel,
		Label:     attr.Label,
		RawValue:  converted,
		UnitCode:  attr.UnitCode,
	}
}

// registerTestData returns files of given folder
func registerTestData(folder string, logger flamingo.Logger) map[string]string {
	testDataFiles := make(map[string]string)
	files, err := os.ReadDir(folder)
	if err != nil {
		logger.Info(err)
		return testDataFiles
	}

	for _, file := range files {
		if !file.IsDir() && strings.HasSuffix(file.Name(), ".json") {
			base := filepath.Base(file.Name())[:len(file.Name())-len(".json")]
			testDataFiles[base] = filepath.Join(folder, file.Name())
		}
	}
	return testDataFiles
}

// unmarshalJSONProduct unmarshals product based on type
func unmarshalJSONProduct(productRaw []byte) (domain.BasicProduct, error) {
	product := &map[string]interface{}{}
	err := json.Unmarshal(productRaw, product)

	if err != nil {
		return nil, err
	}

	productType, ok := (*product)["Type"]

	if !ok {
		return nil, errors.New("type is not specified")
	}

	if productType == domain.TypeConfigurable {
		configurableProduct := &domain.ConfigurableProduct{}
		err = json.Unmarshal(productRaw, configurableProduct)
		if err == nil {
			processSpecificationsAttribute(configurableProduct.Attributes)

			for i := range configurableProduct.Variants {
				processSpecificationsAttribute(configurableProduct.Variants[i].Attributes)
			}

			return *configurableProduct, nil
		}
	}

	if productType == domain.TypeBundle {
		bundleProduct := &domain.BundleProduct{}
		err = json.Unmarshal(productRaw, bundleProduct)
		if err == nil {
			processSpecificationsAttribute(bundleProduct.Attributes)
			return *bundleProduct, nil
		}
	}

	simpleProduct := &domain.SimpleProduct{}
	err = json.Unmarshal(productRaw, simpleProduct)
	if err != nil {
		return nil, err
	}

	processSpecificationsAttribute(simpleProduct.Attributes)
	return *simpleProduct, nil
}

//go:embed testdata/categoryFacetItems.json
var testdata []byte

func loadCategoryFacetItems(givenJSON string) ([]*searchDomain.FacetItem, error) {

	var items []*searchDomain.FacetItem

	if givenJSON != "" {
		fileContent, err := os.ReadFile(givenJSON)
		if err != nil {
			return nil, err
		}

		err = json.Unmarshal(fileContent, &items)
		if err != nil {
			return nil, err
		}
		return items, nil
	}

	err := json.Unmarshal(testdata, &items)
	if err != nil {
		return nil, err
	}

	return items, nil
}
