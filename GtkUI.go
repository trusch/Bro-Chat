package main

import(
	"github.com/mattn/go-gtk/gtk"
	"github.com/mattn/go-gtk/glib" 
	"github.com/mattn/go-gtk/gdk"
	"unsafe"
	"os"
	"log"
	"time"
)

type GtkUI struct {

	BaseUI

	window		*gtk.Window
	outputElem	*gtk.TextView
}

func (ui *GtkUI) Init(sender *BCastSender,receiver *BCastReceiver) error {
	ui.BaseUI.Init(sender,receiver)
	ui.initGtk()
	return nil
}

func (ui *GtkUI) Run() error{
	go gtk.Main()
	go func(){
		for {
			msg := ui.ReadNextMessage()
			if msg!=nil {
				ui.UpdateWhoMap(msg)
				ui.ProcessOutput(msg.String()+"\n")
			}else{
				log.Print("nil msg (timeout?)")
			}
		}
	}()
	go func(){
		for msg := range ui.OutChan {
			log.Print("output to gtk window: ",msg)
			gdk.ThreadsEnter()
			buff := ui.outputElem.GetBuffer()
			buff.InsertAtCursor(msg)
			gdk.ThreadsLeave()			
		}
	}()
	go func(){
		ticker := time.Tick(100*time.Millisecond)
		for{
			<-ticker
			gdk.ThreadsEnter()
			buff := ui.outputElem.GetBuffer()
			var end gtk.TextIter
			buff.GetEndIter(&end)
			ui.outputElem.ScrollToIter(&end,0,false,0,0)
			gdk.ThreadsLeave()		
		}
	}()
	return nil
}

func (ui *GtkUI) initGtk(){
	gtk.Init(&os.Args)
	window := gtk.NewWindow(gtk.WINDOW_TOPLEVEL)
	window.SetTitle("Bro-Chat")
	window.Connect("destroy", func(){
		gtk.MainQuit()
		os.Exit(0)
	})

	vbox := gtk.NewVBox(false, 0)

	scrolledwin := gtk.NewScrolledWindow(nil, nil)
	scrolledwin.SetPolicy(1,1)
	textview := gtk.NewTextView()
	textview.SetEditable(false)
	textview.SetCursorVisible(false)
	scrolledwin.Add(textview)
	vbox.PackStart(scrolledwin,true,true,0)

	ui.outputElem = textview
	
	hbox := gtk.NewHBox(false,0)
	entry := gtk.NewEntry()
	entry.SetText("Hello bro's!")
	hbox.Add(entry)
	button := gtk.NewButtonWithLabel("Send it!")
	button.Clicked(func() {
			userInput := entry.GetText()
			ui.ProcessUserInput(userInput)
			entry.SetText("")
		})
	hbox.Add(button)
	
	align := gtk.NewAlignment(0,1,1,0)
	align.Add(hbox)
	vbox.PackEnd(align,false,true,0)
	
	window.Add(vbox)
	
	window.Connect("key-press-event", func(ctx *glib.CallbackContext) {
		arg := ctx.Args(0)
		kev := *(**gdk.EventKey)(unsafe.Pointer(&arg))
		if !entry.HasFocus(){
			entry.GrabFocus()
			return
		}
		if kev.Keyval==65293 { //Return pressed
			userInput := entry.GetText()
			ui.ProcessUserInput(userInput)
			entry.SetText("")
		}
	})

	window.SetSizeRequest(380, 600)
	window.ShowAll()
	
	glib.ThreadInit(nil)
	gdk.ThreadsInit()
	
	entry.GrabFocus()
	
	ui.window = window
}


func NewGtkUI(nick string) *GtkUI {
	ui :=  new(GtkUI)
	ui.Nickname = nick
	return ui
}
