package bookmarks

import (
	"encoding/xml"
	. "github.com/tbuen/gocmd/internal/global"
	"github.com/tbuen/gocmd/internal/log"
	"os"
	"path/filepath"
)

type Config struct {
	bookmarks []Bookmark
	isChanged bool
}

type XMLBookmark struct {
	Path string `xml:"path,attr"`
	//SortKey   int    `xml:"sortkey,attr"`
	//SortOrder int    `xml:"sortorder,attr"`
	//Hidden    bool   `xml:"hidden,attr"`
}

type XMLBookmarks struct {
	Name      xml.Name      `xml:"bookmarks"`
	Bookmarks []XMLBookmark `xml:"bookmark"`
}

func (c *Config) Get() (bb []Bookmark) {
	bb = make([]Bookmark, len(c.bookmarks))
	copy(bb, c.bookmarks)
	return
}

func (c *Config) Add(b Bookmark) {
	c.bookmarks = append(c.bookmarks, b)
}

func (c *Config) Load(filename string) {
	var xmlBookmarks XMLBookmarks

	file, err := os.Open(filename)
	if err != nil {
		log.Println(log.CONFIG, "Could not open", filename)
		return
	}
	defer file.Close()

	fileInfo, err := file.Stat()
	if err != nil {
		log.Println(log.CONFIG, "Could not stat", filename)
		return
	}
	buffer := make([]byte, fileInfo.Size())
	_, err = file.ReadAt(buffer, 0)
	if err != nil {
		log.Println(log.CONFIG, "Could not read", filename)
		return
	}

	err = xml.Unmarshal(buffer, &xmlBookmarks)
	if err != nil {
		log.Println(log.CONFIG, "Could not parse", filename, err)
		return
	}

	for _, b := range xmlBookmarks.Bookmarks {
		c.bookmarks = append(c.bookmarks, Bookmark{b.Path})
	}
}

func (c *Config) Save(filename string) {
	if !c.isChanged {
		return
	}

	var xmlBookmarks XMLBookmarks

	for _, b := range c.bookmarks {
		xmlBookmarks.Bookmarks = append(xmlBookmarks.Bookmarks, XMLBookmark{b.Path})
	}

	output, err := xml.MarshalIndent(&xmlBookmarks, "", "\t")
	if err != nil {
		log.Println(log.CONFIG, "Error creating bookmarks xml", err)
		return
	}

	dir := filepath.Dir(filename)
	err = os.MkdirAll(dir, os.ModePerm)
	if err != nil {
		log.Println(log.GLOBAL, "Could not mkdir", dir)
		return
	}

	file, err := os.Create(filename)
	if err != nil {
		log.Println(log.GLOBAL, "Could not open", filename, "for writing")
		return
	}
	defer file.Close()

	_, err = file.WriteString(xml.Header)
	if err != nil {
		log.Println(log.GLOBAL, "Could not write", filename)
		return
	}
	_, err = file.Write(output)
	if err != nil {
		log.Println(log.GLOBAL, "Could not write", filename)
		return
	}
}
