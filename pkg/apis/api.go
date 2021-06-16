package apis

import (
	"fmt"
	"math/rand"
	"strings"
	"time"

	"github.com/patrickmn/go-cache"
	"github.com/sirupsen/logrus"

	"github.com/cnrancher/pdf-sender/pkg/types"
	restful "github.com/emicklei/go-restful/v3"
)

var CacheClient = cache.New(10*time.Minute, 20*time.Minute)

func RegisterAPIs() *restful.Container {
	container := restful.NewContainer()

	cors := restful.CrossOriginResourceSharing{
		AllowedHeaders: []string{"Content-Type", "Accept"},
		AllowedMethods: []string{"OPTIONS", "POST"},
		AllowedDomains: []string{"*"},
		CookiesAllowed: false,
		Container:      container}

	docxWs := new(restful.WebService).Path("/").Consumes(restful.MIME_JSON).Produces(restful.MIME_JSON)

	container.Add(docxWs)

	container.Filter(cors.Filter)

	container.Filter(container.OPTIONSFilter)

	docxWs.Route(
		docxWs.POST("/sendCode").
			To(sendCode).
			Param(docxWs.QueryParameter("kind", "request kind(pdf or ent)").Required(true)).Filter(kindFilter).
			Reads(types.User{}))
	docxWs.Route(
		docxWs.POST("/sendEmail").
			To(sendEmail).
			Param(docxWs.QueryParameter("kind", "request kind(pdf or ent)").Required(true)).Filter(kindFilter).
			Reads(types.User{}))
	return container

}

func sendCode(req *restful.Request, resp *restful.Response) {
	var user types.User
	if err := req.ReadEntity(&user); err != nil {
		logrus.Errorf("Get user err:%v", err)
		err = resp.WriteErrorString(400, err.Error())
		if err != nil {
			logrus.Errorf("Failed to write error string err:%v", err)
		}
		return
	}

	if user.Phone == "" {
		logrus.Errorf("Phone number cannot be empty")
		if err := resp.WriteErrorString(400, "手机号不能为空"); err != nil {
			logrus.Errorf("Failed to write error string err:%v", err)
		}
		return
	}

	user.Kind = req.QueryParameter("kind")
	user.Code = GenValidateCode(4)

	logrus.Debugf(user.Code)

	CacheClient.Set(user.Code, &user, cache.DefaultExpiration)

	if err := user.Send(); err != nil {
		logrus.Errorf("Send SMS to phone %s err:%v", user.Phone, err)
		err = resp.WriteErrorString(400, err.Error())
		if err != nil {
			logrus.Errorf("Failed to write error string err:%v", err)
		}
		return
	}

	if err := resp.WriteErrorString(200, "验证码发送成功"); err != nil {
		logrus.Errorf("Failed to write error string err:%v", err)
	}
}

func sendEmail(req *restful.Request, resp *restful.Response) {
	var user types.User
	if err := req.ReadEntity(&user); err != nil {
		logrus.Errorf("Get user err:%v", err)
		err = resp.WriteErrorString(400, err.Error())
		if err != nil {
			logrus.Errorf("Failed to write error string err:%v", err)
		}
		return
	}

	if err := user.Validate(); err != nil {
		logrus.Errorf("Validate user err:%v", err)
		if err := resp.WriteErrorString(400, err.Error()); err != nil {
			logrus.Errorf("Failed to write error string err:%v", err)
		}
		return
	}
	user.Kind = req.QueryParameter("kind")

	entry, found := CacheClient.Get(user.Code)
	if found {
		logrus.Infof("%s get the code in cache successful", user.Phone)
	} else {
		logrus.Errorf("%s failed to get the code in the cache", user.Phone)
		if err := resp.WriteErrorString(400, "验证码超时或不存在"); err != nil {
			logrus.Errorf("Failed to write error string err:%v", err)
		}
		return
	}

	cachedUser := entry.(*types.User)

	if !user.Compare(cachedUser) {
		logrus.Errorf("手机号 %s 校验验证码错误", user.Phone)
		if err := resp.WriteErrorString(400, "验证码错误"); err != nil {
			logrus.Errorf("Failed to write error string err:%v", err)
		}
		return
	}

	go SendEmail(&user)
}

func GenValidateCode(width int) string {
	numeric := [9]byte{1, 2, 3, 4, 5, 6, 7, 8, 9}
	r := len(numeric)
	rand.Seed(time.Now().UnixNano())

	var sb strings.Builder
	for i := 0; i < width; i++ {
		fmt.Fprintf(&sb, "%d", numeric[rand.Intn(r)])
	}

	return sb.String()
}

func kindFilter(req *restful.Request, resp *restful.Response, next *restful.FilterChain) {
	if err := types.IsKindValid(req.QueryParameter("kind")); err != nil {
		if err = resp.WriteErrorString(400, err.Error()); err != nil {
			logrus.Warnf("failed to write error response, %v", err)
		}
		return
	}
	next.ProcessFilter(req, resp)
}
