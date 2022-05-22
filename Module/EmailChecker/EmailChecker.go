package EmailChecker

import (
	emailverifier "github.com/AfterShip/email-verifier"
)

var(
	emailChecker *emailverifier.Verifier
)

func init()  {
	emailChecker = emailverifier.NewVerifier()
	emailChecker.EnableSMTPCheck()
}

//返回false表示格式不合法
func IsValidEmail(emailAddr string)bool  {
	synt := emailChecker.ParseAddress(emailAddr)
	if synt.Valid == false{
		return false
	}
	//_,err := emailChecker.CheckMX(synt.Domain)
	//if err != nil{
	//	return false
	//}
	return true
}

//严格模式
func CheckEmail(emailAddr string)bool  {
	hResult,err := emailChecker.Verify(emailAddr)
	if err != nil{
		return false
	}
	if hResult.Reachable == "yes"{
		return true
	}
	return false
}
