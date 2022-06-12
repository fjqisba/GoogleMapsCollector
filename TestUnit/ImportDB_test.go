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

//key是城市名称,value是编码表
type CityData map[string][]string
//key是省份名称,value是城市数据
type RegionData map[string]CityData

const stmt_createCountry = `create table %s(
		"region"	TEXT NOT NULL,
		"city"		TEXT NOT NULL,
		"zip_codes"	TEXT NOT NULL,
PRIMARY KEY(region,city));`




func TestImportDB(t *testing.T) {

	//填充以下三个字段
	countryName := "UK"
	countryZHName := "英国"
	countryID := "700"

	hRegion,err := os.Open("D:\\工作台\\国外软件\\db\\region.csv")
	if err != nil{
		log.Panicln(err)
	}
	defer hRegion.Close()

	_,err = DataBase.GLocationDB.Sqlx.Exec(fmt.Sprintf(stmt_createCountry,countryName))
	if err != nil{
		if strings.Contains(err.Error(),"already exists") == false{
			log.Panicln(err)
		}
	}

	_,err = DataBase.GLocationDB.Sqlx.Exec("INSERT INTO country(country,countryName,tableName) VALUES(?,?,?)",countryName,countryZHName,countryName)

	//key是CityId,value是ZipCode
	zipCodeMap := make(map[string][]string)
	hZipCode,err := os.Open("D:\\工作台\\国外软件\\db\\zip_codes.csv")
	if err != nil{
		log.Panicln(err)
	}
	defer hZipCode.Close()
	csvZipCode := csv.NewReader(hZipCode)
	csvZipCode.Read()
	for true{
		vec_ZipRecords,_ := csvZipCode.Read()
		if vec_ZipRecords == nil{
			break
		}
		cityId := vec_ZipRecords[0]
		zipCode := vec_ZipRecords[1]
		zipCodeMap[cityId] = append(zipCodeMap[cityId], zipCode)
	}

	//生成区域表,key是省份名称,value是省份数据
	regionMap := make(RegionData)

	//key是省份ID,value是省份名称
	regionCodeMap := make(map[string]string)
	csvRegion := csv.NewReader(hRegion)
	csvRegion.Read()
	vec_RegionRecord,err := csvRegion.Read()
	for vec_RegionRecord != nil{
		if vec_RegionRecord[3] == countryID{
			regionCodeMap[vec_RegionRecord[0]] = vec_RegionRecord[1]
		}
		vec_RegionRecord,_ = csvRegion.Read()
	}

	hCity,err := os.Open("D:\\工作台\\国外软件\\db\\city.csv")
	if err != nil{
		log.Panicln(err)
	}
	defer hCity.Close()
	csvCity := csv.NewReader(hCity)
	csvCity.Read()
	for true{
		vec_CityRecords,_ := csvCity.Read()
		if vec_CityRecords == nil{
			break
		}
		tmpCountryID := vec_CityRecords[2]
		if tmpCountryID != countryID{
			continue
		}
		tmpRegionId := vec_CityRecords[1]
		RegioneName := regionCodeMap[tmpRegionId]
		zipCodes := zipCodeMap[vec_CityRecords[0]]
		if len(zipCodes) == 0{
			//不能为空
			continue
		}
		if regionMap[RegioneName] == nil{
			regionMap[RegioneName] = make(CityData)
		}
		regionMap[RegioneName][vec_CityRecords[3]] = zipCodes
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

	for eProvinceName,eCityMap := range regionMap{
		for eCityName,eZipCodeList := range eCityMap{
			zipCodeBytes, _ := json.Marshal(eZipCodeList)
			_,err = stat.Exec(eProvinceName,eCityName,string(zipCodeBytes))
			if err != nil{
				log.Panicln(err)
			}
		}
	}
	if err = tx.Commit(); err != nil {
		log.Panicln(err)
	}
}