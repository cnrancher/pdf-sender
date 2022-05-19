package types

import (
	"fmt"

	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
	"gorm.io/driver/mysql"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var DBInstance *gorm.DB

func ConnectDB() error {
	DBParams := []interface{}{
		DB.Username,
		DB.Password,
		DB.HostIP,
		DB.Port,
		DB.Name,
	}
	//TODO localization should be able to configure
	logrus.Debugf("Connecting %s with DB %s:%d/%s, user %s, password %s,", DB.Kind, DB.HostIP, DB.Port, DB.Name, DB.Username, DB.Password)
	var err error
	switch DB.Kind {
	case "mysql":
		DBInstance, err = gorm.Open(mysql.Open(fmt.Sprintf(
			"%s:%s@tcp(%s:%d)/%s?charset=utf8&parseTime=true&loc=Local",
			DBParams...,
		)), &gorm.Config{})
	case "pgsql":
		DBInstance, err = gorm.Open(postgres.Open(fmt.Sprintf(
			"user=%s password=%s host=%s port=%d dbname=%s sslmode=disable TimeZone=Asia/Shanghai",
			DBParams...,
		)), &gorm.Config{})
	default:
		return fmt.Errorf("db kind %s is not support", DB.Kind)
	}
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
		return errors.Wrapf(err, "Failed to ping %s", DB.Kind)
	}

	logrus.Infof("Connect %s Successful", DB.Kind)
	logrus.Info("Migrating Database...")
	if err := DBInstance.AutoMigrate(&User{}, &Code{}); err != nil {
		return err
	}
	logrus.Info("Migrating Database succeed")
	return nil
}
