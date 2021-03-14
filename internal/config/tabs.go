package config

import (
	"encoding/xml"
	"fmt"
)

// Directory is the configuration of one directory.
type Directory struct {
	Path string
}

// Panel is the configuration of one panel.
type Panel struct {
	Tabs   []Directory
	Active int
}

// Tabs is the tab configuration.
type Tabs struct {
	Panels []Panel
	Active int
}

type tabXML struct {
	Path      string `xml:"path,attr"`
	SortKey   int    `xml:"sortkey,attr"`
	SortOrder int    `xml:"sortorder,attr"`
	Hidden    bool   `xml:"hidden,attr"`
}

type panelXML struct {
	Active int      `xml:"active,attr"`
	Tabs   []tabXML `xml:"tab"`
}

type tabsXML struct {
	XMLName xml.Name   `xml:"tabs"`
	Active  int        `xml:"active,attr"`
	Panels  []panelXML `xml:"panel"`
}

var tabs Tabs

func GetTabs() *Tabs {
	return &tabs
}

func readTabs(filename string) {
	buf, err := load(filename)
	if err != nil {
		fmt.Println(err)
		return
	}

	var ttx tabsXML
	err = xml.Unmarshal(buf, &ttx)
	if err != nil {
		fmt.Println(err)
		return
	}

	for _, px := range ttx.Panels {
		p := Panel{}
		for _, tx := range px.Tabs {
			t := Directory{tx.Path}
			p.Tabs = append(p.Tabs, t)
		}
		if px.Active < len(p.Tabs) {
			p.Active = px.Active
		}
		tabs.Panels = append(tabs.Panels, p)
	}
	if ttx.Active < len(tabs.Panels) {
		tabs.Active = ttx.Active
	}

	return
}

func writeTabs(filename string) {
	var ttx tabsXML
	for _, p := range tabs.Panels {
		px := panelXML{}
		for _, t := range p.Tabs {
			tx := tabXML{Path: t.Path}
			px.Tabs = append(px.Tabs, tx)
		}
		if p.Active < len(px.Tabs) {
			px.Active = p.Active
		}
		ttx.Panels = append(ttx.Panels, px)
	}
	if tabs.Active < len(ttx.Panels) {
		ttx.Active = tabs.Active
	}

	buf, err := xml.MarshalIndent(ttx, "", "\t")
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
