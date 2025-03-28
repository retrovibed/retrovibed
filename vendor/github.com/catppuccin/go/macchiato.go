package catppuccingo

// Macchiato variant
type macchiato struct{}

// Macchiato flavor variant
var Macchiato Flavor = macchiato{}

// Macchiato
func (macchiato) Name() string { return "macchiato" }

// Rosewater
func (macchiato) Rosewater() Color {
	return Color{
		Hex: "#f4dbd6",
		RGB: [3]uint8{244, 219, 214},
		HSL: [3]float32{10, 0.58, 0.9},
	}
}

// Flamingo
func (macchiato) Flamingo() Color {
	return Color{
		Hex: "#f0c6c6",
		RGB: [3]uint8{240, 198, 198},
		HSL: [3]float32{0, 0.58, 0.86},
	}
}

// Pink
func (macchiato) Pink() Color {
	return Color{
		Hex: "#f5bde6",
		RGB: [3]uint8{245, 189, 230},
		HSL: [3]float32{316, 0.74, 0.85},
	}
}

// Mauve
func (macchiato) Mauve() Color {
	return Color{
		Hex: "#c6a0f6",
		RGB: [3]uint8{198, 160, 246},
		HSL: [3]float32{267, 0.83, 0.8},
	}
}

// Red
func (macchiato) Red() Color {
	return Color{
		Hex: "#ed8796",
		RGB: [3]uint8{237, 135, 150},
		HSL: [3]float32{351, 0.74, 0.73},
	}
}

// Maroon
func (macchiato) Maroon() Color {
	return Color{
		Hex: "#ee99a0",
		RGB: [3]uint8{238, 153, 160},
		HSL: [3]float32{355, 0.71, 0.77},
	}
}

// Peach
func (macchiato) Peach() Color {
	return Color{
		Hex: "#f5a97f",
		RGB: [3]uint8{245, 169, 127},
		HSL: [3]float32{21, 0.86, 0.73},
	}
}

// Yellow
func (macchiato) Yellow() Color {
	return Color{
		Hex: "#eed49f",
		RGB: [3]uint8{238, 212, 159},
		HSL: [3]float32{40, 0.7, 0.78},
	}
}

// Green
func (macchiato) Green() Color {
	return Color{
		Hex: "#a6da95",
		RGB: [3]uint8{166, 218, 149},
		HSL: [3]float32{105, 0.48, 0.72},
	}
}

// Teal
func (macchiato) Teal() Color {
	return Color{
		Hex: "#8bd5ca",
		RGB: [3]uint8{139, 213, 202},
		HSL: [3]float32{171, 0.47, 0.69},
	}
}

// Sky
func (macchiato) Sky() Color {
	return Color{
		Hex: "#91d7e3",
		RGB: [3]uint8{145, 215, 227},
		HSL: [3]float32{189, 0.59, 0.73},
	}
}

// Sapphire
func (macchiato) Sapphire() Color {
	return Color{
		Hex: "#7dc4e4",
		RGB: [3]uint8{125, 196, 228},
		HSL: [3]float32{199, 0.66, 0.69},
	}
}

// Blue
func (macchiato) Blue() Color {
	return Color{
		Hex: "#8aadf4",
		RGB: [3]uint8{138, 173, 244},
		HSL: [3]float32{220, 0.83, 0.75},
	}
}

// Lavender
func (macchiato) Lavender() Color {
	return Color{
		Hex: "#b7bdf8",
		RGB: [3]uint8{183, 189, 248},
		HSL: [3]float32{234, 0.82, 0.85},
	}
}

// Text
func (macchiato) Text() Color {
	return Color{
		Hex: "#cad3f5",
		RGB: [3]uint8{202, 211, 245},
		HSL: [3]float32{227, 0.68, 0.88},
	}
}

// Subtext 1
func (macchiato) Subtext1() Color {
	return Color{
		Hex: "#b8c0e0",
		RGB: [3]uint8{184, 192, 224},
		HSL: [3]float32{228, 0.39, 0.8},
	}
}

// Subtext 0
func (macchiato) Subtext0() Color {
	return Color{
		Hex: "#a5adcb",
		RGB: [3]uint8{165, 173, 203},
		HSL: [3]float32{227, 0.27, 0.72},
	}
}

// Overlay 2
func (macchiato) Overlay2() Color {
	return Color{
		Hex: "#939ab7",
		RGB: [3]uint8{147, 154, 183},
		HSL: [3]float32{228, 0.2, 0.65},
	}
}

// Overlay 1
func (macchiato) Overlay1() Color {
	return Color{
		Hex: "#8087a2",
		RGB: [3]uint8{128, 135, 162},
		HSL: [3]float32{228, 0.15, 0.57},
	}
}

// Overlay 0
func (macchiato) Overlay0() Color {
	return Color{
		Hex: "#6e738d",
		RGB: [3]uint8{110, 115, 141},
		HSL: [3]float32{230, 0.12, 0.49},
	}
}

// Surface 2
func (macchiato) Surface2() Color {
	return Color{
		Hex: "#5b6078",
		RGB: [3]uint8{91, 96, 120},
		HSL: [3]float32{230, 0.14, 0.41},
	}
}

// Surface 1
func (macchiato) Surface1() Color {
	return Color{
		Hex: "#494d64",
		RGB: [3]uint8{73, 77, 100},
		HSL: [3]float32{231, 0.16, 0.34},
	}
}

// Surface 0
func (macchiato) Surface0() Color {
	return Color{
		Hex: "#363a4f",
		RGB: [3]uint8{54, 58, 79},
		HSL: [3]float32{230, 0.19, 0.26},
	}
}

// Base
func (macchiato) Base() Color {
	return Color{
		Hex: "#24273a",
		RGB: [3]uint8{36, 39, 58},
		HSL: [3]float32{232, 0.23, 0.18},
	}
}

// Mantle
func (macchiato) Mantle() Color {
	return Color{
		Hex: "#1e2030",
		RGB: [3]uint8{30, 32, 48},
		HSL: [3]float32{233, 0.23, 0.15},
	}
}

// Crust
func (macchiato) Crust() Color {
	return Color{
		Hex: "#181926",
		RGB: [3]uint8{24, 25, 38},
		HSL: [3]float32{236, 0.23, 0.12},
	}
}
