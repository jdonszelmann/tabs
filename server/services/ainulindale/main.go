package main

import (
	"github.com/jonay2000/ainulindale/server/pkg/server"
	"log"
)

func main() {
	if err := server.StartServer(); err != nil {
		log.Fatal(err)
	}
}
