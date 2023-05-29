package mail

import (
	"fmt"
	"time"

	"github.com/wayne011872/goSterna/mail"
)

func SendMail(mailTitle,mailContent string) {
	myMail := &mail.Mail{}
	myMail.MailInit()
	myMail.SetMailTitle(mailTitle)
	fmt.Printf("[%s] Send E-mail\n",time.Now().Format("2006-01-02 15:04:05"))
	myMail.SetMailBody(mailContent)
	myMail.SendMail()
}