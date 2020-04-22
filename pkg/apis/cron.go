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
	CronJob       = os.Getenv("CRON")

	yesterday string
	today     string
	count     int
)

func CollectInformation() {
	c := cron.New()
	logrus.Infof("Collect information start")
	_, err := c.AddFunc(CronJob, func() {
		DBSelect()
		SendInformation()
	})
	if err != nil {
		logrus.Errorf("Failed cron add function : %v", err)
	}

	c.Start()
}

func SendInformation() {
	m := gomail.NewMessage()

	m.SetAddressHeader("From", "no-reply@rancher.cn", "Rancher Labs 中国")
	m.SetHeader("To", SMTPRancherTo)
	m.SetHeader("Subject", yesterday+"用户信息")
	m.SetBody("text/plain", `截止 `+yesterday+` 00:00 ~ `+today+` 00:00，一共有 `+strconv.Itoa(count)+` 个人`)

	m.Attach("/tmp/" + yesterday + ".xlsx")

	d := gomail.NewDialer(SMTPEndpoint, 587, SMTPUser, SMTPPwd)

	if err := d.DialAndSend(m); err != nil {
		logrus.Errorf("Failed to send collect information email : %v", err)
	} else {
		logrus.Infof("Send collect information email success")
	}

}

func DBSelect() {

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

	stmt, err := DB.Prepare("SELECT * FROM user WHERE date(savetime) = date_sub(curdate(),interval 1 day)")
	if err != nil {
		logrus.Errorf("Failed to prepare SQL statement : %v", err)
	}
	defer stmt.Close()

	rows, err := stmt.Query()
	if nil != err {
		logrus.Errorf("Failed to query : %v", err)
	}

	defer rows.Close()

	count = 0
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

	d, err := time.ParseDuration("-24h")
	if err != nil {
		logrus.Errorf("Failed time parsed duration : %v", err)
	}

	cstSh, err := time.LoadLocation("Asia/Shanghai")
	if err != nil {
		logrus.Errorf("Failed time load location : %v", err)
	}

	yesterday = time.Now().In(cstSh).Add(d).Format("2006-01-02")
	today = time.Now().In(cstSh).Format("2006-01-02")
	err = xlsx.SaveAs("/tmp/" + yesterday + ".xlsx")
	if err != nil {
		logrus.Errorf("Failed to save excel : %v", err)
	}
}
