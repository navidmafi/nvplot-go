package main

import (
	"log"

	"github.com/gotk3/gotk3/gtk"
)

func main() {
	collector, err := NewNVMLCollector()
	gtk.Init(nil)
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

	ui := NewUI(plotter, collector, storage)

	ui.Run()

}
