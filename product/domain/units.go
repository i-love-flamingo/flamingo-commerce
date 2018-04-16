package domain

// Unit describes a unit on an attribute
type Unit struct {
	Code   string
	Symbol string
}

// Constants of the unit codes
const (
	// Piece Units
	PCS = "PCS"

	// Volume Units
	MILLILITER = "MILLILITER"
	CENTILITER = "CENTILITER"
	LITER      = "LITER"
	OUNCE      = "OUNCE"
	PINT       = "PINT"
	BARREL     = "BARREL"
	GALLON     = "GALLON"

	// Weigt Units
	MILLIGRAM = "MILLIGRAM"
	GRAM      = "GRAM"
	KILOGRAM  = "KILOGRAM"
	POUND     = "POUND"
)

// Units provides the unit map
var Units = map[string]Unit{
	// Piece Units
	PCS: {
		Code:   PCS,
		Symbol: "pcs",
	},

	// Volume Units
	MILLILITER: {
		Code:   MILLILITER,
		Symbol: "ml",
	},
	CENTILITER: {
		Code:   CENTILITER,
		Symbol: "cl",
	},
	LITER: {
		Code:   LITER,
		Symbol: "l",
	},
	OUNCE: {
		Code:   OUNCE,
		Symbol: "oz",
	},
	PINT: {
		Code:   PINT,
		Symbol: "pt",
	},
	BARREL: {
		Code:   BARREL,
		Symbol: "bbl",
	},
	GALLON: {
		Code:   GALLON,
		Symbol: "gal",
	},

	// Weight Units
	MILLIGRAM: {
		Code:   MILLIGRAM,
		Symbol: "mg",
	},
	GRAM: {
		Code:   GRAM,
		Symbol: "g",
	},
	KILOGRAM: {
		Code:   KILOGRAM,
		Symbol: "kg",
	},
	POUND: {
		Code:   POUND,
		Symbol: "lb",
	},
}
