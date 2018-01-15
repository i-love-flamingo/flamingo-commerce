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

func getTimestringByOffset(offset int) string {
	return time.Now().Add(time.Hour * 24 * time.Duration(offset)).Format(time.RFC3339)
}

func getDateAttributeWithOffset(attrName string, offset int) Attribute {
	return Attribute{
		RawValue: getTimestringByOffset(offset),
		Code:     attrName,
	}
}

func getDateAttributesWithOffset(newFromDateOffset int, newToDateOffset int) Attributes {
	attributes := make(map[string]Attribute)
	attributes["newFromDate"] = getDateAttributeWithOffset("newFromDate", newFromDateOffset)
	attributes["newToDate"] = getDateAttributeWithOffset("newToDate", newToDateOffset)
	return attributes
}

func getVariantWithOffset(newFromDateOffset int, newToDateOffset int) Variant {
	return Variant{
		BasicProductData: BasicProductData{
			Attributes: getDateAttributesWithOffset(newFromDateOffset, newToDateOffset),
		},
	}
}

func TestIsNew(t *testing.T) {
	productWithoutAttributes := Variant{}
	assert.False(t, productWithoutAttributes.IsNew(), "product without attributes")

	startInThePastNoEnd := getVariantWithOffset(-1, 0)
	delete(startInThePastNoEnd.Attributes, "newToDate")
	assert.True(t, startInThePastNoEnd.IsNew(), "start in the past, no end")

	startInThePastEndInTheFuture := getVariantWithOffset(-1, 1)
	assert.True(t, startInThePastEndInTheFuture.IsNew(), "start in the past, end in the future")

	startInTheFutureEndInTheFuture := getVariantWithOffset(1, 1)
	assert.False(t, startInTheFutureEndInTheFuture.IsNew(), "start in the future, end in the future")

	startInTheFutureNoEnd := getVariantWithOffset(1, 0)
	delete(startInTheFutureNoEnd.Attributes, "newToDate")
	assert.False(t, startInTheFutureNoEnd.IsNew(), "start in future, no end")

	noStartEndInThePast := getVariantWithOffset(0, -1)
	delete(noStartEndInThePast.Attributes, "newFromDate")
	assert.False(t, noStartEndInThePast.IsNew(), "no start, end in the past")

	noStartEndInTheFuture := getVariantWithOffset(0, 1)
	delete(noStartEndInTheFuture.Attributes, "newFromDate")
	assert.True(t, noStartEndInTheFuture.IsNew(), "no start, end in the future")

	invalidDates := getVariantWithOffset(0, 1)
	invalidDates.Attributes["newFromDate"] = Attribute{
		RawValue: "8765..88",
		Code:     "newFromDate",
	}
	assert.False(t, invalidDates.IsNew(), "invalid date should return false")

}
