package gui

import (
	"fmt"
	"github.com/gotk3/gotk3/cairo"
	"regexp"
	"strconv"
	"strings"
	"unicode/utf8"
)

var colorRegexp = regexp.MustCompile("^([0-9a-fA-F]{2})([0-9a-fA-F]{2})([0-9a-fA-F]{2})$")

func setSourceColor(context *cairo.Context, color string) {
	var c [3]float64
	cols := colorRegexp.FindStringSubmatch(color)
	if len(cols) == 4 {
		for i := 0; i < 3; i++ {
			if ii, err := strconv.ParseUint(cols[i+1], 16, 8); err == nil {
				c[i] = float64(ii) / 255
			}
		}
	}
	context.SetSourceRGB(c[0], c[1], c[2])
}

func formatSize(size int64) (result string, length int) {
	if size >= 1_000_000_000 {
		result = fmt.Sprintf("%d.%03d.%03d.%03d", size/1000000000, (size%1000000000)/1000000, (size%1000000)/1000, size%1000)
	} else if size >= 1_000_000 {
		result = fmt.Sprintf("%d.%03d.%03d", size/1000000, (size%1000000)/1000, size%1000)
	} else if size >= 1000 {
		result = fmt.Sprintf("%d.%03d", size/1000, size%1000)
	} else {
		result = fmt.Sprintf("%d", size)
	}
	length = len(result)
	return
}

func appendSpaces(text string, minLen int) (result string) {
	result = text
	if n := minLen - utf8.RuneCountInString(text); n > 0 {
		result += strings.Repeat(" ", n)
	}
	return
}

func prependSpaces(text string, minLen int) (result string) {
	result = text
	if n := minLen - utf8.RuneCountInString(text); n > 0 {
		result = strings.Repeat(" ", n) + result
	}
	return
}

func restrictFront(text string, maxLen int) (result string) {
	result = text
	if n := utf8.RuneCountInString(text) - maxLen; n > 0 {
		for i := 0; i <= n; i++ {
			_, size := utf8.DecodeRuneInString(result)
			result = result[size:]
		}
		result = "\u2026" + result
	}
	return
}

func restrictBack(text string, maxLen int) (result string) {
	result = text
	if n := utf8.RuneCountInString(text) - maxLen; n > 0 {
		for i := 0; i <= n; i++ {
			_, size := utf8.DecodeLastRuneInString(result)
			result = result[:len(result)-size]
		}
		result += "\u2026"
	}
	return
}
