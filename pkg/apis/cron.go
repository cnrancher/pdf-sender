package apis

import (
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/360EntSecGroup-Skylar/excelize"
	"github.com/rancher/pdf-sender/pkg/types"
	cron "github.com/robfig/cron/v3"
	"github.com/sirupsen/logrus"
	"gopkg.in/gomail.v2"
)

var (
	SMTPRancherToDay = os.Getenv("SMTP_RANCHER_TO_DAY")
	SMTPRancherToMon = os.Getenv("SMTP_RANCHER_TO_MON")
	DayCronJob       = os.Getenv("DAY_CRON")
	MonCronJob       = os.Getenv("MON_CRON")
)

const (
	DaySQL = "SELECT * FROM user WHERE date(savetime) = date_sub(curdate(),interval 1 day)"
	MonSQL = "SELECT * FROM user WHERE PERIOD_DIFF(DATE_FORMAT(NOW(),'%Y%m'),DATE_FORMAT(savetime,'%Y%m')) = 1"
)

func CollectInformation() {
	c := cron.New()
	logrus.Infof("Collect information start")

	_, err := c.AddFunc(DayCronJob, func() {
		logrus.Infof("Send Information For Day")

		d, err := time.ParseDuration("-24h")
		if err != nil {
			logrus.Errorf("Failed time parsed duration : %v", err)
		}
		yesterday := time.Now().Add(d).Format("2006-01-02")
		today := time.Now().Format("2006-01-02")
		count := DBSelect(DaySQL, yesterday)

		headMessage := yesterday + "用户信息"
		bodyMessage := yesterday + "08:00 ~ " + today + " 08:00，一共有 " + strconv.Itoa(count) + " 人下载了中文文档。"
		SendInformation(yesterday, headMessage, bodyMessage, SMTPRancherToDay)
	})
	if err != nil {
		logrus.Errorf("Failed cron add function : %v", err)
	}

	_, err = c.AddFunc(MonCronJob, func() {
		logrus.Infof("Send Information For Mon")

		now := time.Now()
		lastMonth := now.AddDate(0, -1, -now.Day()+1).Format("2006-01")
		count := DBSelect(MonSQL, lastMonth)

		headMessage := lastMonth + "月 全部用户信息"
		bodyMessage := lastMonth + "月一共有 " + strconv.Itoa(count) + " 人下载了中文文档。"
		SendInformation(lastMonth, headMessage, bodyMessage, SMTPRancherToMon)
	})
	if err != nil {
		logrus.Errorf("Failed cron add function : %v", err)
	}

	c.Start()
}

func SendInformation(xlsxName, headMessage, bodyMessage, to string) {
	sends := strings.Split(to, ",")

	port, err := strconv.Atoi(SMTPPort)
	if err != nil {
		logrus.Errorf("smtp port err: %v", err)
	}

	m := gomail.NewMessage()
	m.SetAddressHeader("From", SenderEmail, "Rancher Labs 中国")
	m.SetHeader("To", sends...)
	m.SetHeader("Subject", headMessage)
	m.SetBody("text/plain", bodyMessage)
	m.Attach("/tmp/" + xlsxName + ".xlsx")

	d := gomail.NewDialer(SMTPEndpoint, port, SMTPUser, SMTPPwd)

	if err := d.DialAndSend(m); err != nil {
		logrus.Errorf("Failed to send collect information email : %v", err)
	} else {
		logrus.Infof("Send collect information email success")
	}

}

func DBSelect(sql, xlsxName string) int {
	xlsx := excelize.NewFile()
	index := xlsx.GetSheetIndex("Sheet1")
	data := map[string]string{
		"A1": "名字",
		"B1": "公司",
		"C1": "职位",
		"D1": "手机号",
		"E1": "电子邮箱",
		"F1": "保存时间",
		"G1": "邮箱是否有效",
	}

	stmt, err := DB.Prepare(sql)
	if err != nil {
		logrus.Errorf("Failed to prepare SQL statement : %v", err)
	}
	defer stmt.Close()

	rows, err := stmt.Query()
	if nil != err {
		logrus.Errorf("Failed to query : %v", err)
	}
	defer rows.Close()

	count := 0
	for rows.Next() {
		count++
		var user types.User
		row := count + 1
		err := rows.Scan(&user.UID, &user.Name, &user.Company, &user.Position, &user.Phone, &user.Email, &user.SaveTime, &user.Status)
		if err != nil {
			logrus.Errorf("Failed rows scan : %v", err)
		}
		data["A"+strconv.Itoa(row)] = user.Name
		data["B"+strconv.Itoa(row)] = user.Company
		data["C"+strconv.Itoa(row)] = user.Position
		data["D"+strconv.Itoa(row)] = user.Phone
		data["E"+strconv.Itoa(row)] = user.Email
		data["F"+strconv.Itoa(row)] = user.SaveTime.Format("2006-01-02 15:04:05")
		data["G"+strconv.Itoa(row)] = strconv.FormatBool(user.Status)
	}

	for k, v := range data {
		xlsx.SetCellValue("Sheet1", k, v)
	}
	xlsx.SetActiveSheet(index)

	if err = xlsx.SaveAs("/tmp/" + xlsxName + ".xlsx"); err != nil {
		logrus.Errorf("Failed to save excel : %v", err)
	}

	return count
}
