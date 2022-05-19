package TestUnit

import (
	"GoogleMapsCollector/Module/EmailMiner"
	"log"
	"testing"
)

func TestGetEmail(t *testing.T) {

	webSite := "http://www.futurelighting.com/contact"
	emailList := EmailMiner.GetEmail(webSite)
	log.Println(emailList)
}