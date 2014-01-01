package main

import (
	"encoding/json"
	"fmt"
	"github.com/bbodi/tvision"
	"io/ioutil"
	"net/http"
	"os"
)

type WordRegion struct {
	Text          string
	Missing       bool
	AlwaysMissing bool
}

const (
	CmdAddNewPackage tvision.Cmd = tvision.CmdForUserApp + iota
	CmdChoosePkg
	CmdAddToPkg
	CmdEditPkg
	CmdStartPkg
)

const (
	Simple              = iota // a = b
	Double              = iota // a = b && b = a
	Quote               = iota //
	DefinedMissingWords = iota
)

type Card struct {
	Question   string
	WordRegion []WordRegion
	Type       int
	Tags       []string
}

type Package struct {
	Name  string
	Cards []Card
}

type App struct {
	tvision.Desktop
	Packages []Package
}

var app *App

func PackagesJson(w http.ResponseWriter, r *http.Request) {
	js, err := json.Marshal(app.Packages)
	if err != nil {
		fmt.Println(w, "error:", err)
	}
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.Write(js)
}

func persist() {
	js, _ := json.Marshal(app)
	os.Remove("mindy.json")
	err := ioutil.WriteFile("mindy.json", js, os.ModePerm)
	if err != nil {
		tvision.Error(err)
	}
}

func loadDb() *App {
	app := new(App)
	data, err := ioutil.ReadFile("mindy.json")
	if err != nil {
		tvision.Error(err)
	}
	err = json.Unmarshal(data, &app)
	if err != nil {
		tvision.Error(err)
	}
	return app
}

func (self *App) HandleEvent(event *tvision.Event, view *tvision.View) {
	tvision.Trace("Desktop.HandleEvent: ")
	switch event.Type {
	/*	CmdChoosePkg
		CmdAddToPkg
		CmdEditPkg
		CmdStartPkg*/
	case tvision.EvCommand:
		switch event.Cmd {
		case CmdChoosePkg:
			tvision.Trace("	CmdChoosePkg")
			//selectedPkg := event.Data.(Package)
		case CmdAddNewPackage:
			tvision.Trace("	CmdAddNewPackage")
		case CmdAddToPkg:
			showAddToPkgWindow()
		case tvision.CmdQuit:
			event.SetProcessed()
			self.Running = false
			desktopView.StopExecuting(tvision.ExecutingResult{tvision.CmdQuit, nil})
		}
	case tvision.EvKey:
		switch event.Key {
		case tvision.KeyF2:
			tvision.ToogleLogWindow()
		case tvision.KeyF3:
			tvision.LoggingTextArea.View.Rect = tvision.LoggingTextArea.View.Rect.Grow(1, 1)
		case tvision.KeyF4:
			tvision.LoggingTextArea.View.Rect = tvision.LoggingTextArea.View.Rect.Grow(-1, -1)
		}
	case tvision.EvResize:
		//app.ReDraw(view)
		event.SetProcessed()
	}
}

func showAddToPkgWindow() {
	winView, win := tvision.CreateWindow(50, 20, "Add new Card")
	winView.SetFontSize(12)

	comboView1, combo := tvision.CreateComboBox(20)
	combo.AddItem("Egy")
	combo.AddItem("Kettő")
	combo.AddItem("Három")
	combo.AddItem("Négy")
	combo.Editable = true
	comboView2, _ := tvision.CreateComboBox(20)
	winView.AddView(25, 10, comboView1)
	winView.AddView(25, 12, comboView2)

	editorView, _ := CreateAskEditor(20, 10)
	//win.Closeable = true
	win.Resizeable = true
	winView.AddView(1, 1, editorView)

	desktopView.AddView(28, 10, winView)
}

func CreateApp() (*App, *tvision.View) {
	group := new(tvision.View)
	app := loadDb()
	group.Widget = app
	app.Running = true
	return app, group
}

var desktopView *tvision.View

func main() {
	tvision.InitSdlEngine(140, 60, tvision.Pixel(21))

	app, desktopView = CreateApp()
	persist()

	//tvision.InitTermBoxEngine()
	//tvision.InitMultipleEngine()
	app.Init(desktopView)
	defer app.Close()

	view, box := tvision.CreateSelectBox("Csomagok")
	for _, pkg := range app.Packages {
		selectItem := box.AddItem(pkg.Name, CmdChoosePkg, pkg)
		selectItem.AddSubItem("Kezdés", CmdStartPkg, pkg)
		selectItem.AddSubItem("Hozzáadás", CmdAddToPkg, pkg)
		selectItem.AddSubItem("Szerkesztés", CmdEditPkg, pkg)
	}
	box.AddEventItem("Új hozzáadása", CmdAddNewPackage)
	box.AddEventItem("Kilépés", tvision.CmdQuit)
	desktopView.AddView(10, 10, view)

	comboView, combo := tvision.CreateComboBox(20)
	comboView.SetFontSize(18)
	combo.AddItem("Egy")
	combo.AddItem("Kettő")
	combo.AddItem("Három")
	combo.AddItem("Négy")
	combo.Editable = true
	desktopView.AddView(0, 0, comboView)

	desktopView.Execute()
}
