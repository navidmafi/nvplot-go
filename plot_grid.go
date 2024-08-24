package main

import (
	"fmt"

	"github.com/diamondburned/gotk4/pkg/cairo"
)

func (vp *VRAMPlotter) Grid(cr *cairo.Context, width, height int) {
	fWidth := float64(width)
	fHeight := float64(height)
	fPadding := float64(vp.padding)
	fVertLines := float64(vp.vertLines)
	fHorizLines := float64(vp.horizLines)

	ContainerRectMin := (fPadding) / 2
	ContainerXMax := fWidth - (fPadding)/2
	ContainerYMax := fHeight - (fPadding)/2
	usableContainerHeight := ContainerYMax - ContainerRectMin
	usableContainerWidth := ContainerXMax - ContainerRectMin

	cr.SetSourceRGB(0.3, 0.3, 0.3)
	cr.SetLineWidth(1)

	// vertical grid lines
	for i := 0; i <= vp.vertLines; i++ {
		x := (float64(i)/fVertLines)*usableContainerWidth + ContainerRectMin
		cr.MoveTo(x, ContainerRectMin)
		cr.LineTo(x, ContainerYMax)
	}
	cr.Stroke()

	// horizontal grid lines
	for i := 0; i <= vp.horizLines; i++ {
		y := (float64(i)/fHorizLines)*usableContainerHeight + ContainerRectMin
		cr.MoveTo(ContainerRectMin, y)
		cr.LineTo(ContainerXMax, y)
	}
	cr.Stroke()

	// Label horizontal lines
	for i := 0; i <= vp.horizLines; i++ {
		y := (float64(i)/fHorizLines)*usableContainerHeight + ContainerRectMin + 5
		cr.MoveTo(ContainerRectMin-30, y)
		cr.ShowText(fmt.Sprintf("%d", int((1-float64(i)/fHorizLines)*float64(vp.maxValue))))
		cr.Stroke()
	}

	// // Label vertical lines
	// for i := 0; i <= vp.vertLines; i++ {
	// 	x := float64(i)/fVertLines*(fWidth-2*fPadding) + fPadding
	// 	cr.MoveTo(x, fHeight-fPadding+10)
	// 	cr.ShowText(fmt.Sprintf("%d", i*10))
	// 	cr.Stroke()
	// }
}
