package TaskManager

import (
	"GoogleMapsCollector/ConfigManager"
	"GoogleMapsCollector/Model"
	"GoogleMapsCollector/Module/CsvResult"
	"GoogleMapsCollector/Module/GooglePageScraper"
	"GoogleMapsCollector/TaskManager/TaskSignal"
	"fmt"
	"log"
	"strconv"
	"sync"
	"time"
)

var(
	GTaskManager TaskManager
)

type TaskManager struct {
	//等待结束
	wgFinish sync.WaitGroup
	//任务队列
	TaskList []Model.CollectionTask
	//任务完成后的界面处理函数
	TaskFinishCB func()
}

//传入任务,采集指定的页面地址,返回数据列表
func (this *TaskManager)startCollectTask(task *Model.CollectionTask)(ret []string)  {

	//要完成的目标个数
	taskCount, _ := strconv.Atoi(ConfigManager.Instance.GetEmailPerZipCode())
	//用于去重
	insertMap := make(map[string]struct{})

	//执行第一次采集任务
	tmpTaskList := task.CollectLocationIDList1()
	for _,eTask := range tmpTaskList{
		if _,bExists := insertMap[eTask];bExists == false{
			insertMap[eTask] = struct{}{}
			ret = append(ret, eTask)
		}
	}
	if len(ret) >= taskCount{
		return ret[0:taskCount]
	}
	if TaskSignal.GetTaskStatus() == Model.TASK_STOP{
		return ret
	}

	log.Println("采集地址数量小于预期,扩大搜索范围:",len(ret))

	//执行第二次采集任务
	tmpTaskList = task.CollectLocationIDList2()
	for _,eTask := range tmpTaskList{
		if _,bExists := insertMap[eTask];bExists == false{
			insertMap[eTask] = struct{}{}
			ret = append(ret, eTask)
		}
	}
	if len(ret) >= taskCount{
		return ret[0:taskCount]
	}
	if TaskSignal.GetTaskStatus() == Model.TASK_STOP{
		return ret
	}

	log.Println("采集地址数量小于预期,继续扩大搜索范围:",len(ret))
	//执行第三次采集任务
	tmpTaskList = task.CollectLocationIDList3()
	for _,eTask := range tmpTaskList{
		if _,bExists := insertMap[eTask];bExists == false{
			insertMap[eTask] = struct{}{}
			ret = append(ret, eTask)
		}
	}
	if len(ret) >= taskCount{
		return ret[0:taskCount]
	}
	return ret
}

func (this *TaskManager)startScrapeTask(task *Model.CollectionTask,url string) {
	defer func() {
		this.wgFinish.Done()
	}()
	GooglePageScraper.GetData(task,url)
	return
	//defer func() {
	//	this.wgFinish.Done()
	//}()
	//
	//go func(tmpTask *Model.CollectionTask) {
	//
	//	ch <- struct{}{}
	//}(task)
	//
	//select {
	//case <- ch:
	//	return
	//case <- time.After(60 * time.Second):
	//	return
	//}
}

//返回true表示任务停止,返回false表示任务继续
func (this *TaskManager)handleZipCodeTask(eZipCodeTask *Model.CollectionTask)bool  {

	defer func() {
		this.wgFinish.Wait()
	}()

	fileName := eZipCodeTask.Country + "_" + time.Now().Format("20060102150405") + ".csv"
	err := CsvResult.Instance.OpenCsv(fileName)
	if err != nil{
		log.Println("创建csv文件失败")
		return false
	}
	log.Println("创建csv任务:",fileName)
	defer CsvResult.Instance.CloseCsv()

	idList := this.startCollectTask(eZipCodeTask)
	log.Println("生成地址(条):",len(idList))
	if TaskSignal.GetTaskStatus() == Model.TASK_STOP{
		log.Println("用户点击停止任务,返回")
		return true
	}
	if len(idList) == 0{
		return false
	}
	//开始批量提取任务
	threadCount := 0
	for _,tmpIdKey := range idList{
		this.wgFinish.Add(1)
		threadCount = threadCount + 1
		go this.startScrapeTask(eZipCodeTask,fmt.Sprintf("https://maps.google.com/?cid=0x0:0x%s",tmpIdKey))
		if threadCount >= 5{
			this.wgFinish.Wait()
			threadCount = 0
		}
		if TaskSignal.GetTaskStatus() == Model.TASK_STOP{
			log.Println("用户点击停止任务")
			this.wgFinish.Wait()
			return true
		}
	}
	this.wgFinish.Wait()
	return false
}

func (this *TaskManager)handleTask()  {
	for _,eZipCodeTask := range this.TaskList{
		if this.handleZipCodeTask(&eZipCodeTask) == true{
			return
		}
	}
}

func (this *TaskManager)Thread_ExecuteTask()  {
	defer func() {
		this.wgFinish.Wait()
		this.TaskFinishCB()
	}()
	this.handleTask()
}
