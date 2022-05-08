package FyneApp

import (
	"GoogleMapsCollector/DataBase"
	"GoogleMapsCollector/Model"
	"GoogleMapsCollector/Utils/ProjectPath"
	"errors"
	"fmt"
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/widget"
)

type FyneApp struct {
	app fyne.App
	mainWindow fyne.Window

	select_Country *widget.Select
	select_State *widget.Select
	select_City *widget.Select

	countryList []string
	countryNameList []string
}

func NewFyneApp()*FyneApp  {
	return &FyneApp{
		app:app.New(),
	}
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
}

//省份被选择
func (this *FyneApp)OnStateSelected(state string) {

	if state == "全部省份"{
		this.select_City.Disable()
		return
	}

	selectCountry := this.countryList[this.select_Country.SelectedIndex()]
	var cityList []string
	stmt := fmt.Sprintf("SELECT distinct city FROM %s where region=?",selectCountry)
	err := DataBase.GLocationDB.Sqlx.Select(&cityList,stmt,state)
	if err != nil{
		return
	}
	this.select_City.Enable()
	this.select_City.Options = []string{"全部城市"}
	this.select_City.Options = append(this.select_City.Options, cityList...)
	this.select_City.SetSelectedIndex(0)
}

//城市被选择
func (this *FyneApp)OnCitySelected(city string) {

}

func (this *FyneApp)makeMainMenu()*fyne.MainMenu  {

	quitAction := fyne.NewMenuItem("退出", func() {
		this.app.Quit()
	})
	quitAction.IsQuit = true

	menu_Program:= fyne.NewMenu("程序", quitAction)
	menu_about:= fyne.NewMenu("关于", fyne.NewMenuItem("关于采集器", func() {
		info := dialog.NewInformation("关于谷歌地图采集器","企业内部定制版,禁止外部分享",this.mainWindow)
		info.SetDismissText("好的")
		info.Show()
	}))

	mainMenu := fyne.NewMainMenu(menu_Program,menu_about)

	return mainMenu
}

func (this *FyneApp)InitializeComponent()error  {


	//初始化窗口
	this.mainWindow = this.app.NewWindow("谷歌地图采集器")
	this.mainWindow.Resize(fyne.NewSize(800,600))
	this.mainWindow.CenterOnScreen()
	this.mainWindow.SetMaster()

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
	customerLayout.Add(this.select_Country )

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
	label_selectCity.Move(fyne.NewPos(20,105))
	this.select_City = widget.NewSelect([]string{"全部城市"},this.OnCitySelected)
	this.select_City.PlaceHolder = "请选择"
	this.select_City.Move(fyne.NewPos(110,110))
	this.select_City.Resize(fyne.NewSize(130,25))
	this.select_City.Disable()

	customerLayout.Add(label_selectCity)
	customerLayout.Add(this.select_City)




	this.mainWindow.SetContent(customerLayout)

	return nil
}

func (this *FyneApp)InitApp()error  {

	t := &myTheme{}
	t.SetFonts(ProjectPath.GProjectBinPath + "\\rsrc\\simsun.ttc","")
	this.app.Settings().SetTheme(t)

	err := this.InitializeComponent()
	if err != nil{
		return err
	}

	return nil
}

func (this *FyneApp)Run()  {
	this.mainWindow.ShowAndRun()
}