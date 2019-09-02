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

func TestSources_ReduceInNeg(t *testing.T) {
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

	var sourcesreduce1 Sources
	sourcesreduce1 = append(sourcesreduce1, Source{
		LocationCode: "loc2",
		Qty:          5,
	})
	var sourcesreduce2 Sources
	sourcesreduce2 = append(sourcesreduce2, Source{
		LocationCode: "loc2",
		Qty:          2,
	})

	sources = sources.Reduce(sourcesreduce1)
	sources = sources.Reduce(sourcesreduce2)
	assert.Equal(t, -4, sources.QtySum())
}
