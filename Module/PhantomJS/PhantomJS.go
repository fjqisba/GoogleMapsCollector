package PhantomJS

import (
	"GoogleMapsCollector/Utils/ProjectPath"
	"os/exec"
)


func GetPageHtml(url string)string  {
	cmd := exec.Command(ProjectPath.GProjectBinPath + "\\rsrc\\phantomjs.exe",
		ProjectPath.GProjectBinPath + "\\rsrc\\page.js",url)
	outPut,err := cmd.Output()
	if err != nil{
		return ""
	}
	return string(outPut)
}

func ScrapeGoogleMapHtml(url string)string  {
	cmd := exec.Command(ProjectPath.GProjectBinPath + "\\rsrc\\phantomjs.exe",
		ProjectPath.GProjectBinPath + "\\rsrc\\index.js",url)
	outPut,err := cmd.Output()
	if err != nil{
		return ""
	}
	return string(outPut)
}