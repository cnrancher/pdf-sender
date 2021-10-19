package apis

import (
	"fmt"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/360EntSecGroup-Skylar/excelize"
	"github.com/cnrancher/pdf-sender/pkg/email"
	"github.com/cnrancher/pdf-sender/pkg/types"
	"github.com/pkg/errors"
	cronv3 "github.com/robfig/cron/v3"
	"github.com/sirupsen/logrus"
	"github.com/urfave/cli"
	"gorm.io/gorm"
)

var rows = []string{
	"name",
	"company",
	"position",
	"phone",
	"email",
	"savetime",
	"status",
	"kind",
}

var (
	DaySQL = "SELECT " + strings.Join(rows, ",") + " FROM user WHERE date(savetime) = date_sub(curdate(),interval 1 day)"
	MonSQL = "SELECT " + strings.Join(rows, ",") + " FROM user WHERE PERIOD_DIFF(DATE_FORMAT(NOW(),'%Y%m'),DATE_FORMAT(savetime,'%Y%m')) = 1"
)

const (
	dayFormat   = "2006-01-02"
	monthFormat = "2006-01"
)

type cronRunner struct {
	cronDB *gorm.DB
	cron   *cronv3.Cron
}

func StartCorn(ctx *cli.Context) error {
	var err error
	runner := cronRunner{
		cronDB: types.DBInstance.Session(&gorm.Session{
			PrepareStmt: true,
		}),
		cron: cronv3.New(),
	}
	logrus.Infof("Collect information start")

	if len(types.Email.DailyReceiver) != 0 {
		_, err = runner.cron.AddFunc(types.Email.CRONDaily, func() {
			logrus.Infof("Send Information For Day")

			d, err := time.ParseDuration("-24h")
			if err != nil {
				logrus.Errorf("Failed time parsed duration : %v", err)
			}
			yesterday := time.Now().Add(d).Format(dayFormat)
			today := time.Now().Format(dayFormat)
			count, filename := runner.DBSelect(DaySQL, yesterday)
			defer os.Remove(filename)

			headMessage := fmt.Sprintf("%s 用户信息", yesterday)
			bodyMessage := fmt.Sprintf("%s 08:00 ~ %s 08:00，一共有 %d 人下载了中文文档以及白皮书。\n", yesterday, today, count)
			toMap := map[string]string{}
			for _, addr := range types.Email.MonthlyReceiver {
				toMap[addr] = ""
			}
			if err := email.SendEmail(&email.Content{
				Subject:   headMessage,
				Body:      bodyMessage,
				Attach:    filename,
				To:        toMap,
				From:      types.Email.Sender,
				FromAlias: "Rancher Labs 中国",
			}); err != nil {
				logrus.Warnf("failed to send monthly email")
			}
		})
		if err != nil {
			return errors.Wrap(err, "Failed cron add function")
		}
	}

	if len(types.Email.MonthlyReceiver) != 0 {
		_, err = runner.cron.AddFunc(types.Email.CRONMonthly, func() {
			logrus.Infof("Send Information For Mon")

			now := time.Now()
			lastMonth := now.AddDate(0, -1, -now.Day()+1).Format(monthFormat)
			count, filename := runner.DBSelect(MonSQL, lastMonth)
			defer os.Remove(filename)

			headMessage := fmt.Sprintf("%s月 全部用户信息", lastMonth)
			bodyMessage := fmt.Sprintf("%s月一共有 %d 人下载了中文文档以及白皮书。\n", lastMonth, count)
			toMap := map[string]string{}
			for _, addr := range types.Email.MonthlyReceiver {
				toMap[addr] = ""
			}
			if err := email.SendEmail(&email.Content{
				Subject:   headMessage,
				Body:      bodyMessage,
				Attach:    filename,
				To:        toMap,
				From:      types.Email.Sender,
				FromAlias: "Rancher Labs 中国",
			}); err != nil {
				logrus.Warnf("failed to send monthly email")
			}
		})
		if err != nil {
			return errors.Wrap(err, "Failed cron add function")
		}
	}

	return runner.Start()
}

func (r *cronRunner) Start() error {
	r.cron.Start()
	return nil
}

func (r *cronRunner) DBSelect(sql, xlsxName string) (int, string) {
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
		"H1": "文档类型",
	}

	rows, err := r.cronDB.Raw(sql).Rows()
	if nil != err {
		logrus.Errorf("Failed to query : %v", err)
	}

	kinds := types.GetKindDescription()

	count := 0
	for rows.Next() {
		count++
		var user types.User
		row := count + 1
		err := rows.Scan(&user.Name, &user.Company, &user.Position, &user.Phone, &user.Email, &user.SaveTime, &user.Status, &user.Kind)
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
		data["H"+strconv.Itoa(row)] = kinds[user.Kind]
	}

	for k, v := range data {
		xlsx.SetCellValue("Sheet1", k, v)
	}
	xlsx.SetActiveSheet(index)
	filename := "/tmp/" + xlsxName + ".xlsx"
	if err = xlsx.SaveAs(filename); err != nil {
		logrus.Errorf("Failed to save excel : %v", err)
	}

	return count, filename
}
