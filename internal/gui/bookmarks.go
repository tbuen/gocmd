package gui

import (
	"github.com/gotk3/gotk3/cairo"
	"github.com/gotk3/gotk3/pango"
	"github.com/tbuen/gocmd/internal/backend"
)

func drawBookmarks(context *cairo.Context, layout *pango.Layout, width, height, sx1, sx2 float64, active bool, bm *backend.Bookmarks) {
	const scrollbarWidth = 8.0

	ch := 15.0
	cw := 6.0

	setSourceColor(context, "F6F5F4")
	context.Rectangle(0, 0, width, height)
	context.Fill()

	setSourceColor(context, "000000")
	context.Rectangle(0, 0, width, height)
	context.Stroke()
	setSourceColor(context, "F6F5F4")
	context.MoveTo(sx1, 0)
	context.LineTo(sx2, 0)
	context.Stroke()

	if bm == nil {
		return
	}

	//layout.SetText(".");
	//layout.GetPixelSize(cw, ch);
	lines := int((height - ch - 27) / ch)
	columns := int((width - 10) / cw)

	if active {
		//setSourceColor(context, "3584E4")
		setSourceColor(context, "E48435")
	} else {
		setSourceColor(context, "707070")
	}
	context.Rectangle(2, 2, width-5, ch+2)
	context.Fill()

	setSourceColor(context, "FFFFFF")
	context.MoveTo(5, 3)
	layout.SetText("Bookmarks", -1)
	pango.CairoShowLayout(context, layout)

	width -= scrollbarWidth
	columns = int((width - 10) / cw)

	bookmarks := bm.Bookmarks()

	selection := bm.Selection()
	offset := bm.DispOffset()

	if len(bookmarks) <= lines {
		offset = 0
	} else if offset > len(bookmarks)-lines {
		offset = len(bookmarks) - lines
	}
	if selection >= offset+lines {
		offset = selection - lines + 1
	}
	if selection < offset {
		offset = selection
	}
	bm.SetDispOffset(offset)

	lenPath := columns

	for i := 0; i <= lines && offset+i < len(bookmarks); i++ {
		bookmark := bookmarks[offset+i]
		if i%2 == 0 {
			setSourceColor(context, "FFFFFF")
			context.Rectangle(2, 6+(float64(i)+1)*ch, width-5, ch)
			context.Fill()
		} else {
			setSourceColor(context, "F6F5F4")
			context.Rectangle(2, 6+(float64(i)+1)*ch, width-5, ch)
			context.Fill()
		}
		if active && offset+i == selection {
			setSourceColor(context, "000000")
			context.Rectangle(3, 6+(float64(i)+1)*ch, width-6, ch)
			context.Stroke()
		}
		setSourceColor(context, "000000")
		context.MoveTo(5, 5+(float64(i)+1)*ch)

		var line string
		line = bookmark.Path()
		line = restrictBack(line, lenPath)

		layout.SetText(line, -1)
		pango.CairoShowLayout(context, layout)
	}
	context.Save()
	context.Translate(width-3, ch+5)
	drawScrollbar(context, scrollbarWidth, height-ch-ch-11, len(bookmarks), lines, offset)
	context.Restore()

	width += scrollbarWidth
	columns = int((width - 10) / cw)
	setSourceColor(context, "B0B0B0")
	context.Rectangle(2, height-20, width-5, ch+2)
	context.Fill()
}
