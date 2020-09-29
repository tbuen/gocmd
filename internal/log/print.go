package log

import (
	"fmt"
	"os"
)

func Println(module int, vars ...interface{}) {
	if modules&module != 0 {
		var s string
		for _, x := range vars {
			s += fmt.Sprint(x)
		}
		fmt.Println(s)
	}
}

func Fatalln(vars ...interface{}) {
	var s string
	for _, x := range vars {
		s += fmt.Sprint(x)
	}
	fmt.Println(s)
	os.Exit(1)
}
