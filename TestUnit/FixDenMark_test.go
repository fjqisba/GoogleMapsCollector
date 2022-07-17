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
	"os"
	"strings"
	"testing"
)

type cityTask struct {
	CityName string
	StateName string
	ZipData []string
	Url string
}

func getPage(u string)string  {
	resp,err := http.Get(u)
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

func getCityData(task *cityTask)error  {

	htmlContent := getPage(task.Url)
	doc,err := goquery.NewDocumentFromReader(strings.NewReader(htmlContent))
	if err != nil{
		return err
	}

	sel := doc.Find("tbody")
	if len(sel.Nodes) != 1{
		return errors.New("no node")
	}
	selData := sel.Find("td")
	for i:=0;i<len(selData.Nodes);i=i+6{
		task.ZipData = append(task.ZipData, selData.Nodes[i].FirstChild.Data)
		task.StateName = selData.Nodes[i+2].FirstChild.Data
	}
	return nil
}

func TestImportDenMarkDB(t *testing.T) {

	//填充以下三个字段
	countryName := "Greece"
	countryZHName := "希腊"

	hFile,err := os.Open("D:\\希腊.txt")
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
	for{
		vec_Records,err := 	hCsv.Read()
		if err != nil{
			break
		}
		_,err = stat.Exec(vec_Records[0],vec_Records[1],vec_Records[2])
		if err != nil{
			continue
		}
	}
	if err = tx.Commit(); err != nil {
		log.Panicln(err)
	}
	return

}

func getIceLandData()error {

	htmlContent := getPage("https://www.17tr.com/is/410659.shtml")

	doc,err := goquery.NewDocumentFromReader(strings.NewReader(htmlContent))
	if err != nil{
		return err
	}

	sel := doc.Find("tbody")
	if len(sel.Nodes) != 1{
		return errors.New("no node")
	}
	selData := sel.Find("td")

	type addrKey struct {
		StateName string
		CityName string
	}
	finalMap := make(map[addrKey][]string)

	hFile,err := os.OpenFile("D:\\Ice.csv",os.O_RDWR|os.O_APPEND|os.O_CREATE, 0666)
	if err != nil{
		log.Panicln(err)
	}
	hCsvWriter := csv.NewWriter(hFile)
	defer hFile.Close()

	for i:=0;i<len(selData.Nodes);i=i+5{
		zipCode := selData.Nodes[i].FirstChild.Data
		cityName := selData.Nodes[i+1].FirstChild.Data
		tmpKey := addrKey{
			StateName: "Iceland",
			CityName: cityName,
		}
		finalMap[tmpKey] = append(finalMap[tmpKey], zipCode)
	}

	for eKey,eValue:= range finalMap{
		zipBytes, _ := json.Marshal(eValue)
		err = hCsvWriter.Write([]string{eKey.StateName,eKey.CityName,string(zipBytes)})
		if err != nil {
			log.Panicln(err)
		}
		hCsvWriter.Flush()
	}
	return nil
}

func TestFixDenMarkDB(t *testing.T) {

	getIceLandData()
	return

	stateContent := getPage("https://www.17tr.com/no/")
	stateDoc,err := goquery.NewDocumentFromReader(strings.NewReader(stateContent))
	if err != nil{
		log.Panicln(err)
	}

	hFile,err := os.OpenFile("D:\\NO.csv",os.O_RDWR|os.O_APPEND, 0666)
	if err != nil{
		log.Panicln(err)
	}
	hCsvWriter := csv.NewWriter(hFile)
	defer hFile.Close()

	selHref := stateDoc.Find(".link-item.link-image.no-desc")

	for iState:=1;iState<len(selHref.Nodes);iState++{
		stateUrl := selHref.Nodes[iState].FirstChild.Attr[0].Val
		htmlContent := getPage(stateUrl)
		doc,err := goquery.NewDocumentFromReader(strings.NewReader(htmlContent))
		if err != nil{
			log.Panicln(err)
		}
		sel := doc.Find("a[target]")
		var vec_cityTask []cityTask
		for _,eNode := range sel.Nodes{
			if len(eNode.Attr)!=2{
				log.Panicln("错误的数据")
			}
			if eNode.Attr[0].Val == "https://www.17tr.com/post/"{
				continue
			}
			if eNode.Attr[0].Val == "https://www.17tr.com/feedbacks/"{
				continue
			}
			if len(eNode.FirstChild.Attr) != 2{
				log.Panicln(err)
			}
			vec_cityTask = append(vec_cityTask, cityTask{CityName:eNode.FirstChild.Attr[1].Val,Url:eNode.Attr[1].Val})
		}
		log.Println("任务总数:",len(vec_cityTask))
		for i:=0;i<len(vec_cityTask);i++{
			err = getCityData(&vec_cityTask[i])
			if err != nil{
				log.Panicln(err)
			}
			zipBytes,_ := json.Marshal(vec_cityTask[i].ZipData)
			err = hCsvWriter.Write([]string{vec_cityTask[i].StateName,vec_cityTask[i].CityName,string(zipBytes)})
			if err != nil{
				log.Panicln(err)
			}
			log.Println("完成任务:",vec_cityTask[i])
			hCsvWriter.Flush()
		}
	}
}