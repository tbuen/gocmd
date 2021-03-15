package config

import (
	"encoding/xml"
	"fmt"
)

// Bookmark is a bookmark.
type Bookmark struct {
	Path      string
	SortKey   int
	SortOrder int
	Hidden    bool
}

type bookmarks struct {
	bookmarks []Bookmark
	isChanged bool
}

type bookmarkXML struct {
	Path      string `xml:"path,attr"`
	SortKey   int    `xml:"sortkey,attr"`
	SortOrder int    `xml:"sortorder,attr"`
	Hidden    bool   `xml:"hidden,attr"`
}

type bookmarksXML struct {
	Name      xml.Name      `xml:"bookmarks"`
	Bookmarks []bookmarkXML `xml:"bookmark"`
}

var bookm bookmarks

// Bookmarks returns the bookmarks.
func Bookmarks() []Bookmark {
	return bookm.bookmarks
}

// AddBookmark adds a bookmark.
func AddBookmark(b Bookmark) {
	bookm.bookmarks = append(bookm.bookmarks, b)
}

func readBookmarks(filename string) {
	buf, err := load(filename)
	if err != nil {
		fmt.Println(err)
		return
	}

	var bx bookmarksXML
	err = xml.Unmarshal(buf, &bx)
	if err != nil {
		fmt.Println(err)
		return
	}

	for _, b := range bx.Bookmarks {
		bookm.bookmarks = append(bookm.bookmarks, Bookmark{b.Path, b.SortKey, b.SortOrder, b.Hidden})
	}
}

func writeBookmarks(filename string) {
	if !bookm.isChanged {
		return
	}

	var bx bookmarksXML
	for _, b := range bookm.bookmarks {
		bx.Bookmarks = append(bx.Bookmarks, bookmarkXML{b.Path, b.SortKey, b.SortOrder, b.Hidden})
	}

	buf, err := xml.MarshalIndent(bx, "", "\t")
	if err != nil {
		fmt.Println(err)
		return
	}
	buf = append([]byte(xml.Header), buf...)

	save(filename, buf)
	if err != nil {
		fmt.Println(err)
		return
	}

	return
}
