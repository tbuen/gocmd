package gui

import (
	"fmt"
	"github.com/gotk3/gotk3/cairo"
	"github.com/gotk3/gotk3/pango"
	"github.com/tbuen/gocmd/internal/backend"
	"github.com/tbuen/gocmd/internal/config"
	"github.com/tbuen/gocmd/internal/log"
	"strings"
	"unicode/utf8"
)

const (
	B = 1 << (10 * iota)
	KB
	MB
	GB
)

func setSourceColor(context *cairo.Context, color config.Color) {
	context.SetSourceRGB(color[0], color[1], color[2])
}

func drawPanel(context *cairo.Context, layout *pango.Layout, width, height float64, active bool, dir backend.Directory) {
	const scrollbarWidth = 8.0

	//    int cw, ch;
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
	log.Println(log.GUI, "lines:", lines)
	columns := int((width - 19) / cw)
	log.Println(log.GUI, "columns:", columns)

	if active {
		//context.SetSourceRGB(0x00, 0x40/255.0, 0xb0/255.0)
		//context.SetSourceRGB(0, 0x50/255.0, 0x70/255.0)
		context.SetSourceRGB(0x35/255.0, 0x84/255.0, 0xe4/255.0)
	} else {
		context.SetSourceRGB(0x70/255.0, 0x70/255.0, 0x70/255.0)
	}
	context.Rectangle(7, 7, width-13, ch+2)
	context.Fill()

	state := dir.State()
	path := dir.Path()
	if removeChars := utf8.RuneCountInString(path) - columns; removeChars > 0 {
		for i := 0; i <= removeChars; i++ {
			_, size := utf8.DecodeLastRuneInString(path)
			path = path[:len(path)-size]
		}
		path += "\u2026"
	}

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
		log.Println(log.GUI, "columns:", columns)

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
		for i := 0; i <= lines && offset+i < len(files); i++ {
			file := files[offset+i]
			if file.Marked() {
				context.SetSourceRGB(0xFF/255.0, 0xA0/255.0, 0x90/255.0)
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
			color := file.Color()
			//context.SetSourceRGB(color[0], color[1], color[2])
			setSourceColor(context, color)
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
			var size string
			if file.Dir() {
				size = "   " + "      <DIR>"
			} else {
				/*
					size = fmt.Sprintf("%11d", file.Size())
					switch s := file.Size(); {
					case s >= GB:
						size += fmt.Sprintf(" %4dG", s/GB)
					case s >= MB:
						size += fmt.Sprintf(" %4dM", s/MB)
					case s >= KB:
						size += fmt.Sprintf(" %4dK", s/KB)
					default:
						size += "      "//fmt.Sprintf(" %4dB", s)
					}*/
				switch s := file.Size(); {
				case s >= 1000000000:
					size += fmt.Sprintf("%2d.%03d.%03d.%03d", s/1000000000, (s%1000000000)/1000000, (s%1000000)/1000, s%1000)
				case s >= 1000000:
					size += fmt.Sprintf("   %3d.%03d.%03d", s/1000000, (s%1000000)/1000, s%1000)
				case s >= 1000:
					size += fmt.Sprintf("       %3d.%03d", s/1000, s%1000)
				default:
					size += fmt.Sprintf("           %3d", s)
				}
			}
			time := file.Time().Format("02.01.2006 15:04")
			neededSpace := 29 /*+ 6*/ + 3 + 10 + 13
			desiredNameLen := 10
			if columns-neededSpace > desiredNameLen {
				desiredNameLen = columns - neededSpace
			}
			if removeChars := utf8.RuneCountInString(line) - desiredNameLen; removeChars > 0 {
				for i := 0; i <= removeChars; i++ {
					_, size := utf8.DecodeLastRuneInString(line)
					line = line[:len(line)-size]
				}
				line += "\u2026"
			}
			if addSpaces := desiredNameLen - utf8.RuneCountInString(line); addSpaces > 0 {
				line += strings.Repeat(" ", addSpaces)
			}
			user, group := file.UserGroup()
			line += " " + size + " " + time + " " + user + "." + group + " " + file.Perm()

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
	/*long sel = directory.selection;
	  for (long i = 0; i <= lines && directory.offset + i < directory.files.length; ++i) {
	     auto file = directory.files[directory.offset + i];
	     if (file.marked) {
	        context.setSourceRgb(0xFF/255.0, 0xA0/255.0, 0x90/255.0);
	        context.rectangle(7, 11 + (i + 1) * ch, width - 13, ch);
	       context.fill;
	    }
	    if (directory.focus && directory.offset + i == sel) {
	        context.setSourceRgb(0, 0, 0);
	        context.rectangle(8, 11 + (i + 1) * ch, width - 14, ch);
	        context.stroke;
	     }
	     context.setSourceRgb(file.color[0]/255.0, file.color[1]/255.0, file.color[2]/255.0);
	     context.moveTo(10, 10 + (i + 1) * ch);
	     if (file.isLink && file.isDir) {
	        layout.setText("[" ~ file.name ~ "] -> " ~ file.link);
	     } else if (file.isLink) {
	        layout.setText(file.name ~ " -> " ~ file.link);
	     } else if (file.isDir) {
	        layout.setText("[" ~ file.name ~ "]");
	     } else {
	        layout.setText(file.name);
	     }
	     PgCairo.showLayout(context, layout);
	     auto time = file.time;
	     context.moveTo(width - 10 - 19 * cw, 10 + (i + 1) * ch);
	     layout.setText(format("%02d.%02d.%04d %02d:%02d:%02d", time.day, time.month, time.year, time.hour, time.minute, time.second));
	     PgCairo.showLayout(context, layout);
	  }*/

}
