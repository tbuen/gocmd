package gui

import (
	"fmt"
	"strings"
	"unicode/utf8"
)

func formatSize(size int64) (result string, length int) {
	if size >= 1_000_000_000 {
		result = fmt.Sprintf("%3d.%03d.%03d.%03d", size/1000000000, (size%1000000000)/1000000, (size%1000000)/1000, size%1000)
	} else if size >= 1_000_000 {
		result = fmt.Sprintf("%3d.%03d.%03d", size/1000000, (size%1000000)/1000, size%1000)
	} else if size >= 1000 {
		result = fmt.Sprintf("%3d.%03d", size/1000, size%1000)
	} else {
		result = fmt.Sprintf("%3d", size)
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
