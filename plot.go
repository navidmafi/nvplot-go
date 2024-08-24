package main

import (
	"github.com/diamondburned/gotk4/pkg/cairo"
	"github.com/diamondburned/gotk4/pkg/gtk/v4"
)

type Plotter interface {
	Plot(area *gtk.DrawingArea, cr *cairo.Context, width, height int)
}

type VRAMPlotter struct {
	dataStorage *DataStorage
	padding     int
	vertLines   int
	horizLines  int
	maxValue    int
	minValue    int
}

func NewVRAMPlotter(dataStorage *DataStorage, minValue, maxValue int) *VRAMPlotter {
	return &VRAMPlotter{
		dataStorage: dataStorage,
		padding:     100,
		vertLines:   10,
		horizLines:  20,
		maxValue:    maxValue,
		minValue:    minValue,
	}
}

func (vp *VRAMPlotter) Plot(area *gtk.DrawingArea, cr *cairo.Context, width, height int) {
	data := vp.dataStorage.GetData()

	// TODO : refactor

	fWidth := float64(width)
	fHeight := float64(height)
	fPadding := float64(vp.padding)

	ContainerRectMin := (fPadding) / 2
	ContainerXMax := fWidth - (fPadding)/2
	ContainerYMax := fHeight - (fPadding)/2
	usableContainerHeight := ContainerYMax - ContainerRectMin
	usableContainerWidth := ContainerXMax - ContainerRectMin

	dsSize := vp.dataStorage.maxPoints
	if len(data) == 0 {
		return
	}
	cr.SetSourceRGB(0, 0, 0)
	cr.Paint()

	vp.Grid(cr, width, height)

	cr.SetSourceRGB(0, 1, 0)
	cr.SetLineWidth(2)

	cr.MoveTo(ContainerRectMin, ContainerYMax)

	normalizedY := func(Y float64) float64 {
		return (ContainerYMax - Y)
	}

	for i, dp := range data {

		x := ContainerRectMin + ((float64(i) / float64(dsSize-1)) * usableContainerWidth)
		y := normalizedY((float64(dp.UsedMB) / float64(dp.TotalMB)) * usableContainerHeight)
		if i == 0 {
			cr.MoveTo(x, y)
			continue
		}
		cr.LineTo(x, y)
	}

	cr.Stroke()
}
