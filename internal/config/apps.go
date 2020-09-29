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
}

var extcfgs map[string]extcfg

func Read() {
	type App struct {
		Descr string   `xml:"descr,attr"`
		Color string   `xml:"color,attr"`
		Ext   []string `xml:"ext"`
		Cmd   string   `xml:"cmd"`
	}
	type AppConfig struct {
		XMLName xml.Name `xml:"apps"`
		Apps    []App    `xml:"app"`
	}

	file, err := os.Open(filenameApps)
	if err != nil {
		log.Println(log.GLOBAL, "Could not open ", filenameApps)
		return
	}
	defer file.Close()

	fileInfo, err := file.Stat()
	if err != nil {
		log.Println(log.GLOBAL, "Could not stat ", filenameApps)
		return
	}
	buffer := make([]byte, fileInfo.Size())
	_, err = file.ReadAt(buffer, 0)
	if err != nil {
		log.Println(log.GLOBAL, "Could not read ", filenameApps)
		return
	}

	appcfg := AppConfig{}
	err = xml.Unmarshal(buffer, &appcfg)
	if err != nil {
		log.Println(log.GLOBAL, "Could not parse ", filenameApps, ": ", err)
		return
	}
	log.Println(log.CONFIG, "Number of apps: ", len(appcfg.Apps))

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
			log.Println(log.GLOBAL, "Invalid color: ", a.Color, " ("+filenameApps+")")
		}
		for _, e := range a.Ext {
			extcfgs[e] = extcfg{a.Descr, c, a.Cmd}
		}
		log.Println(log.CONFIG, i, ". Descr: ", a.Descr)
		log.Println(log.CONFIG, i, ". Color: ", a.Color)
		log.Println(log.CONFIG, i, ". Color: ", c)
		log.Println(log.CONFIG, i, ". Ext: ", a.Ext)
		log.Println(log.CONFIG, i, ". Cmd: ", a.Cmd)
	}
	log.Println(log.CONFIG, "Number of extensions: ", len(extcfgs))
	log.Println(log.CONFIG, "Extensions: ", extcfgs)
}

func FileColor(ext string) (c Color) {
	if e, ok := extcfgs[ext]; ok {
		c = e.color
	}
	return
}
