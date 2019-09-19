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

func TestSources_Next(t *testing.T) {
	var sources Sources
	sources = append(sources, Source{
		LocationCode: "loc",
		Qty:          2,
	},
		Source{
			LocationCode: "loc2",
			Qty:          2,
		})

	var source Source
	var err error
	for i := 1; i <= 2; i++ {
		source, sources, err = sources.Next()
		assert.NoError(t, err)
		assert.Equal(t, "loc", source.LocationCode)
		assert.Equal(t, 1, source.Qty)
		assert.Equal(t, 4-i, sources.QtySum())
	}
	for i := 1; i <= 2; i++ {
		source, sources, err = sources.Next()
		assert.NoError(t, err)
		assert.Equal(t, "loc2", source.LocationCode)
		assert.Equal(t, 1, source.Qty)
		assert.Equal(t, 2-i, sources.QtySum())
	}
	assert.Equal(t, 0, sources.QtySum())
	_, _, err = sources.Next()
	assert.Error(t, err)
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
