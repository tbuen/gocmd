package gui

import (
	"github.com/gotk3/gotk3/cairo"
)

func drawScrollbar(context *cairo.Context, width, height float64, total, visible, offset int) {
	context.SetSourceRGB(0x70/255.0, 0x70/255.0, 0x70/255.0)
	//context.Rectangle(0, 0, width, height)
	//context.Fill()
	ratio := 1.0
	if total > visible {
		ratio = float64(visible) / float64(total)
	}
	start := float64(offset) / float64(total) * height
	context.Rectangle(2, start, width-2, height*ratio)
	context.Fill()
}
