package TestUnit

import (
	"GoogleMapsCollector/DataBase"
	"encoding/csv"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/PuerkitoBio/goquery"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"os"
	"strings"
	"testing"
)

func doRequest(uri string)string  {
	proxyFunc := func(req *http.Request) (*url.URL, error){
		return url.Parse("http://127.0.0.1:1080")
	}
	hClient := http.Client{
		Transport: &http.Transport{Proxy: proxyFunc},
	}
	hReq,err := http.NewRequest("GET",uri,nil)
	if err != nil{
		return ""
	}
	hReq.Header.Set("User-Agent","Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/102.0.5005.124 Safari/537.36 Edg/102.0.1245.44")
	resp,err := hClient.Do(hReq)
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


func getHtmlContent(uri string)string {
	for i:=0;i<5;i++{
		retHtml := doRequest(uri)
		if retHtml != ""{
			return retHtml
		}
	}
	return ""
}

type StateLink struct {
	StateName string
	LinkAddr string
}

func handleStateData(hCsv *csv.Writer,region string,stateInfo *StateLink)error  {

	stateContent := getHtmlContent(stateInfo.LinkAddr)
	if stateContent == ""{
		return errors.New("获取stateContent失败")
	}

	//说明数据已到达最底层
	if strings.Index(stateContent,"Locality") != -1 && strings.Index(stateContent,"Elevation") != -1{
		log.Println("写入地区:",region +","+ stateInfo.StateName)
		hCsv.Write([]string{region,stateInfo.StateName})
		return nil
	}

	doc,err := goquery.NewDocumentFromReader(strings.NewReader(stateContent))
	if err != nil{
		return err
	}

	sel := doc.Find("tbody")
	if len(sel.Nodes) != 1{
		log.Println("写入地区:",region + "," + stateInfo.StateName)
		hCsv.Write([]string{region,stateInfo.StateName})
		return nil
	}

	var stateList []StateLink
	selData := sel.Find("td")
	if len(selData.Nodes) == 0{
		log.Println("写入地区:",region + "," + stateInfo.StateName)
		hCsv.Write([]string{region,stateInfo.StateName})
		return nil
	}

	for i:=0;i<len(selData.Nodes);i=i+1{
		pFirstChild := selData.Nodes[i].FirstChild
		if pFirstChild == nil{
			log.Println("异常节点",stateInfo.LinkAddr,stateInfo.StateName)
			continue
		}
		if pFirstChild.Data == "a"{
			tmpStateName := pFirstChild.FirstChild.Data
			cIndex := strings.LastIndex(tmpStateName,"(")
			if cIndex != -1{
				tmpStateName = tmpStateName[0:cIndex]
			}
			tmpStateName = strings.TrimSpace(tmpStateName)
			stateList = append(stateList, StateLink{
				StateName:tmpStateName,
				LinkAddr: "https://www.azpostalcodes.com" + pFirstChild.Attr[0].Val,
			})
		}else{
			hCsv.Write([]string{stateInfo.StateName,pFirstChild.Data})
			log.Println("写入地区:",stateInfo.StateName + "," + pFirstChild.Data)
			continue
		}
	}

	for _,eState := range stateList{
		err = handleStateData(hCsv,stateInfo.StateName,&eState)
		if err != nil{
			log.Println("处理数据失败:",eState.StateName,err)
		}
		hCsv.Flush()
	}
	return nil
}

func getCountryData()  {

	countryZhName := "美国"
	countryName := "United States"
	countryZM := "us"

	hFile,err := os.Create("D:\\" + countryZhName + ".csv")
	if err != nil{
		log.Panicln(err)
	}
	defer hFile.Close()
	hCsvWriter := csv.NewWriter(hFile)

	mainPageContent := getHtmlContent("https://www.azpostalcodes.com/" + countryZM)
	doc,err := goquery.NewDocumentFromReader(strings.NewReader(mainPageContent))
	if err != nil{
		log.Panicln(err)
		return
	}
	sel := doc.Find("tbody")
	if len(sel.Nodes) != 1{
		return
	}

	var stateList []StateLink
	selData := sel.Find("a[href]")
	for i:=0;i<len(selData.Nodes);i=i+1{
		tmpStateName := selData.Nodes[i].FirstChild.Data
		cIndex := strings.LastIndex(tmpStateName,"(")
		if cIndex != -1{
			tmpStateName = tmpStateName[0:cIndex]
		}
		tmpStateName = strings.TrimSpace(tmpStateName)
		stateList = append(stateList, StateLink{
			StateName:tmpStateName,
			LinkAddr: "https://www.azpostalcodes.com" + selData.Nodes[i].Attr[0].Val,
		})
	}

	for _,eState := range stateList{
		err = handleStateData(hCsvWriter,countryName,&eState)
		if err != nil{
			log.Println("处理数据失败:",eState.StateName,err)
		}
		hCsvWriter.Flush()
	}
	hCsvWriter.Flush()
}

func TestImportCountryTxt(t *testing.T) {

	//getCountryData()
	//return

	//填充以下三个字段
	countryName := "USA_FAST"
	countryZHName := "美国极速版"

	hFile,err := os.Open("D:\\美国.csv")
	if err != nil{
		log.Panicln(err)
	}
	defer hFile.Close()

	_,err = DataBase.GLocationDB.Sqlx.Exec(fmt.Sprintf(stmt_createCountry,countryName))
	if err != nil{
		if strings.Contains(err.Error(),"already exists") == false{
			log.Panicln(err)
		}
	}

	_,err = DataBase.GLocationDB.Sqlx.Exec("INSERT INTO country(country,countryName,tableName) VALUES(?,?,?)",countryName,countryZHName,countryName)

	hCsv := csv.NewReader(hFile)
	//hCsv.Comma = []rune("\t")[0]
	type addrKey struct {
		StateName string
		CityName string
	}

	finalMap := make(map[addrKey][]string)

	var filterMap = make(map[string]struct{})
	for{
		vec_Records,err := 	hCsv.Read()
		if err != nil{
			break
		}
		tmpKey := addrKey{
			StateName: vec_Records[0],
			CityName: vec_Records[1],
		}
		zipCode := ""

		hash := tmpKey.StateName+tmpKey.CityName+zipCode
		_,bExists := filterMap[hash]
		if bExists == false{
			filterMap[hash] = struct{}{}
			finalMap[tmpKey] = append(finalMap[tmpKey], zipCode)
		}
	}

	//准备批量提交数据
	tx,err := DataBase.GLocationDB.Sqlx.Begin()
	if err != nil{
		log.Panicln(err)
	}
	//准备插入语句
	stat, err := tx.Prepare(fmt.Sprintf("insert into %s(region,city,zip_codes) values(?,?,?)",countryName))
	if err != nil {
		log.Panicln(err)
	}
	defer stat.Close()

	for eAddr,eZipCodeList := range finalMap{
		zipCodeBytes, _ := json.Marshal(eZipCodeList)
		_,err = stat.Exec(eAddr.StateName,eAddr.CityName,string(zipCodeBytes))
		if err != nil{
			log.Panicln(err)
		}
	}
	if err = tx.Commit(); err != nil {
		log.Panicln(err)
	}
	return

}