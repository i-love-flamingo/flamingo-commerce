package domain

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestSources_MainLocation(t *testing.T) {
	var sources Sources
	assert.Equal(t, "", sources.MainLocation())

	sources = append(sources, Source{
		LocationCode: "loc",
		Qty:          1,
	})

	assert.Equal(t, "loc", sources.MainLocation())

}

func TestSources_QtySum(t *testing.T) {
	var sources Sources
	assert.Equal(t, "", sources.MainLocation())

	sources = append(sources, Source{
		LocationCode: "loc1",
		Qty:          1,
	})

	sources = append(sources, Source{
		LocationCode: "loc2",
		Qty:          2,
	})

	assert.Equal(t, 3, sources.QtySum())
}

func TestSources_Reduce(t *testing.T) {
	var sources Sources
	assert.Equal(t, "", sources.MainLocation())

	sources = append(sources, Source{
		LocationCode: "loc1",
		Qty:          1,
	})

	sources = append(sources, Source{
		LocationCode: "loc2",
		Qty:          2,
	})

	var sourcesreduce Sources
	sourcesreduce = append(sourcesreduce, Source{
		LocationCode: "loc2",
		Qty:          1,
	})

	sources = sources.Reduce(sourcesreduce)
	assert.Equal(t, 2, sources.QtySum())
}
