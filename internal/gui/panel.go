package gui

import (
	"github.com/gotk3/gotk3/cairo"
	"github.com/gotk3/gotk3/pango"
	"github.com/tbuen/gocmd/internal/backend"
	"github.com/tbuen/gocmd/internal/config"
	"unicode/utf8"
)

func setSourceColor(context *cairo.Context, color config.Color) {
	context.SetSourceRGB(color[0], color[1], color[2])
}

func drawPanel(context *cairo.Context, layout *pango.Layout, width, height float64, active bool, dir backend.Directory) {
	const scrollbarWidth = 8.0

	ch := 15.0
	cw := 6.0

	context.SetSourceRGB(0, 0, 0)
	context.Rectangle(5, 5, width-8, height-9)
	context.Stroke()
	context.Rectangle(6, 6, width-12, height-11)
	context.Clip()

	if dir == nil {
		return
	}

	//layout.SetText(".");
	//layout.GetPixelSize(cw, ch);
	lines := int((height - 19 - ch) / ch)
	columns := int((width - 19) / cw)

	if active {
		context.SetSourceRGB(0x35/255.0, 0x84/255.0, 0xe4/255.0)
	} else {
		context.SetSourceRGB(0x70/255.0, 0x70/255.0, 0x70/255.0)
	}
	context.Rectangle(7, 7, width-13, ch+2)
	context.Fill()

	state := dir.State()
	path := restrictFront(dir.Path(), columns)

	switch state {
	case backend.STATE_IDLE:
		context.SetSourceRGB(1, 1, 1)
	case backend.STATE_ERROR:
		context.SetSourceRGB(1, 0, 0)
	case backend.STATE_RELOAD:
		context.SetSourceRGB(1, 1, 0)
	}
	context.MoveTo(10, 8)
	layout.SetText(path, -1)
	pango.CairoShowLayout(context, layout)

	if state == backend.STATE_IDLE {
		width -= scrollbarWidth
		columns = int((width - 19) / cw)

		selection := dir.Selection()
		offset := dir.DispOffset()

		if selection >= offset+lines {
			offset = selection - lines + 1
		}
		if selection < offset {
			offset = selection
		}
		dir.SetDispOffset(offset)

		files := dir.Files()

		minLenName := 15
		extraLen := 0
		var lenName, lenSize, lenTime, lenOwner, lenPerm int
		showSize, showTime, showOwner, showPerm := view.size, view.time, view.owner, view.perm
		if showSize {
			var maxSize int64
			for _, f := range files {
				if s := f.Size(); s > maxSize {
					maxSize = s
				}
			}
			_, l := formatSize(maxSize)
			if l < 5 {
				l = 5
			}
			lenSize = l + 2
			if minLenName+extraLen+lenSize <= columns {
				extraLen += lenSize
			} else {
				showSize, showTime, showOwner, showPerm = false, false, false, false
			}
		}
		if showTime {
			lenTime = 18
			if minLenName+extraLen+lenTime <= columns {
				extraLen += lenTime
			} else {
				showTime, showOwner, showPerm = false, false, false
			}
		}
		var maxlenUser, maxlenGroup int
		if showOwner {
			for _, f := range files {
				user, group := f.UserGroup()
				lu, lg := utf8.RuneCountInString(user), utf8.RuneCountInString(group)
				if lu > maxlenUser {
					maxlenUser = lu
				}
				if lg > maxlenGroup {
					maxlenGroup = lg
				}
			}
			lenOwner = maxlenUser + maxlenGroup + 3
			if minLenName+extraLen+lenOwner <= columns {
				extraLen += lenOwner
			} else {
				showOwner, showPerm = false, false
			}
		}
		if showPerm {
			lenPerm = 11
			if minLenName+extraLen+lenPerm <= columns {
				extraLen += lenPerm
			} else {
				showPerm = false
			}
		}
		lenName = columns - extraLen

		for i := 0; i <= lines && offset+i < len(files); i++ {
			file := files[offset+i]
			if file.Marked() {
				context.SetSourceRGB(0xF2/255.0, 0x6B/255.0, 0x3A/255.0)
				//context.SetSourceRGB(0xFF/255.0, 0xA0/255.0, 0x90/255.0)
				context.Rectangle(7, 11+(float64(i)+1)*ch, width-13, ch)
				context.Fill()
			} else if i%2 == 0 {
				context.SetSourceRGB(1.0, 1.0, 1.0)
				context.Rectangle(7, 11+(float64(i)+1)*ch, width-13, ch)
				context.Fill()
			} else {
				context.SetSourceRGB(0xF6/255.0, 0xF5/255.0, 0xF4/255.0)
				context.Rectangle(7, 11+(float64(i)+1)*ch, width-13, ch)
				context.Fill()
			}
			if active && offset+i == selection {
				context.SetSourceRGB(0, 0, 0)
				context.Rectangle(8, 11+(float64(i)+1)*ch, width-14, ch)
				context.Stroke()
			}
			setSourceColor(context, file.Color())
			context.MoveTo(10, 10+(float64(i)+1)*ch)

			var line string
			link, linkOk, linkTarget := file.Link()
			if link {
				if linkOk {
					//if file.Dir() {
					//	line = "[" + file.Name() + "] -> " + linkTarget
					//} else {
					line = file.Name() + " -> " + linkTarget
					//}
				} else {
					line = file.Name() + " -| " + linkTarget
				}
			} else {
				//if file.Dir() {
				//	line = "[" + file.Name() + "]"
				//} else {
				line = file.Name()
				//}
			}
			if file.Dir() {
				line = "/" + line
			} else if file.Pipe() {
				line = "|" + line
			} else if file.Socket() {
				line = "=" + line
			} else {
				line = " " + line
			}

			line = restrictBack(line, lenName)
			line = appendSpaces(line, lenName)

			if showSize {
				var size string
				if file.Dir() {
					size = "<DIR>"
				} else {
					size, _ = formatSize(file.Size())
				}
				line += prependSpaces(size, lenSize)
			}
			if showTime {
				time := file.Time().Format("02.01.2006 15:04")
				line += "  " + time
			}
			if showOwner {
				user, group := file.UserGroup()
				line += "  " + appendSpaces(user, maxlenUser) + " " + appendSpaces(group, maxlenGroup)
			}
			if showPerm {
				line += "  " + file.Perm()
			}

			layout.SetText(line, -1)
			pango.CairoShowLayout(context, layout)
		}
		context.Save()
		context.Translate(width+7-13, 7+ch+2+1)
		drawScrollbar(context, scrollbarWidth, height-11-ch-2-3, len(files), lines, offset)
		context.Restore()
	} else {
		width -= scrollbarWidth
		context.Save()
		context.Translate(width+7-13, 7+ch+2+1)
		drawScrollbar(context, scrollbarWidth, height-11-ch-2-3, 0, 0, 0)
		context.Restore()
	}
}
