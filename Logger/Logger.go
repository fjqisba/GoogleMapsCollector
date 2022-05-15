package Logger

import (
	"GoogleMapsCollector/Utils/ProjectPath"
	"github.com/sirupsen/logrus"
	"os"
)

var(
	ErrorLogger *logrus.Logger
	InfoLogger *logrus.Logger
)

func newLogger(loggerName string) (ret *logrus.Logger) {
	ret = logrus.New()
	flag := os.O_RDWR | os.O_CREATE | os.O_APPEND
	src, err := os.OpenFile(ProjectPath.GProjectBinPath + "\\logs\\"+loggerName+".txt", flag, 0666)
	if err != nil {
		panic("fatal error:can not open logs file")
	}
	ret.Out = src
	return ret
}

func init()  {

	os.Mkdir(ProjectPath.GProjectBinPath  + "\\logs", 0666)

	ErrorLogger = newLogger("error")
	InfoLogger = newLogger("info")
}