package gui

import (
	"github.com/gotk3/gotk3/cairo"
	"github.com/gotk3/gotk3/pango"
	"github.com/tbuen/gocmd/internal/backend/panel"
)

func drawTabs(context *cairo.Context, layout *pango.Layout, width, height float64, header *panel.Header) (sx1, sx2 float64) {
	cw := 6.0
	space := 30.0
	offset := header.Offset

	var text string
	var widths []float64
	var activeX1 float64
	var activeX2 float64
	var sumWidth float64
	for i, title := range header.Titles {
		w := (float64(len(title)) + 2) * cw
		text += " " + title + " "
		widths = append(widths, w)
		sumWidth += w
		if i < header.Active {
			activeX1 += w
		}
		if i <= header.Active {
			activeX2 += w
		}
	}

	if sumWidth < (width-2*space)+1 {
		offset = 0.0
	} else if sumWidth-offset < (width-2*space)+1 {
		offset = sumWidth - (width - 2*space) + 1
	}
	if offset > activeX1 {
		offset = activeX1
	}
	if offset < activeX2-(width-2*space)+1 {
		offset = activeX2 - (width - 2*space) + 1
	}
	header.Offset = offset

	if offset > 0 {
		setSourceColor(context, "000000")
	} else {
		setSourceColor(context, "B0B0B0")
	}
	context.MoveTo(space/3, height/2)
	context.LineTo(space/3+10, height/2-5)
	context.LineTo(space/3+10, height/2+5)
	context.Fill()

	if offset < sumWidth-(width-2*space)+1 {
		setSourceColor(context, "000000")
	} else {
		setSourceColor(context, "B0B0B0")
	}
	context.MoveTo(width-space/3, height/2)
	context.LineTo(width-space/3-10, height/2-5)
	context.LineTo(width-space/3-10, height/2+5)
	context.Fill()

	context.Rectangle(space-1, -1, width-(2*space)+1, height+1)
	context.Clip()

	x := space - offset
	for i, w := range widths {
		if i == header.Active {
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
	context.MoveTo(space-offset, 5)
	layout.SetMarkup(text, -1)
	pango.CairoShowLayout(context, layout)
	layout.SetMarkup("", -1)

	return
}
