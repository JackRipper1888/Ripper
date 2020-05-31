package httpkit

import (
	"bytes"
	"encoding/json"
	"github.com/Ripper/errkit"
	"github.com/astaxie/beego/logs"
	"github.com/bitly/go-simplejson"
	"io/ioutil"
	"net/http"
	"time"
)

func Get(url, errName string, param map[string]string) error {
	// 动态更新token
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
		logs.Error(errName, err.Error())
		return errkit.New(-1, err.Error())
	}
	if resp == nil {
		return errkit.New(-1, "请求体为空")
	}
	//返回的状态码
	statusCode := resp.StatusCode
	defer resp.Body.Close()
	if statusCode != 200 {
		logs.Error(errName, statusCode)
		return errkit.New(-1, "请求体为空")
	}
	respBody, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		logs.Error(errName, statusCode)
		return errkit.New(-1, "请求错误"+err.Error())
	}
	bodyJson, err := simplejson.NewJson(respBody)
	if err != nil {
		logs.Error(errName, err.Error())
		return errkit.New(-1, "下行json解析错误"+err.Error())
	}
	code := bodyJson.Get("code").MustInt(0)
	msg := bodyJson.Get("msg").MustString("")
	if code != 1 {
		logs.Error(errName, msg)
		return errkit.New(-1, msg)
	}
	logs.Info(errName, "suc")
	return nil
}

func Post(url, errNema string, Body map[string]interface{}) error {
	//logPre := "|#controller|Post|id=%s,controllerIps=%s|msg:%s"
	// 超时时间：5秒
	client := &http.Client{Timeout: 5 * time.Second}
	jsonStr, _ := json.Marshal(Body)
	resp, err := client.Post(url, "application/x-www-form-urlencoded", bytes.NewBuffer(jsonStr))
	if err != nil {
		logs.Error(errNema, err.Error())
		return errkit.New(-1, errNema+err.Error())
	}
	if resp == nil {
		return errkit.New(-1, "请求体为空")
	}
	defer resp.Body.Close()
	statusCode := resp.StatusCode
	if statusCode != 200 {
		logs.Error(errNema, statusCode)
		return errkit.New(-1, "请求错误")
	}
	respBody, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		logs.Error(errNema, statusCode)
		return errkit.New(-1, err.Error())
	}
	bodyJson, err := simplejson.NewJson(respBody)
	if err != nil {
		logs.Error(errNema, err.Error())
		return errkit.New(-1, err.Error())
	}
	code := bodyJson.Get("code").MustInt(0)
	msg := bodyJson.Get("msg").MustString("")
	if code != 1 {
		logs.Error(errNema, msg)
		return errkit.New(-1, msg)
	}
	return nil
}
