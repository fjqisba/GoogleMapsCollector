package TestUnit

import (
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"os"
	"strings"
	"testing"
)

func TestFixDB(t *testing.T) {

	hFile, _ := os.Open("D:\\list.txt")
	fileBytes,_ := ioutil.ReadAll(hFile)
	vec_list := strings.Split(string(fileBytes),"\r\n")
	log.Println(len(vec_list))

	for _,eCityName := range vec_list{

		resp,err := http.Get("https://www.nowmsg.com/findzip/be_postal_code.asp?CityName="+  url.QueryEscape(eCityName))
		if err != nil{
			log.Println("爬取出错:",eCityName)
			continue
		}
		respBytes,_ := ioutil.ReadAll(resp.Body)
		html := string(respBytes)
		resp.Body.Close()
		aIndex := strings.Index(html,"<tbody>")
		if aIndex == -1{
			log.Println("爬取出错",eCityName)
			continue
		}
		bIndex := strings.Index(html[aIndex:],"</tbody>")
		if bIndex == -1{
			log.Println("寻找尾部失败",eCityName)
			continue
		}
		keyHtml := html[aIndex+7:aIndex+bIndex]
		keyHtml = strings.ReplaceAll(keyHtml," ","")
		vec_keyList := strings.Split(keyHtml,"</td><td>")
		if len(vec_keyList) != 8{
			log.Println("爬取错误:",eCityName)
			continue
		}
		fmt.Println(eCityName+"," + vec_keyList[3])

	}

}