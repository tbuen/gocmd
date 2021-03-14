package config

import (
	//"encoding/xml"
	. "github.com/tbuen/gocmd/internal/global"
	//"github.com/tbuen/gocmd/internal/log"
	//"os"
	//"path/filepath"
)

type Bookmarks struct {
	bookmarks []Bookmark
	isChanged bool
}

func (c *Bookmarks) Get() (bb []Bookmark) {
	bb = make([]Bookmark, len(c.bookmarks))
	copy(bb, c.bookmarks)
	return
}

/*
func (c *Config) Add(b Bookmark) {
	c.bookmarks = append(c.bookmarks, b)
}

func (c *Config) Load(filename string) {
	type XmlBookmark struct {
		Path string `xml:"path,attr"`
		//SortKey   int    `xml:"sortkey,attr"`
		//SortOrder int    `xml:"sortorder,attr"`
		//Hidden    bool   `xml:"hidden,attr"`
	}
	type Xml struct {
		Name      xml.Name      `xml:"bookmarks"`
		Bookmarks []XmlBookmark `xml:"bookmark"`
	}

	buf, err := os.ReadFile(filename)
	if err != nil {
		log.Println(log.GLOBAL, err)
		return
	}

	var cfg Xml
	err = xml.Unmarshal(buf, &cfg)
	if err != nil {
		log.Println(log.GLOBAL, err)
		return
	}

	for _, b := range cfg.Bookmarks {
		c.bookmarks = append(c.bookmarks, Bookmark{b.Path})
	}
}

func (c *Config) Save(filename string) {
	type XmlBookmark struct {
		Path string `xml:"path,attr"`
		//SortKey   int    `xml:"sortkey,attr"`
		//SortOrder int    `xml:"sortorder,attr"`
		//Hidden    bool   `xml:"hidden,attr"`
	}
	type Xml struct {
		Name      xml.Name      `xml:"bookmarks"`
		Bookmarks []XmlBookmark `xml:"bookmark"`
	}

	if !c.isChanged {
		return
	}

	var cfg Xml
	for _, b := range c.bookmarks {
		cfg.Bookmarks = append(cfg.Bookmarks, XmlBookmark{b.Path})
	}

	buf, err := xml.MarshalIndent(&cfg, "", "\t")
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
}*/
