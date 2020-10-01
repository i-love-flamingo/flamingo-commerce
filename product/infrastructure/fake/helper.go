package fake

import (
	"encoding/json"
	"errors"
	"flamingo.me/flamingo-commerce/v3/product/domain"
	"io/ioutil"
	"path/filepath"
	"strings"
)

// registerTestData
func registerTestData(folder string) map[string]string {
	testDataFiles := make(map[string]string)
	files, err := ioutil.ReadDir(folder)
	if err == nil {
		for _, file := range files {
			if !file.IsDir() && strings.HasSuffix(file.Name(), ".json") {
				base := filepath.Base(file.Name())[:len(file.Name())-5]
				testDataFiles[base] = filepath.Join(folder, file.Name())
			}
		}
	}
	return testDataFiles
}

//unmarshalJSONProduct
func unmarshalJSONProduct(productRaw []byte) (domain.BasicProduct, error) {
	product := &map[string]interface{}{}
	err := json.Unmarshal(productRaw, product)

	if err != nil {
		return nil, err
	}

	productType, ok := (*product)["Type"]

	if !ok {
		return nil, errors.New("Type is not specified")
	}

	if productType == domain.TypeConfigurable {
		configurableProduct := &domain.ConfigurableProduct{}
		err = json.Unmarshal(productRaw, configurableProduct)
		if err == nil {
			return *configurableProduct, nil
		}
	}

	simpleProduct := &domain.SimpleProduct{}
	err = json.Unmarshal(productRaw, simpleProduct)
	if err != nil {
		return nil, err
	}

	return *simpleProduct, nil
}
