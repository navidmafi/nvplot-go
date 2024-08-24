package main

import (
	"log"

	"github.com/diamondburned/gotk4/pkg/gtk/v4"
)

func main() {
	collector, err := NewNVMLCollector()

	if err != nil {
		log.Fatal("Error initializing NVML collector:", err)
	}
	defer collector.Close()

	storage := NewDataStorage(100)
	info, err := collector.GetVRAMUsage()
	if err != nil {
		panic(err)
	}

	plotter := NewVRAMPlotter(storage, 0, info.TotalMB)

	app := gtk.NewApplication("com.navidmafi.vram-usage", 0)
	app.ConnectActivate(func() {
		ui := NewUI(plotter, collector, storage)
		ui.Run()
		app.AddWindow(ui.window)
	})

	if code := app.Run(nil); code > 0 {
		log.Fatal("Error running application:", code)
	}
}
