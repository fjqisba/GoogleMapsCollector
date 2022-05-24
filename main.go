package main

import (
	"GoogleMapsCollector/FyneApp"
	"GoogleMapsCollector/Utils"
	"GoogleMapsCollector/Utils/ProjectPath"
	"log"
	"os"
)

func main() {
	_, err := Utils.CreateMutex("GoogleMapCollector")
	if err != nil{
		log.Println("程序禁止多开")
		return
	}
	os.Setenv("google-chrome",ProjectPath.GProjectBinPath + "\\chrome\\chrome.exe")
	fyne := FyneApp.NewFyneApp()
	err = fyne.InitApp()
	if err != nil{
		log.Panicln(err)
	}
	fyne.Run()
	fyne.RunServer()
}