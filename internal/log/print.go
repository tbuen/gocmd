package log

import (
	"fmt"
	"os"
)

func Println(module int, vars ...interface{}) {
	if modules&module != 0 {
		fmt.Println(vars...)
	}
}

func Fatalln(vars ...interface{}) {
	fmt.Println(vars...)
	os.Exit(1)
}
