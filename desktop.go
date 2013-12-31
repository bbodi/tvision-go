package tvision

import (
	"fmt"
)

const logW = 20
const logH = 20

type Desktop struct {
	Running    bool
	drawBuffer *DrawBuffer
}

func (self *Desktop) Init(view *View) {
	_, loggingTextArea := CreateTextArea(logW, logH, "Logs")
	LoggingTextArea = loggingTextArea
	desktop = self
	desktopView = view
	w, h := engine.Size()
	self.drawBuffer = CreateDrawBuffer(w, h)
}

func (self *Desktop) Close() {
	engine.Close()
}

func (self *Desktop) String() string {
	return "Desktop"
}

func (self *Desktop) Draw(view *View) *DrawBuffer {
	self.drawBuffer.Clear(ColorWhite, ColorGray)
	return self.drawBuffer
}

func ToogleLogWindow() {
	w, _ := engine.Size()
	loggingView := LoggingTextArea.View
	if loggingView.Owner == nil {
		desktopView.AddView(w/2-logW/2, 1, loggingView)
		//LoggingTextArea.View.Show()
	} else {
		desktopView.RemoveView(loggingView)
		//LoggingTextArea.View.Hide()
	}
}

var desktop *Desktop
var desktopView *View

var LoggingWindow *Window
var LoggingTextArea *TextArea

//var view, LoggingWindow = CreateTextArea(1, 1, 20, 10, "Logs")

func Info(str ...interface{}) {
	if LoggingTextArea == nil {
		return
	}
	LoggingTextArea.PushFront(TextRegion{ColorWhite, ColorNone, fmt.Sprint(str)})
}

func Warn(str ...interface{}) {
	if LoggingTextArea == nil {
		return
	}
	LoggingTextArea.PushFront(TextRegion{ColorYellow, ColorNone, fmt.Sprint(str)})
}

func Error(str ...interface{}) {
	if LoggingTextArea == nil {
		return
	}
	LoggingTextArea.PushFront(TextRegion{ColorRed, ColorNone, fmt.Sprint(str)})
}

func Trace(str ...interface{}) {
	if LoggingTextArea == nil {
		return
	}
	LoggingTextArea.PushFront(TextRegion{ColorCyan, ColorNone, fmt.Sprint(str)})
}
