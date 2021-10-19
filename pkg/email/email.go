package email

import (
	"strings"

	"github.com/cnrancher/pdf-sender/pkg/types"
	"github.com/sirupsen/logrus"
	"gopkg.in/gomail.v2"
)

type Content struct {
	From,
	FromAlias,
	Subject,
	Body,
	BodyFormat string
	To, CC map[string]string
	Attach string
}

func SendEmail(content *Content) error {
	if types.Config.Debug {
		logrus.Debug("going to send email content")
		logrus.Debugf("%+v\n", *content)
	}
	if types.Config.DryRun {
		return nil
	}
	m := gomail.NewMessage()
	if content.FromAlias != "" {
		m.SetAddressHeader("From", content.From, content.FromAlias)
	} else {
		m.SetHeader("From", content.From)
	}
	var toList, ccList []string
	for addr, alias := range content.To {
		toList = append(toList, addr)
		if alias != "" {
			m.SetAddressHeader("To", addr, alias)
		}
	}
	if len(toList) != 0 {
		m.SetHeader("To", toList...)
	}
	for addr, alias := range content.CC {
		ccList = append(ccList, addr)
		if alias != "" {
			m.SetAddressHeader("Cc", addr, alias)
		}
	}
	if len(ccList) != 0 {
		m.SetHeader("Cc", ccList...)
	}

	m.SetHeader("Subject", content.Subject)
	if content.BodyFormat == "" {
		m.SetBody("text/plain", content.Body)
	} else {
		m.SetBody(content.BodyFormat, content.Body)
	}

	if content.Attach != "" {
		m.Attach(content.Attach)
	}

	d := gomail.NewDialer(
		types.Email.Endpoint,
		types.Email.Port,
		types.Email.User,
		types.Email.Password)
	if err := d.DialAndSend(m); err != nil {
		logrus.Errorf("Send Email to %s err:%v", strings.Join(toList, ","), err)
		return err
	}
	logrus.Infof("Send Email to %s success", strings.Join(toList, ","))
	return nil
}
