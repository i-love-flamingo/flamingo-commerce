package application_test

import (
	"context"
	"net/url"
	"testing"

	"flamingo.me/flamingo-commerce/v3/breadcrumbs"
	"flamingo.me/flamingo/v3/framework/web"
	"github.com/stretchr/testify/assert"

	"flamingo.me/flamingo-commerce/v3/category/application"
	"flamingo.me/flamingo-commerce/v3/category/domain"
)

type (
	MockRouter struct{}
)

func (router *MockRouter) URL(name string, params map[string]string) *url.URL {
	return &url.URL{
		Path: "/foo",
	}
}

func TestBreadcrumbService_AddBreadcrumb(t *testing.T) {
	type args struct {
		category domain.Category
	}

	controller := new(breadcrumbs.Controller)

	tests := []struct {
		name string
		args args
		want interface{}
	}{
		{
			name: "no category active",
			args: args{
				category: getCategoryTreeWithoutActive(),
			},
			want: []breadcrumbs.Crumb{},
		},
		{
			name: "one category active",
			args: args{
				category: getCategoryTreeWithSingleActive(),
			},
			want: []breadcrumbs.Crumb{
				{
					Title: "Root",
					Url:   "/foo", // hardcoded value, we test only for the title to be correct
					Code:  "root",
				},
			},
		},
		{
			name: "my funny full tree",
			args: args{
				category: getFullCategoryTree(),
			},
			want: []breadcrumbs.Crumb{
				{
					Title: "Root",
					Url:   "/foo", // hardcoded value, we test only for the title to be correct
					Code:  "root",
				},
				{
					Title: "Sub1 Active",
					Url:   "/foo", // hardcoded value, we test only for the title to be correct
					Code:  "root-sub1-active",
				},
				{
					Title: "Sub2 Active",
					Url:   "/foo", // hardcoded value, we test only for the title to be correct
					Code:  "root-sub2-active",
				},
				{
					Title: "Sub3 Active",
					Url:   "/foo", // hardcoded value, we test only for the title to be correct
					Code:  "root-sub3-active",
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			// Init context and add breadcrumb
			bs := &application.BreadcrumbService{}
			bs.Inject(&MockRouter{})
			request := web.CreateRequest(nil, nil)
			ctx := web.ContextWithRequest(context.Background(), request)
			bs.AddBreadcrumb(ctx, tt.args.category)

			// get breadcrumb and validate it
			breadcrumb := controller.Data(ctx, nil, nil)

			assert.Equal(t, tt.want, breadcrumb)

		})
	}
}

func getCategoryTreeWithoutActive() domain.Category {
	categoryRoot := domain.CategoryData{
		CategoryCode: "root",
		CategoryName: "Root",
		IsActive:     false,
		Children:     nil,
	}
	return categoryRoot
}

func getCategoryTreeWithSingleActive() domain.Category {
	categoryRoot := domain.CategoryData{
		CategoryCode: "root",
		CategoryName: "Root",
		IsActive:     true,
		Children:     nil,
	}
	return categoryRoot
}

func getFullCategoryTree() domain.Category {
	categoryRoot := domain.CategoryData{
		CategoryCode: "root",
		CategoryName: "Root",
		IsActive:     true,
		Children: []*domain.CategoryData{
			{
				CategoryCode: "root-sub1-inactive",
				CategoryName: "Sub1 Inactive",
				IsActive:     false,
				Children:     nil,
			},
			{
				CategoryCode: "root-sub1-active",
				CategoryName: "Sub1 Active",
				IsActive:     true,
				Children: []*domain.CategoryData{
					{
						CategoryCode: "root-sub2-active",
						CategoryName: "Sub2 Active",
						IsActive:     true,
						Children: []*domain.CategoryData{
							{
								CategoryCode: "root-sub3-active",
								CategoryName: "Sub3 Active",
								IsActive:     true,
								Children:     nil,
							},
							{
								CategoryCode: "root-sub3-inactive",
								CategoryName: "Sub3 Inactive",
								IsActive:     false,
								Children:     nil,
							},
						},
					},
					{
						CategoryCode: "root-sub2-inactive",
						CategoryName: "Sub2 Inactive",
						IsActive:     false,
						Children:     nil,
					},
				},
			},
		},
	}
	return categoryRoot
}
