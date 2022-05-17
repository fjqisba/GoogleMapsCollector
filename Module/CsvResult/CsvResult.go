package CsvResult

import (
	"GoogleMapsCollector/Model"
	"GoogleMapsCollector/Utils/ProjectPath"
	"encoding/csv"
	"os"
	"sync"
)

var(
	Instance CsvResult
	title = []string{
		"关键字","行业","名称","地址","城市","省份","邮编","电话","邮箱","官网","谷歌链接"}
)

type CsvResult struct {
	hFile *os.File
	writer *csv.Writer
	mutex sync.Mutex
}

func (this *CsvResult)OpenCsv(filename string)error  {
	this.mutex.Lock()
	defer this.mutex.Unlock()
	if this.hFile != nil{
		this.writer.Flush()
		this.hFile.Close()
	}
	var err error
	this.hFile,err = os.Create(ProjectPath.GProjectBinPath + "\\csv\\" + filename)
	if err != nil{
		return err
	}
	this.writer = csv.NewWriter(this.hFile)
	this.writer.Write(title)
	return nil
}

func (this *CsvResult)CloseCsv() {
	this.mutex.Lock()
	defer this.mutex.Unlock()
	this.hFile.Close()
	this.hFile = nil
}

func (this *CsvResult)WriteResult(geoData *Model.ScraperData)  {

	this.mutex.Lock()
	defer this.mutex.Unlock()

	var writeData []string
	writeData = append(writeData, geoData.Category)
	writeData = append(writeData, geoData.RealCategory)
	writeData = append(writeData, geoData.BusinessName)
	writeData = append(writeData, geoData.Address)
	writeData = append(writeData, geoData.City)
	writeData = append(writeData, geoData.State)
	writeData = append(writeData, geoData.PostalCode)
	writeData = append(writeData, geoData.Phone)
	writeData = append(writeData, geoData.Email)
	writeData = append(writeData, geoData.Website)
	writeData = append(writeData, geoData.GoogleUrl)
	this.writer.Write(writeData)
	this.writer.Flush()
}

func init()  {
	os.Mkdir(ProjectPath.GProjectBinPath + "\\csv",0666)
}