package gui

import (
	"github.com/gotk3/gotk3/cairo"
	"github.com/gotk3/gotk3/pango"
	"github.com/tbuen/gocmd/internal/fs"
)

func drawPanel(context *cairo.Context, layout *pango.Layout, width, height float64, dir fs.Directory) {

	//    int cw, ch;
	//lines := 10

	context.SetSourceRGB(0, 0, 0)
	context.Rectangle(5, 5, width-8, height-9)
	context.Stroke()
	context.Rectangle(6, 6, width-10, height-10)
	context.Clip()

	//layout.SetText(".");
	//layout.GetPixelSize(cw, ch);
	//lines = (height - 19 - ch) / ch;
	//debug writeln("lines ", lines);

	/*if (directory.selection >= directory.offset + lines) {
	     directory.offset = directory.selection - lines + 1;
	  }
	  if (directory.selection < directory.offset) {
	     directory.offset = directory.selection;
	  }*/

	//////context.setSourceRgb(0x00, 0x40/255.0, 0xb0/255.0);
	//if (directory.focus) {
	//context.setSourceRgb(0x4a/255.0, 0x90/255.0, 0xd9/255.0);
	context.SetSourceRGB(0, 0x50/255.0, 0x70/255.0)
	//} else {
	//  context.setSourceRgb(0x70/255.0, 0x70/255.0, 0x70/255.0);
	//}
	ch := 15.0
	context.Rectangle(7, 7, width-13, ch+2)
	context.Fill()

	context.SetSourceRGB(1, 1, 1)
	context.MoveTo(10, 8)
	layout.SetText(dir.Path(), -1)
	pango.CairoShowLayout(context, layout)

	if dir.State() == fs.STATE_IDLE {
		for i, file := range dir.Files() {
			context.SetSourceRGB(0, 0, 0)
			context.MoveTo(10, 10+(float64(i)+1)*ch)
			layout.SetText(file.Name(), -1)
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
