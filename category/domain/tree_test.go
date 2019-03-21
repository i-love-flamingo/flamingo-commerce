package domain

import (
	"testing"
)

func TestTreeData_GetSubTree(t *testing.T) {
	type args struct {
		tree Tree
	}

	var (
		inactiveCategory = TreeData{IsActive: false}
		activeCategory   = TreeData{IsActive: true}
	)
	tests := []struct {
		name string
		args args
		want int
	}{

		{
			name: "one subtree",
			args: args{
				tree: TreeData{
					SubTreesData: []*TreeData{
						&inactiveCategory,
					},
					IsActive: false,
				},
			},
			want: 1,
		},
		{
			name: "two subtree",
			args: args{
				tree: TreeData{
					SubTreesData: []*TreeData{
						&inactiveCategory,
						&activeCategory,
					},
					IsActive: false,
				},
			},
			want: 2,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := len(tt.args.tree.SubTrees()); got != tt.want {
				t.Errorf("SubTrees() = %v, want %v", got, tt.want)
			}
		})
	}
}
