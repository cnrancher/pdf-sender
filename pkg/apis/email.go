package apis

import (
	"github.com/cnrancher/pdf-sender/pkg/types"
	"github.com/sirupsen/logrus"
	"gopkg.in/gomail.v2"
)

func SendEmail(user *types.User) {
	subject, sender, err := types.GetSenderNameAndSubjectByKind(user.Kind)
	if err != nil {
		logrus.Warnf("failed to send email of kind %s for user %v, %v", user.Kind, *user, err)
		return
	}
	body, err := types.GetBodyByKind(user.Kind)
	if err != nil {
		logrus.Warnf("failed to send email of kind %s for user %v, %v", user.Kind, *user, err)
		return
	}
	m := gomail.NewMessage()
	m.SetAddressHeader("From", types.Email.Sender, sender)
	m.SetHeader("To", user.Email)
	m.SetHeader("Subject", subject)
	m.SetBody("text/plain", body)

	d := gomail.NewDialer(
		types.Email.Endpoint,
		types.Email.Port,
		types.Email.User,
		types.Email.Password)

	if err := d.DialAndSend(m); err != nil {
		logrus.Errorf("Send Email to %s err:%v", user.Email, err)
		user.Status = false
	} else {
		logrus.Infof("Send Email to %s success", user.Email)
		user.Status = true
	}

	DBSave(user)
}
