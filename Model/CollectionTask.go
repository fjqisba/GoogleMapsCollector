package Model

import (
	"fmt"
	"net/url"
)

type TaskState int32
const(
	//无任务情况,按钮显示开始任务
	TASK_START TaskState = 0
	//正在执行任务情况,按钮显示停止任务
	TASK_EXECUTE TaskState = 1
	//用户点击停止任务,按钮显示正在停止任务中
	TASK_STOP TaskState = 2
)

type CollectionTask struct {
	//任务
	TaskId int
	//关键字
	Category string
	//国家
	Country string
	//省份
	State string
	//城市
	City string
	//邮政编码
	ZipCode string
}

func (this *CollectionTask)BuildSearchRequest()string{
	str := url.QueryEscape(this.Category) + ","
	if this.ZipCode != ""{
		str = str + "," + this.ZipCode
	}
	if this.City != ""{
		str = str + "," + url.QueryEscape(this.City)
	}
	if this.State != ""{
		str = str + "," + url.QueryEscape(this.State)
	}
	if this.Country != ""{
		str = str + "," + url.QueryEscape(this.Country)
	}

	return fmt.Sprintf("https://www.google.com/maps/search/%s?force=tt",str)
}