package EmailMiner

import (
	"GoogleMapsCollector/Logger"
	"GoogleMapsCollector/Module/EmailChecker"
	"fmt"
	"github.com/weppos/publicsuffix-go/publicsuffix"
	"io/ioutil"
	"net/http"
	"net/url"
	"regexp"
	"strings"
	"sync"
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

func findCorrectEmail(items [][]string)string  {
	for _,array := range items{
		str := strings.ToLower(array[0])
		if strings.Index(str,"@mail.com") != -1{
			continue
		}
		if strings.Index(str,"example") != -1{
			continue
		}
		if strings.Index(str,".jpg") != -1{
			continue
		}
		if strings.Index(str,".wix") != -1{
			continue
		}
		if strings.Index(str,".png") != -1{
			continue
		}
		return array[0]
	}
	return ""
}

func detectUrl(url string,items [][]string)string  {
	for _,array := range items{
		for _,value := range contactPage{
			str := strings.ToLower(array[2])
			if strings.Index(str,value) == -1{
				continue
			}
			text := array[2]
			if strings.Index(text,"http") == -1{
				if text[0] == '/'{
					text = strings.TrimRight(url,"/") + text
				}else{
					tmpUrl := strings.TrimRight(url,"/")
					tmpUrl = strings.ReplaceAll(tmpUrl,"http://","")
					text = tmpUrl + "/" + text
				}
			}else{
				text2 := strings.ReplaceAll(url,"http://","")
				text2 = strings.ReplaceAll(text2,"https://","")
				text2 = strings.ReplaceAll(text2,"www","")
				vec_Split := strings.Split(text2,"/")
				if len(vec_Split) == 0{
					continue
				}
				text2 = vec_Split[0]
				if strings.Index(text,text2) == -1{
					continue
				}
			}
			if text == ""{
				continue
			}
			html := getPage(text)
			matchList := regex_Email.FindAllStringSubmatch(html,-1)
			if len(matchList) > 0{
				return findCorrectEmail(matchList)
			}
		}
	}
	return ""
}


func getPage(website string)string  {
	resp,err := http.Get(website)
	if err != nil{
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
	this.wg.Add(1)
	defer this.wg.Done()

	html := getPage(expUrl)
	matchList := regex_Email.FindAllStringSubmatch(html,20)
	if len(matchList) > 0{
		//To do...
	}
	vec_EmailList := regex_GenEmail.FindAllString(html,20)
	for _,eCheckEmail := range vec_EmailList{
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
