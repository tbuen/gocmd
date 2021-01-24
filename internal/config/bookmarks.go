package config

import (
	"encoding/xml"
	"github.com/tbuen/gocmd/internal/log"
	"os"
	"path/filepath"
)

type Bookmark struct {
	Path string `xml:"path,attr"`
	//SortKey   int    `xml:"sortkey,attr"`
	//SortOrder int    `xml:"sortorder,attr"`
	//Hidden    bool   `xml:"hidden,attr"`
}

type BookmarkConfig struct {
	XMLName   xml.Name   `xml:"bookmarks"`
	Bookmarks []Bookmark `xml:"bookmark"`
}

var bookmarkcfg BookmarkConfig

func ReadBookmarks() {
	file, err := os.Open(filenameBookmarks)
	if err != nil {
		log.Println(log.CONFIG, "Could not open", filenameBookmarks)
		return
	}
	defer file.Close()

	fileInfo, err := file.Stat()
	if err != nil {
		log.Println(log.CONFIG, "Could not stat", filenameBookmarks)
		return
	}
	buffer := make([]byte, fileInfo.Size())
	_, err = file.ReadAt(buffer, 0)
	if err != nil {
		log.Println(log.CONFIG, "Could not read", filenameBookmarks)
		return
	}

	err = xml.Unmarshal(buffer, &bookmarkcfg)
	if err != nil {
		log.Println(log.CONFIG, "Could not parse", filenameBookmarks, err)
		return
	}
}

func WriteBookmarks() {
	output, err := xml.MarshalIndent(&bookmarkcfg, "", "\t")
	if err != nil {
		log.Println(log.CONFIG, "Error creating bookmarks xml", err)
		return
	}

	dir := filepath.Dir(filenameBookmarks)
	err = os.MkdirAll(dir, os.ModePerm)
	if err != nil {
		log.Println(log.GLOBAL, "Could not mkdir", dir)
		return
	}

	file, err := os.Create(filenameBookmarks)
	if err != nil {
		log.Println(log.GLOBAL, "Could not open", filenameBookmarks, "for writing")
		return
	}
	defer file.Close()

	_, err = file.WriteString(xml.Header)
	if err != nil {
		log.Println(log.GLOBAL, "Could not write", filenameBookmarks)
		return
	}
	_, err = file.Write(output)
	if err != nil {
		log.Println(log.GLOBAL, "Could not write", filenameBookmarks)
		return
	}
}
