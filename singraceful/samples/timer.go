package main

import (
	"fmt"
	"github.com/sin-z/sin-common/singraceful"
	"time"
)

func task() {
	fmt.Println("task")
}
func main() {
	singraceful.Timer(time.Second*5, task)
	singraceful.Timer(time.Second*5, task)
	singraceful.Timer(time.Second*5, task)
	singraceful.RunBackend()
	fmt.Printf("backend done\n")
}
