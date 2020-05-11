package main

import (
	"encoding/json"
	"fmt"
	"github.com/astaxie/beego/logs"
	"github.com/bitly/go-simplejson"
	"io/ioutil"
	"net"
	"net/http"
	"os"
	"tracker/utils"
)

//go语言的interface类型 本质是一个指针类型
type AnimalIF interface {
	Sleep()        //让动物睡觉
	Color() string //获取动物颜色
	Type() string  //获取动物类型
}

// ===========================
type cat struct {
	color string
}

func (this *cat) Sleep() {
	fmt.Println("Cat is Sleep")
}

func (this *cat) Color() string {
	return this.color
}

func (this *cat) Type() string {
	return "Cat"
}

// =================
type dog struct {
	color string
}

func (this *dog) Sleep() {
	fmt.Println("Dog is Sleep")
}

func (this *dog) Color() string {
	return this.color
}

func (this *dog) Type() string {
	return "Dog"
}

//创建一个工厂
func Factory(kind string) AnimalIF {
	switch kind {
	case "dog":
		//生产一个狗
		return &dog{"yellow"}
	case "cat":
		//生产一个猫
		cat := new(cat)
		cat.color = "green"
		return cat
	default:
		fmt.Println("kind is wrong")
		return nil
	}
}

func showAnimal(animal AnimalIF) {
	animal.Sleep()              //多态现象
	fmt.Println(animal.Color()) //多态
	fmt.Println(animal.Type())  //多态
}

func main() {
	//var animal AnimalIF
	//animal = Factory("cat")
	//count := 0
	//for {
	//	count++
	//	Get(
	//		"http://183.60.143.82:3030/peer/select/tracker","0001","0008",map[string]string{
	//		"streamid":"us",
	//	})
	//	if count == 10000{
	//		return
	//	}
	//}
	//animal.Sleep()
	//showAnimal(animal)
	//t2,_ := time.Parse("2016-01-02 15:05:05", "2018-04-23 00:00:06")
	//t1 := time.Date(2018, 1, 2, 15, 5, 0, 0, time.Local)
	//t2 := time.Date(2018, 1, 2, 15, 0, 0, 0, time.Local)
	//
	//logs.Error(t1.Unix()%(24*3600)%300,t2.Unix()%300)

}

func peerclient() {
	logs.Info("peerclient() start")
	conn, err := net.Dial("udp", "183.60.143.82:8001")
	if err != nil {
		os.Exit(1)
		return
	}
	defer conn.Close()
	msgNum := 0
	for msgNum < 10000 {
		msgNum++
		//select {
		//case <- time.After(1*time.Second):
		data := map[string]interface{}{
			"cmd":  "login",
			"data": "12371892",
		}
		body, _ := json.Marshal(data)
		conn.Write([]byte(body))

		var msg [1024]byte
		_, err = conn.Read(msg[0:])
		if err != nil {
			fmt.Println(err)
			return
		}
		fmt.Println("msgcount:", msgNum)
	}
}

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
		return utils.New(-1, err.Error())
	}
	if resp == nil {
		return utils.New(-1, "请求体为空")
	}
	//返回的状态码
	statusCode := resp.StatusCode
	defer resp.Body.Close()
	if statusCode != 200 {
		logs.Error("@"+logPre, confId, controllerIps, statusCode)
		return utils.New(-1, "请求体为空")
	}
	respBody, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		logs.Error("@"+logPre, confId, controllerIps, statusCode)
		return utils.New(-1, "请求错误"+err.Error())
	}
	bodyJson, err := simplejson.NewJson(respBody)
	if err != nil {
		logs.Error("@"+logPre, confId, controllerIps, err.Error())
		return utils.New(-1, "下行json解析错误"+err.Error())
	}
	code := bodyJson.Get("code").MustInt(0)
	msg := bodyJson.Get("msg").MustString("")
	if code != 1 {
		logs.Error("@"+logPre, confId, controllerIps, msg)
		return utils.New(-1, msg)
	}
	logs.Info(logPre, confId, controllerIps, "suc")
	return nil
}
