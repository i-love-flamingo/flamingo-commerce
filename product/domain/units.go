package domain

// Unit describes a unit on an attribute
type Unit struct {
	Code   string
	Symbol string
}

// Constants of the unit codes
// See also the PIM unit codes on which these are based (and no other will get into this system)
const (
	// Area Units
	SQUARE_MILLIMETER = "SQUARE_MILLIMETER"
	SQUARE_CENTIMETER = "SQUARE_CENTIMETER"
	SQUARE_DECIMETER  = "SQUARE_DECIMETER"
	SQUARE_METER      = "SQUARE_METER"
	CENTIARE          = "CENTIARE"
	SQUARE_DEKAMETER  = "SQUARE_DEKAMETER"
	ARE               = "ARE"
	SQUARE_HECTOMETER = "SQUARE_HECTOMETER"
	HECTARE           = "HECTARE"
	SQUARE_KILOMETER  = "SQUARE_KILOMETER"
	SQUARE_MIL        = "SQUARE_MIL"
	SQUARE_INCH       = "SQUARE_INCH"
	SQUARE_FOOT       = "SQUARE_FOOT"
	SQUARE_YARD       = "SQUARE_YARD"
	ARPENT            = "ARPENT"
	ACRE              = "ACRE"
	SQUARE_FURLONG    = "SQUARE_FURLONG"
	SQUARE_MILE       = "SQUARE_MILE"

	// Binary Units
	BIT      = "BIT"
	BYTE     = "BYTE"
	KILOBYTE = "KILOBYTE"
	MEGABYTE = "MEGABYTE"
	GIGABYTE = "GIGABYTE"
	TERABYTE = "TERABYTE"

	// Loudness Units
	DECIBEL = "DECIBEL"

	// Frequency Units
	HERTZ     = "HERTZ"
	KILOHERTZ = "KILOHERTZ"
	MEGAHERTZ = "MEGAHERTZ"
	GIGAHERTZ = "GIGAHERTZ"
	TERAHERTZ = "TERAHERTZ"

	// Length Unit
	MILLIMETER = "MILLIMETER"
	CENTIMETER = "CENTIMETER"
	DECIMETER  = "DECIMETER"
	METER      = "METER"
	DEKAMETER  = "DEKAMETER"
	HECTOMETER = "HECTOMETER"
	KILOMETER  = "KILOMETER"
	MIL        = "MIL"
	INCH       = "INCH"
	FEET       = "FEET"
	YARD       = "YARD"
	CHAIN      = "CHAIN"
	FURLONG    = "FURLONG"
	MILE       = "MILE"

	// Power Units
	WATT     = "WATT"
	KILOWATT = "KILOWATT"
	MEGAWATT = "MEGAWATT"
	GIGAWATT = "GIGAWATT"
	TERAWATT = "TERAWATT"

	// Voltage Units
	MILLIVOLT = "MILLIVOLT"
	CENTIVOLT = "CENTIVOLT"
	DECIVOLT  = "DECIVOLT"
	VOLT      = "VOLT"
	DEKAVOLT  = "DEKAVOLT"
	HECTOVOLT = "HECTOVOLT"
	KILOVOLT  = "KILOVOLT"

	// Intensity
	MILLIAMPERE = "MILLIAMPERE"
	CENTIAMPERE = "CENTIAMPERE"
	DECIAMPERE  = "DECIAMPERE"
	AMPERE      = "AMPERE"
	DEKAMPERE   = "DEKAMPERE"
	HECTOAMPERE = "HECTOAMPERE"
	KILOAMPERE  = "KILOAMPERE"

	// Resistance
	MILLIOHM = "MILLIOHM"
	CENTIOHM = "CENTIOHM"
	DECIOHM  = "DECIOHM"
	OHM      = "OHM"
	DEKAOHM  = "DEKAOHM"
	HECTOHM  = "HECTOHM"
	KILOHM   = "KILOHM"
	MEGOHM   = "MEGOHM"

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
	// Area Units
	SQUARE_MILLIMETER: {
		Code:   SQUARE_MILLIMETER,
		Symbol: "mm²",
	},

	SQUARE_CENTIMETER: {
		Code:   SQUARE_CENTIMETER,
		Symbol: "cm²",
	},

	SQUARE_DECIMETER: {
		Code:   SQUARE_DECIMETER,
		Symbol: "dm²",
	},

	SQUARE_METER: {
		Code:   SQUARE_METER,
		Symbol: "m²",
	},

	CENTIARE: {
		Code:   CENTIARE,
		Symbol: "ca",
	},

	SQUARE_DEKAMETER: {
		Code:   SQUARE_DEKAMETER,
		Symbol: "dam²",
	},

	ARE: {
		Code:   ARE,
		Symbol: "a",
	},

	SQUARE_HECTOMETER: {
		Code:   SQUARE_HECTOMETER,
		Symbol: "hm²",
	},

	HECTARE: {
		Code:   HECTARE,
		Symbol: "ha",
	},

	SQUARE_KILOMETER: {
		Code:   SQUARE_KILOMETER,
		Symbol: "km²",
	},

	SQUARE_MIL: {
		Code:   "SQUARE_MIL",
		Symbol: "sq mil",
	},

	SQUARE_INCH: {
		Code:   "SQUARE_INCH",
		Symbol: "in²",
	},

	SQUARE_FOOT: {
		Code:   "SQUARE_FOOT",
		Symbol: "ft²",
	},

	SQUARE_YARD: {
		Code:   "SQUARE_YARD",
		Symbol: "yd²",
	},

	ARPENT: {
		Code:   "ARPENT",
		Symbol: "arpent",
	},

	ACRE: {
		Code:   "ACRE",
		Symbol: "A",
	},

	SQUARE_FURLONG: {
		Code:   "SQUARE_FURLONG",
		Symbol: "fur²",
	},

	SQUARE_MILE: {
		Code:   "SQUARE_MILE",
		Symbol: "mi²",
	},

	// Binary Units
	BIT: {
		Code:   BIT,
		Symbol: "b",
	},
	BYTE: {
		Code:   BYTE,
		Symbol: "B",
	},
	KILOBYTE: {
		Code:   KILOBYTE,
		Symbol: "kB",
	},
	MEGABYTE: {
		Code:   MEGABYTE,
		Symbol: "MB",
	},
	GIGABYTE: {
		Code:   GIGABYTE,
		Symbol: "GB",
	},
	TERABYTE: {
		Code:   TERABYTE,
		Symbol: "TB",
	},

	// Loudness Units
	DECIBEL: {
		Code:   DECIBEL,
		Symbol: "dB",
	},

	// Frequency Units
	HERTZ: {
		Code:   HERTZ,
		Symbol: "Hz",
	},
	KILOHERTZ: {
		Code:   KILOHERTZ,
		Symbol: "kHz",
	},
	MEGAHERTZ: {
		Code:   MEGAHERTZ,
		Symbol: "MHz",
	},
	GIGAHERTZ: {
		Code:   GIGAHERTZ,
		Symbol: "GHz",
	},
	TERAHERTZ: {
		Code:   TERAHERTZ,
		Symbol: "THz",
	},

	// Length Units
	MILLIMETER: {
		Code:   MILLIMETER,
		Symbol: "mm",
	},
	CENTIMETER: {
		Code:   CENTIMETER,
		Symbol: "cm",
	},
	DECIMETER: {
		Code:   DECIMETER,
		Symbol: "dm",
	},
	METER: {
		Code:   METER,
		Symbol: "m",
	},
	DEKAMETER: {
		Code:   DEKAMETER,
		Symbol: "dam",
	},
	HECTOMETER: {
		Code:   HECTOMETER,
		Symbol: "hm",
	},
	KILOMETER: {
		Code:   KILOMETER,
		Symbol: "km",
	},
	MIL: {
		Code:   MIL,
		Symbol: "mil",
	},
	INCH: {
		Code:   INCH,
		Symbol: "in",
	},
	FEET: {
		Code:   FEET,
		Symbol: "ft",
	},
	YARD: {
		Code:   YARD,
		Symbol: "yd",
	},
	CHAIN: {
		Code:   CHAIN,
		Symbol: "ch",
	},
	FURLONG: {
		Code:   FURLONG,
		Symbol: "fur",
	},
	MILE: {
		Code:   MILE,
		Symbol: "mi",
	},

	// Power Units
	WATT: {
		Code:   WATT,
		Symbol: "W",
	},
	KILOWATT: {
		Code:   KILOWATT,
		Symbol: "mW",
	},
	MEGAWATT: {
		Code:   MEGAWATT,
		Symbol: "MW",
	},
	GIGAWATT: {
		Code:   GIGAWATT,
		Symbol: "GW",
	},
	TERAWATT: {
		Code:   TERAWATT,
		Symbol: "TW",
	},

	// Voltage Units
	MILLIVOLT: {
		Code:   MILLIVOLT,
		Symbol: "mV",
	},
	CENTIVOLT: {
		Code:   CENTIVOLT,
		Symbol: "cV",
	},
	DECIVOLT: {
		Code:   DECIVOLT,
		Symbol: "dV",
	},
	VOLT: {
		Code:   VOLT,
		Symbol: "V",
	},
	DEKAVOLT: {
		Code:   DEKAVOLT,
		Symbol: "daV",
	},
	HECTOVOLT: {
		Code:   HECTOVOLT,
		Symbol: "hV",
	},
	KILOVOLT: {
		Code:   KILOVOLT,
		Symbol: "kV",
	},

	// Intensity
	MILLIAMPERE: {
		Code:   MILLIAMPERE,
		Symbol: "mA",
	},
	CENTIAMPERE: {
		Code:   CENTIAMPERE,
		Symbol: "cA",
	},
	DECIAMPERE: {
		Code:   DECIAMPERE,
		Symbol: "dA",
	},
	AMPERE: {
		Code:   AMPERE,
		Symbol: "A",
	},
	DEKAMPERE: {
		Code:   "DEKAMPERE",
		Symbol: "daA",
	},
	HECTOAMPERE: {
		Code:   "HECTOAMPERE",
		Symbol: "hA",
	},
	KILOAMPERE: {
		Code:   "KILOAMPERE",
		Symbol: "kA",
	},

	// Resistance
	MILLIOHM: {
		Code:   MILLIOHM,
		Symbol: "mΩ",
	},
	CENTIOHM: {
		Code:   CENTIOHM,
		Symbol: "cΩ",
	},
	DECIOHM: {
		Code:   DECIOHM,
		Symbol: "dΩ",
	},
	OHM: {
		Code:   OHM,
		Symbol: "Ω",
	},
	DEKAOHM: {
		Code:   DEKAOHM,
		Symbol: "daΩ",
	},
	HECTOHM: {
		Code:   HECTOHM,
		Symbol: "hΩ",
	},
	KILOHM: {
		Code:   KILOHM,
		Symbol: "kΩ",
	},
	MEGOHM: {
		Code:   MEGOHM,
		Symbol: "MΩ",
	},

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
