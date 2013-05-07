package gui

import (
	"fmt"
	"github.com/mattn/go-gtk/gdk"
	"github.com/mattn/go-gtk/glib"
	"github.com/mattn/go-gtk/gtk"
)

type GUI struct {
	width    int
	height   int
	window   *gtk.Window
	notebook *gtk.Notebook
	pages    map[string]*Page
	menuItem *gtk.MenuItem
}

type Page struct {
	textView *gtk.TextView
	nickTV   *gtk.TextView
	entry    *gtk.Entry
}

func NewGUI(title string, width, height int) *GUI {
	glib.ThreadInit(nil)
	gdk.ThreadsInit()
	gdk.ThreadsEnter()
	gtk.Init(nil)

	window := gtk.NewWindow(gtk.WINDOW_TOPLEVEL)
	window.SetPosition(gtk.WIN_POS_CENTER)
	window.SetTitle(title)
	window.SetIconName("gtk-dialog-info")
	window.Connect("destroy", func(ctx *glib.CallbackContext) {
		println("got destroy!", ctx.Data().(string))
		gtk.MainQuit()
	}, "foo")

	vbox := gtk.NewVBox(false, 0)

	notebook := gtk.NewNotebook()

	vbox.Add(notebook)
	window.Add(vbox)
	window.SetSizeRequest(width, height)

	return &GUI{window: window, notebook: notebook, pages: make(map[string]*Page),
		width: width, height: height}
}

func (gui *GUI) StartMain() {
	gui.window.ShowAll()
	gtk.Main()
}

func (gui *GUI) CreateChannelWindow(context string, buttonFunc func()) {
	var page *gtk.Frame

	if context == "" {
		page = gtk.NewFrame("Server")
		gui.notebook.AppendPage(page, gtk.NewLabel("Server"))
	} else {
		page = gtk.NewFrame(context)
		gui.notebook.AppendPage(page, gtk.NewLabel(context))
	}

	vbox := gtk.NewVBox(false, 1)
	hbox1 := gtk.NewHBox(false, 1)

	var nickTV *gtk.TextView
	var textView *gtk.TextView

	if context != "" {
		swin := gtk.NewScrolledWindow(nil, nil)
		swin.SetPolicy(gtk.POLICY_AUTOMATIC, gtk.POLICY_AUTOMATIC)
		swin.SetShadowType(gtk.SHADOW_IN)
		textView = gtk.NewTextView()
		textView.SetEditable(false)
		textView.SetCursorVisible(false)
		textView.SetWrapMode(gtk.WRAP_WORD)
		textView.SetSizeRequest(600, 500)
		swin.Add(textView)
		hbox1.Add(swin)

		swin2 := gtk.NewScrolledWindow(nil, nil)
		swin2.SetPolicy(gtk.POLICY_AUTOMATIC, gtk.POLICY_AUTOMATIC)
		swin2.SetShadowType(gtk.SHADOW_IN)
		nickTV := gtk.NewTextView()
		nickTV.SetEditable(false)
		nickTV.SetCursorVisible(false)
		nickTV.SetWrapMode(gtk.WRAP_WORD)
		nickTV.SetSizeRequest(200, 500)
		swin2.Add(nickTV)
		hbox1.Add(swin2)
	} else {
		swin := gtk.NewScrolledWindow(nil, nil)
		swin.SetPolicy(gtk.POLICY_AUTOMATIC, gtk.POLICY_AUTOMATIC)
		swin.SetShadowType(gtk.SHADOW_IN)
		textView = gtk.NewTextView()
		textView.SetEditable(false)
		textView.SetCursorVisible(false)
		textView.SetWrapMode(gtk.WRAP_WORD)
		textView.SetSizeRequest(800, 500)
		swin.Add(textView)
		hbox1.Add(swin)
	}

	vbox.Add(hbox1)
	hbox2 := gtk.NewHBox(false, 1)

	// entry
	entry := gtk.NewEntry()
	entry.SetSizeRequest(700, 40)
	hbox2.Add(entry)

	button := gtk.NewButtonWithLabel("Send")
	button.Clicked(buttonFunc)
	hbox2.Add(button)

	vbox.Add(hbox2)

	page.Add(vbox)

	newPage := &Page{textView: textView, nickTV: nickTV, entry: entry}
	gui.pages[context] = newPage
}

func (gui *GUI) DeleteCurrentWindow() {
	gui.notebook.RemovePage(nil, gui.notebook.GetCurrentPage())
}

func (gui *GUI) WriteToChannel(s, context string) (err error) {
	var endIter gtk.TextIter
	tmp, ok := gui.pages[context]
	if !ok {
		return fmt.Errorf("No such window %s", context)
	}
	textBuffer := tmp.textView.GetBuffer()
	textBuffer.GetEndIter(&endIter)
	textBuffer.Insert(&endIter, s+"\n")

	gui.AutoScroll(textBuffer, &endIter)
	return nil
}

func (gui *GUI) WriteToNicks(s, context string) {
	var endIter gtk.TextIter
	textBuffer := gui.pages[context].nickTV.GetBuffer()
	textBuffer.GetEndIter(&endIter)
	textBuffer.Insert(&endIter, s+"\n")
}

func (gui *GUI) EmptyNicks(s, context string) {
	textBuffer := gui.pages[context].nickTV.GetBuffer()
	textBuffer.SetText("")
}

func (gui *GUI) AutoScroll(textbuffer *gtk.TextBuffer, endIter *gtk.TextIter) {
	// TODO
}
