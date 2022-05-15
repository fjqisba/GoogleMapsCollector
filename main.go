package main

import (
	"GoogleMapsCollector/FyneApp"
	"GoogleMapsCollector/Utils/ProjectPath"
	"log"
	"os"
)

func main() {
	os.Setenv("google-chrome",ProjectPath.GProjectBinPath + "\\chrome\\chrome.exe")

	fyne := FyneApp.NewFyneApp()
	err := fyne.InitApp()
	if err != nil{
		log.Panicln(err)
	}
	fyne.Run()
}