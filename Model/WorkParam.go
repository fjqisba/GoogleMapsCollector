package Model


type WorkParam struct {
	//关键字
	Category []string	`json:"category"`
	//选择的国家
	CountryName string `json:"country"`
	//选择的省份
	StateName string	`json:"state"`
	//选择的城市
	CityList []string	`json:"city"`
}