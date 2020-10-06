package fake

import (
	"encoding/json"
	"io/ioutil"
	"path/filepath"
	"strings"

	"flamingo.me/flamingo/v3/framework/flamingo"

	"flamingo.me/flamingo-commerce/v3/category/domain"
)

// LoadCategoryTree returns tree data from file
func LoadCategoryTree(testDataFiles map[string]string, logger flamingo.Logger) []*domain.TreeData {
	var tree []*domain.TreeData
	if categoryTreeFile, ok := testDataFiles["categoryTree"]; ok {
		data, err := ioutil.ReadFile(categoryTreeFile)
		if err != nil {
			logger.Warn(err)
			return tree
		}
		err = json.Unmarshal(data, &tree)
		if err != nil {
			logger.Warn(err)
		}
	} else {
		jsonFile, err := Asset("categoryTree.json")
		if err != nil {
			logger.Warn(err)
			return tree
		}
		err = json.Unmarshal(jsonFile, &tree)
		if err != nil {
			logger.Warn(err)
		}
	}
	return tree
}

// LoadCategory returns category data from file
func LoadCategory(categoryCode string, testDataFiles map[string]string, logger flamingo.Logger) domain.Category {
	var categoryData *domain.CategoryData

	if categoryTreeFile, ok := testDataFiles[categoryCode]; ok {
		data, err := ioutil.ReadFile(categoryTreeFile)
		if err != nil {
			logger.Warn(err)
			return nil
		}
		err = json.Unmarshal(data, &categoryData)
		if err != nil {
			logger.Warn(err)
			return nil
		}
	} else {
		jsonFile, err := Asset(categoryCode + ".json")
		if err != nil {
			logger.Warn(err)
			return nil
		}
		err = json.Unmarshal(jsonFile, &categoryData)
		if err != nil {
			logger.Warn(err)
			return nil
		}
	}
	return categoryData
}

// RegisterTestData returns files of given folder
func RegisterTestData(folder string, logger flamingo.Logger) map[string]string {
	testDataFiles := make(map[string]string)
	files, err := ioutil.ReadDir(folder)
	if err != nil {
		logger.Warn(err)
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
