package tvision

type Color struct {
	R, G, B uint8
	visible bool
}

var (
	ColorNone    = Color{0, 0, 0, false}
	ColorBlack   = Color{0, 0, 0, true}
	ColorRed     = Color{255, 0, 0, true}
	ColorGreen   = Color{0, 255, 0, true}
	ColorYellow  = Color{255, 255, 0, true}
	ColorBlue    = Color{0, 0, 255, true}
	ColorMagenta = Color{255, 0, 255, true}
	ColorCyan    = Color{0, 255, 255, true}
	ColorWhite   = Color{255, 255, 255, true}
	ColorGray    = Color{100, 100, 100, true}
)
