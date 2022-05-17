package PageExtractor

import (
	"GoogleMapsCollector/ConfigManager"
	"GoogleMapsCollector/Logger"
	"GoogleMapsCollector/Module/PhantomJS"
	"io/ioutil"
	"net/http"
	"regexp"
	"strconv"
	"strings"
	"time"
)

const reg_CompanyID string = `:0x(\w*?)\\\"]]],`

var(
	regex_MatchCompany *regexp.Regexp
)

type PageExtractor struct {

}

func getPage(url string)string  {
	xClient := http.Client{
		Transport: &http.Transport{
			DisableKeepAlives:true,
		},
		Timeout: 10 * time.Second,
	}
	hReq,err := http.NewRequest("GET",url,nil)
	if err != nil{
		return ""
	}
	hReq.Header.Set("User-Agent","Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/87.0.4280.88 Safari/537.36")
	resp,err := xClient.Do(hReq)
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

//返回不重复的页面列表

func ExtractPage(url string)(ret []string)  {
	html := PhantomJS.Instance.ScrapeGoogleMapHtml(url)
	index := strings.Index(html,"<<FirstPageEnd>>")
	if index == -1{
		return ret
	}
	result := regex_MatchCompany.FindAllStringSubmatch(html[0:index],-1)
	tmpMap := make(map[string]struct{})
	for _,eResult := range result{
		if len(eResult) != 2{
			Logger.ErrorLogger.Error("解析Html失败",eResult)
			continue
		}
		if _,bExists := tmpMap[eResult[1]];bExists == false{
			ret = append(ret, eResult[1])
			tmpMap[eResult[1]] = struct{}{}
		}
	}

	zipCodeCount, _ := strconv.Atoi(ConfigManager.Instance.GetEmailPerZipCode())
	if len(ret) >= zipCodeCount{
		return ret
	}
	//继续补充节点
	var vec_Url []string
	tmpArray := strings.Split(html[index:],"\r")
	for _,eTmp := range tmpArray{
		if strings.Index(eTmp,"http") > -1{
			vec_Url = append(vec_Url,strings.Replace(eTmp,"\n","",-1))
		}
	}
	for _,eUrl := range vec_Url{
		Logger.InfoLogger.Info("其它链接:",eUrl)
		//tmpHtml := getPage(eUrl)
	}
	return ret
}

func init()  {
	regex_MatchCompany = regexp.MustCompile(reg_CompanyID)
}
