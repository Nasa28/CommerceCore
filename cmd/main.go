package main

import (
	"log"

	"github.com/Nasa28/CommerceCore/cmd/api"
)

func main() {
	server := api.NewAPIServer(":4000", nil)
	if err := server.Run(); err != nil {
		log.Fatal(err)
	}
}
