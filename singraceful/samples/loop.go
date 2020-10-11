package main

import (
	"fmt"
	"github.com/sin-z/sin-common/singraceful"
	"time"
)

func task() {
	fmt.Println("task")
	time.Sleep(time.Millisecond * 500)
}
func main() {
	singraceful.Loop(task)
	singraceful.RunBackend()
	fmt.Printf("backend done\n")
}
