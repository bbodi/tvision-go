package tvision

import (
	"container/list"
	"log"
	"os"
)

var focusedView *View

type Rect struct {
	X, Y int
	W, H int
}

func (rect Rect) Grow(addX, addY int) Rect {
	rect.W += addX
	rect.H += addY
	return rect
}

func (rect Rect) Move(addX, addY int) Rect {
	rect.X += addX
	rect.Y += addY
	return rect
}

func (rect Rect) Contains(x, y int) bool {
	containsX := x >= rect.X && x < rect.X+rect.W
	containsY := y >= rect.Y && y < rect.Y+rect.H
	return containsX && containsY
}

type View struct {
	Rect            Rect
	Owner           *View
	Widget          Widget
	ExecutingResult ExecutingResult
	executing       bool
	views           list.List
	dirty           bool
	hasDirtyChild   bool
	hidden          bool
	focused         bool
	focusedElement  *list.Element
	fontSize        Pixel
}

func (self *View) FontSize() Pixel {
	if self.fontSize != 0 || self.Owner == nil {
		return self.fontSize
	}
	return self.Owner.FontSize()
}

func (self *View) SetFontSize(p int) {
	self.fontSize = Pixel(p)
}

type ExecutingResult struct {
	Cmd  Cmd
	Data interface{}
}

type Widget interface {
	Draw(view *View) *DrawBuffer
	HandleEvent(*Event, *View)
}

func (self *View) Hidden() bool {
	return self.hidden
}

func (self *View) Hide() {
	self.hidden = true
	self.Modified()
}

func (self *View) Show() {
	self.hidden = false
	self.Modified()
}

func (self *View) GetExecutingResult() ExecutingResult {
	return self.ExecutingResult
}

func (self *View) StartExecuting() {
	self.executing = true
}

func (self *View) StopExecuting(executingResult ExecutingResult) {
	if self.isExecuting() == false {
		Error("View is not in Executing: ", self)
	}
	self.ExecutingResult = executingResult
	self.executing = false
}

func (self *View) isExecuting() bool {
	return self.executing
}

func (self *View) Modified() {
	self.dirty = true
	for parent := self.Owner; parent != nil; parent = parent.Owner {
		parent.hasDirtyChild = true
	}
}

func (self *View) clickedInMe(event *Event) bool {
	if self.Rect.X == 25 && self.Rect.Y == 10 {
		self.Rect = self.Rect
	}
	mx, my := self.makeLocal(event.MouseX, event.MouseY)
	ok := mx < self.Rect.W && my < self.Rect.H && mx >= 0 && my >= 0
	log.Println("Click: ", ok, ": ", event.MouseX, ", ", event.MouseY, " -> ", mx, ", ", my)
	return mx < self.Rect.W && my < self.Rect.H && mx >= 0 && my >= 0
}

func (self *View) makeLocal(x, y Pixel) (int, int) {
	offX := self.calcOffsetX()
	offY := self.calcOffsetY()
	x -= offX
	y -= offY
	chW := engine.LoadFont(self.FontSize()).CharWidth()
	chH := engine.LoadFont(self.FontSize()).CharHeight()
	var localX int
	if x < 0 {
		localX = -1
	} else {
		localX = int(x / chW)
	}
	var localY int
	if y < 0 {
		localY = -1
	} else {
		localY = int(y / chH)
	}
	return localX, localY
}

func (self *View) HandleEvent(event *Event) {
	if self.hidden {
		return
	}
	self.doOnChildrenBackToFront(func(v *View) bool {
		v.HandleEvent(event)
		if event.Type == EvNothing {
			return false
		}
		return true
	})
	if event.Type == EvSystemMouse {
		if self.clickedInMe(event) {
			setFocusedView(self)
			self.Modified()
			event.Type = EvMouse
			mx, my := self.makeLocal(event.MouseX, event.MouseY)
			event.LocalMouseX = mx
			event.LocalMouseY = my
		}
	}
	Trace("View.HandleEvent: ", self.Widget)
	self.Widget.HandleEvent(event, self)
}

func (self *View) doOnChildrenBackToFront(action func(v *View) bool) {
	for e := self.views.Back(); e != nil; e = e.Prev() {
		v := e.Value.(*View)
		if cont := action(v); cont == false {
			break
		}
	}
}

func (self *View) doOnChildrenFrontToBack(action func(v *View) bool) {
	for e := self.views.Front(); e != nil; e = e.Next() {
		v := e.Value.(*View)
		if cont := action(v); cont == false {
			break
		}
	}
}

func (self *View) calcOffsetX() Pixel {
	var fontSize Pixel
	var parentOffset Pixel
	if self.Owner == nil {
		fontSize = DefaultFontSize
		parentOffset = Pixel(0)
	} else {
		fontSize = self.Owner.FontSize()
		parentOffset = self.Owner.calcOffsetX()
	}
	selfOffset := engine.LoadFont(fontSize).CharWidth().Mult(self.Rect.X)
	return parentOffset + selfOffset
}

func (self *View) calcOffsetY() Pixel {
	var fontSize Pixel
	var parentOffset Pixel
	if self.Owner == nil {
		fontSize = DefaultFontSize
		parentOffset = Pixel(0)
	} else {
		fontSize = self.Owner.FontSize()
		parentOffset = self.Owner.calcOffsetY()
	}
	selfOffset := engine.LoadFont(fontSize).CharHeight().Mult(self.Rect.Y)
	return parentOffset + selfOffset
}

func (self *View) ExecuteView(x, y int, view *View) ExecutingResult {
	if self == view {
		panic("self.ExecuteView(self)")
	}
	Trace("ExecuteView: ", view)
	self.AddView(x, y, view)
	view.StartExecuting()
	savedFocusedView := focusedView
	for view.isExecuting() {
		if view.dirty || view.hasDirtyChild {
			regions := make([]*screenRegion, 0)
			self.Draw(&regions)
			DrawToBackBuffer(regions)
			engine.SwapBackBuffer()
		}
		event := pollEvent()
		view.HandleEvent(event)
	}
	self.RemoveView(view)
	setFocusedView(savedFocusedView)
	return view.GetExecutingResult()
}

func (self *View) Execute() ExecutingResult {
	self.Modified()
	self.StartExecuting()
	for self.isExecuting() {
		if self.dirty || self.hasDirtyChild {
			regions := make([]*screenRegion, 0)
			self.Draw(&regions)
			DrawToBackBuffer(regions)
			engine.SwapBackBuffer()
		}
		event := pollEvent()
		self.HandleEvent(event)
	}
	return self.ExecutingResult
}

func (self *View) X() int {
	if self.Owner != nil {
		return self.Owner.X() + self.Rect.X
	}
	return self.Rect.X
}

func (self *View) Y() int {
	if self.Owner != nil {
		return self.Owner.Y() + self.Rect.Y
	}
	return self.Rect.Y
}

func (self *View) Draw(regions *[]*screenRegion) {
	if self.hidden {
		return
	}
	buff := self.Widget.Draw(self)
	*regions = append(*regions, &screenRegion{buff, self})

	self.doOnChildrenFrontToBack(func(v *View) bool {
		v.Draw(regions)
		return true
	})
	self.dirty = false
	self.hasDirtyChild = false
}

func pollEvent() *Event {
	event := engine.PollEvent()
	return event
}

func (self *View) BroadcastCommand(cmd Cmd, data interface{}) {
	self.Broadcast(EvCommand, cmd, data)
}

func (self *View) Broadcast(typ EventType, cmd Cmd, data interface{}) {
	if self.Owner == nil {
		event := new(Event)
		event.Type = typ
		event.Cmd = cmd
		event.Data = data
		self.HandleEvent(event)
	} else {
		self.Owner.BroadcastCommand(cmd, data)
	}
}

func (self *View) AddView(x, y int, view *View) {
	if view.Owner != nil {
		Warn(view, " has a Parent!")
	}
	view.Rect.X = x
	view.Rect.Y = y
	self.views.PushBack(view)
	view.Owner = self
	setFocusedView(view)
	view.Modified()
}

func (self *View) RemoveView(view *View) {
	e := self.findViewElement(view)
	v := e.Value.(*View)
	self.views.Remove(e)
	setFocusedView(self.LastView())
	v.Owner = nil
	self.Modified()
}

func (self *View) findViewElement(view *View) *list.Element {
	for e := self.views.Front(); e != nil; e = e.Next() {
		v := e.Value.(*View)
		if view == v {
			return e
		}
	}
	return nil
}

func (self *View) LastView() *View {
	if self.views.Len() == 0 {
		return nil
	}
	return self.views.Back().Value.(*View)
}

func setFocusedView(newFocusedView *View) {
	if focusedView != nil {
		focusedView.setFocused(false)
		focusedView = nil
	}
	focusedView = newFocusedView
	if newFocusedView != nil {
		newFocusedView.setFocused(true)
		if newFocusedView.Owner != nil {
			newFocusedView.BringToFront()
		}
	}
}

func (self *View) setFocused(focused bool) {
	self.focused = focused
	if focused {
		self.Widget.HandleEvent(&Event{Type: EvGetFocus}, self)
	} else {
		self.Widget.HandleEvent(&Event{Type: EvLostFocus}, self)
	}
	if self.Owner != nil {
		self.Owner.setFocused(focused)
	}
}

func (self *View) BringToFront() {
	if self.Owner == nil {
		return
	}
	self.Owner.makeLast(self)
	self.Owner.BringToFront()
}

func (self *View) makeLast(child *View) {
	e := self.findViewElement(child)
	self.views.MoveToBack(e)
}

func (self *View) NextView() {
	if self.views.Len() == 0 {
		os.Exit(1)
		return
	}
	if self.focusedElement == nil {
		self.focusedElement = self.views.Front()
	} else if self.focusedElement.Next() == nil {
		self.focusedElement = self.views.Front()
	} else {
		self.focusedElement = self.focusedElement.Next()
	}
	setFocusedView(self.focusedElement.Value.(*View))
}
