package Model

import (
	"GoogleMapsCollector/Module/PageExtractor"
	"context"
	"fmt"
	"github.com/chromedp/chromedp"
	"log"
	"net/url"
	"time"
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

func (this *CollectionTask)buildSearchRequest1()string{
	str := url.QueryEscape(this.Category) + ","
	if this.ZipCode != ""{
		str = str + "+" + this.ZipCode
	}
	if this.City != ""{
		str = str + "+" + url.QueryEscape(this.City)
	}
	if this.State != ""{
		str = str + "+" + url.QueryEscape(this.State)
	}
	if this.Country != ""{
		str = str + "+" + url.QueryEscape(this.Country)
	}

	return fmt.Sprintf("https://www.google.com/maps/search/%s?force=tt",str)
}

func (this *CollectionTask)buildSearchRequest2()string {
	str := url.QueryEscape(this.Category) + ","
	if this.ZipCode != ""{
		str = str + "+" + this.ZipCode
	}
	if this.City != ""{
		str = str + "+" + url.QueryEscape(this.City)
	}
	if this.State != ""{
		str = str + "+" + url.QueryEscape(this.State)
	}
	return fmt.Sprintf("https://www.google.com/maps/search/%s?force=tt",str)
}

func (this *CollectionTask)buildSearchRequest3()string {
	str := url.QueryEscape(this.Category) + ","
	if this.ZipCode != ""{
		str = str + "+" + this.ZipCode
	}
	if this.City != ""{
		str = str + "+" + url.QueryEscape(this.City)
	}
	return fmt.Sprintf("https://www.google.com/maps/search/%s?force=tt",str)
}

func (this *CollectionTask)CollectCompanyListByChrome()string {
	requestUrl := this.buildSearchRequest1()
	opts := append(chromedp.DefaultExecAllocatorOptions[:],
		chromedp.Flag("headless",false),
	)
	allocCtx, _ := chromedp.NewExecAllocator(context.Background(), opts...)
	ctx, cancel := chromedp.NewContext(allocCtx)
	defer cancel()
	var outHtml string
	err := chromedp.Run(ctx,chromedp.Navigate(requestUrl),
		chromedp.Sleep(10 * time.Second),
		chromedp.OuterHTML("head",&outHtml))
	if err != nil{
		log.Println(err)
	}
	return ""
}

func (this *CollectionTask)CollectLocationIDList3()(ret []string) {
	return PageExtractor.ExtractPage(this.buildSearchRequest3())
}

func (this *CollectionTask)CollectLocationIDList2()(ret []string) {
	return PageExtractor.ExtractPage(this.buildSearchRequest2())
}

func (this *CollectionTask)CollectLocationIDList1()(ret []string) {
	return PageExtractor.ExtractPage(this.buildSearchRequest1())
}