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
	singraceful.Go(task)
	singraceful.Go(task)
	singraceful.Go(task)
	singraceful.Go(task)
	singraceful.RunBackend()
	fmt.Printf("backend done\n")
}
