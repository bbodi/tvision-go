package tvision

/*
import (
	"github.com/nsf/termbox-go"
)

type TermBoxEngine struct {
	cellBuffer []Cell
	w, h       int
}

func (self *TermBoxEngine) LoadFont(int) *Font {
	return &Font{1, 1}
}
func (self *TermBoxEngine) Clear(fg, bg Attribute) {
	termbox.Clear(toTermBoxColor(fg), toTermBoxColor(bg))
	self.syncCells()
}

func (self *TermBoxEngine) DrawToBackBuffer(drawBuffer *DrawBuffer, x, y, fontSize int) {
	termBoxCellBuffer := termbox.CellBuffer()
	srcW := drawBuffer.w
	dstW := self.w
	for srcX := 0; srcX < drawBuffer.w; srcX++ {
		for srcY := 0; srcY < drawBuffer.h; srcY++ {
			srcIndex := srcY*srcW + srcX
			dstX := x + srcX
			dstY := y + srcY
			if self.outOfRange(dstX, dstY) {
				continue
			}
			dstIndex := dstY*dstW + dstX
			srcCell := drawBuffer.cells[srcIndex]
			dstCell := &termBoxCellBuffer[dstIndex]

			if srcCell.Fg != ColorNone {
				dstCell.Fg = toTermBoxColor(srcCell.Fg)
			}
			if srcCell.Bg != ColorNone {
				dstCell.Bg = toTermBoxColor(srcCell.Bg)
			}
			if srcCell.Ch != 0 {
				dstCell.Ch = srcCell.Ch
			}
		}
	}
}

func (self *TermBoxEngine) SwapBackBuffer() {
	termbox.Flush()
}

func (self *TermBoxEngine) outOfRange(x, y int) bool {
	return x >= self.w || y >= self.h || y < 0 || x < 0
}

func (self *TermBoxEngine) Init() {
	err := termbox.Init()
	if err != nil {
		panic(err)
	}
	termbox.SetInputMode(termbox.InputEsc | termbox.InputMouse)
	self.w, self.h = self.Size()
	self.cellBuffer = make([]Cell, self.w*self.h)
}

func (self *TermBoxEngine) Close() {
	termbox.Close()
}

func (self *TermBoxEngine) Size() (int, int) {
	return termbox.Size()
}

func (self *TermBoxEngine) PollEvent() *Event {
	teve := termbox.PollEvent()
	event := new(Event)
	switch teve.Type {
	case termbox.EventKey:
		event.Type = EvKey
		event.Ch = teve.Ch
		event.Key = self.convertKey(teve.Key)
		if teve.Mod&0x01 != 0 {
			event.Mod.Alt = true
		}
	case termbox.EventMouse:
		event.Type = EvSystemMouse
		//event.MouseX = teve.MouseX
		//event.MouseY = teve.MouseY
	}
	return event
}

func toTermBoxColor(color Attribute) termbox.Attribute {
	switch color {
	case ColorBlack:
		return termbox.ColorBlack
	case ColorRed:
		return termbox.ColorRed
	case ColorGreen:
		return termbox.ColorGreen
	case ColorYellow:
		return termbox.ColorYellow
	case ColorBlue:
		return termbox.ColorBlue
	case ColorMagenta:
		return termbox.ColorMagenta
	case ColorCyan:
		return termbox.ColorCyan
	case ColorWhite:
		return termbox.ColorWhite
	}
	return termbox.ColorDefault

}

func (self *TermBoxEngine) convertKey(tkey termbox.Key) Key {
	switch tkey {
	case termbox.KeyF1:
		return KeyF1
	case termbox.KeyF2:
		return KeyF2
	case termbox.KeyF3:
		return KeyF3
	case termbox.KeyF4:
		return KeyF4
	case termbox.KeyF5:
		return KeyF5
	case termbox.KeyF6:
		return KeyF6
	case termbox.KeyF7:
		return KeyF7
	case termbox.KeyF8:
		return KeyF8
	case termbox.KeyF9:
		return KeyF9
	case termbox.KeyF10:
		return KeyF10
	case termbox.KeyF11:
		return KeyF11
	case termbox.KeyF12:
		return KeyF12

	case termbox.KeyInsert:
		return KeyInsert
	case termbox.KeyDelete:
		return KeyDelete
	case termbox.KeyHome:
		return KeyHome
	case termbox.KeyEnd:
		return KeyEnd
	case termbox.KeyPgup:
		return KeyPgup
	case termbox.KeyPgdn:
		return KeyPgdn
	case termbox.KeyArrowUp:
		return KeyArrowUp
	case termbox.KeyArrowDown:
		return KeyArrowDown
	case termbox.KeyArrowLeft:
		return KeyArrowLeft
	case termbox.KeyArrowRight:
		return KeyArrowRight
	case termbox.KeyTab:
		return KeyTab
	case termbox.KeyEnter:
		return KeyEnter
	case termbox.KeyBackspace:
		return KeyBackspace
	case termbox.KeySpace:
		return KeySpace
	case termbox.KeyEsc:
		return KeyEsc
	default:
		return 0
	}
}

func (self *TermBoxEngine) syncCells() {
	for i, cell := range termbox.CellBuffer() {
		self.cellBuffer[i] = Cell{cell.Ch, Attribute(cell.Fg), Attribute(cell.Bg)}
	}
}
*/
