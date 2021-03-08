package apis

import (
	"database/sql"
	"fmt"
	"time"

	"github.com/cnrancher/pdf-sender/pkg/types"
	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
)

var DB *sql.DB

var (
	dbhostip,
	dbusername,
	dbpassword,
	dbname string
	dbport int
)

func ConnectMysql() error {
	dbinfo := fmt.Sprintf(
		"%s:%s@tcp(%s:%d)/%s?charset=utf8&parseTime=true",
		dbusername,
		dbpassword,
		dbhostip,
		dbport,
		dbname,
	)
	logrus.Infof("Connecting mysql with %s", dbinfo)
	var err error

	DB, err = sql.Open("mysql", dbinfo)
	if nil != err {
		return errors.Wrap(err, "Failed to open Database")
	}

	DB.SetConnMaxLifetime(100)
	DB.SetMaxIdleConns(10)

	if err = DB.Ping(); nil != err {
		return errors.Wrap(err, "Failed to ping MySQL")
	}

	logrus.Infof("Connect MySQL Successful")
	return nil
}

func DBSave(user *types.User) {
	stmt, err := DB.Prepare("INSERT INTO user(name, company, position, phone, email, savetime, status) values(?,?,?,?,?,?,?)")
	if nil != err {
		logrus.Errorf("Failed to prepare SQL statement : %v", err)
	}

	defer stmt.Close()

	_, err = stmt.Exec(user.Name, user.Company, user.Position, user.Phone, user.Email, time.Now(), user.Status)
	if nil != err {
		logrus.Errorf("Failed to executes SQL : %v", err)
	}
}
