package apis

import (
	"strings"

	"github.com/cnrancher/pdf-sender/pkg/email"
	"github.com/cnrancher/pdf-sender/pkg/types"
	"github.com/emicklei/go-restful/v3"
	"github.com/sirupsen/logrus"
	"gorm.io/gorm"
)

func register(req *restful.Request, resp *restful.Response) {
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

	content, err := getContentFromRegisterConfig(&user)
	if err != nil {
		logrus.Warnf("failed to get content from register config, %v", err)
		return
	}
	if content == nil {
		logrus.Infof("no email content should be send for user %s(%s)", user.Name, user.Phone)
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

func getContentFromRegisterConfig(user *types.User) (*email.Content, error) {
	writer := &strings.Builder{}
	if err := email.RegisterTemplate.Execute(writer, map[string]interface{}{
		"kind": types.Register.Kinds[user.Kind],
		"user": user,
	}); err != nil {
		logrus.Error(err)
		return nil, err
	}
	var rtn email.Content
	rtn.Body = writer.String()
	rtn.To = map[string]string{}
	for _, addr := range types.Register.Receivers {
		rtn.To[addr] = ""
	}
	rtn.Subject = types.Register.Subject
	rtn.From = types.Email.Sender
	rtn.FromAlias = types.Register.SenderName
	return &rtn, nil
}
