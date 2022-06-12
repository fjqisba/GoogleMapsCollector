package EmailMiner

import (
	"GoogleMapsCollector/ConfigManager"
	"GoogleMapsCollector/Logger"
	"GoogleMapsCollector/Module/EmailChecker"
	"github.com/weppos/publicsuffix-go/publicsuffix"
	"log"

	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"regexp"
	"strings"
	"sync"
	"time"
)

const strRegex_Email = `(mailto\\:|)([\\w\\.\\-]+)@((([\\-\\w]+\\.)+[a-zA-Z]{2,4})|(([0-9]{1,3}\\.){3}[0-9]{1,3}))`
const strRegex_Href = `href=(\"|'|)(.*?)(\"|'|)[>|\\s]`
const strRegex_GenEmail = `[a-zA-Z0-9_-]+@[a-zA-Z0-9_-]+(\.[a-zA-Z0-9_-]+)+`

type EmailMiner struct {
	wg sync.WaitGroup
	mutex sync.Mutex
	result []string
}

var(
	regex_GenEmail *regexp.Regexp
	regex_Email *regexp.Regexp
	regex_Href *regexp.Regexp
	contactPage = []string{
		"contacty", "kontakt", "conta",
	}
	pathList = []string{
		"contact","contact-us",
	}
	guessEmailList = []string{"info","support","buiro",
	"contacto","general","sales","purchase","order","commercial","office","gm"}
)

//返回true表示过滤
func filterEmail(emailAddr string)bool  {
	str := strings.ToLower(emailAddr)
	if strings.HasSuffix(str,".jpg") == true{
		return true
	}
	if strings.HasSuffix(str,".gif") == true{
		return true
	}
	if strings.Index(str,"@mail.com") != -1{
		return true
	}
	if strings.Index(str,"example") != -1{
		return true
	}
	if strings.Index(str,".wix") != -1{
		return true
	}
	if strings.Index(str,".png") != -1{
		return true
	}
	if strings.Index(str,"sentry.io") != -1{
		return true
	}
	return false
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
		Timeout: 30 * time.Second,
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

func (this *EmailMiner)insertEmail(email string)  {
	this.mutex.Lock()
	defer this.mutex.Unlock()
	this.result = append(this.result, email)
}

func (this *EmailMiner)guessEmail(domain string){
	for _,eGuess := range guessEmailList{
		this.wg.Add(1)
		go func(checkEmail string) {
			defer this.wg.Done()
			email := checkEmail + "@" + domain
			if EmailChecker.CheckEmail(email) == true{
				this.insertEmail(email)
			}
		}(eGuess)
	}
	return
}

func (this *EmailMiner)exoploreUrl(expUrl string)  {
	defer this.wg.Done()
	this.wg.Add(1)
	html := getHtmlContent(expUrl)
	vec_EmailList := regex_GenEmail.FindAllString(html,20)
	for _,eCheckEmail := range vec_EmailList{
		if filterEmail(eCheckEmail) == true{
			continue
		}
		if EmailChecker.IsValidEmail(eCheckEmail) == false{
			continue
		}
		this.insertEmail(eCheckEmail)
	}
}

func GetEmail(website string)[]string  {
	var emailMiner EmailMiner
	return emailMiner.DetectEmail(website)
}

//采集Email
func (this *EmailMiner)DetectEmail(website string)(retList []string)  {
	eUrl,err := url.Parse(website)
	if err != nil {
		Logger.ErrorLogger.Error("[GetEmail]:", website)
		return retList
	}
	host := eUrl.Host
	tIndex := strings.IndexByte(eUrl.Host,':')
	if tIndex != -1{
		host = eUrl.Host[0:tIndex]
	}
	strDomain,err := publicsuffix.Domain(host)
	if err == nil{
		this.guessEmail(strDomain)
	}

	go this.exoploreUrl(website)
	for _,ePath := range pathList{
		expUrl := fmt.Sprintf("%s://%s/%s",eUrl.Scheme,eUrl.Host,ePath)
		go this.exoploreUrl(expUrl)
	}
	this.wg.Wait()

	tmpMap := make(map[string]struct{})
	for _,eResult := range this.result{
		if _,bOk := tmpMap[eResult];bOk == false{
			tmpMap[eResult] = struct{}{}
			retList = append(retList, eResult)
		}
	}
	return retList
}

func init()  {
	regex_Email = regexp.MustCompile(strRegex_Email)
	regex_Href = regexp.MustCompile(strRegex_Href)
	regex_GenEmail = regexp.MustCompile(strRegex_GenEmail)
}
