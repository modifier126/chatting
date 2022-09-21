package main

import (
	"chatDemo/router"
	"chatDemo/service"
	"fmt"
)

func main() {
	fmt.Println("main...")

	go service.Manager.Start()

	r := router.NewRouter()
	r.Run(":9080")
}
