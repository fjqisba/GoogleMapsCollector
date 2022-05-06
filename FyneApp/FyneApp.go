package FyneApp

import (
	"GoogleMapsCollector/Utils/ProjectPath"
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
}

func NewFyneApp()*FyneApp  {
	return &FyneApp{
		app:app.New(),
	}
}

//国家被选择
func (this *FyneApp)OnCountrySelected(country string)  {
	if country == "全部国家"{
		this.select_State.Disable()
		this.select_City.Disabled()
		this.select_State.ClearSelected()
		this.select_City.ClearSelected()
		return
	}
	this.select_State.Enable()
}

//省份被选择
func (this *FyneApp)OnStateSelected(state string) {

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

func (this *FyneApp)InitializeComponent()  {


	//初始化窗口
	this.mainWindow = this.app.NewWindow("谷歌地图采集器")
	this.mainWindow.Resize(fyne.NewSize(800,600))
	this.mainWindow.CenterOnScreen()
	this.mainWindow.SetMaster()

	//设置主菜单
	this.mainWindow.SetMainMenu(this.makeMainMenu())

	//建立自定义布局
	customerLayout := container.NewWithoutLayout()

	//添加国家选择
	label_selectCountry := widget.NewLabel("选择国家:")
	label_selectCountry.Move(fyne.NewPos(20,25))
	this.select_Country = widget.NewSelect([]string{"全部国家","法国","英国"}, this.OnCountrySelected)
	this.select_Country .Move(fyne.NewPos(110,30))
	this.select_Country .Resize(fyne.NewSize(130,25))
	this.select_Country.PlaceHolder = "请选择"
	customerLayout.Add(label_selectCountry)
	customerLayout.Add(this.select_Country )

	//添加省份选择
	label_selectState := widget.NewLabel("选择省份:")
	label_selectState.Move(fyne.NewPos(20,65))
	this.select_State = widget.NewSelect([]string{"全部省份"},this.OnStateSelected)
	this.select_State.Move(fyne.NewPos(110,70))
	this.select_State.Resize(fyne.NewSize(130,25))
	this.select_State.Disable()
	this.select_State.PlaceHolder = "请选择"
	customerLayout.Add(label_selectState)
	customerLayout.Add(this.select_State)

	//添加城市选择
	label_selectCity := widget.NewLabel("选择城市:")
	label_selectCity.Move(fyne.NewPos(20,105))
	this.select_City = widget.NewSelect([]string{"全部城市"},this.OnCitySelected)
	this.select_City.Move(fyne.NewPos(110,110))
	this.select_City.Resize(fyne.NewSize(130,25))
	this.select_City.Disable()
	this.select_City.PlaceHolder = "请选择"
	customerLayout.Add(label_selectCity)
	customerLayout.Add(this.select_City)

	this.mainWindow.SetContent(customerLayout)
}

func (this *FyneApp)InitApp()error  {

	t := &myTheme{}
	t.SetFonts(ProjectPath.GProjectBinPath + "\\rsrc\\simsun.ttc","")
	this.app.Settings().SetTheme(t)

	this.InitializeComponent()

	return nil
}

func (this *FyneApp)Run()  {
	this.mainWindow.ShowAndRun()
}