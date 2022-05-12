package TaskManager

import (
	"GoogleMapsCollector/Model"
	"context"
	"log"
	"sync"
	"time"
)

type TaskManager struct {
	Ctx context.Context
	CancelFunc context.CancelFunc
	//等待结束
	wgFinish sync.WaitGroup
	//任务队列
	TaskList []Model.CollectionTask
	TaskFinish func()
}

func (this *TaskManager)Thread_ExecuteTask()  {

	defer func() {
		this.wgFinish.Wait()
		this.TaskFinish()
	}()

	for _,eTask := range this.TaskList{
		log.Println(eTask)
	}
	time.Sleep(10 * time.Second)
}

func NewTaskManager()*TaskManager  {
	return &TaskManager{}
}

