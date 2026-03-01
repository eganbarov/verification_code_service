package sender

import "fmt"

type CodeSender interface {
	SendCode(code string) error
}

type SmsSender struct{}

func (s *SmsSender) SendCode(code string) error {
	fmt.Println("Code is sent by sms gateway: " + code)
	return nil
}
