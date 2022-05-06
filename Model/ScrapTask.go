package Model

type ScrapTask struct {
	//任务
	TaskId int
	//关键字
	Category string
	//地址
	Location string
	//国家
	Country string
	//省份
	State string
	//城市
	City string
	//邮政编码
	ZipCode string
}