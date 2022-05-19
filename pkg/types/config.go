package types

import (
	"fmt"
	"io/ioutil"
	"strings"

	"github.com/imdario/mergo"
	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
	"github.com/urfave/cli"
	"gopkg.in/yaml.v2"
)

var (
	Aliyun   = aliyun{}
	DB       = db{}
	Email    = email{}
	Register = register{
		Kinds: map[string]string{},
	}
	Config = config{
		DB:       &DB,
		SMTP:     &Email,
		Aliyun:   &Aliyun,
		Register: &Register,
	}
	ErrMissingRequired = errors.New("required field is missing")
)

type config struct {
	Debug     bool            `json:"debug,omitempty" yaml:"debug"`
	DryRun    bool            `json:"dryRun,omitempty" yaml:"dryRun"`
	Port      int             `json:"port,omitempty" yaml:"port"`
	Kinds     map[string]Kind `json:"kinds,omitempty" yaml:"kinds"`
	SMTP      *email          `json:"smtp,omitempty" yaml:"smtp"`
	DB        *db             `json:"db,omitempty" yaml:"db"`
	Aliyun    *aliyun         `json:"aliyun,omitempty" yaml:"aliyun"`
	Documents []Document      `json:"documents,omitempty" yaml:"documents"`
	Register  *register       `json:"register,omitempty" yaml:"register"`
	CodeClean bool            `json:"codeClean,omitempty" yaml:"codeClean"`
}

type email struct {
	User            string   `json:"user,omitempty" yaml:"user"`
	Password        string   `json:"password,omitempty" yaml:"password"`
	Endpoint        string   `json:"endpoint,omitempty" yaml:"endpoint"`
	Port            int      `json:"port,omitempty" yaml:"port"`
	Sender          string   `json:"sender,omitempty" yaml:"sender"`
	DailyReceiver   []string `json:"dailyReceiver,omitempty" yaml:"dailyReceiver"`
	MonthlyReceiver []string `json:"monthlyReceiver,omitempty" yaml:"monthlyReceiver"`
	CRONDaily       string   `json:"cronDaily,omitempty" yaml:"cronDaily"`
	CRONMonthly     string   `json:"cronMonthly,omitempty" yaml:"cronMonthly"`
}

func (e *email) Validate() error {
	if e.User != "" &&
		e.Password != "" &&
		e.Endpoint != "" &&
		e.Port != 0 &&
		e.Sender != "" &&
		e.CRONDaily != "" &&
		e.CRONMonthly != "" {
		return nil
	}
	return ErrMissingRequired
}

type register struct {
	SenderName string            `json:"senderName,omitempty" yaml:"senderName"`
	Subject    string            `json:"subject,omitempty" yaml:"subject"`
	Receivers  []string          `json:"receivers,omitempty" yaml:"receivers"`
	Template   string            `json:"template,omitempty" yaml:"template"`
	Kinds      map[string]string `json:"kinds,omitempty" yaml:"kinds"`
}

func (r *register) Validate() error {
	if len(r.Receivers) == 0 {
		return nil
	}
	if r.Subject != "" && r.Template != "" {
		return nil
	}
	return ErrMissingRequired
}

type db struct {
	Kind     string `json:"kind,omitempty" yaml:"kind"`
	HostIP   string `json:"hostIp,omitempty" yaml:"hostIp"`
	Port     int    `json:"port,omitempty" yaml:"port"`
	Username string `json:"username,omitempty" yaml:"username"`
	Password string `json:"password,omitempty" yaml:"password"`
	Name     string `json:"name,omitempty" yaml:"name"`
}

func (e *db) Validate() error {
	if e.HostIP != "" &&
		e.Port >= 0 && e.Port <= 65535 &&
		e.Username != "" &&
		e.Password != "" &&
		e.Name != "" {
		return nil
	}
	return errors.Wrap(ErrMissingRequired, "db fields are missed")
}

type aliyun struct {
	Region                  string `json:"region,omitempty" yaml:"region"`
	AccessKey               string `json:"accessKey,omitempty" yaml:"accessKey"`
	AccessSecret            string `json:"accessSecret,omitempty" yaml:"accessSecret"`
	OSSScheme               string `json:"ossScheme,omitempty" yaml:"ossScheme"`
	OSSBucket               string `json:"ossBucket,omitempty" yaml:"ossBucket"`
	OSSPathPrefix           string `json:"ossPathPrefix,omitempty" yaml:"ossPathPrefix"`
	OSSSignURLExpiresSecond int64  `json:"ossSignURLExpiresSecond,omitempty" yaml:"ossSignURLExpiresSecond"`
	OSSCnameEndpoint        string `json:"ossCnameEndpoint,omitempty" yaml:"ossCnameEndpoint"`
	SMSSignName             string `json:"smsSignName,omitempty" yaml:"smsSignName"`
	SMSTemplateCode         string `json:"smsTempateCode,omitempty" yaml:"smsTempateCode"`
}

func (e *aliyun) Validate() error {
	if e.Region != "" &&
		e.AccessKey != "" &&
		e.AccessSecret != "" &&
		e.OSSScheme != "" &&
		e.OSSBucket != "" &&
		e.OSSSignURLExpiresSecond > 0 &&
		e.SMSSignName != "" &&
		e.SMSTemplateCode != "" {
		return nil
	}
	if e.OSSScheme != "https" && e.OSSScheme != "http" {
		return errors.New("aliyun oss scheme is not valid")
	}
	return errors.Wrap(ErrMissingRequired, "")
}

type Kind struct {
	Header      string `json:"header,omitempty" yaml:"header"`
	Footer      string `json:"footer,omitempty" yaml:"footer"`
	Subject     string `json:"subject,omitempty" yaml:"subject"`
	SenderName  string `json:"senderName,omitempty" yaml:"senderName"`
	Description string `json:"description,omitempty" yaml:"description"`
}

func (e *Kind) Validate() error {
	if e.Header != "" &&
		e.Footer != "" &&
		e.Subject != "" &&
		e.SenderName != "" {
		return nil
	}
	return ErrMissingRequired
}

type Document struct {
	Name     string   `json:"name,omitempty" yaml:"name"`
	Title    string   `json:"title,omitempty" yaml:"title"`
	URL      string   `json:"url,omitempty" yaml:"url"`
	PWD      string   `json:"pwd,omitempty" yaml:"pwd"`
	Filename string   `json:"filename,omitempty" yaml:"filename"`
	Kind     []string `json:"kind,omitempty" yaml:"kind"`

	OSSPathPrefixOverride string `json:"ossPathPrefixOverride,omitempty" yaml:"ossPathPrefixOverride"`
}

func (d Document) Validate() error {
	switch {
	case d.Name == "":
		return fmt.Errorf("Document name is missing")
	case d.URL != "": // Using direct URL for download
		return nil
	case d.Filename == "":
		return fmt.Errorf("the filename for Document %s is missing", d.Name)
	case OSSBucketEndpoint != "": // Using oss to ge
		exists, err := IsObjectExists(d.OSSPathPrefixOverride, d.Filename)
		if err != nil {
			return err
		}
		if !exists {
			return fmt.Errorf("the Document file %s does not exists in OSS ", d.Name)
		}
		return nil
	}
	return nil
}

// When generating Document download link for email,
// we will use following priority for configuration.
// 1. Direct URL from configuration with password.
// 2. OSS gLobal Configuration.
func (d Document) GetLine(isPublic bool) string {
	var fileURL string
	switch {
	case d.URL != "":
		fileURL = d.URL
		if d.PWD != "" {
			fileURL += "     访问密码： " + d.PWD
		}
	case OSSBucketEndpoint != "": // using oss as file hoster
		fileURL = GetOSSFileDownloadURL(d.Filename, d.OSSPathPrefixOverride, isPublic)
	}

	title := d.Title
	if title == "" {
		title = d.Name
	}
	return fmt.Sprintf("%s:     %s\n", title, fileURL)
}

func (c *config) Validate() error {
	if err := c.DB.Validate(); err != nil {
		return err
	}
	var supported []string
	for k, v := range c.Kinds {
		if err := v.Validate(); err != nil {
			return errors.Wrapf(err, "failed to validate kind %s configuration", k)
		}
		supported = append(supported, k)
	}
	if err := c.SMTP.Validate(); err != nil {
		return err
	}
	if err := c.Aliyun.Validate(); err != nil {
		return err
	}
	if err := c.Register.Validate(); err != nil {
		return err
	}
	for _, doc := range c.Documents {
		for _, kind := range doc.Kind {
			if _, ok := c.Kinds[kind]; !ok {
				return fmt.Errorf("kind %s of docs %s is in supported kinds %s", kind, doc.Name, strings.Join(supported, ","))
			}
		}
	}
	return nil
}

func GetKindDescription() map[string]string {
	rtn := map[string]string{}
	for k, v := range Config.Kinds {
		if v.Description == "" {
			rtn[k] = k
		} else {
			rtn[k] = v.Description
		}
	}
	return rtn
}

func MergeConfig(ctx *cli.Context) error {
	filename := ctx.GlobalString("config-file")
	data, err := ioutil.ReadFile(filename)
	if err != nil {
		return err
	}

	Config.SMTP.DailyReceiver = ctx.StringSlice("smtp-rancher-to-day")
	Config.SMTP.MonthlyReceiver = ctx.StringSlice("smtp-rancher-to-mon")

	fileConfig := config{}
	if err := yaml.Unmarshal(data, &fileConfig); err != nil {
		return err
	}

	if err := mergo.Map(&Config, fileConfig); err != nil {
		return err
	}

	return nil
}

func Validate() error {
	return Config.Validate()
}

func IsKindValid(requester, kind string) error {
	if kind == "" {
		return errors.New("query kind is empty")
	}
	var ok bool
	switch requester {
	case "pdf":
		_, ok = Config.Kinds[kind]
	case "register":
		_, ok = Config.Register.Kinds[kind]
	default:
		_, ok1 := Config.Kinds[kind]
		_, ok2 := Config.Register.Kinds[kind]
		ok = ok1 || ok2
	}
	if ok {
		return nil
	}
	return fmt.Errorf("kind %s is not valid", kind)
}

func GetConfigCommand() cli.Command {
	cmd := cli.Command{
		Name:   "config",
		Usage:  "Print out the merged configuration",
		Action: printConfig,
		Flags:  GetFlags(),
	}
	return cmd
}

func printConfig(ctx *cli.Context) error {
	if err := MergeConfig(ctx); err != nil {
		return err
	}

	if ctx.GlobalBool("debug") {
		logrus.SetLevel(logrus.DebugLevel)
	}
	data, err := yaml.Marshal(Config)
	if err != nil {
		return errors.Wrap(err, "failed to convert config to yaml")
	}
	println(string(data))
	return nil
}

func SetRunStatus(debug, dryRun bool) {
	Config.DryRun = dryRun
	Config.Debug = debug || dryRun
	if Config.Debug {
		logrus.SetLevel(logrus.DebugLevel)
	}
}
