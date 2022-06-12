package TestUnit

import (
	"GoogleMapsCollector/DataBase"
	"encoding/csv"
	"encoding/json"
	"fmt"
	"log"
	"os"
	"strings"
	"testing"
)

func TestImportCountryTxt(t *testing.T) {

	//填充以下三个字段
	countryName := "NewZealand"
	countryZHName := "新西兰"

	hFile,err := os.Open("D:\\jp.txt")
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
	hCsv.Comma = []rune("\t")[0]
	type addrKey struct {
		StateName string
		CityName string
	}

	finalMap := make(map[addrKey][]string)

	for{
		vec_Records,err := 	hCsv.Read()
		if err != nil{
			break
		}

		tmpState := vec_Records[5]
		if tmpState == ""{
			tmpState = vec_Records[3]
		}
		if tmpState == ""{
			tmpState = countryName
		}
		tmpKey := addrKey{
			StateName: tmpState,
			CityName: vec_Records[2],
		}
		finalMap[tmpKey] = append(finalMap[tmpKey], vec_Records[1])
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