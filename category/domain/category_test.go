package domain

import (
	"reflect"
	"testing"
)

func TestCategoryData_Attribute(t *testing.T) {
	type fields struct {
		CategoryCode       string
		CategoryName       string
		CategoryPath       string
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

func TestAttributes_mapToAdditionalAttributes(t *testing.T) {
	tests := []struct {
		name  string
		field Attributes
		want  interface{}
	}{
		{
			name: "empty attributes",
			want: AdditionalAttributes{},
		},
		{
			name: "no matching attributes",
			field: Attributes{
				"test": "ok",
			},
			want: AdditionalAttributes{},
		},
		{
			name: "match title",
			field: Attributes{
				"title": "Title",
			},
			want: AdditionalAttributes{Title: "Title"},
		},
		{
			name: "match all",
			field: Attributes{
				"title":            "Title",
				"marketingTitle":   "MarketingTitle",
				"shortDescription": "ShortDescription",
				"content":          "Content",
			},
			want: AdditionalAttributes{
				Title:            "Title",
				MarketingTitle:   "MarketingTitle",
				ShortDescription: "ShortDescription",
				Content:          "Content",
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			a := tt.field
			if got := a.mapToAdditionalAttributes(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Attributes.mapToAdditionalAttributes() = %v, want %v", got, tt.want)
			}
		})
	}
}
