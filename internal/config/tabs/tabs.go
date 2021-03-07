package tabs

import (
	"encoding/xml"
	"github.com/tbuen/gocmd/internal/log"
	"os"
	"path/filepath"
)

type Tab struct {
	Path      string `xml:"path,attr"`
	SortKey   int    `xml:"sortkey,attr"`
	SortOrder int    `xml:"sortorder,attr"`
	Hidden    bool   `xml:"hidden,attr"`
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

type Config struct {
	tc TabConfig
}

func (c *Config) Load(filename string) {
	buf, err := os.ReadFile(filename)
	if err != nil {
		log.Println(log.GLOBAL, err)
		return
	}

	err = xml.Unmarshal(buf, c.tc)
	if err != nil {
		log.Println(log.GLOBAL, err)
		return
	}
}

func (c *Config) Save(filename string) {
	buf, err := xml.MarshalIndent(c.tc, "", "\t")
	if err != nil {
		log.Println(log.GLOBAL, err)
		return
	}

	buf = append([]byte(xml.Header), buf...)

	err = os.MkdirAll(filepath.Dir(filename), 0777)
	if err != nil {
		log.Println(log.GLOBAL, err)
		return
	}
	err = os.WriteFile(filename, buf, 0666)
	if err != nil {
		log.Println(log.GLOBAL, err)
		return
	}
}
