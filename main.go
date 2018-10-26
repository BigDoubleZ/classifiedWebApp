package main

import (
	"classfiedWebApp/app"
	"fmt"
	"os"
)

func main() {
	app, err := app.Init()
	if err != nil {
		fmt.Println("Failed to init app")
		os.Exit(1)
	}

	app.Run()
}
