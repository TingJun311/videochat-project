package main

import (
	"log"
	"videochat-project/internal/servers"
)

func main() {
	if err := servers.Run(); err != nil {
		log.Fatalln(err.Error())
	}	
}