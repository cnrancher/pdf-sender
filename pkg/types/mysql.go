package types

import (
	"fmt"

	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

var DBInstance *gorm.DB

func ConnectMysql() error {
	dbinfo := fmt.Sprintf(
		"%s:%s@tcp(%s:%d)/%s?charset=utf8&parseTime=true",
		DB.Username,
		DB.Password,
		DB.HostIP,
		DB.Port,
		DB.Name,
	)
	logrus.Debugf("Connecting mysql with %s", dbinfo)
	var err error

	DBInstance, err = gorm.Open(mysql.Open(dbinfo), &gorm.Config{})
	if nil != err {
		return errors.Wrap(err, "Failed to open Database")
	}
	sqlDB, err := DBInstance.DB()
	if err != nil {
		return err
	}
	sqlDB.SetConnMaxLifetime(100)
	sqlDB.SetMaxIdleConns(10)

	if err = sqlDB.Ping(); nil != err {
		return errors.Wrap(err, "Failed to ping MySQL")
	}

	logrus.Infof("Connect MySQL Successful")
	logrus.Info("Migrating Database...")
	if err := DBInstance.AutoMigrate(&User{}, &Code{}); err != nil {
		return err
	}
	logrus.Info("Migrating Database succeed")
	return nil
}
