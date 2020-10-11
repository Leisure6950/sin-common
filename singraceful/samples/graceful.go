package main

import (
	"fmt"
	"github.com/sin-z/sin-common/singraceful"
	"time"
)

func main() {
	singraceful.Add(1)
	go func() {
		time.Sleep(time.Second * 10)
		singraceful.Done()
		fmt.Printf("task done\n")
	}()

	singraceful.Go(func() error {
		return nil
	}).Error()
	fmt.Printf("start 1 service\n")
	singraceful.RunServer(func() error {
		time.Sleep(time.Second * 3)
		return nil
	})
	fmt.Printf("service done\n")
}
