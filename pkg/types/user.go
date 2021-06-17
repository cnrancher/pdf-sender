package types

import (
	"encoding/json"
	"fmt"
	"strconv"
	"time"

	"github.com/sirupsen/logrus"

	"github.com/aliyun/alibaba-cloud-sdk-go/services/dysmsapi"
)

type User struct {
	UID      int       `json:"uid"`
	Name     string    `json:"name"`
	Company  string    `json:"company"`
	Position string    `json:"position"`
	Email    string    `json:"email"`
	Phone    string    `json:"phone"`
	Code     string    `json:"code"`
	SaveTime time.Time `json:"saveTime"`
	Status   bool      `json:"status"`
	Kind     string    `json:"-"`
}

func (u *User) Send() error {
	client, err := dysmsapi.NewClientWithAccessKey(Aliyun.Region, Aliyun.AccessKey, Aliyun.AccessSecret)
	if err != nil {
		return fmt.Errorf("使用密钥创建客户端失败:%v", err)
	}

	intCode, err := strconv.Atoi(u.Code)
	if err != nil {
		return fmt.Errorf("验证码转换类型失败:%v", err)
	}

	request := dysmsapi.CreateSendSmsRequest()

	request.Scheme = "https"
	request.PhoneNumbers = u.Phone
	request.SignName = Aliyun.SMSSignName
	request.TemplateCode = Aliyun.SMSTemplateCode
	request.TemplateParam = fmt.Sprintf(`{"code":"%d"}`, intCode)

	response, err := client.SendSms(request)
	if err != nil {
		return fmt.Errorf("发送短信错误:%v", err)
	}

	if response.Code != "OK" {
		return fmt.Errorf("Failed to send Aliyun SMS message:%s", response.Message)
	}

	logrus.Infof("发送成功")
	return nil
}

func New(bodyData []byte) (*User, error) {
	var user User
	if err := json.Unmarshal(bodyData, &user); err != nil {
		return nil, err
	}

	return &user, nil
}

func (u *User) Validate() error {
	if u.Code == "" {
		return fmt.Errorf("验证码不能为空")
	}
	if u.Phone == "" {
		return fmt.Errorf("手机号不能为空")
	}
	if u.Email == "" {
		return fmt.Errorf("电子邮箱不能为空")
	}

	return nil
}

func (u *User) Compare(target *User) bool {
	return u.Code == target.Code &&
		u.Phone == target.Phone &&
		u.Kind == target.Kind
}
