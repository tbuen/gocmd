package config

import (
	"encoding/xml"
	"github.com/tbuen/gocmd/internal/log"
	"os"
	"regexp"
	"strconv"
)

type Color [3]float64

type extcfg struct {
	descr string
	color Color
	cmd   string
	args  []string
}

type progcfg struct {
	cmd  string
	args []string
}

var extcfgs map[string]extcfg
var view progcfg
var edit progcfg

func Read() {
	type App struct {
		Descr string   `xml:"descr,attr"`
		Color string   `xml:"color,attr"`
		Exts  []string `xml:"ext"`
		Cmd   string   `xml:"cmd"`
		Args  []string `xml:"arg"`
	}
	type AppConfig struct {
		XMLName xml.Name `xml:"apps"`
		View    App      `xml:"view"`
		Edit    App      `xml:"edit"`
		Apps    []App    `xml:"app"`
	}

	file, err := os.Open(filenameApps)
	if err != nil {
		log.Println(log.GLOBAL, "Could not open", filenameApps)
		return
	}
	defer file.Close()

	fileInfo, err := file.Stat()
	if err != nil {
		log.Println(log.GLOBAL, "Could not stat", filenameApps)
		return
	}
	buffer := make([]byte, fileInfo.Size())
	_, err = file.ReadAt(buffer, 0)
	if err != nil {
		log.Println(log.GLOBAL, "Could not read", filenameApps)
		return
	}

	appcfg := AppConfig{}
	err = xml.Unmarshal(buffer, &appcfg)
	if err != nil {
		log.Println(log.GLOBAL, "Could not parse", filenameApps, err)
		return
	}
	log.Println(log.CONFIG, "Number of apps:", len(appcfg.Apps))

	view.cmd = appcfg.View.Cmd
	view.args = appcfg.View.Args
	edit.cmd = appcfg.Edit.Cmd
	edit.args = appcfg.Edit.Args

	extcfgs = make(map[string]extcfg)
	re := regexp.MustCompile("^#([0-9a-fA-F]{2})([0-9a-fA-F]{2})([0-9a-fA-F]{2})$")
	for i, a := range appcfg.Apps {
		c := Color{}
		cols := re.FindStringSubmatch(a.Color)
		if len(cols) == 4 {
			for i := 0; i < 3; i++ {
				if ii, err := strconv.ParseUint(cols[i+1], 16, 8); err == nil {
					c[i] = float64(ii) / 255
				}
			}
		} else {
			log.Println(log.GLOBAL, "Invalid color:", a.Color, "("+filenameApps+")")
		}
		for _, e := range a.Exts {
			extcfgs[e] = extcfg{a.Descr, c, a.Cmd, a.Args}
		}
		log.Println(log.CONFIG, i, "Descr:", a.Descr)
		log.Println(log.CONFIG, i, "Color:", a.Color)
		log.Println(log.CONFIG, i, "Color:", c)
		log.Println(log.CONFIG, i, "Exts: ", a.Exts)
		log.Println(log.CONFIG, i, "Cmd:  ", a.Cmd)
		log.Println(log.CONFIG, i, "Args: ", a.Args)
	}
	log.Println(log.CONFIG, "Number of extensions:", len(extcfgs))
	log.Println(log.CONFIG, "Extensions:", extcfgs)
}

func FileColor(ext string) (c Color) {
	if e, ok := extcfgs[ext]; ok {
		c = e.color
	}
	return
}

func FileCmd(ext string) (c string, a []string) {
	if e, ok := extcfgs[ext]; ok {
		c = e.cmd
		a = e.args
	}
	return
}

func View() (string, []string) {
	return view.cmd, view.args
}

func Edit() (string, []string) {
	return edit.cmd, edit.args
}
