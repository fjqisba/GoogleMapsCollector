package Utils

import "os"

func IsPathExists(path string)bool  {
	_,err := os.Stat(path)
	if err == nil{
		return true
	}
	return false
}

