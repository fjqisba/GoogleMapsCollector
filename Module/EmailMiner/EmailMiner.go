package EmailMiner

import (
	"io/ioutil"
	"net/http"
	"regexp"
	"strings"
)

const strRegex_Email = `(mailto\\:|)([\\w\\.\\-]+)@((([\\-\\w]+\\.)+[a-zA-Z]{2,4})|(([0-9]{1,3}\\.){3}[0-9]{1,3}))`
const strRegex_Href = `href=(\"|'|)(.*?)(\"|'|)[>|\\s]`

var(
	regex_Email *regexp.Regexp
	regex_Href *regexp.Regexp

	contactPage = []string{
		"contacty", "kontakt", "conta",
	}
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

func GetEmail(website string)string  {
	html := getPage(website)
	matchList := regex_Email.FindAllStringSubmatch(html,-1)
	if len(matchList) > 0{
		return findCorrectEmail(matchList)
	}
	//else{
	//	matchList = regex_Href.FindAllStringSubmatch(html,-1)
	//	return detectUrl(website,matchList)
	//}
	return ""
}

func init()  {
	regex_Email = regexp.MustCompile(strRegex_Email)
	regex_Href = regexp.MustCompile(strRegex_Href)
}
