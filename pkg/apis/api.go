package apis

import (
	"fmt"
	"math/rand"
	"strings"
	"time"

	"github.com/sirupsen/logrus"
	"gorm.io/gorm"

	"github.com/cnrancher/pdf-sender/pkg/email"
	"github.com/cnrancher/pdf-sender/pkg/types"
	restful "github.com/emicklei/go-restful/v3"
)

type validateFunc func(string) error

func RegisterAPIs() *restful.Container {
	container := restful.NewContainer()

	cors := restful.CrossOriginResourceSharing{
		AllowedHeaders: []string{"Content-Type", "Accept"},
		AllowedMethods: []string{"OPTIONS", "POST"},
		AllowedDomains: []string{"*"},
		CookiesAllowed: false,
		Container:      container}

	docxWs := new(restful.WebService).Path("/").Consumes(restful.MIME_JSON).Produces(restful.MIME_JSON)

	container.Add(docxWs)

	container.Filter(cors.Filter)

	container.Filter(container.OPTIONSFilter)

	docxWs.Route(
		docxWs.POST("/sendCode").
			To(sendCode).
			Param(docxWs.QueryParameter("kind", "request kind(pdf/ent/demo/pricing/contact)").Required(true)).Filter(kindFilter("")).
			Reads(types.Code{}))
	docxWs.Route(
		docxWs.POST("/sendEmail").
			To(sendEmail).
			Param(docxWs.QueryParameter("kind", "request kind(pdf or ent)").Required(true)).Filter(kindFilter("pdf")).
			Reads(types.User{}))
	docxWs.Route(
		docxWs.POST("/register").
			Param(docxWs.QueryParameter("kind", "request kind(demo, pricing or contact)").Required(true)).
			Filter(kindFilter("register")).
			Reads(types.User{}).
			To(register))
	return container

}

func sendCode(req *restful.Request, resp *restful.Response) {
	var code types.Code
	if err := req.ReadEntity(&code); err != nil {
		logrus.Errorf("Get user err:%v", err)
		err = resp.WriteErrorString(400, err.Error())
		if err != nil {
			logrus.Errorf("Failed to write error string err:%v", err)
		}
		return
	}

	if code.Phone == "" {
		logrus.Errorf("Phone number cannot be empty")
		if err := resp.WriteErrorString(400, "手机号不能为空"); err != nil {
			logrus.Errorf("Failed to write error string err:%v", err)
		}
		return
	}

	code.Kind = req.QueryParameter("kind")
	code.Code = GenValidateCode(4)

	logrus.Debugf(code.Code)

	if err := code.SaveAndSend(); err != nil {
		logrus.Errorf("Send SMS to phone %s err:%v", code.Phone, err)
		err = resp.WriteErrorString(400, err.Error())
		if err != nil {
			logrus.Errorf("Failed to write error string err:%v", err)
		}
		return
	}

	if err := resp.WriteErrorString(200, "验证码发送成功"); err != nil {
		logrus.Errorf("Failed to write error string err:%v", err)
	}
}

func sendEmail(req *restful.Request, resp *restful.Response) {
	var user types.User
	if err := req.ReadEntity(&user); err != nil {
		logrus.Errorf("Get user err:%v", err)
		err = resp.WriteErrorString(400, err.Error())
		if err != nil {
			logrus.Errorf("Failed to write error string err:%v", err)
		}
		return
	}
	user.Kind = req.QueryParameter("kind")
	code, err := user.Validate()
	if err != nil {
		logrus.Errorf("Validate user err:%v", err)
		if err := resp.WriteErrorString(400, err.Error()); err != nil {
			logrus.Errorf("Failed to write error string err:%v", err)
		}
		return
	}

	if err := types.DBInstance.Transaction(func(tx *gorm.DB) error {
		if err := user.Save(tx).Error; err != nil {
			return err
		}
		if err := tx.Model(code).UpdateColumn("state", "used").Error; err != nil {
			return err
		}
		return nil
	}); err != nil {
		resp.WriteErrorString(400, "数据库更新失败")
		return
	}

	content, err := getContentFromUser(&user)
	if err != nil {
		logrus.Warnf("failed to get email content for user %s(%s), %v", user.Name, user.Phone, err)
		return
	}
	if err := email.SendEmail(content); err != nil {
		logrus.Warnf("failed to send email for user %s(%s), %v", user.Name, user.Phone, err)
		return
	}

	user.Status = true
	if err := user.Save(nil).Error; err != nil {
		logrus.Warnf("failed to update user %d state, %v", user.UID, err)
	}
}

func GenValidateCode(width int) string {
	numeric := [9]byte{1, 2, 3, 4, 5, 6, 7, 8, 9}
	r := len(numeric)
	rand.Seed(time.Now().UnixNano())

	var sb strings.Builder
	for i := 0; i < width; i++ {
		fmt.Fprintf(&sb, "%d", numeric[rand.Intn(r)])
	}

	return sb.String()
}

func kindFilter(requestType string) func(*restful.Request, *restful.Response, *restful.FilterChain) {
	return func(req *restful.Request, resp *restful.Response, next *restful.FilterChain) {
		if err := types.IsKindValid(requestType, req.QueryParameter("kind")); err != nil {
			if err = resp.WriteErrorString(400, err.Error()); err != nil {
				logrus.Warnf("failed to write error response, %v", err)
			}
			return
		}
		next.ProcessFilter(req, resp)
	}
}

func getContentFromUser(user *types.User) (*email.Content, error) {
	rtn := &email.Content{
		To: map[string]string{},
		CC: map[string]string{},
	}
	subject, sender, err := email.GetSenderNameAndSubjectByKind(user.Kind)
	if err != nil {
		logrus.Warnf("failed to send email of kind %s for user %v, %v", user.Kind, *user, err)
		return nil, err
	}
	rtn.Subject = subject
	rtn.To[user.Email] = ""
	rtn.From = types.Email.Sender
	rtn.FromAlias = sender
	body, err := email.GetBodyByKind(user.Kind)
	if err != nil {
		logrus.Warnf("failed to send email of kind %s for user %v, %v", user.Kind, *user, err)
		return nil, err
	}
	rtn.Body = body
	return rtn, nil
}
