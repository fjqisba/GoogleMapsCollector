package TestUnit

import (
	"GoogleMapsCollector/Module/EmailMiner"
	"GoogleMapsCollector/Utils/ProjectPath"
	"log"
	"os"
	"testing"
)

func TestGetEmail(t *testing.T) {

	os.Setenv("google-chrome",ProjectPath.GProjectBinPath + "\\chrome\\chrome.exe")
	emailList := EmailMiner.GetEmail("https://urbanlightsdenver.com/about")
	log.Println(emailList)
}