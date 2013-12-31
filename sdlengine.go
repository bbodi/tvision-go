package tvision

import (
	//	"fmt"
	"github.com/DeedleFake/sdl"
	"github.com/DeedleFake/sdl/ttf"
	"time"
	"unicode/utf8"
)

var DefaultFontSize Pixel

type SdlEngine struct {
	renderer                           *sdl.Renderer
	window                             *sdl.Window
	w, h                               int
	screenSurface                      *sdl.Surface
	font                               *sdlFont
	pixelFormat                        *sdl.PixelFormat
	nextKeyEventTick                   time.Time
	fonts                              map[Pixel]*sdlFont
	doubleClickDelay                   time.Duration
	nextTickWhileDoubleClickCouldOccur *time.Time
	lastEvent                          *sdl.Event
}

type sdlFont struct {
	Font
	font *ttf.Font
}

func (self *SdlEngine) Init(w, h int, fontSize Pixel) {
	DefaultFontSize = fontSize
	self.doubleClickDelay, _ = time.ParseDuration("500ms")
	self.fonts = make(map[Pixel]*sdlFont)
	err := ttf.Init()
	if err != nil {
		panic(err)
	}
	self.LoadFont(DefaultFontSize)
	self.font = self.fonts[DefaultFontSize]

	self.w = w
	self.h = h
	screenW, screenH := self.w*int(self.font.fontSizeX), self.h*int(self.font.fontSizeY)

	err = sdl.Init(sdl.INIT_EVERYTHING)
	if err != nil {
		panic(err)
	}
	win, err := sdl.CreateWindow("Hello World!", 100, 100, screenW, screenH, sdl.WINDOW_SHOWN)
	if err != nil {
		panic(err)
	}
	renderer, err := win.CreateRenderer(-1, sdl.RENDERER_ACCELERATED|sdl.RENDERER_PRESENTVSYNC)
	_ = renderer
	if err != nil {
		panic(err)
	}

	self.renderer = renderer
	self.window = win
	self.pixelFormat, err = sdl.AllocFormat(sdl.PIXELFORMAT_RGBA8888)
	if err != nil {
		panic(err)
	}
	self.screenSurface, err = sdl.CreateRGBSurface(
		screenH, screenH,
		32,
		self.pixelFormat.Rmask,
		self.pixelFormat.Gmask,
		self.pixelFormat.Bmask,
		self.pixelFormat.Amask)
	if err != nil {
		panic(err)
	}
}

func (self *SdlEngine) LoadFont(size Pixel) *Font {
	if size == 0 {
		size = DefaultFontSize
	}
	font, ok := self.fonts[size]
	if ok == false {
		ttfFont, err := ttf.OpenFont("DejaVuSansMono.ttf", int(size))
		if err != nil {
			panic(err)
		}
		fontSizeX, fontSizeY, _ := ttfFont.SizeText("A")
		font := &sdlFont{Font{Pixel(fontSizeX), Pixel(fontSizeY)}, ttfFont}
		self.fonts[size] = font
		return &font.Font
	}
	return &font.Font
}

func (self *SdlEngine) Clear(fg, bg Color) {
	self.renderer.Clear()
}

func (self *SdlEngine) DrawToBackBuffer(drawBuffer *DrawBuffer, view *View) {
	self.drawBackground(drawBuffer, view)
	self.drawCharacters(drawBuffer, view)
}

func (self *SdlEngine) ttfFont(fontSize Pixel) *ttf.Font {
	if fontSize == 0 {
		fontSize = DefaultFontSize
	}
	self.LoadFont(fontSize)
	return self.fonts[fontSize].font
}

func (self *SdlEngine) drawCharacters(drawBuffer *DrawBuffer, view *View) {
	x, y := 0, 0
	var offsetX Pixel = view.calcOffsetX()
	var offsetY Pixel = view.calcOffsetY()
	font := self.LoadFont(view.FontSize())
	fw := font.CharWidth()
	fh := font.CharHeight()
	for y < drawBuffer.h {
		newX, word, color := getContinousCharacters(drawBuffer, x, y)
		if utf8.RuneCountInString(word) > 0 {
			w := newX - x
			xpos := offsetX + fw.Mult(x)
			ypos := offsetY + fh.Mult(y)
			self.drawString(word, color, xpos, ypos, fw.Mult(w), fh, self.ttfFont(view.FontSize()))
		}
		x = newX
		if x >= drawBuffer.w {
			y++
			x = 0
		}
	}
}

func (self *SdlEngine) drawString(str string, color Color, x, y, w, h Pixel, font *ttf.Font) {
	surf, err := font.RenderUTF8Solid(str, toSdlColor(color, 255))
	if err != nil {
		panic(err)
	}
	defer surf.Free()
	texture, err := self.renderer.CreateTextureFromSurface(surf)
	if err != nil {
		panic(err)
	}
	self.renderer.Copy(texture, nil, &sdl.Rect{int32(x), int32(y), int32(w), int32(h)})
}

func (self *SdlEngine) SwapBackBuffer() {
	self.renderer.Present()
}

func toSdlRect(r Rect, font *Font) *sdl.Rect {
	x := int32(r.X * int(font.fontSizeX))
	y := int32(r.Y * int(font.fontSizeY))
	w := int32(r.W * int(font.fontSizeX))
	h := int32(r.H * int(font.fontSizeY))
	return &sdl.Rect{x, y, w, h}
}

func (self *SdlEngine) drawRect(color Color, x, y, w, h Pixel) {
	self.renderer.SetDrawColor(toSdlRgba(color, 255))
	self.renderer.FillRect(&sdl.Rect{int32(x), int32(y), int32(w), int32(h)})
}

func toSdlColor(attr Color, alpha uint8) sdl.Color {
	r, g, b, _ := toSdlRgba(attr, alpha)
	var color sdl.Color
	color.R = r
	color.G = g
	color.B = uint(b)
	return color
}

func toSdlRgba(color Color, alpha uint8) (uint8, uint8, uint8, uint8) {
	return color.R, color.G, color.B, alpha
}

func (self *SdlEngine) drawBackground(drawBuffer *DrawBuffer, view *View) {
	srcX, srcY := 0, 0
	var offsetX Pixel = view.calcOffsetX()
	var offsetY Pixel = view.calcOffsetY()
	font := self.LoadFont(view.FontSize())
	fw := font.CharWidth()
	fh := font.CharHeight()
	for {
		region := getSameColorRegions(drawBuffer, srcX, srcY)

		if region.color != ColorNone {
			x := fw.Mult(srcX) + offsetX
			y := fh.Mult(srcY) + offsetY
			w := fw.Mult(region.rect.W)
			h := fh.Mult(region.rect.H)
			self.drawRect(region.color, x, y, w, h)
		}
		srcX += region.rect.W
		srcY += region.rect.H - 1

		if srcX >= drawBuffer.w {
			srcX = 0
			srcY += 1
		}
		if srcY >= drawBuffer.h {
			break
		}
	}

}

type colorRegion struct {
	color Color
	rect  Rect
}

func getSameColorRegions(drawBuffer *DrawBuffer, x, y int) colorRegion {
	var region colorRegion
	region.rect.W = 0
	region.rect.H = 1
	srcW := drawBuffer.w
	region.color = drawBuffer.Cell(x, y).Bg

	srcX := x
	srcY := y
	for srcY < drawBuffer.h {
		srcCell := drawBuffer.Cell(srcX, srcY)
		if srcCell.Bg != region.color {
			notSameColorLine := srcX != srcW
			if notSameColorLine {
				hasMoreRow := region.rect.H != 1
				if hasMoreRow {
					region.rect.H -= 1
					region.rect.W = srcW
				}
			}
			return region
		}
		srcX++
		region.rect.W += 1
		if srcX >= drawBuffer.w {
			if x > 0 {
				return region
			}
			srcX = 0
			region.rect.W = 0
			region.rect.H += 1
			srcY++
		}
	}

	region.rect.H -= 1
	region.rect.W = srcW
	return region
}

func (self *SdlEngine) outOfRange(x, y int) bool {
	return x >= self.w || y >= self.h || y < 0 || x < 0
}

func getContinousCharacters(drawBuffer *DrawBuffer, offX, offY int) (int, string, Color) {
	i := offY*drawBuffer.w + offX
	lineEnd := i + drawBuffer.w - offX
	cells := drawBuffer.cells[i:lineEnd]
	fg := cells[0].Fg
	var word string
	for _, cell := range cells {
		if cell.Fg != fg {
			break
		}
		var ch string
		if cell.Ch == 0 {
			ch = " "
		} else {
			ch = string(cell.Ch)
		}
		word = word + ch
		offX++
	}
	return offX, word, fg
}

func (self *SdlEngine) Close() {
	sdl.Quit()
}

func (self *SdlEngine) Size() (int, int) {
	return self.w, self.h
}

func (self *SdlEngine) PollEvent() *Event {
	for {
		sdlEvent := self.readSdlEvent()
		self.lastEvent = sdlEvent
		event := self.processSdlEvent(sdlEvent)
		if event != nil {
			return event
		}
	}
}

func (self *SdlEngine) readSdlEvent() *sdl.Event {
	var sdlEvent sdl.Event
	ok := false
	for ok == false {
		ok = sdl.PollEvent(&sdlEvent)
		switch t := sdlEvent.(type) {
		case *sdl.KeyboardEvent:
			if t.State == 0 {
				ok = false
			}
		case *sdl.MouseButtonEvent:
			if t.State == 0 {
				ok = false
			}
		}
		sdl.Delay(10)
	}
	return &sdlEvent
}

func (self *SdlEngine) processSdlEvent(sdlEvent *sdl.Event) *Event {
	event := new(Event)
	switch t := (*sdlEvent).(type) {
	case *sdl.TextInputEvent:
		ch, _ := utf8.DecodeRune(t.Text[:])
		event.Type = EvKey
		event.Ch = ch
	case *sdl.KeyboardEvent:
		if t.Repeat == 1 {
			event.Type = EvKeyRepeat
		} else {
			event.Type = EvKey
		}
		event.Key = self.convertKey(t)
		if event.Key == KeyNormal {
			return nil
		}
		event.Mod = self.convertModificationButtons(t)
	case *sdl.MouseButtonEvent:
		event.Type = EvSystemMouse
		event.MouseX = Pixel(t.X)
		event.MouseY = Pixel(t.Y)
		if self.nextTickWhileDoubleClickCouldOccur != nil && self.nextTickWhileDoubleClickCouldOccur.After(time.Now()) {
			event.DoubleClick = true
			self.nextTickWhileDoubleClickCouldOccur = nil
		} else {
			t := time.Now().Add(self.doubleClickDelay)
			self.nextTickWhileDoubleClickCouldOccur = &t
		}
	case *sdl.QuitEvent:
		event.Type = EvCommand
		event.Cmd = CmdQuit
	}
	return event
}

func (self *SdlEngine) convertModificationButtons(event *sdl.KeyboardEvent) Modifier {
	mod := uint16(event.Keysym.Mod)
	return Modifier{
		Lalt:     mod&uint16(sdl.KMOD_LALT) == 0,
		Ralt:     mod&uint16(sdl.KMOD_RALT) == 0,
		Alt:      mod&uint16(sdl.KMOD_ALT) == 0,
		Lctrl:    mod&uint16(sdl.KMOD_LCTRL) == 0,
		Rctrl:    mod&uint16(sdl.KMOD_RCTRL) == 0,
		Ctrl:     mod&uint16(sdl.KMOD_CTRL) == 0,
		Lshift:   mod&uint16(sdl.KMOD_LSHIFT) == 0,
		Rshift:   mod&uint16(sdl.KMOD_RSHIFT) == 0,
		Shift:    mod&uint16(sdl.KMOD_SHIFT) == 0,
		CapsLock: mod&uint16(sdl.KMOD_CAPS) == 0,
		NumLock:  mod&uint16(sdl.KMOD_NUM) == 0,
	}
}

func (self *SdlEngine) convertKey(event *sdl.KeyboardEvent) Key {
	switch event.Keysym.Sym {
	case sdl.K_F1:
		return KeyF1
	case sdl.K_F2:
		return KeyF2
	case sdl.K_F3:
		return KeyF3
	case sdl.K_F4:
		return KeyF4
	case sdl.K_F5:
		return KeyF5
	case sdl.K_F6:
		return KeyF6
	case sdl.K_F7:
		return KeyF7
	case sdl.K_F8:
		return KeyF8
	case sdl.K_F9:
		return KeyF9
	case sdl.K_F10:
		return KeyF10
	case sdl.K_F11:
		return KeyF11
	case sdl.K_F12:
		return KeyF12

	case sdl.K_INSERT:
		return KeyInsert
	case sdl.K_DELETE:
		return KeyDelete
	case sdl.K_HOME:
		return KeyHome
	case sdl.K_END:
		return KeyEnd
	case sdl.K_PAGEUP:
		return KeyPgup
	case sdl.K_PAGEDOWN:
		return KeyPgdn
	case sdl.K_UP:
		return KeyArrowUp
	case sdl.K_DOWN:
		return KeyArrowDown
	case sdl.K_LEFT:
		return KeyArrowLeft
	case sdl.K_RIGHT:
		return KeyArrowRight
	case sdl.K_TAB:
		return KeyTab
	case sdl.K_RETURN:
		return KeyEnter
	case sdl.K_BACKSPACE:
		return KeyBackspace
	case sdl.K_SPACE:
		return KeySpace
	case sdl.K_ESCAPE:
		return KeyEsc
	default:
		return KeyNormal
	}
}
