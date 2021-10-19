package types

import (
	"fmt"
	"strconv"
	"time"

	"github.com/aliyun/alibaba-cloud-sdk-go/services/dysmsapi"
	"github.com/sirupsen/logrus"
	"gorm.io/gorm"
)

type Code struct {
	UID         uint      `json:"-" gorm:"primaryKey"`
	Phone       string    `json:"phone,omitempty" gorm:"type:VARCHAR(20);index:phone"`
	Code        string    `json:"-" gorm:"type:VARCHAR(20);index:phone"`
	Kind        string    `json:"-" gorm:"type:VARCHAR(20);index:phone"`
	RequestTime time.Time `json:"-" gorm:"autoCreateTime:milli"`
	UpdateTime  time.Time `json:"-" gorm:"autoUpdateTime:milli"`
	State       string    `json:"-" gorm:"size:20;default:active;index:phone"`
}

type CodeContent struct {
}

func (c *Code) SaveAndSend() error {
	if err := c.Save().Error; err != nil {
		return err
	}
	if Config.DryRun {
		return nil
	}
	return c.Send()
}

func (c *Code) Send() error {
	client, err := dysmsapi.NewClientWithAccessKey(Aliyun.Region, Aliyun.AccessKey, Aliyun.AccessSecret)
	if err != nil {
		return fmt.Errorf("使用密钥创建客户端失败:%v", err)
	}

	intCode, err := strconv.Atoi(c.Code)
	if err != nil {
		return fmt.Errorf("验证码转换类型失败:%v", err)
	}

	request := dysmsapi.CreateSendSmsRequest()

	request.Scheme = "https"
	request.PhoneNumbers = c.Phone
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

func (c *Code) Save() *gorm.DB {
	if c.UID == 0 {
		return DBInstance.Create(c)
	}
	return DBInstance.Save(c)
}

func (*Code) TableName() string {
	return "code"
}
