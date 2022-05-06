package ProjectPath

import (
	"log"
	"os"
	"path/filepath"
)

var(
	//项目运行路径
	GProjectBinPath string
)

func init()  {
	var err error
	GProjectBinPath = filepath.Dir(os.Args[0])
	GProjectBinPath, err = filepath.Abs(GProjectBinPath)
	if err != nil {
		log.Fatalln(err)
	}
}