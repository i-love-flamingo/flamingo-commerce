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
			return *configurableProduct, nil
		}
	}

	if productType == domain.TypeBundle {
		bundleProduct := &domain.BundleProduct{}
		err = json.Unmarshal(productRaw, bundleProduct)
		if err == nil {
			return *bundleProduct, nil
		}
	}

	simpleProduct := &domain.SimpleProduct{}
	err = json.Unmarshal(productRaw, simpleProduct)
	if err != nil {
		return nil, err
	}

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
