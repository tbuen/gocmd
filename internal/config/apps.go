package config

import (
	"encoding/xml"
	"fmt"
	"regexp"
	"strings"
)

type ext struct {
	descr string
	color string
	cmd   string
	args  []string
}

type prog struct {
	cmd  string
	args []string
}

type appl struct {
	ext  map[string]ext
	view prog
	edit prog
}

type appXML struct {
	Descr string   `xml:"descr,attr"`
	Color string   `xml:"color,attr"`
	Exts  []string `xml:"ext"`
	Cmd   string   `xml:"cmd"`
	Args  []string `xml:"arg"`
}
type appsXML struct {
	XMLName xml.Name `xml:"apps"`
	View    appXML   `xml:"view"`
	Edit    appXML   `xml:"edit"`
	Apps    []appXML `xml:"app"`
}

var apps appl

// FileColor returns the color of a file.
func FileColor(ext string) (col string) {
	col = "000000"
	if e, ok := apps.ext[strings.ToUpper(ext)]; ok {
		col = e.color
	}
	return
}

// FileCmd returns the command and arguments associated with a file.
func FileCmd(ext string) (cmd string, args []string) {
	if e, ok := apps.ext[strings.ToUpper(ext)]; ok {
		cmd = e.cmd
		args = e.args
	}
	return
}

// View returns the command and arguments for viewing a file.
func View() (cmd string, args []string) {
	cmd, args = apps.view.cmd, apps.view.args
	return
}

// Edit returns the command and arguments for editing a file.
func Edit() (cmd string, args []string) {
	cmd, args = apps.edit.cmd, apps.edit.args
	return
}

func readApps(filename string) {
	buf, err := load(filename)
	if err != nil {
		fmt.Println(err)
		return
	}

	var ax appsXML
	err = xml.Unmarshal(buf, &ax)
	if err != nil {
		fmt.Println(err)
		return
	}

	apps.view.cmd = ax.View.Cmd
	apps.view.args = ax.View.Args
	apps.edit.cmd = ax.Edit.Cmd
	apps.edit.args = ax.Edit.Args

	apps.ext = make(map[string]ext)
	re := regexp.MustCompile("^#[0-9a-fA-F]{6}$")
	for _, app := range ax.Apps {
		col := "000000"
		if re.MatchString(app.Color) {
			col = app.Color[1:]
		}
		for _, e := range app.Exts {
			apps.ext[strings.ToUpper(e)] = ext{app.Descr, col, app.Cmd, app.Args}
		}
	}

	return
}
