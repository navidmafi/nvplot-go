package main

import (
	"fmt"
	"time"

	"github.com/diamondburned/gotk4/pkg/cairo"
	"github.com/diamondburned/gotk4/pkg/gtk/v4"
	"github.com/diamondburned/gotk4/pkg/pango"
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
	ui := &UI{
		plotter:   plotter,
		collector: collector,
		storage:   storage,
	}

	ui.window = gtk.NewWindow()
	ui.window.SetTitle("GPU Usage")
	ui.window.SetDefaultSize(800, 600)

	vbox := gtk.NewBox(gtk.OrientationVertical, 0)

	titleLabel := gtk.NewLabel("GPU Usage")
	titleLabel.SetHAlign(gtk.AlignCenter)
	titleLabel.SetMarginTop(10)
	titleLabel.SetMarginBottom(10)
	boldAttrlist := pango.NewAttrList()
	// arrts := []*pango.Attribute{pango.NewAttrWeight(pango.WeightBold), pango.NewAttrScale(1.5)}
	boldAttrlist.Insert(pango.NewAttrWeight(pango.WeightBold))
	boldAttrlist.Insert(pango.NewAttrScale(1.5))
	titleLabel.SetAttributes(boldAttrlist)

	ui.infoLabel = gtk.NewLabel("")
	ui.infoLabel.SetHAlign(gtk.AlignCenter)
	ui.infoLabel.SetMarginBottom(10)

	ui.drawing = gtk.NewDrawingArea()
	ui.drawing.SetVExpand(true)
	ui.drawing.SetDrawFunc(ui.draw)

	vbox.Append(titleLabel)
	vbox.Append(ui.infoLabel)
	vbox.Append(ui.drawing)

	ui.window.SetChild(vbox)

	return ui
}

func (ui *UI) draw(area *gtk.DrawingArea, cr *cairo.Context, width, height int) {
	ui.plotter.Plot(area, cr, width, height)
}

func (ui *UI) Run() {
	ui.window.SetVisible(true)
	go ui.updateData()
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
		ui.updateInfoLabel(memInfo)
		ui.drawing.QueueDraw()
	}
}

func (ui *UI) updateInfoLabel(memInfo MemoryInfo) {
	ui.infoLabel.SetMarkup(fmt.Sprintf(
		"<b>Total:</b> %d MB | <b>Used:</b> %d MB | <b>Free:</b> %d MB",
		memInfo.TotalMB, memInfo.UsedMB, memInfo.FreeMB,
	))
}
