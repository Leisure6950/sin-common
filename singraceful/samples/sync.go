package main

import (
	"fmt"
	"github.com/sin-z/sin-common/singraceful"
	"time"
)

func task() error {
	fmt.Println("task")
	time.Sleep(time.Second * 10)
	return nil
}
func main() {
	go func() {
		singraceful.Sync(task)
		singraceful.Sync(task)
	}()
	singraceful.RunBackend()
	fmt.Printf("backend done\n")
}
