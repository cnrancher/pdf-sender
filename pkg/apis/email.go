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

	PDFUrl = os.Getenv("PDF_URL")
	PDFPwd = os.Getenv("PDF_PWD")

	body = `您好，
	
您可以通过下面的链接和密码下载 Rancher 中文文档。
	
下载链接：
	
` + PDFUrl + ` 
	
访问密码：
	
` + PDFPwd + `
	
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
		logrus.Errorf("Send Email to %s err:%v", user.Email, err)
		user.Status = false
	} else {
		logrus.Infof("Send Email to %s success", user.Email)
		user.Status = true
	}

	DBSave(user)
}
