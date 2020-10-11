package main

import (
	"fmt"
	"github.com/sin-z/sin-common/singraceful"
	"time"
)

func task() {
	fmt.Println("task")
	return
}
func main() {
	singraceful.Ticker(time.Second*1, task)
	singraceful.Ticker(time.Second*1, task)
	singraceful.Ticker(time.Second*1, task)
	singraceful.RunBackend()

	fmt.Printf("backend done\n")
}
