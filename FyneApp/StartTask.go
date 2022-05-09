package FyneApp

import (
	"errors"
	"fyne.io/fyne/v2/dialog"
	"log"
	"strings"
)

func (this *FyneApp)StartTask(vec_KeyWords []string)error  {

	this.button_StartTask.Text = "停止任务"

	//开始正式处理任务
	log.Println("开始采集任务")

	targetCountry := this.countryList[this.select_Country.SelectedIndex()]
	log.Println("采集国家:",targetCountry)

	//To do...生成爬虫任务所需要的数据
	

	//var targetCitys []string
	//if this.select_State.Selected == "全部省份"{
	//	stmt := fmt.Sprintf("SELECT distinct region FROM %s ",targetCountry)
	//	err := DataBase.GLocationDB.Sqlx.Select(&targetStates,stmt)
	//	if err != nil{
	//		return err
	//	}
	//}else{
	//	targetStates = append(targetStates, this.select_State.Selected)
	//}
	//log.Println(targetStates)


	return nil
}

func (this *FyneApp)StopTask()  {

	this.button_StartTask.Text = "开始任务"


}

func (this *FyneApp)TaskHandlerEntry() {
	if this.button_StartTask.Text == "开始任务"{
		//检查选择目标
		if this.select_Country.SelectedIndex() == -1{
			errWnd := dialog.NewError(errors.New("请先选择国家"),this.mainWindow)
			errWnd.SetDismissText("好的")
			errWnd.Show()
			return
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
			return
		}
		this.StartTask(vec_KeyWords)
	}else{
		this.StopTask()
	}
}