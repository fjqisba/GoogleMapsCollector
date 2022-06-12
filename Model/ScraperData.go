package Model

import "time"

type ScraperData struct {

	//关键字
	Category string

	//谷歌地图链接
	GoogleUrl string	`json:"google_url"`

	State string		`json:"state"`

	//类型
	RealCategory string  `json:"real_category"`

	//公司
	BusinessName string	 `json:"business_name"`

	//地址
	Address string		`json:"address"`

	//城市
	City string			`json:"city"`

	PostalCode string	`json:"postal_code"`

	Latitude string		`json:"latitude"`
	Longitude string	`json:"longitude"`

	//网址
	Website string		`json:"website"`

	//电话
	Phone string		`json:"phone"`

	//邮箱
	Email string		`json:"email"`

	CreateTime time.Time  `json:"create_time"`
}