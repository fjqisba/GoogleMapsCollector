package FyneApp

import (
	"GoogleMapsCollector/Model"
	"GoogleMapsCollector/TaskManager/TaskSignal"
	"encoding/json"
	"github.com/gin-gonic/gin"
	"io/ioutil"
	"net/http"
)

func (this *FyneApp)addWorkHandler(c *gin.Context) {
	var retJson map[string]interface{} = make(map[string]interface{})
	respBody, err := ioutil.ReadAll(c.Request.Body)
	if err != nil {
		retJson["code"] = 201
		retJson["msg"] = "参数错误"
		c.JSON(http.StatusOK, retJson)
		return
	}
	var workMsg Model.WorkParam
	err = json.Unmarshal(respBody,&workMsg)
	if err != nil{
		retJson["code"] = 202
		retJson["msg"] = "参数错误"
		c.JSON(http.StatusOK, retJson)
		return
	}
	if TaskSignal.GetTaskStatus() != Model.TASK_START{
		retJson["code"] = 203
		retJson["msg"] = "已有任务正在执行"
		c.JSON(http.StatusOK, retJson)
		return
	}
	err = this.WorkApi(&workMsg)
	if err != nil{
		retJson["code"] = 204
		retJson["msg"] = err.Error()
		c.JSON(http.StatusOK, retJson)
		return
	}
	retJson["code"] = 200
	retJson["msg"] = "添加任务成功"
	c.JSON(http.StatusOK, retJson)
}