package domain

import (
	"reflect"
	"testing"
)

func TestGetActive(t *testing.T) {
	type args struct {
		c Category
	}

	var (
		inactiveCategory = CategoryData{IsActive: false}
		activeCategory   = CategoryData{IsActive: true}
	)
	tests := []struct {
		name string
		args args
		want Category
	}{
		{
			name: "Empty category returns nil",
			args: args{
				c: nil,
			},
			want: nil,
		},
		{
			name: "Inactive tree retuns nil",
			args: args{
				c: CategoryData{
					Children: []*CategoryData{
						&inactiveCategory,
						&inactiveCategory,
					},
					IsActive: false,
				},
			},
			want: nil,
		},
		{
			name: "Active tree returns non nil active category",
			args: args{
				c: CategoryData{
					Children: []*CategoryData{
						&inactiveCategory,
						&activeCategory,
					},
					IsActive: false,
				},
			},
			want: &activeCategory,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := GetActive(tt.args.c); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("GetActive() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestCategoryData_Attribute(t *testing.T) {
	type fields struct {
		CategoryCode       string
		CategoryName       string
		CategoryPath       string
		Children           []*CategoryData
		IsActive           bool
		IsPromoted         bool
		CategoryMedia      Medias
		CategoryTypeCode   string
		CategoryAttributes Attributes
	}
	type args struct {
		code string
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		want   interface{}
	}{
		{
			name: "empty attributes",
			args: args{
				code: "test",
			},
			want: nil,
		},
		{
			name: "not found",
			args: args{
				code: "invalid",
			},
			fields: fields{
				CategoryAttributes: Attributes{
					"test": "ok",
				},
			},
			want: nil,
		},
		{
			args: args{
				code: "test",
			},
			want: nil,
		},
		{
			name: "found",
			args: args{
				code: "test",
			},
			fields: fields{
				CategoryAttributes: Attributes{
					"test": "ok",
				},
			},
			want: "ok",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := CategoryData{
				CategoryCode:       tt.fields.CategoryCode,
				CategoryName:       tt.fields.CategoryName,
				CategoryPath:       tt.fields.CategoryPath,
				Children:           tt.fields.Children,
				IsActive:           tt.fields.IsActive,
				IsPromoted:         tt.fields.IsPromoted,
				CategoryMedia:      tt.fields.CategoryMedia,
				CategoryTypeCode:   tt.fields.CategoryTypeCode,
				CategoryAttributes: tt.fields.CategoryAttributes,
			}
			if got := c.Attribute(tt.args.code); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("CategoryData.Attribute() = %v, want %v", got, tt.want)
			}
		})
	}
}
