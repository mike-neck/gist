package main

import (
	"log"
	"os"
)

func main() {
	app := GetApplication(os.Args)
	err := app.Start()
	if err != nil {
		log.Fatalln("error", err)
	}
}
