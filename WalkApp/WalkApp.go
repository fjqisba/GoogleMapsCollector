package WalkApp

import (
	"github.com/lxn/walk"
	. "github.com/lxn/walk/declarative"
	"log"
)

type WalkApp struct {
	privteHwnd *walk.MainWindow
	hwndMain *MainWindow
}

func NewWalkApp()*WalkApp  {
	return &WalkApp{}
}

func (this *WalkApp)Run() {
	_,err := this.hwndMain.Run()
	if err != nil{
		log.Panicln(err)
	}
}


func (this *WalkApp)InitializeComponent()error {
	label_Country := TextLabel{Text: "hello world"}
	label_Country.Row = 1
	label_Country.Column = 1
	this.hwndMain.Children = append(this.hwndMain.Children, label_Country)
	return nil
}

func (this *WalkApp)InitWalkApp()error  {

	this.hwndMain = &MainWindow{
		AssignTo: &this.privteHwnd,
		Title:    "谷歌地图采集器",
		MinSize:  Size{240, 320},
		Size:     Size{300, 400},
		Layout:   Grid{Columns: 2,MarginsZero: true},
	}

	err := this.InitializeComponent()
	if err != nil{
		return err
	}
	return nil
}