package FyneApp

import (
	"GoogleMapsCollector/ConfigManager"
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/widget"
	"strconv"
)

type SettingWindow struct {
	dialog.Dialog
	entry_PerZip *widget.Entry
	entry_SytemProxy *widget.Entry
}

func (this *SettingWindow)makeSettingWindow()fyne.CanvasObject  {

	totalLayOut := container.NewVBox()
	perZipLayOut := container.NewHBox()
	label_PerZip := widget.NewLabel("每条邮政编码提取地址数:")
	this.entry_PerZip = widget.NewEntry()
	zipCodeCount, _:= strconv.Atoi(ConfigManager.Instance.GetEmailPerZipCode())
	if zipCodeCount == 0{
		zipCodeCount = 30
	}
	this.entry_PerZip.Text = strconv.Itoa(zipCodeCount)
	perZipLayOut.Add(label_PerZip)
	perZipLayOut.Add(this.entry_PerZip)

	label_SytemProxy := widget.NewLabel("使用HTTP代理:")
	this.entry_SytemProxy = widget.NewEntry()
	this.entry_SytemProxy.Text = ConfigManager.Instance.GetSystemProxy()
	totalLayOut.Add(perZipLayOut)
	totalLayOut.Add(label_SytemProxy)
	totalLayOut.Add(this.entry_SytemProxy)
	return totalLayOut
}

func (this *SettingWindow)onConfirm(bConfirm bool)  {
	if bConfirm == false{
		return
	}
	zipCodeCount, _ := strconv.Atoi(this.entry_PerZip.Text)
	if zipCodeCount != 0{
		ConfigManager.Instance.SetEmailPerZipCode(this.entry_PerZip.Text)
	}
	ConfigManager.Instance.SetSystemProxy(this.entry_SytemProxy.Text)
	ConfigManager.Instance.Save()
	return
}

func NewSettingWindow(parent fyne.Window)(ret *SettingWindow)  {
	ret = &SettingWindow{}
	diag := dialog.NewCustomConfirm("设置","完成","取消",ret.makeSettingWindow(),ret.onConfirm,parent)
	ret.Dialog = diag
	return ret
}
