package types

import (
	"encoding/json"
	"fmt"
	"time"

	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
	"gorm.io/gorm"
)

type User struct {
	UID        int       `json:"uid" gorm:"primaryKey"`
	Name       string    `json:"name" gorm:"size:255"`
	Company    string    `json:"company" gorm:"size:255"`
	Position   string    `json:"position" gorm:"size:255"`
	Email      string    `json:"email" gorm:"size:255"`
	Phone      string    `json:"phone" gorm:"size:255"`
	Code       string    `json:"code" gorm:"-"`
	SaveTime   time.Time `json:"saveTime" gorm:"column:savetime;autoUpdateTime:milli"`
	Status     bool      `json:"status"`
	Kind       string    `json:"-" gorm:"size:20"`
	City       string    `json:"city,omitempty" gorm:"size:255"`
	Department string    `json:"department,omitempty" gorm:"size:255"`
}

func New(bodyData []byte) (*User, error) {
	var user User
	if err := json.Unmarshal(bodyData, &user); err != nil {
		return nil, err
	}

	return &user, nil
}

func (u *User) Validate() (*Code, error) {
	if u.Code == "" {
		return nil, fmt.Errorf("验证码不能为空")
	}
	if u.Phone == "" {
		return nil, fmt.Errorf("手机号不能为空")
	}
	if u.Email == "" {
		return nil, fmt.Errorf("电子邮箱不能为空")
	}
	if u.Kind == "" {
		return nil, fmt.Errorf("未支持表单类型")
	}

	rows, err := DBInstance.Model(&Code{}).Where("phone = ? AND code = ? AND kind = ? AND state = ?", u.Phone, u.Code, u.Kind, "active").Order("request_time desc").Rows()
	if err != nil {
		return nil, errors.Wrapf(err, "failed to query codes")
	}
	defer rows.Close()
	for rows.Next() {
		var code Code
		if err := DBInstance.ScanRows(rows, &code); err != nil {
			return nil, errors.Wrap(err, "failed to decode row")
		}
		if time.Now().Sub(code.RequestTime) < 10*time.Minute {
			logrus.Debugf("code uid %d %s matching user %d", code.UID, code.Code, u.UID)
			return &code, nil
		}
	}
	return nil, fmt.Errorf("验证码不合法")
}

func (u *User) Compare(target *User) bool {
	return u.Code == target.Code &&
		u.Phone == target.Phone &&
		u.Kind == target.Kind
}

func (u *User) Save(tx *gorm.DB) *gorm.DB {
	db := DBInstance
	if tx != nil {
		db = tx
	}
	if u.UID == 0 {
		return db.Create(u)
	}
	return db.Save(u)
}

func (*User) TableName() string {
	return "user"
}
