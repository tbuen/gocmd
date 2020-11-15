package config

import (
	"encoding/xml"
	"github.com/tbuen/gocmd/internal/log"
	"os"
	"path/filepath"
)

type Tab struct {
	Path string `xml:"path,attr"`
}
type Panel struct {
	Active int   `xml:"active,attr"`
	Tabs   []Tab `xml:"tab"`
}
type TabConfig struct {
	XMLName xml.Name `xml:"tabs"`
	Panels  [2]Panel `xml:"panel"`
}

func WriteTabs(tabcfg *TabConfig) {
	output, err := xml.MarshalIndent(tabcfg, "", "   ")
	if err != nil {
		log.Println(log.CONFIG, "Error creating tabs xml", err)
		return
	}

	dir := filepath.Dir(filenameTabs)
	err = os.MkdirAll(dir, os.ModePerm)
	if err != nil {
		log.Println(log.GLOBAL, "Could not mkdir", dir)
		return
	}

	file, err := os.Create(filenameTabs)
	if err != nil {
		log.Println(log.GLOBAL, "Could not open", filenameTabs, "for writing")
		return
	}
	defer file.Close()

	_, err = file.WriteString(xml.Header)
	if err != nil {
		log.Println(log.GLOBAL, "Could not write", filenameTabs)
		return
	}
	_, err = file.Write(output)
	if err != nil {
		log.Println(log.GLOBAL, "Could not write", filenameTabs)
		return
	}
}
