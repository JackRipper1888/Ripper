package httpkit

import (
	"github.com/Ripper/errkit"
	"github.com/astaxie/beego/logs"
	"github.com/bitly/go-simplejson"
	"io/ioutil"
	"net/http"
)

func Get(url, confId, controllerIps string, param map[string]string) error {
	// 动态更新token
	logPre := "|#controller|Get|id=%s,controllerIps=%s|msg:%s"
	if len(param) > 0 {
		url += "?"
		paramStr := ""
		for k, v := range param {
			paramStr += "&" + k + "=" + v
		}
		url += paramStr[1:]
	}
	resp, err := http.Get(url)
	if err != nil {
		logs.Error("@"+logPre, confId, controllerIps, err.Error())
		return errkit.New(-1, err.Error())
	}
	if resp == nil {
		return errkit.New(-1, "请求体为空")
	}
	//返回的状态码
	statusCode := resp.StatusCode
	defer resp.Body.Close()
	if statusCode != 200 {
		logs.Error("@"+logPre, confId, controllerIps, statusCode)
		return errkit.New(-1, "请求体为空")
	}
	respBody, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		logs.Error("@"+logPre, confId, controllerIps, statusCode)
		return errkit.New(-1, "请求错误"+err.Error())
	}
	bodyJson, err := simplejson.NewJson(respBody)
	if err != nil {
		logs.Error("@"+logPre, confId, controllerIps, err.Error())
		return errkit.New(-1, "下行json解析错误"+err.Error())
	}
	code := bodyJson.Get("code").MustInt(0)
	msg := bodyJson.Get("msg").MustString("")
	if code != 1 {
		logs.Error("@"+logPre, confId, controllerIps, msg)
		return errkit.New(-1, msg)
	}
	logs.Info(logPre, confId, controllerIps, "suc")
	return nil
}
