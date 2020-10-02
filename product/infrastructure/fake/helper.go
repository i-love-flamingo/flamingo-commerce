package fake

import (
	"encoding/json"
	"errors"
	"flamingo.me/flamingo/v3/framework/flamingo"

	"flamingo.me/flamingo-commerce/v3/product/domain"
	"io/ioutil"
	"path/filepath"
	"strings"
)

// registerTestData returns files of given folder
func registerTestData(folder string, logger flamingo.Logger) map[string]string {
	testDataFiles := make(map[string]string)
	files, err := ioutil.ReadDir(folder)
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

	simpleProduct := &domain.SimpleProduct{}
	err = json.Unmarshal(productRaw, simpleProduct)
	if err != nil {
		return nil, err
	}

	return *simpleProduct, nil
}
