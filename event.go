package tvision

type Cmd int

const (
	CmdNothing Cmd = iota
	CmdQuit
	CmdCancel
	CmdOk

	CmdForUserApp
)

type EventType int

const (
	EvNothing EventType = iota
	EvResize
	EvCommand
	EvMouse
	EvSystemMouse
	EvKey
	EvKeyRepeat
	EvLostFocus
	EvGetFocus
	EvTick
)

type Key uint16

const (
	KeyNormal Key = iota
	KeyF1
	KeyF2
	KeyF3
	KeyF4
	KeyF5
	KeyF6
	KeyF7
	KeyF8
	KeyF9
	KeyF10
	KeyF11
	KeyF12
	KeyInsert
	KeyDelete
	KeyHome
	KeyEnd
	KeyPgup
	KeyPgdn
	KeyArrowUp
	KeyArrowDown
	KeyArrowLeft
	KeyArrowRight
	KeyTab
	KeyEnter
	KeyCtrl
	KeyAlt
	KeyBackspace
	KeySpace
	KeyEsc
)

type Modifier struct {
	Lalt, Ralt, Alt       bool
	Lctrl, Rctrl, Ctrl    bool
	Lshift, Rshift, Shift bool
	CapsLock, NumLock     bool
}

type Event struct {
	Type                     EventType
	Cmd                      Cmd
	MouseX, MouseY           Pixel
	DoubleClick              bool
	LocalMouseX, LocalMouseY int
	Ch                       rune
	Key                      Key
	Mod                      Modifier
	Data                     interface{}
}

func (event *Event) SetProcessed() {
	if event.Type == EvGetFocus || event.Type == EvLostFocus || event.Type == EvTick {
		panic("ClearEvent(focus)")
	}
	event.Type = EvNothing
}
