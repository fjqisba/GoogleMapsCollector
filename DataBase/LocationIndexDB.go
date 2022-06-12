package DataBase

import (
	"GoogleMapsCollector/Model"
	"GoogleMapsCollector/Utils/ProjectPath"
	"encoding/json"
	"github.com/syndtr/goleveldb/leveldb"
	"log"
	"time"
)

var(
	GLocationIndex LocationIndexDB
)

type LocationIndexDB struct {
	db *leveldb.DB
}

func (this *LocationIndexDB)SetLocationData(locationUrl string,geoData *Model.ScraperData)  {
	geoData.CreateTime = time.Now()
	geoBytes,err := json.Marshal(geoData)
	if err != nil{
		return
	}
	this.db.Put([]byte(locationUrl),geoBytes,nil)
}

func (this *LocationIndexDB)GetLocationData(locationUrl string)(retData Model.ScraperData)  {
	geoBytes,err := this.db.Get([]byte(locationUrl),nil)
	if err != nil{
		return retData
	}
	json.Unmarshal(geoBytes,&retData)
	return retData
}


func init()  {
	var err error
	GLocationIndex.db, err = leveldb.OpenFile(ProjectPath.GProjectBinPath + "\\db\\LocationIndex", nil)
	if err != nil {
		log.Panicln("fatal error:打开索引数据库失败")
	}
}