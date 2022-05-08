package DataBase

import (
	"GoogleMapsCollector/Utils/ProjectPath"
	"github.com/jmoiron/sqlx"
	_ "github.com/mattn/go-sqlite3"
	"log"
	"strings"
)

var(
	GLocationDB LocationDB
)

type LocationDB struct {
	Sqlx *sqlx.DB
}

const stmt_createCountryTable = `create table country(
		"country"		TEXT PRIMARY KEY NOT NULL,
		"countryName"	TEXT NOT NULL
);`

func init()  {
	var err error
	GLocationDB.Sqlx,err = sqlx.Open("sqlite3",ProjectPath.GProjectBinPath +"\\db\\location.db")
	if err != nil{
		log.Panicln(err)
	}

	_,err = GLocationDB.Sqlx.Exec(stmt_createCountryTable)
	if err != nil{
		if strings.Contains(err.Error(),"already exists") == false{
			log.Panicln(err)
		}
	}
}