package apis

import (
	"bytes"
	"errors"
	"fmt"
	"net/url"
	"os"
	"strings"

	"github.com/cnrancher/pdf-sender/pkg/types"
	"github.com/sirupsen/logrus"
	"github.com/urfave/cli"
	"gopkg.in/gomail.v2"
)

type document struct {
	Name     string `json:"name"`
	Title    string `json:"title"`
	URL      string `json:"url"`
	PWD      string `json:"pwd"`
	Filename string `json:"filename"`
}

func (d document) validate() error {
	switch {
	case d.Name == "":
		return fmt.Errorf("document name is missing")
	case d.URL != "": // Using direct URL for download
		return nil
	case d.Filename == "":
		return fmt.Errorf("the filename for document %s is missing", d.Name)
	case types.OSSBucketEndpoint != "": // Using oss to ge
		exists, err := types.IsObjectExists(d.Filename)
		if err != nil {
			return err
		}
		if !exists {
			return errors.New("the document file does not exists in OSS")
		}
		return nil
	case fileURLPrefix == "":
		return fmt.Errorf("document %s is missing url configuration or file url prefix is not set", d.Name)
	}
	return nil
}

// When generating document download link for email,
// we will use following priority for configuration.
// 1. Direct URL from configuration with password.
// 2. OSS gLobal Configuration.
// 3. Global url prefix and filename with password.
func (d document) getLine(isPublic bool) string {
	var fileURL string
	switch {
	case d.URL != "":
		fileURL = d.URL
		if d.PWD != "" {
			fileURL += "     访问密码： " + d.PWD
		}
	case types.OSSBucketEndpoint != "": // using oss as file hoster
		fileURL = types.GetOSSFileDownloadURL(d.Filename, isPublic)
	default:
		fileURL = fileURLPrefix
		if !strings.HasSuffix(fileURL, "/") {
			fileURL += "/"
		}
		fileURL += url.PathEscape(d.Filename)
		if d.PWD != "" {
			fileURL += "     访问密码： " + d.PWD
		}
	}

	title := d.Title
	if title == "" {
		title = d.Name
	}
	return fmt.Sprintf("%s:     %s\n", title, fileURL)
}

var (
	fileURLPrefix,
	senderEmail,
	smtpUser,
	smtpPWD,
	smtpEndpoint string
	smtpPort int

	documents = map[string]document{
		"k3s": {
			Name:     "k3s",
			Title:    "K3s",
			Filename: "K3s.pdf",
			URL:      os.Getenv("K3s_PDF_URL"),
			PWD:      os.Getenv("K3s_PWD"),
		},
		"harvester": {
			Name:     "harvester",
			Title:    "Harvester",
			Filename: "Harvester.pdf",
			URL:      os.Getenv("Harvester_PDF_URL"),
			PWD:      os.Getenv("Harvester_PWD"),
		},
		"octopus": {
			Name:     "octopus",
			Title:    "Octopus",
			Filename: "Octopus_CN_Doc.pdf",
			URL:      os.Getenv("Octopus_PDF_URL"),
			PWD:      os.Getenv("Octopus_PWD"),
		},
		"rancher2": {
			Name:     "rancher2",
			Title:    "Rancher2.x",
			Filename: "Rancher2.x_CN_Doc.pdf",
			URL:      os.Getenv("Rancher2_PDF_URL"),
			PWD:      os.Getenv("Rancher2_PWD"),
		},
		"rancher2.5": {
			Name:     "rancher2.5",
			Title:    "Rancher2.5",
			Filename: "Rancher2.5_CN_Doc.pdf",
		},
		"rancher1": {
			Name:     "rancher1",
			Title:    "Rancher1.6",
			Filename: "rancher1.6.pdf",
		},
		"rke": {
			Name:     "rke",
			Title:    "RKE",
			Filename: "rke.pdf",
			URL:      os.Getenv("RKE_PDF_URL"),
			PWD:      os.Getenv("RKE_PWD"),
		},
	}
	appliedDocuments []string
)

const (
	header = `您好，

您可以通过下面的链接下载 Rancher 中文文档。

`
	footer = `Best Regards,
Rancher Labs 源澈科技
`
)

func InitEmailBody(ctx *cli.Context) error {
	if types.OSSBucketEndpoint != "" {
		logrus.Info("Using oss as file host url")
	}
	for _, kind := range ctx.StringSlice("published-docs") {
		d, ok := documents[kind]
		if !ok {
			return fmt.Errorf("document kind %s is not configured", kind)
		}
		if err := d.validate(); err != nil {
			return err
		}
		appliedDocuments = append(appliedDocuments, kind)
	}

	return nil
}

func SendEmail(user *types.User) {
	m := gomail.NewMessage()
	m.SetAddressHeader("From", senderEmail, "Rancher Labs 中国")
	m.SetHeader("To", user.Email)
	m.SetHeader("Subject", "Rancher 2.x 中文文档")
	m.SetBody("text/plain", getBody())

	d := gomail.NewDialer(smtpEndpoint, smtpPort, smtpUser, smtpPWD)

	if err := d.DialAndSend(m); err != nil {
		logrus.Errorf("Send Email to %s err:%v", user.Email, err)
		user.Status = false
	} else {
		logrus.Infof("Send Email to %s success", user.Email)
		user.Status = true
	}

	DBSave(user)
}

func getBody() string {
	isPublic := true
	if types.OSSBucketEndpoint != "" {
		info, err := types.GetBucketInfo()
		if err != nil {
			logrus.Warnf("failed to get bucket info, assuming public access, error: %v", err)
		} else {
			isPublic = types.IsBucketPrivate(info.ACL)
		}
	}

	writer := bytes.NewBufferString(header)
	for _, kind := range appliedDocuments {
		d := documents[kind]
		line := d.getLine(isPublic) + "\n"
		writer.WriteString(line)
	}
	writer.WriteString(footer)
	return writer.String()
}
