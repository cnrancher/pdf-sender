package apis

import (
	"os"

	"github.com/rancher/pdf-sender/pkg/types"
	"github.com/sirupsen/logrus"
	"gopkg.in/gomail.v2"
)

var (
	SMTPUser     = os.Getenv("SMTP_USER")
	SMTPPwd      = os.Getenv("SMTP_PWD")
	SMTPEndpoint = os.Getenv("SMTP_ENDPOINT")
	SMTPPort     = os.Getenv("SMTP_PORT")

	body = `您好，
	
您可以通过下面的链接和密码下载 Rancher 中文文档。
	
下载链接：
	
https://v2.fangcloud.com/share/2bcac9426816768baa179a8435 
	
访问密码：
	
957f1e
	
Best Regards,
Rancher Labs 源澈科技`
)

func SendEmail(user *types.User) {

	m := gomail.NewMessage()
	m.SetAddressHeader("From", "no-reply@rancher.cn", "Rancher Labs 中国")
	m.SetHeader("To", user.Email)
	m.SetHeader("Subject", "Rancher 2.x 中文文档")
	m.SetBody("text/plain", body)

	d := gomail.NewDialer(SMTPEndpoint, 587, SMTPUser, SMTPPwd)

	err := d.DialAndSend(m)

	if err != nil {
		logrus.Errorf("Send Email err:%v", err)
		user.Status = false
	} else {
		logrus.Infof("Send Email success")
		user.Status = true
	}

	DBSave(user)
}
