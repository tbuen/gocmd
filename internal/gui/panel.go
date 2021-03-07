package gui

import (
	"github.com/gotk3/gotk3/cairo"
	"github.com/gotk3/gotk3/pango"
	"github.com/tbuen/gocmd/internal/backend/dir"
	"strconv"
	"unicode/utf8"
)

func drawPanel(context *cairo.Context, layout *pango.Layout, width, height, sx1, sx2 float64, active bool, d *dir.Directory) {
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

	if d == nil {
		return
	}

	//layout.SetText(".");
	//layout.GetPixelSize(cw, ch);
	lines := int((height - ch - 27) / ch)
	columns := int((width - 10) / cw)

	state := d.State()
	path := restrictFront(d.Path(), columns)

	if active {
		setSourceColor(context, "3584E4")
	} else {
		setSourceColor(context, "707070")
	}
	context.Rectangle(2, 2, width-5, ch+2)
	context.Fill()

	switch state {
	case dir.STATE_IDLE:
		setSourceColor(context, "FFFFFF")
	case dir.STATE_ERROR:
		setSourceColor(context, "FF0000")
	case dir.STATE_RELOAD:
		setSourceColor(context, "FFFF00")
	}
	context.MoveTo(5, 3)
	layout.SetText(path, -1)
	pango.CairoShowLayout(context, layout)

	if state == dir.STATE_IDLE {
		width -= scrollbarWidth
		columns = int((width - 10) / cw)

		files := d.Files()

		selection := d.Selection()
		offset := d.DispOffset()

		if len(files) <= lines {
			offset = 0
		} else if offset > len(files)-lines {
			offset = len(files) - lines
		}
		if selection >= offset+lines {
			offset = selection - lines + 1
		}
		if selection < offset {
			offset = selection
		}
		d.SetDispOffset(offset)

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
				//setSourceColor(context, "F26B3A")
				setSourceColor(context, "B0B0B0")
				context.Rectangle(2, 6+(float64(i)+1)*ch, width-5, ch)
				context.Fill()
			} else if i%2 == 0 {
				if active {
					setSourceColor(context, "FFFFFF")
				} else {
					setSourceColor(context, "EEEEEE")
				}
				context.Rectangle(2, 6+(float64(i)+1)*ch, width-5, ch)
				context.Fill()
			} else {
				if active {
					setSourceColor(context, "F6F5F4")
				} else {
					setSourceColor(context, "E6E5E4")
				}
				context.Rectangle(2, 6+(float64(i)+1)*ch, width-5, ch)
				context.Fill()
			}
			if offset+i == selection {
				setSourceColor(context, "000000")
				context.Rectangle(3, 6+(float64(i)+1)*ch, width-6, ch)
				context.Stroke()
			}
			setSourceColor(context, file.Color())
			context.MoveTo(5, 5+(float64(i)+1)*ch)

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
		context.Translate(width-3, ch+5)
		drawScrollbar(context, scrollbarWidth, height-ch-ch-11, len(files), lines, offset)
		context.Restore()
	} else {
		width -= scrollbarWidth
		context.Save()
		context.Translate(width-3, ch+5)
		drawScrollbar(context, scrollbarWidth, height-ch-ch-11, 0, 0, 0)
		context.Restore()
	}

	width += scrollbarWidth
	columns = int((width - 10) / cw)
	setSourceColor(context, "B0B0B0")
	context.Rectangle(2, height-20, width-5, ch+2)
	context.Fill()
	if state == dir.STATE_IDLE {
		info := d.Info()
		var text string
		if d.Hidden() {
			text += ". "
		}
		text += strconv.Itoa(info.NumDirs)
		if info.NumDirs == 1 {
			text += " directory - "
		} else {
			text += " directories - "
		}
		text += strconv.Itoa(info.NumFiles)
		if info.NumFiles == 1 {
			text += " file - "
		} else {
			text += " files - "
		}
		sizeStr, _ := formatSize(info.SizeFiles)
		text += sizeStr
		if info.SizeFiles == 1 {
			text += " byte"
		} else {
			text += " bytes"
		}
		if info.NumSelectedDirs+info.NumSelectedFiles > 0 {
			text += ", selected: " + strconv.Itoa(info.NumSelectedDirs)
			if info.NumSelectedDirs == 1 {
				text += " directory - "
			} else {
				text += " directories - "
			}
			text += strconv.Itoa(info.NumSelectedFiles)
			if info.NumSelectedFiles == 1 {
				text += " file - "
			} else {
				text += " files - "
			}
			sizeStr, _ := formatSize(info.SizeSelectedFiles)
			text += sizeStr
			if info.SizeSelectedFiles == 1 {
				text += " byte"
			} else {
				text += " bytes"
			}
		}
		text = restrictBack(text, columns)
		setSourceColor(context, "000000")
		context.MoveTo(5, height-19)
		layout.SetText(text, -1)
		pango.CairoShowLayout(context, layout)
	}
}
