package domain

// Unit describes a unit on an attribute
type Unit struct {
	Code   string
	Symbol string
}

// Constants of the unit codes
const (
	PCS = "PCS"

	MILLILITER = "MILLILITER"
	CENTILITER = "CENTILITER"
	LITER      = "LITER"
	OUNCE      = "OUNCE"
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
	"OUNCE": {
		Code:   "OUNCE",
		Symbol: "oz",
	},
	"PINT": {
		Code:   "PINT",
		Symbol: "pt",
	},
	"BARREL": {
		Code:   "BARREL",
		Symbol: "bbl",
	},
	"GALLON": {
		Code:   "GALLON",
		Symbol: "gal",
	},

	// Weight Units
}
