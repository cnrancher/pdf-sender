package apis

import (
	"fmt"
	"io/ioutil"
	"math/rand"
	"strings"
	"time"

	"github.com/patrickmn/go-cache"
	"github.com/sirupsen/logrus"

	"github.com/emicklei/go-restful"
	"github.com/rancher/pdf-sender/pkg/types"
)

func RegisterAPIs() *restful.Container {
	container := restful.NewContainer()

	cors := restful.CrossOriginResourceSharing{
		AllowedHeaders: []string{"Content-Type", "Accept"},
		AllowedMethods: []string{"OPTIONS", "POST"},
		AllowedDomains: []string{},
		CookiesAllowed: false,
		Container:      container}

	docxWs := new(restful.WebService).Path("/").Consumes(restful.MIME_JSON).Produces(restful.MIME_JSON)

	container.Add(docxWs)

	container.Filter(cors.Filter)

	container.Filter(container.OPTIONSFilter)

	docxWs.Route(docxWs.POST("/sendCode").To(sendCode))
	docxWs.Route(docxWs.POST("/sendEmail").To(sendEmail))
	return container

}

func sendCode(req *restful.Request, resp *restful.Response) {
	bodyData, err := ioutil.ReadAll(req.Request.Body)
	if err != nil {
		logrus.Errorf("Read req body err:%v", err)
		err = resp.WriteErrorString(400, err.Error())
		if err != nil {
			logrus.Errorf("Failed to write error string err:%v", err)
		}
		return
	}

	user, err := types.New(bodyData)
	if err != nil {
		logrus.Errorf("Get user err:%v", err)
		err = resp.WriteErrorString(400, err.Error())
		if err != nil {
			logrus.Errorf("Failed to write error string err:%v", err)
		}
		return
	}

	if user.Phone == "" {
		logrus.Errorf("Phone number cannot be empty")
		err = resp.WriteErrorString(400, "手机号不能为空")
		if err != nil {
			logrus.Errorf("Failed to write error string err:%v", err)
		}
		return
	}

	user.Code = GenValidateCode(4)

	logrus.Infof(user.Code)

	CacheClient.Set(user.Phone, user.Code, cache.DefaultExpiration)

	err = user.Send()
	if err != nil {
		logrus.Errorf("Send SMS err:%v", err)
		err = resp.WriteErrorString(400, err.Error())
		if err != nil {
			logrus.Errorf("Failed to write error string err:%v", err)
		}
		return
	}

	err = resp.WriteErrorString(200, "验证码发送成功")
	if err != nil {
		logrus.Errorf("Failed to write error string err:%v", err)
	}

	return
}

func sendEmail(req *restful.Request, resp *restful.Response) {
	bodyData, err := ioutil.ReadAll(req.Request.Body)
	if err != nil {
		logrus.Errorf("Read req body err:%v", err)
		err = resp.WriteErrorString(400, err.Error())
		if err != nil {
			logrus.Errorf("Failed to write error string err:%v", err)
		}
		return
	}

	user, err := types.New(bodyData)
	if err != nil {
		logrus.Errorf("Get user err:%v", err)
		err = resp.WriteErrorString(400, err.Error())
		if err != nil {
			logrus.Errorf("Failed to write error string err:%v", err)
		}
		return
	}

	err = user.Validate()
	if err != nil {
		logrus.Errorf("Validate user err:%v", err)
		err = resp.WriteErrorString(400, err.Error())
		if err != nil {
			logrus.Errorf("Failed to write error string err:%v", err)
		}
		return
	}

	code, found := CacheClient.Get(user.Phone)
	if found {
		logrus.Infof("%s get the code in cache successful", user.Phone)
	} else {
		logrus.Errorf("%s failed to get the code in the cache", user.Phone)
		err = resp.WriteErrorString(400, "验证码超时或不存在")
		if err != nil {
			logrus.Errorf("Failed to write error string err:%v", err)
		}
		return

	}

	if user.Code != code {
		logrus.Errorf("验证码错误")
		err = resp.WriteErrorString(400, "验证码错误")
		if err != nil {
			logrus.Errorf("Failed to write error string err:%v", err)
		}
		return
	}

	go SendEmail(user)

	return
}

func GenValidateCode(width int) string {
	numeric := [10]byte{0, 1, 2, 3, 4, 5, 6, 7, 8, 9}
	r := len(numeric)
	rand.Seed(time.Now().UnixNano())

	var sb strings.Builder
	for i := 0; i < width; i++ {
		fmt.Fprintf(&sb, "%d", numeric[rand.Intn(r)])
	}

	return sb.String()
}
