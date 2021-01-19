package gui

import (
	"github.com/gotk3/gotk3/cairo"
	"github.com/gotk3/gotk3/pango"
)

func drawTabs(context *cairo.Context, layout *pango.Layout, width, height float64, tabs []string, active int) (sx1, sx2 float64) {
	cw := 6.0

	var text string
	var widths []float64
	for _, title := range tabs {
		text += " " + title + " "
		widths = append(widths, (float64(len(title))+2)*cw)
	}

	x := 5.0
	for i, w := range widths {
		if i == active {
			setSourceColor(context, "F6F5F4")
			context.Rectangle(x, 0, w, height)
			context.Fill()
			context.SetLineWidth(2)
			setSourceColor(context, "000000")
			context.MoveTo(x, height)
			context.LineTo(x, 0)
			context.LineTo(x+w, 0)
			context.LineTo(x+w, height)
			context.Stroke()
			sx1 = x + 1
			sx2 = x + w - 1
		} else {
			context.SetLineWidth(1)
			setSourceColor(context, "000000")
			context.MoveTo(x, height)
			context.LineTo(x, 2)
			context.LineTo(x+w, 2)
			context.LineTo(x+w, height)
			context.Stroke()
		}
		x += w
	}
	context.MoveTo(5, 5)
	layout.SetMarkup(text, -1)
	pango.CairoShowLayout(context, layout)
	layout.SetMarkup("", -1)
	return
}
