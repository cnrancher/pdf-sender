package apis

import (
	"os"
	"strconv"

	"github.com/rancher/pdf-sender/pkg/types"
	"github.com/sirupsen/logrus"
	"gopkg.in/gomail.v2"
)

var (
	SMTPUser     = os.Getenv("SMTP_USER")
	SMTPPwd      = os.Getenv("SMTP_PWD")
	SMTPEndpoint = os.Getenv("SMTP_ENDPOINT")
	SMTPPort     = os.Getenv("SMTP_PORT")

	Rancher2PDFUrl = os.Getenv("Rancher2_PDF_URL")
	Rancher2PDFPwd = os.Getenv("Rancher2_PWD")

	RKEPDFUrl = os.Getenv("RKE_PDF_URL")
	RKEPDFPwd = os.Getenv("RKE_PWD")

	K3sPDFUrl = os.Getenv("K3s_PDF_URL")
	K3sPDFPwd = os.Getenv("K3s_PWD")

	OctopusPDFUrl = os.Getenv("Octopus_PDF_URL")
	OctopusPDFPwd = os.Getenv("Octopus_PWD")

	HarvesterPDFUrl = os.Getenv("Harvester_PDF_URL")
	HarvesterPDFPwd = os.Getenv("Harvester_PWD")

	SenderEmail = os.Getenv("SENDER_EMAIL")

	body = `您好，
	
您可以通过下面的链接和密码下载 Rancher 中文文档。

Rancher2.x：    ` + Rancher2PDFUrl + `     访问密码： ` + Rancher2PDFPwd + `

RKE：    ` + RKEPDFUrl + `     访问密码： ` + RKEPDFPwd + `

K3s：    ` + K3sPDFUrl + `     访问密码： ` + K3sPDFPwd + `

Octopus：    ` + OctopusPDFUrl + `     访问密码： ` + OctopusPDFPwd + `

Harvester：    ` + HarvesterPDFUrl + `     访问密码： ` + HarvesterPDFPwd + `
	
Best Regards,
Rancher Labs 源澈科技`
)

func SendEmail(user *types.User) {
	port, err := strconv.Atoi(SMTPPort)
	if err != nil {
		logrus.Errorf("smtp port err: %v", err)
	}

	m := gomail.NewMessage()
	m.SetAddressHeader("From", SenderEmail, "Rancher Labs 中国")
	m.SetHeader("To", user.Email)
	m.SetHeader("Subject", "Rancher 2.x 中文文档")
	m.SetBody("text/plain", body)

	d := gomail.NewDialer(SMTPEndpoint, port, SMTPUser, SMTPPwd)
	err = d.DialAndSend(m)
	if err != nil {
		logrus.Errorf("Send Email to %s err:%v", user.Email, err)
		user.Status = false
	} else {
		logrus.Infof("Send Email to %s success", user.Email)
		user.Status = true
	}

	DBSave(user)
}
