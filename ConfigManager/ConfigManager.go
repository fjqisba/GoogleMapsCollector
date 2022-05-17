package ConfigManager

import (
	"GoogleMapsCollector/Utils"
	"GoogleMapsCollector/Utils/ProjectPath"
	"gopkg.in/ini.v1"
	"log"
	"os"
	"strconv"
)

var(
	Instance ConfigManager
)

type mainConfig struct {
	//每个邮编提取多少邮件
	EmailPerZipCode int
	//系统代理
	SystemProxy string
}

type ConfigManager struct {
	ini *ini.File
	FilePath string
	mainConfig mainConfig
}


func (this *ConfigManager)GetSystemProxy()string  {
	return this.ini.Section("main").Key("SystemProxy").Value()
}

func (this *ConfigManager)SetSystemProxy(proxy string)  {
	this.ini.Section("main").Key("SystemProxy").SetValue(proxy)
}

func (this *ConfigManager)GetEmailPerZipCode()string  {
	return this.ini.Section("main").Key("EmailPerZipCode").Value()
}

func (this *ConfigManager)SetEmailPerZipCode(value string)  {
	mainSection := this.ini.Section("main")
	if mainSection == nil{
		return
	}
	hKey := mainSection.Key("EmailPerZipCode")
	if hKey == nil{
		return
	}
	hKey.SetValue(value)
}

func (this *ConfigManager)Save()  {
	this.ini.SaveTo(this.FilePath)
}

func (this *ConfigManager)initConfigManager(settingPath string)error  {
	var err error
	this.ini,err = ini.Load(settingPath)
	if err != nil{
		return err
	}
	this.FilePath = settingPath
	this.mainConfig.EmailPerZipCode, _ = strconv.Atoi(this.GetEmailPerZipCode())
	return nil
}

func init()  {
	settingPath := ProjectPath.GProjectBinPath  + "\\setting.ini"
	if Utils.IsPathExists(settingPath) == false{
		hFile,_ := os.Create(settingPath)
		if hFile != nil{
			hFile.Close()
		}
	}
	err := Instance.initConfigManager(settingPath)
	if err != nil{
		log.Panicln(err)
	}
}