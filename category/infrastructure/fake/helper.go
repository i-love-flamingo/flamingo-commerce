package fake

import (
	"embed"
	"encoding/json"
	"os"
	"path/filepath"
	"strings"

	"flamingo.me/flamingo/v3/framework/flamingo"

	"flamingo.me/flamingo-commerce/v3/category/domain"
)

//go:embed mock
var mock embed.FS

// LoadCategoryTree returns tree data from file
func LoadCategoryTree(testDataFiles map[string]string, logger flamingo.Logger) []*domain.TreeData {
	var tree []*domain.TreeData
	if categoryTreeFile, ok := testDataFiles["categoryTree"]; ok {
		data, err := os.ReadFile(categoryTreeFile)
		if err != nil {
			logger.Warn(err)
			return tree
		}
		err = json.Unmarshal(data, &tree)
		if err != nil {
			logger.Warn(err)
		}
	} else {
		jsonFile, err := mock.ReadFile("mock/categoryTree.json")
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
		data, err := os.ReadFile(categoryTreeFile)
		if err != nil {
			logger.Warn(err)
			return nil
		}
		categoryData, err = unmarshalCategoryData(data)
		if err != nil {
			logger.Warn(err)
			return nil
		}
	} else {
		jsonFile, err := mock.ReadFile("mock/" + categoryCode + ".json")
		if err != nil {
			logger.Warn(err)
			return nil
		}
		categoryData, err = unmarshalCategoryData(jsonFile)
		if err != nil {
			logger.Warn(err)
			return nil
		}
	}
	return categoryData
}

type jsonCategoryData struct {
	CategoryCode       string
	CategoryName       string
	CategoryPath       string
	CategoryTypeCode   string
	IsPromoted         bool
	IsActive           bool
	CategoryMedia      []domain.MediaData
	CategoryAttributes map[string]domain.Attribute
	Promotion          jsonCategoryPromotion
}

type jsonCategoryPromotion struct {
	LinkType   string
	LinkTarget string
	Media      []domain.MediaData
}

func unmarshalCategoryData(jsonData []byte) (*domain.CategoryData, error) {
	var categoryData *jsonCategoryData
	if err := json.Unmarshal(jsonData, &categoryData); err != nil {
		return nil, err
	}

	categoryMedia := make([]domain.Media, len(categoryData.CategoryMedia))
	for i, media := range categoryData.CategoryMedia {
		categoryMedia[i] = media
	}

	promotionMedia := make([]domain.Media, len(categoryData.Promotion.Media))
	for i, media := range categoryData.Promotion.Media {
		promotionMedia[i] = media
	}

	return &domain.CategoryData{
		CategoryCode:       categoryData.CategoryCode,
		CategoryName:       categoryData.CategoryName,
		CategoryPath:       categoryData.CategoryPath,
		IsPromoted:         categoryData.IsPromoted,
		IsActive:           categoryData.IsActive,
		CategoryMedia:      categoryMedia,
		CategoryTypeCode:   categoryData.CategoryTypeCode,
		CategoryAttributes: categoryData.CategoryAttributes,
		Promotion: domain.Promotion{
			LinkType:   categoryData.Promotion.LinkType,
			LinkTarget: categoryData.Promotion.LinkTarget,
			Media:      promotionMedia,
		},
	}, nil
}

// RegisterTestData returns files of given folder
func RegisterTestData(folder string, logger flamingo.Logger) map[string]string {
	testDataFiles := make(map[string]string)
	files, err := os.ReadDir(folder)
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
