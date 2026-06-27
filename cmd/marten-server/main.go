package main

import (
	"fmt"

	"github.com/tosterabgx/marten/internal/server"
)

func main() {
	fmt.Println("Starting Marten server")
	if err := server.RunControlServer(); err != nil {
		panic(err)
	}
}
