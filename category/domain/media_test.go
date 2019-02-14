package domain

import (
	"reflect"
	"testing"
)

func TestMedias_Has(t *testing.T) {
	type args struct {
		usage string
	}

	var (
		detailMedia Media = MediaData{MediaUsage: MediaUsageDetail}
		teaserMedia Media = MediaData{MediaUsage: MediaUsageTeaser}
	)
	tests := []struct {
		name string
		m    Medias
		args args
		want bool
	}{
		{
			name: "Test empty media has no list item",
			m:    Medias{},
			args: args{
				usage: MediaUsageTeaser,
			},
			want: false,
		},
		{
			name: "Test empty media has no detail item",
			m:    Medias{},
			args: args{
				usage: MediaUsageDetail,
			},
			want: false,
		},
		{
			name: "Test detail only media has no list item",
			m: Medias{
				detailMedia,
			},
			args: args{
				usage: MediaUsageTeaser,
			},
			want: false,
		},
		{
			name: "Test detail only media has detail item",
			m: Medias{
				teaserMedia,
				detailMedia,
			},
			args: args{
				usage: MediaUsageDetail,
			},
			want: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.m.Has(tt.args.usage); got != tt.want {
				t.Errorf("Medias.Has() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestMedias_Get(t *testing.T) {
	type args struct {
		usage string
	}

	var (
		detailMedia Media = MediaData{MediaUsage: MediaUsageDetail}
		teaserMedia Media = MediaData{MediaUsage: MediaUsageTeaser}
	)
	tests := []struct {
		name string
		m    Medias
		args args
		want Media
	}{
		{
			name: "Empty media returns nil",
			m:    Medias{},
			args: args{
				usage: MediaUsageDetail,
			},
			want: MediaData{},
		},
		{
			name: "List without list media returns nil",
			m: Medias{
				detailMedia,
			},
			args: args{
				usage: MediaUsageTeaser,
			},
			want: MediaData{},
		},
		{
			name: "List with list media returns correct media",
			m: Medias{
				detailMedia,
				teaserMedia,
			},
			args: args{
				usage: MediaUsageDetail,
			},
			want: detailMedia,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.m.Get(tt.args.usage); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Medias.Get() = %v, want %v", got, tt.want)
			}
		})
	}
}
