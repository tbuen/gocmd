package gui

import (
	"github.com/gotk3/gotk3/cairo"
	"github.com/gotk3/gotk3/pango"
)

func drawTabs(context *cairo.Context, layout *pango.Layout, width, height float64, tabs []string, active int) {
	//setSourceColor(context, "F26B3A")
	//context.Rectangle(5, 5, width-8, height-9)
	//context.Fill()
	context.Rectangle(6, 6, width-12, height-11)
	context.Clip()

	var text string
	for i, title := range tabs {
		if i == 0 {
			text = title
		} else {
			text += "  " + title
		}
	}

	setSourceColor(context, "000000")
	context.MoveTo(10, 8)
	layout.SetMarkup(text, -1)
	pango.CairoShowLayout(context, layout)
	layout.SetMarkup("", -1)
}
