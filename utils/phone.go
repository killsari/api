package utils

import (
	"fmt"
	"github.com/sirupsen/logrus"
	"regexp"
)

//发短信
func SendShortMessage(shortMessage string, tel string) (err error) {
	logrus.Info(fmt.Sprintf("send short message[%s] to tel[%s]", shortMessage, tel))
	return
}

// ValidMobilePhoneNum validates the mobile phone number in China.
//
// Return:
//     true for valid or false for invalid.
func ValidMobilePhoneNum(number string) bool {
	p := `^\d{11}$`
	r := regexp.MustCompile(p)

	if r.MatchString(number) {
		return true
	}
	return false
}
