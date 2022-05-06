package DataBase

import (
	"GoogleMapsCollector/Utils/ProjectPath"
	"database/sql"
	_ "github.com/mattn/go-sqlite3"
	"log"
)

var(
	GLocationDB LocationDB
)

type LocationDB struct {
	Sqlx *sql.DB
}

func init()  {
	var err error
	GLocationDB.Sqlx,err = sql.Open("sqlite3",ProjectPath.GProjectBinPath +"\\db\\location.db")
	if err != nil{
		log.Panicln(err)
	}
}