package domain

import (
	"testing"

	"time"

	"github.com/stretchr/testify/assert"
)

func TestIsSalable(t *testing.T) {
	salable := Saleable{
		IsSaleable:   true,
		SaleableTo:   time.Now().Add(time.Hour * time.Duration(1)),
		SaleableFrom: time.Now().Add(time.Hour * time.Duration(-1)),
	}
	assert.True(t, salable.IsSaleableNow(), "salable test")

	salable2 := Saleable{
		IsSaleable:   true,
		SaleableFrom: time.Now().Add(time.Hour * time.Duration(1)),
		SaleableTo:   time.Now().Add(time.Hour * time.Duration(-1)),
	}
	assert.False(t, salable2.IsSaleableNow(), "salable2 test")

	salable3 := Saleable{
		IsSaleable: true,
		SaleableTo: time.Now().Add(time.Hour * time.Duration(1)),
	}
	assert.True(t, salable3.IsSaleableNow(), "salable3 test")

	salable4 := Saleable{
		IsSaleable: true,
	}
	assert.True(t, salable4.IsSaleableNow(), "salable4 test")
}
