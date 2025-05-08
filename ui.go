package main

import (
	"fmt"
	"time"

	"github.com/gotk3/gotk3/cairo"
	"github.com/gotk3/gotk3/glib"
	"github.com/gotk3/gotk3/gtk"
)

type UI struct {
	window    *gtk.Window
	drawing   *gtk.DrawingArea
	infoLabel *gtk.Label
	plotter   Plotter
	collector *NVMLCollector
	storage   *DataStorage
}

func NewUI(plotter Plotter, collector *NVMLCollector, storage *DataStorage) *UI {
	// 1) Initialize GTK (must be first)
	gtk.Init(nil)

	ui := &UI{
		plotter:   plotter,
		collector: collector,
		storage:   storage,
	}

	// 2) Create top‐level window
	w, err := gtk.WindowNew(gtk.WINDOW_TOPLEVEL)
	if err != nil {
		panic("Failed to create window: " + err.Error())
	}
	ui.window = w
	ui.window.SetTitle("GPU Usage")
	ui.window.SetDefaultSize(800, 600)

	// 3) Vertical box container
	vbox, err := gtk.BoxNew(gtk.ORIENTATION_VERTICAL, 0)
	if err != nil {
		panic(err)
	}

	// 4) Title label with bold + scale attributes
	titleLabel, _ := gtk.LabelNew("GPU Usage")
	titleLabel.SetHAlign(gtk.ALIGN_CENTER)
	titleLabel.SetMarginTop(10)
	titleLabel.SetMarginBottom(10)

	titleLabel.SetMarkup(`<span weight="bold" size="15000">GPU Usage</span>`)

	// 5) Info label (initially empty)
	infoLabel, _ := gtk.LabelNew("")
	infoLabel.SetHAlign(gtk.ALIGN_CENTER)
	infoLabel.SetMarginBottom(10)
	ui.infoLabel = infoLabel

	// 6) Drawing area and connect draw signal
	drawing, _ := gtk.DrawingAreaNew()
	drawing.SetVExpand(true)
	drawing.Connect("draw", func(da *gtk.DrawingArea, cr *cairo.Context) {
		// width/height can be queried or passed if needed
		w, h := da.GetAllocatedWidth(), da.GetAllocatedHeight()
		ui.plotter.Plot(da, cr, w, h)
	})
	ui.drawing = drawing

	// 7) Pack widgets into vbox
	vbox.PackStart(titleLabel, false, false, 0)
	vbox.PackStart(ui.infoLabel, false, false, 0)
	vbox.PackStart(ui.drawing, true, true, 0)

	// 8) Add vbox to window
	ui.window.Add(vbox)

	// 9) Setup destroy handler
	ui.window.Connect("destroy", func() {
		gtk.MainQuit()
	})

	return ui
}

func (ui *UI) Run() {
	// Show all widgets
	ui.window.ShowAll()

	// Start background updates
	go ui.updateData()

	// Enter GTK main loop
	gtk.Main()
}

func (ui *UI) updateData() {
	ticker := time.NewTicker(time.Second)
	defer ticker.Stop()

	for range ticker.C {
		memInfo, err := ui.collector.GetVRAMUsage()
		if err != nil {
			fmt.Println("Error getting VRAM usage:", err)
			continue
		}

		ui.storage.AddDataPoint(memInfo)

		// UI updates must run in GTK main thread
		glib.IdleAdd(func() {
			ui.updateInfoLabel(memInfo)
			ui.drawing.QueueDraw()
		})
	}
}

func (ui *UI) updateInfoLabel(memInfo MemoryInfo) {
	// Use Pango markup
	ui.infoLabel.SetMarkup(fmt.Sprintf(
		"<b>Total:</b> %d MB | <b>Used:</b> %d MB | <b>Free:</b> %d MB",
		memInfo.TotalMB, memInfo.UsedMB, memInfo.FreeMB,
	))
}
