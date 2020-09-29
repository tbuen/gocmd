package gui

import (
	"github.com/gotk3/gotk3/cairo"
	"github.com/gotk3/gotk3/pango"
	"github.com/tbuen/gocmd/internal/fs"
	"github.com/tbuen/gocmd/internal/log"
)

func drawPanel(context *cairo.Context, layout *pango.Layout, width, height float64, active bool, dir fs.Directory) {

	//    int cw, ch;
	ch := 15.0

	state := dir.State()
	selection := dir.Selection()
	offset := dir.DispOffset()

	context.SetSourceRGB(0, 0, 0)
	context.Rectangle(5, 5, width-8, height-9)
	context.Stroke()
	context.Rectangle(6, 6, width-10, height-10)
	context.Clip()

	//layout.SetText(".");
	//layout.GetPixelSize(cw, ch);
	lines := int((height - 19 - ch) / ch)
	log.Println(log.GUI, "lines:", lines)

	if selection >= offset+lines {
		offset = selection - lines + 1
	}
	if selection < offset {
		offset = selection
	}
	dir.SetDispOffset(offset)

	//////context.setSourceRgb(0x00, 0x40/255.0, 0xb0/255.0);
	if active {
		//context.setSourceRgb(0x4a/255.0, 0x90/255.0, 0xd9/255.0);
		context.SetSourceRGB(0, 0x50/255.0, 0x70/255.0)
	} else {
		context.SetSourceRGB(0x70/255.0, 0x70/255.0, 0x70/255.0)
	}
	context.Rectangle(7, 7, width-13, ch+2)
	context.Fill()

	switch state {
	case fs.STATE_IDLE:
		context.SetSourceRGB(1, 1, 1)
	case fs.STATE_ERROR:
		context.SetSourceRGB(1, 0, 0)
	case fs.STATE_RELOAD:
		context.SetSourceRGB(1, 1, 0)
	}
	context.MoveTo(10, 8)
	layout.SetText(dir.Path(), -1)
	pango.CairoShowLayout(context, layout)

	if state == fs.STATE_IDLE {
		files := dir.Files()
		for i := 0; i <= lines && offset+i < len(files); i++ {
			file := files[offset+i]
			if file.IsMarked() {
				context.SetSourceRGB(0xFF/255.0, 0xA0/255.0, 0x90/255.0)
				context.Rectangle(7, 11+(float64(i)+1)*ch, width-13, ch)
				context.Fill()
			}
			if active && offset+i == selection {
				context.SetSourceRGB(0, 0, 0)
				context.Rectangle(8, 11+(float64(i)+1)*ch, width-14, ch)
				context.Stroke()
			}
			context.SetSourceRGB(0, 0, 0)
			context.MoveTo(10, 10+(float64(i)+1)*ch)
			if file.IsDir() {
				layout.SetText("["+file.Name()+"]", -1)
			} else {
				layout.SetText(file.Name(), -1)
			}
			pango.CairoShowLayout(context, layout)
		}
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
