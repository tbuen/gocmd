package apps

import (
	"encoding/xml"
	"github.com/tbuen/gocmd/internal/log"
	"os"
	"regexp"
	"strings"
)

type extcfg struct {
	descr string
	color string
	cmd   string
	args  []string
}

type progcfg struct {
	cmd  string
	args []string
}

type Config struct {
	ext  map[string]extcfg
	view progcfg
	edit progcfg
}

func (c *Config) FileColor(ext string) (col string) {
	col = "000000"
	if e, ok := c.ext[strings.ToUpper(ext)]; ok {
		col = e.color
	}
	return
}

func (c *Config) FileCmd(ext string) (cmd string, args []string) {
	if e, ok := c.ext[strings.ToUpper(ext)]; ok {
		cmd = e.cmd
		args = e.args
	}
	return
}

func (c *Config) View() (string, []string) {
	return c.view.cmd, c.view.args
}

func (c *Config) Edit() (string, []string) {
	return c.edit.cmd, c.edit.args
}

func (c *Config) Load(filename string) {
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

	buf, err := os.ReadFile(filename)
	if err != nil {
		log.Println(log.GLOBAL, err)
		return
	}

	appcfg := AppConfig{}
	err = xml.Unmarshal(buf, &appcfg)
	if err != nil {
		log.Println(log.GLOBAL, err)
		return
	}
	log.Println(log.CONFIG, "Number of apps:", len(appcfg.Apps))

	c.view.cmd = appcfg.View.Cmd
	c.view.args = appcfg.View.Args
	c.edit.cmd = appcfg.Edit.Cmd
	c.edit.args = appcfg.Edit.Args

	c.ext = make(map[string]extcfg)
	re := regexp.MustCompile("^#[0-9a-fA-F]{6}$")
	for i, a := range appcfg.Apps {
		col := "000000"
		if re.MatchString(a.Color) {
			col = a.Color[1:]
		} else {
			log.Println(log.GLOBAL, "Invalid color:", a.Color, "("+filename+")")
		}
		for _, e := range a.Exts {
			c.ext[strings.ToUpper(e)] = extcfg{a.Descr, col, a.Cmd, a.Args}
		}
		log.Println(log.CONFIG, i, "Descr:", a.Descr)
		log.Println(log.CONFIG, i, "Color: #"+col)
		log.Println(log.CONFIG, i, "Exts: ", a.Exts)
		log.Println(log.CONFIG, i, "Cmd:  ", a.Cmd)
		log.Println(log.CONFIG, i, "Args: ", a.Args)
	}
	log.Println(log.CONFIG, "Number of extensions:", len(c.ext))
	log.Println(log.CONFIG, "Extensions:", c.ext)
}
