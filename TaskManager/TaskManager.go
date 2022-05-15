package TaskManager

import (
	"GoogleMapsCollector/ConfigManager"
	"GoogleMapsCollector/Model"
	"GoogleMapsCollector/Module/GooglePageScraper"
	"GoogleMapsCollector/TaskManager/TaskSignal"
	"fmt"
	"log"
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
	taskCount := ConfigManager.GConfigManager.MainConfig.EmailPerZipCode
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

	ch := make(chan struct{},1)
	go func(tmpTask *Model.CollectionTask) {
		GooglePageScraper.GetData(tmpTask,url)
		ch <- struct{}{}
	}(task)

	select {
	case <- ch:
		return
	case <- time.After(30 * time.Second):
		return
	}
}

func (this *TaskManager)handleTask()  {

	for _,eZipCodeTask := range this.TaskList{
		idMap := this.startCollectTask(&eZipCodeTask)
		log.Println("生成地址(条):",len(idMap))
		if TaskSignal.GetTaskStatus() == Model.TASK_STOP{
			log.Println("用户点击停止任务,返回")
			return
		}
		if len(idMap) == 0{
			continue
		}
		//开始批量提取任务
		for tmpIdKey,_ := range idMap{
			this.wgFinish.Add(1)
			go this.startScrapeTask(&eZipCodeTask,fmt.Sprintf("https://maps.google.com/?cid=0x0:0x%s",tmpIdKey))
			if TaskSignal.GetTaskStatus() == Model.TASK_STOP{
				log.Println("用户点击停止任务,直接返回")
				return
			}
		}
		this.wgFinish.Wait()
	}
}

func (this *TaskManager)Thread_ExecuteTask()  {
	defer func() {
		this.wgFinish.Wait()
		this.TaskFinishCB()
	}()
	this.handleTask()
}
