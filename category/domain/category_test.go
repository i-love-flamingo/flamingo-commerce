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
		inactiveCategory Category = CategoryData{IsActive: false}
		activeCategory   Category = CategoryData{IsActive: true}
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
					Children: []Category{
						inactiveCategory,
						inactiveCategory,
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
					Children: []Category{
						inactiveCategory,
						activeCategory,
					},
					IsActive: false,
				},
			},
			want: activeCategory,
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
