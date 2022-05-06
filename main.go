package main

import (
	"GoogleMapsCollector/FyneApp"
	"log"
)

func main() {
	fyne := FyneApp.NewFyneApp()
	err := fyne.InitApp()
	if err != nil{
		log.Panicln(err)
	}
	fyne.Run()
}