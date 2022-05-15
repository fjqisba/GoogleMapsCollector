package FyneApp

import (
	"GoogleMapsCollector/DataBase"
	"GoogleMapsCollector/Model"
	"GoogleMapsCollector/TaskManager"
	"GoogleMapsCollector/TaskManager/TaskSignal"
	"GoogleMapsCollector/Utils/ProjectPath"
	"errors"
	"fmt"
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/data/binding"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/widget"
)



type FyneApp struct {
	app fyne.App
	mainWindow fyne.Window
	//界面信息
	select_Country *widget.Select
	select_State *widget.Select
	select_City *widget.List
	button_SelectAllCity *widget.Button
	button_StartTask *widget.Button
	entry_KeyWord *widget.Entry
	//国家列表,英文
	countryList []string
	//国家列表,中文
	countryNameList []string
	//城市列表,[]*CitySelectData
	cityList binding.UntypedList
}

func NewFyneApp()*FyneApp  {
	return &FyneApp{
		app:app.New(),
	}
}

func (this *FyneApp)makeMainMenu()*fyne.MainMenu  {

	quitAction := fyne.NewMenuItem("退出", func() {
		this.app.Quit()
	})
	quitAction.IsQuit = true

	setAction := fyne.NewMenuItem("设置", func() {
		setWnd := NewSettingWindow(this.mainWindow)
		setWnd.Show()
	})

	menu_Program:= fyne.NewMenu("程序",setAction, quitAction)
	menu_about:= fyne.NewMenu("关于", fyne.NewMenuItem("关于采集器", func() {
		info := dialog.NewInformation("关于谷歌地图采集器","企业内部定制版,禁止外部分享",this.mainWindow)
		info.SetDismissText("好的")
		info.Show()
	}))

	mainMenu := fyne.NewMainMenu(menu_Program,menu_about)
	return mainMenu
}


func (this *FyneApp)onConfirmClose(bConfirm bool) {
	if bConfirm == true{
		this.mainWindow.Close()
	}
}

func (this *FyneApp)onCloseWindow()  {

	if TaskSignal.GetTaskStatus() == Model.TASK_START{
		this.mainWindow.Close()
		return
	}

	//需要等待任务结束
	errHnd := dialog.NewConfirm("提示","请先停止正在运行中的任务!",this.onConfirmClose,this.mainWindow)
	errHnd.SetConfirmText("强制退出")
	errHnd.SetDismissText("取消")
	errHnd.Show()
}

func (this *FyneApp)InitializeComponent()error  {

	//初始化窗口
	this.mainWindow = this.app.NewWindow("谷歌地图采集器")
	this.mainWindow.Resize(fyne.NewSize(800,600))
	this.mainWindow.CenterOnScreen()
	this.mainWindow.SetMaster()

	this.mainWindow.SetCloseIntercept(this.onCloseWindow)
	//设置主菜单
	this.mainWindow.SetMainMenu(this.makeMainMenu())

	//建立自定义布局
	customerLayout := container.NewWithoutLayout()

	var tmpCountryMapping []Model.CountryNameMapping
	//添加国家选择
	err := DataBase.GLocationDB.Sqlx.Select(&tmpCountryMapping,"SELECT country,countryName from country")
	if err != nil{
		return err
	}
	if len(tmpCountryMapping) == 0{
		return errors.New("no country")
	}
	for _,eCountryMapping := range tmpCountryMapping{
		this.countryList = append(this.countryList, eCountryMapping.Country)
		this.countryNameList = append(this.countryNameList, eCountryMapping.CountryName)
	}
	label_selectCountry := widget.NewLabel("选择国家:")
	this.select_Country = widget.NewSelect(this.countryNameList, this.OnCountrySelected)
	this.select_Country.PlaceHolder = "请选择"
	label_selectCountry.Move(fyne.NewPos(20,25))
	this.select_Country.Move(fyne.NewPos(110,30))
	this.select_Country.Resize(fyne.NewSize(130,25))

	customerLayout.Add(label_selectCountry)
	customerLayout.Add(this.select_Country)

	//添加省份选择
	label_selectState := widget.NewLabel("选择省份:")
	label_selectState.Move(fyne.NewPos(20,65))
	this.select_State = widget.NewSelect([]string{"全部省份"},this.OnStateSelected)
	this.select_State.PlaceHolder = "请选择"
	this.select_State.Move(fyne.NewPos(110,70))
	this.select_State.Resize(fyne.NewSize(130,25))
	this.select_State.Disable()

	customerLayout.Add(label_selectState)
	customerLayout.Add(this.select_State)

	//添加城市选择
	label_selectCity := widget.NewLabel("选择城市:")
	label_selectCity.Move(fyne.NewPos(270,25))
	this.button_SelectAllCity = widget.NewButton("全选",this.OnSelectAllCity)
	this.button_SelectAllCity.Move(fyne.NewPos(370,32))
	this.button_SelectAllCity.Resize(fyne.NewSize(100,20))

	this.cityList = binding.NewUntypedList()
	this.select_City = widget.NewListWithData(this.cityList, func() fyne.CanvasObject {
		return widget.NewCheck("", nil)
	} , func(item binding.DataItem, object fyne.CanvasObject) {
		untypeData, _ := item.(binding.Untyped).Get()
		object.(*widget.Check).Bind(untypeData.(*Model.CitySelectData).CitySwitch)
		object.(*widget.Check).Text = untypeData.(*Model.CitySelectData).CityName
		object.(*widget.Check).Refresh()
	})
	if this.select_City == nil{
		return errors.New("no cityList")
	}
	this.select_City.Move(fyne.NewPos(350,70))
	this.select_City.Resize(fyne.NewSize(400,300))
	customerLayout.Add(label_selectCity)
	customerLayout.Add(this.select_City)
	customerLayout.Add(this.button_SelectAllCity)

	//添加关键字搜索
	label_keyword := widget.NewLabel("采集关键字(换行符分割多个):")
	label_keyword.Move(fyne.NewPos(20,200))
	this.entry_KeyWord = widget.NewMultiLineEntry()
	this.entry_KeyWord.Move(fyne.NewPos(28,230))
	this.entry_KeyWord.Resize(fyne.NewSize(200,120))
	customerLayout.Add(this.entry_KeyWord)
	customerLayout.Add(label_keyword)

	//添加任务按钮
	this.button_StartTask = widget.NewButton("开始任务", this.TaskHandlerEntry)
	this.button_StartTask.Move(fyne.NewPos(270,430))
	this.button_StartTask.Resize(fyne.NewSize(200,120))
	customerLayout.Add(this.button_StartTask)

	this.mainWindow.SetContent(customerLayout)
	return nil
}

//国家被选择
func (this *FyneApp)OnCountrySelected(country string)  {

	var selectCountry string
	for index,eCountryName := range this.countryNameList{
		if country == eCountryName{
			selectCountry = this.countryList[index]
			break
		}
	}

	var stateList []string
	stmt := fmt.Sprintf("SELECT distinct region FROM %s ",selectCountry)
	err := DataBase.GLocationDB.Sqlx.Select(&stateList,stmt)
	if err != nil{
		return
	}

	this.select_State.Options = []string{"全部省份"}
	this.select_State.Options = append(this.select_State.Options, stateList...)
	this.select_State.Enable()
	this.select_State.SetSelectedIndex(0)

	this.cityList.Set([]interface{}{})
}

//省份被选择
func (this *FyneApp)OnStateSelected(state string) {

	if state == "全部省份"{
		this.cityList.Set([]interface{}{})
		return
	}

	selectCountry := this.countryList[this.select_Country.SelectedIndex()]
	var cityList []string
	stmt := fmt.Sprintf("SELECT distinct city FROM %s where region=?",selectCountry)
	err := DataBase.GLocationDB.Sqlx.Select(&cityList,stmt,state)
	if err != nil{
		return
	}

	var cityBoundList []interface{}
	for _,eCityName := range cityList{
		cityBoundList = append(cityBoundList, &Model.CitySelectData{
			CityName:eCityName,
			CitySwitch:binding.NewBool(),
		})
	}

	this.cityList.Set(cityBoundList)
}

//选择所有的城市
func (this *FyneApp)OnSelectAllCity()  {

	if this.button_SelectAllCity.Text == "全选"{
		this.button_SelectAllCity.Text = "清空"
		vec_AllCity, _  := this.cityList.Get()
		for _,eCityList := range vec_AllCity{
			bTrue := binding.NewBool()
			bTrue.Set(true)
			eCityList.(*Model.CitySelectData).CitySwitch = bTrue
		}
	}else{
		this.button_SelectAllCity.Text = "全选"
		vec_AllCity, _  := this.cityList.Get()
		for _,eCityList := range vec_AllCity{
			eCityList.(*Model.CitySelectData).CitySwitch = binding.NewBool()
		}
	}
	this.button_SelectAllCity.Refresh()
	this.select_City.Refresh()
	return
}


func (this *FyneApp)InitApp()error  {

	//设置程序规模
	//os.Setenv("FYNE_SCALE", "0.9")
	t := &myTheme{}
	t.SetFonts(ProjectPath.GProjectBinPath + "\\rsrc\\simsun.ttc","")
	this.app.Settings().SetTheme(t)
	err := this.InitializeComponent()
	if err != nil{
		return err
	}

	TaskManager.GTaskManager.TaskFinishCB = this.onTaskFinished
	TaskSignal.SetTaskStatus(Model.TASK_START)
	return nil
}

func (this *FyneApp)Run()  {
	this.mainWindow.ShowAndRun()
}