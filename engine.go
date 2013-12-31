package tvision

type Pixel int

func (self Pixel) Mult(a int) Pixel {
	return Pixel(int(self) * a)
}

type Font struct {
	fontSizeX, fontSizeY Pixel
}

func (self *Font) CharWidth() Pixel {
	return self.fontSizeX
}

func (self *Font) CharHeight() Pixel {
	return self.fontSizeY
}

type Engine interface {
	Clear(fg, bg Color)
	DrawToBackBuffer(drawBuffer *DrawBuffer, view *View)
	SwapBackBuffer()
	Init(w, h int, fontSize Pixel)
	Close()
	Size() (int, int)
	PollEvent() *Event
	LoadFont(fontSize Pixel) *Font
}

var engine Engine

/*
func InitTermBoxEngine() {
	engine = new(TermBoxEngine)
	engine.Init()
}*/

func InitSdlEngine(w, h int, fontSize Pixel) {
	engine = new(SdlEngine)
	engine.Init(w, h, fontSize)
}
