package FyneApp

import (
	"GoogleMapsCollector/DataBase"
	"GoogleMapsCollector/Model"
	"GoogleMapsCollector/TaskManager"
	"GoogleMapsCollector/TaskManager/TaskSignal"
	"encoding/json"
	"errors"
	"fmt"
	"fyne.io/fyne/v2/dialog"
	"log"
	"strings"
)

func (this *FyneApp)StartTask(vec_KeyWords []string)error  {

	this.button_StartTask.Text = "停止任务"
	this.button_StartTask.Refresh()

	//开始正式处理任务
	log.Println("开始生成任务")
	targetCountry := this.countryList[this.select_Country.SelectedIndex()]
	log.Println("采集国家:",targetCountry)

	type ZipCodeData struct {
		Region string		`db:"region"`
		City string			`db:"city"`
		ZipCodes string		`db:"zip_codes"`
	}

	//生成ZipCode临时数据
	var zipCodeList []ZipCodeData
	if this.select_State.Selected == "全部省份" {
		stmt := fmt.Sprintf("SELECT region,city,zip_codes FROM %s",targetCountry)
		err := DataBase.GLocationDB.Sqlx.Select(&zipCodeList,stmt)
		if err != nil{
			return err
		}
	}else{
		citySelectList,err := this.cityList.Get()
		if err != nil{
			return err
		}
		stmt := fmt.Sprintf("select region,city,zip_codes FROM %s where region=? and city=?",targetCountry)
		for _,eCitySelect := range citySelectList{
			bSelect,_ := eCitySelect.(*Model.CitySelectData).CitySwitch.Get()
			if bSelect == false {
				continue
			}
			cityName := eCitySelect.(*Model.CitySelectData).CityName
			rows,_ := DataBase.GLocationDB.Sqlx.Query(stmt,this.select_State.Selected,cityName)
			if rows == nil{
				continue
			}
			for rows.Next(){
				var tmpZipCodeData ZipCodeData
				err = rows.Scan(&tmpZipCodeData.Region,&tmpZipCodeData.City,&tmpZipCodeData.ZipCodes)
				if err != nil{
					continue
				}
				zipCodeList = append(zipCodeList, tmpZipCodeData)
			}
		}
	}

	//生成任务集合
	gTaskId := 1
	var CollectTaskList []Model.CollectionTask
	for _,eKeyWord := range vec_KeyWords{
		for _,eZipCodeData := range zipCodeList{
			var vec_ZipCode []string
			err := json.Unmarshal([]byte(eZipCodeData.ZipCodes),&vec_ZipCode)
			if err != nil{
				continue
			}
			for _,eZipCode := range vec_ZipCode{
				CollectTaskList = append(CollectTaskList, Model.CollectionTask{
					TaskId : gTaskId,
					Category : eKeyWord,
					Country:targetCountry,
					State:eZipCodeData.Region,
					City:eZipCodeData.City,
					ZipCode:eZipCode})
				gTaskId = gTaskId + 1
			}
		}
	}

	log.Println("生成任务完成,总数:",len(CollectTaskList))
	TaskManager.GTaskManager.TaskList = CollectTaskList
	go TaskManager.GTaskManager.Thread_ExecuteTask()
	return nil
}

func (this *FyneApp)onTaskFinished()  {
	log.Println("任务结束")
	//恢复按钮使用
	this.button_StartTask.Text = "开始任务"
	this.button_StartTask.Enable()
	this.button_StartTask.Refresh()
	TaskSignal.SetTaskStatus(Model.TASK_START)
}

func (this *FyneApp)StopTask()  {

	//先禁用按钮
	this.button_StartTask.Text = "等待任务结束......"
	this.button_StartTask.Disable()

	//传递停止信号
	TaskSignal.SetTaskStatus(Model.TASK_STOP)
}

//检查任务执行参数,返回false表示检查失败
func (this *FyneApp)preCheckTaskParam()([]string,bool)  {

	if this.select_Country.SelectedIndex() == -1{
		errWnd := dialog.NewError(errors.New("请先选择国家"),this.mainWindow)
		errWnd.SetDismissText("好的")
		errWnd.Show()
		return nil,false
	}

	//检查关键字
	var vec_KeyWords []string
	tmpKeyWords := strings.Split(this.entry_KeyWord.Text,"\n")
	for _,eKeyWord := range tmpKeyWords{
		if eKeyWord == ""{
			continue
		}
		vec_KeyWords = append(vec_KeyWords, eKeyWord)
	}
	if len(vec_KeyWords) == 0{
		errWnd := dialog.NewError(errors.New("请填入关键字"),this.mainWindow)
		errWnd.SetDismissText("好的")
		errWnd.Show()
		return nil,false
	}
	return vec_KeyWords,true
}

func (this *FyneApp)TaskHandlerEntry() {

	currentState := TaskSignal.GetTaskStatus()


	//开始任务
	if currentState == Model.TASK_START{
		//检查选择目标
		keyWordList,bCheckResult := this.preCheckTaskParam()
		if bCheckResult == false{
			return
		}
		TaskSignal.SetTaskStatus(Model.TASK_EXECUTE)
		this.StartTask(keyWordList)
	}

	//结束任务
	if currentState == Model.TASK_EXECUTE {
		TaskSignal.SetTaskStatus(Model.TASK_STOP)
		go this.StopTask()
	}
}