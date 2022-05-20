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
	"city",
	"company",
	"department",
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
	dayFormat      = "2006-01-02"
	monthFormat    = "2006-01"
	codeCleanCron  = "0 * * * *"
	codeDeleteCron = "30 0 * * *"
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
			if !types.Config.Debug {
				defer os.Remove(filename)
			}

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
			if !types.Config.Debug {
				defer os.Remove(filename)
			}

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
	if types.Config.CodeClean {
		logrus.Infof("adding code clean cron")
		// code cron
		if _, err := runner.cron.AddFunc(codeCleanCron, func() {
			result := runner.cronDB.Model(&types.Code{}).Where("state = ? and request_time < "+getDateFromNowString("1 hour"), "active").Update("state", "used")
			if result.Error != nil {
				logrus.Warnf("failed to clean useless code, %v", err)
				return
			}
			logrus.Infof("%d codes have been cleaned", result.RowsAffected)
		}); err != nil {
			return errors.Wrap(err, "failed to add code clean cron function")
		}

		logrus.Infof("adding code delete cron")
		if _, err := runner.cron.AddFunc(codeDeleteCron, func() {
			result := runner.cronDB.Delete(&types.Code{}, "request_time < "+getDateFromNowString("1 day"))
			if result.Error != nil {
				logrus.Warnf("failed to delete useless code, %v", err)
				return
			}
			logrus.Infof("%d codes have been deleted", result.RowsAffected)
		}); err != nil {
			return errors.Wrap(err, "failed to add code delete cron function")
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
	data := tableWithHeader(0, 0, []string{
		"名字",
		"城市",
		"公司",
		"部门",
		"职位",
		"手机号",
		"电子邮箱",
		"保存时间",
		"邮件已发送",
		"访问功能",
		"文档类型",
	})

	rows, err := r.cronDB.Raw(sql).Rows()
	if nil != err {
		logrus.Errorf("Failed to query : %v", err)
	}

	kinds := types.GetKindDescription()

	count := 0
	for rows.Next() {
		count++
		var user types.User
		err := r.cronDB.ScanRows(rows, &user)
		if err != nil {
			logrus.Errorf("Failed rows scan : %v", err)
		}
		function := "pdf"
		if _, ok := kinds[user.Kind]; !ok {
			function = types.Config.Register.Kinds[user.Kind]
		}
		setRow(data, 0, count, []string{
			user.Name,
			user.City,
			user.Company,
			user.Department,
			user.Position,
			user.Phone,
			user.Email,
			user.SaveTime.Format("2006-01-02 15:04:05"),
			strconv.FormatBool(user.Status),
			function,
			kinds[user.Kind],
		})
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

// getTableHeader row and column begins with 0. Only support input len <= 26
func tableWithHeader(row, column int, headers []string) map[string]string {
	columnBegin := 'A'
	rtn := map[string]string{}
	currentColumn := column
	for _, header := range headers {
		key := fmt.Sprintf("%c%d", columnBegin+rune(currentColumn), row+1)
		rtn[key] = header
		currentColumn++
	}
	return rtn
}

func setRow(table map[string]string, column, row int, datas []string) {
	columnBegin := 'A'
	currentColumn := column
	for _, data := range datas {
		key := fmt.Sprintf("%c%d", columnBegin+rune(currentColumn), row+1)
		table[key] = data
		currentColumn++
	}
}

func getDateFromNowString(interval string) string {
	if types.DB.Kind == "pgsql" {
		return fmt.Sprintf("NOW() - interval '%s'", interval)
	}
	// default to mysql
	return fmt.Sprintf("DATE_SUB(NOW(),INTERVAL %s)", interval)
}
