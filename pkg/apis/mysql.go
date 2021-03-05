package apis

import (
	"database/sql"
	"errors"
	"os"
	"strings"
	"time"

	"github.com/cnrancher/pdf-sender/pkg/types"
	"github.com/sirupsen/logrus"

	_ "github.com/go-sql-driver/mysql"
)

var DB *sql.DB

var (
	dbhostsip  = os.Getenv("DB_HOST_IP")
	dbusername = os.Getenv("DB_USERNAME")
	dbpassword = os.Getenv("DB_PASSWORD")
	dbname     = os.Getenv("DB_NAME")
)

func ConnectMysql() {

	dbinfo := strings.Join([]string{dbusername, ":", dbpassword, "@tcp(", dbhostsip, ")/", dbname, "?charset=utf8&parseTime=true"}, "")

	logrus.Infof(dbinfo)
	err := errors.New("")

	DB, err = sql.Open("mysql", dbinfo)
	if nil != err {
		logrus.Fatalf("Failed to open Database : %v", err)
	}

	DB.SetConnMaxLifetime(100)
	DB.SetMaxIdleConns(10)

	if err = DB.Ping(); nil != err {
		logrus.Fatalf("Failed to ping MySQL : %v", err)
	}

	logrus.Infof("Connect MySQL Successful")

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
