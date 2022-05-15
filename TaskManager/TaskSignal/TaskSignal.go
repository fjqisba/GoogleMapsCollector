package TaskSignal

import (
	"GoogleMapsCollector/Model"
	"sync/atomic"
)

var(
	taskStatus atomic.Value
)

func GetTaskStatus()Model.TaskState  {
	return taskStatus.Load().(Model.TaskState)
}

func SetTaskStatus(state Model.TaskState)  {
	taskStatus.Store(state)
}