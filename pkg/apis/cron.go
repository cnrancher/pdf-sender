package apis

import (
	"os"
	"strconv"
	"time"

	"github.com/360EntSecGroup-Skylar/excelize"
	"github.com/rancher/pdf-sender/pkg/types"
	cron "github.com/robfig/cron/v3"
	"github.com/sirupsen/logrus"
	"gopkg.in/gomail.v2"
)

var (
	SMTPRancherTo = os.Getenv("SMTP_RANCHER_TO")
	CronJob = os.Getenv("CRON")
)

func CollectInformation() {
	c := cron.New()
	logrus.Infof("Collect information start")
	c.AddFunc(CronJob, func() {
		excelName := DBSelect()
		SendInformation(excelName)
	})
	c.Start()
}

func SendInformation(excelName string) {
	m := gomail.NewMessage()

	m.SetAddressHeader("From", "no-reply@rancher.cn", "Rancher Labs 中国")
	m.SetHeader("To", SMTPRancherTo)
	m.SetHeader("Subject", excelName+"用户信息")
	m.Attach("/tmp/" + excelName + ".xlsx")

	d := gomail.NewDialer(SMTPEndpoint, 587, SMTPUser, SMTPPwd)

	if err := d.DialAndSend(m); err != nil {
		logrus.Errorf("Failed to send collect information email : %v", err)
	} else {
		logrus.Infof("Send collect information email success")
	}

}

func DBSelect() string {

	xlsx := excelize.NewFile()

	index := xlsx.NewSheet("用户信息表")

	data := map[string]string{
		"A1": "名字",
		"B1": "公司",
		"C1": "职位",
		"D1": "手机号",
		"E1": "电子邮箱",
		"F1": "保存时间",
		"G1": "邮箱是否有效",
	}

	rows, err := DB.Query("SELECT * FROM " + dbtable + " WHERE date(savetime) = date_sub(curdate(),interval 1 day)")
	if nil != err {
		logrus.Errorf("Failed to query : %v", err)
	}

	defer rows.Close()

	var count int = 1
	for rows.Next() {
		count++
		var user types.User
		err := rows.Scan(&user.UID, &user.Name, &user.Company, &user.Position, &user.Phone, &user.Email, &user.SaveTime, &user.Status)
		if err != nil {
			logrus.Errorf("Failed rows scan : %v", err)
		}
		data["A"+strconv.Itoa(count)] = user.Name
		data["B"+strconv.Itoa(count)] = user.Company
		data["C"+strconv.Itoa(count)] = user.Position
		data["D"+strconv.Itoa(count)] = user.Phone
		data["E"+strconv.Itoa(count)] = user.Email
		data["F"+strconv.Itoa(count)] = user.SaveTime.Format("2006-01-02 15:04:05")
		data["G"+strconv.Itoa(count)] = strconv.FormatBool(user.Status)
	}

	for k, v := range data {
		xlsx.SetCellValue("用户信息表", k, v)
	}

	xlsx.SetActiveSheet(index)

	d, _ := time.ParseDuration("-24h")

	excelName := time.Now().Add(d).Format("2006-01-02")
	err = xlsx.SaveAs("/tmp/" + excelName + ".xlsx")
	if err != nil {
		logrus.Errorf("Failed to save excel : %v", err)
	}

	return excelName
}
