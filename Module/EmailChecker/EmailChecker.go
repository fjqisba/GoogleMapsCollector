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
