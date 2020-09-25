package log

import (
	"fmt"
	"log"
)

func Println(module int, text string, vars ...interface{}) {
	if config&module != 0 {
		s := text
		for _, x := range vars {
			s += " " + fmt.Sprint(x)
		}
		log.Println(s)
	}
}

func Fatal(text string, vars ...interface{}) {
	s := text
	for _, x := range vars {
		s += " " + fmt.Sprint(x)
	}
	log.Fatal(s)
}
