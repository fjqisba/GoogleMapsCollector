package GooglePageScraper

import (
	"GoogleMapsCollector/ConfigManager"
	"GoogleMapsCollector/DataBase"
	"GoogleMapsCollector/Model"
	"GoogleMapsCollector/Module/CsvResult"
	"GoogleMapsCollector/Module/EmailMiner"
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"regexp"
	"strings"
	"time"
)

const strRegex_RealCategory = `null,null,null,null,null,null,null,null,null,null,null,null,null,null,null,null,null,null,null,null,null,null,null,null,null,null,null,null,null,null,null,null,null,null,null,null,null,null,null,null,null,null,null,null,null,null,null,null,null,null,null,null,null,null,null,null,null,null,null,null,null,null,null,null,null,null,null,null,>\"(.*?)>\"-`
const strRegex_BusinessName = `---7,-->\"(.*?)\"-`
const strRegex_Address = `-1,-->\"(.*?)>\"`
const strRegex_City = `-4,-->\"(.*?)>\"`

const strRegex_PostalCode1 = `(\\d{5,6}) ([^,]+),`
const strRegex_PostalCode2 = `(\\d{5,6}) (.*)`
const strRegex_PostalCode3 = `, ([^\\d]+)\\s(\\d{5,6}), (.*)`
const strRegex_PostalCode4 = `, (.*), (.*) (\\d{5,6})`

const strRegex_Landtitude = `https://www.google.com/maps/preview/place/([^/]+)/@(.*?),(.*?),`

const strRegex_Website = `null,->\"(http.+?)>\",`
const strRegex_Phone = `tel:(.*?)>\"`

var(
	regex_RealCategory *regexp.Regexp
	regex_BusinessName *regexp.Regexp
	regex_Address *regexp.Regexp
	regex_City *regexp.Regexp
	regex_PostalCode1 *regexp.Regexp
	regex_PostalCode2 *regexp.Regexp
	regex_PostalCode3 *regexp.Regexp
	regex_PostalCode4 *regexp.Regexp
	regex_Landtitude *regexp.Regexp
	regex_Website *regexp.Regexp
	regex_Phone *regexp.Regexp
)

//尝试从地址中解析出邮政编码和城市
func parseAddress(scrapeData *Model.ScraperData)  {
	tmpMatchList := regex_PostalCode1.FindStringSubmatch(scrapeData.Address)
	if len(tmpMatchList) > 0{
		scrapeData.PostalCode = tmpMatchList[1]
		scrapeData.City = tmpMatchList[2]
		return
	}
	tmpMatchList = regex_PostalCode2.FindStringSubmatch(scrapeData.Address)
	if len(tmpMatchList) > 0{
		scrapeData.PostalCode = tmpMatchList[1]
		scrapeData.City = tmpMatchList[2]
		return
	}
	tmpMatchList = regex_PostalCode3.FindStringSubmatch(scrapeData.Address)
	if len(tmpMatchList) > 0{
		scrapeData.PostalCode = tmpMatchList[1]
		scrapeData.City = tmpMatchList[2]
		return
	}
	tmpMatchList = regex_PostalCode4.FindStringSubmatch(scrapeData.Address)
	if len(tmpMatchList) > 0{
		scrapeData.PostalCode = tmpMatchList[1]
		scrapeData.City = tmpMatchList[3]
		return
	}
}

func DetectWebsite(matchUrl string)string{
	if strings.Index(matchUrl,"googleusercontent.com") != -1{
		return ""
	}
	if strings.Index(matchUrl,"www.google.com") != -1{
		return ""
	}
	matchUrl = strings.ReplaceAll(matchUrl,">>u003d","=")
	matchUrl = strings.ReplaceAll(matchUrl,">>u0026","&")
	return matchUrl
}

func getHtmlContent(pageUrl string)string  {
	proxyFunc := http.ProxyFromEnvironment
	proxyUrl := ConfigManager.Instance.GetSystemProxy()
	if proxyUrl != ""{
		proxyFunc = func(req *http.Request) (*url.URL, error){
			return url.Parse("http://" + proxyUrl)
		}
	}
	xClient := http.Client{
		Transport: &http.Transport{
			DisableKeepAlives:true,
			Proxy: proxyFunc,
		},
		Timeout: 60 * time.Second,
	}
	hReq,err := http.NewRequest("GET",pageUrl,nil)
	if err != nil{
		return ""
	}
	hReq.Header.Set("User-Agent","Mozilla / 5.0(Windows NT 10.0; Win64; x64) AppleWebKit / 537.36(KHTML, like Gecko) Chrome / 96.0.4664.45 Safari / 537.36")
	resp,err := xClient.Do(hReq)
	if err != nil{
		log.Println("access url failed",pageUrl,":",err)
		return ""
	}
	defer resp.Body.Close()
	respBytes,err := ioutil.ReadAll(resp.Body)
	if err != nil{
		return ""
	}
	return string(respBytes)
}

func GetData(task* Model.CollectionTask,pageUrl string)  {

	log.Println("start analyze url:",pageUrl)

	tmpScraperData := DataBase.GLocationIndex.GetLocationData(pageUrl)
	tmpScraperData.Category = task.Category
	if tmpScraperData.BusinessName != ""{
		if time.Now().Sub(tmpScraperData.CreateTime) < 7 * 24 * time.Hour{
			log.Println("read csv from cache:",pageUrl)
			CsvResult.Instance.WriteResult(&tmpScraperData)
			return
		}
	}
	tmpScraperData.State = task.State
	tmpScraperData.GoogleUrl = pageUrl
	html := getHtmlContent(pageUrl)
	if html == ""{
		return
	}
	html = strings.ReplaceAll(html,"\\n","A")
	html = strings.ReplaceAll(html,"[","-")
	html = strings.ReplaceAll(html,"]","-")
	html = strings.ReplaceAll(html,"\\",">")

	tmpMatchList := regex_RealCategory.FindStringSubmatch(html)
	if len(tmpMatchList) > 0{
		tmpScraperData.RealCategory = tmpMatchList[1]
	}else{
		tmpScraperData.RealCategory = task.Category
	}

	tmpMatchList = regex_BusinessName.FindStringSubmatch(html)
	if len(tmpMatchList) > 0{
		tmpScraperData.BusinessName = strings.ReplaceAll(tmpMatchList[1],">>u0026","")
		tmpScraperData.BusinessName = strings.ReplaceAll(tmpScraperData.BusinessName,">","")
	}else{
		tmpScraperData.BusinessName = "N/A"
	}

	//解析地址
	tmpMatchList = regex_Address.FindStringSubmatch(html)
	if len(tmpMatchList) > 0{
		tmpScraperData.Address = strings.ReplaceAll(tmpMatchList[1],">>u0026","")
	}

	//解析城市
	tmpMatchList = regex_City.FindStringSubmatch(html)
	if len(tmpMatchList) > 0{
		tmpScraperData.City = strings.ReplaceAll(tmpMatchList[1],">>u0026","")
	}

	parseAddress(&tmpScraperData)

	//解析经纬度
	tmpMatchList = regex_Landtitude.FindStringSubmatch(html)
	if len(tmpMatchList) > 0{
		tmpScraperData.Latitude = tmpMatchList[2]
		tmpScraperData.Longitude = tmpMatchList[3]
	}

	//解析网站
	tmpMatchList = regex_Website.FindStringSubmatch(html)
	if len(tmpMatchList) > 0{
		tmpScraperData.Website = DetectWebsite(tmpMatchList[1])
	}

	//解析电话
	tmpMatchList = regex_Phone.FindStringSubmatch(html)
	if len(tmpMatchList) > 0{
		tmpScraperData.Phone = tmpMatchList[1]
	}

	if tmpScraperData.Website != ""{
		emailList := EmailMiner.GetEmail(tmpScraperData.Website)
		if len(emailList) > 0{
			emailBytes,_ := json.Marshal(emailList)
			tmpScraperData.Email = string(emailBytes)
		}
	}

	DataBase.GLocationIndex.SetLocationData(pageUrl,&tmpScraperData)

	//写出结果
	log.Println("collect finish,write to csv:",pageUrl)
	CsvResult.Instance.WriteResult(&tmpScraperData)
}

func init()  {
	regex_RealCategory = regexp.MustCompile(strRegex_RealCategory)
	regex_BusinessName = regexp.MustCompile(strRegex_BusinessName)
	regex_Address = regexp.MustCompile(strRegex_Address)
	regex_City = regexp.MustCompile(strRegex_City)

	regex_PostalCode1 = regexp.MustCompile(strRegex_PostalCode1)
	regex_PostalCode2 = regexp.MustCompile(strRegex_PostalCode2)
	regex_PostalCode3 = regexp.MustCompile(strRegex_PostalCode3)
	regex_PostalCode4 = regexp.MustCompile(strRegex_PostalCode4)

	regex_Landtitude = regexp.MustCompile(strRegex_Landtitude)
	regex_Website = regexp.MustCompile(strRegex_Website)

	regex_Phone = regexp.MustCompile(strRegex_Phone)
}
