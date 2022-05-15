package PhantomJS

import (
	"GoogleMapsCollector/Utils/ProjectPath"
	"os/exec"
)

var(
	Instance PhantomJS
)

type PhantomJS struct {

}

func (this *PhantomJS)ScrapeGoogleMapHtml(url string)string  {
	cmd := exec.Command(ProjectPath.GProjectBinPath + "\\rsrc\\phantomjs.exe",
		ProjectPath.GProjectBinPath + "\\rsrc\\index.js",url)
	outPut,err := cmd.Output()
	if err != nil{
		return ""
	}
	return string(outPut)
}