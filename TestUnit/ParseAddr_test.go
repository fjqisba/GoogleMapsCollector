package TestUnit

import (
	"GoogleMapsCollector/Module/GooglePageScraper"
	"log"
	"regexp"
	"strings"
	"testing"
)

const strRegex_Website = `null,->\"(http.+?)>\",`


func TestParseAddr(t *testing.T)  {

	addrUrl := "https://maps.google.com/?cid=0x0:0x21db3bd2d0b69879"
	html := GooglePageScraper.GetGooglePageHtml(addrUrl)
	if html == ""{
		return
	}
	html = strings.ReplaceAll(html,"\\n","A")
	html = strings.ReplaceAll(html,"[","-")
	html = strings.ReplaceAll(html,"]","-")
	html = strings.ReplaceAll(html,"\\",">")

	regex_Website := regexp.MustCompile(strRegex_Website)
	tmpMatchList := regex_Website.FindStringSubmatch(html)
	if len(tmpMatchList) > 0{
		webSite := GooglePageScraper.DetectWebsite(tmpMatchList[1])
		log.Println(webSite)
	}

}
