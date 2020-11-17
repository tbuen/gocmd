package config

import (
	"encoding/xml"
	"github.com/tbuen/gocmd/internal/log"
	"os"
	"path/filepath"
)

const (
	SORT_BY_NAME = iota
	SORT_BY_EXT
	SORT_BY_SIZE
	SORT_BY_TIME
)

const (
	SORT_ASCENDING = iota
	SORT_DESCENDING
)

type Tab struct {
	Path      string `xml:"path,attr"`
	SortKey   int    `xml:"sortkey,attr"`
	SortOrder int    `xml:"sortorder,attr"`
}

type Panel struct {
	Active int   `xml:"active,attr"`
	Tabs   []Tab `xml:"tab"`
}

type TabConfig struct {
	XMLName xml.Name `xml:"tabs"`
	Active  int      `xml:"active,attr"`
	Panels  []Panel  `xml:"panel"`
}

func ReadTabs() (tabcfg *TabConfig, err error) {
	file, err := os.Open(filenameTabs)
	if err != nil {
		log.Println(log.CONFIG, "Could not open", filenameTabs)
		return
	}
	defer file.Close()

	fileInfo, err := file.Stat()
	if err != nil {
		log.Println(log.CONFIG, "Could not stat", filenameTabs)
		return
	}
	buffer := make([]byte, fileInfo.Size())
	_, err = file.ReadAt(buffer, 0)
	if err != nil {
		log.Println(log.CONFIG, "Could not read", filenameTabs)
		return
	}

	tabcfg = new(TabConfig)
	err = xml.Unmarshal(buffer, tabcfg)
	if err != nil {
		log.Println(log.CONFIG, "Could not parse", filenameTabs, err)
		tabcfg = nil
		return
	}
	return
}

func WriteTabs(tabcfg *TabConfig) {
	output, err := xml.MarshalIndent(tabcfg, "", "\t")
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
