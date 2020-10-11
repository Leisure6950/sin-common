package main

import (
	"fmt"
	"github.com/sin-z/sin-common/sintime"
	"time"
)

func main() {
	fmt.Println(sintime.Now().FormatDef())
	fmt.Println(sintime.Now().DayBegin().FormatDef())
	fmt.Println(sintime.Now().DayEnd().FormatDef())
	fmt.Println(sintime.FormatDef(time.Now()))
	fmt.Println(sintime.ParseDef("2020-05-14 17:00:00"))
}
