package backstage

import (
	"bytes"
	"encoding/json"
	"io/ioutil"
	"net/http"
	"time"
)

//机器人信息
type mycontent struct {
	Content string `json:"content"`
}
type isAtdata struct {
	IsAtAll bool `json:"isAtAll"`
}

type senddata struct {
	Msgtype string    `json:"msgtype"`
	Text    mycontent `json:"text"`
	At      isAtdata  `json:"at"`
}

func Testsend(src string) {

	var data senddata
	data.Msgtype = "text"
	data.Text.Content = src
	data.At.IsAtAll = true
	b, err := json.Marshal(data)
	if err != nil {
		gLog.Error("json:%s", err)
		return
	}
	gLog.Info("send content :%s", string(b))

	body := bytes.NewBuffer(b)

	transport := http.Transport{
		DisableKeepAlives: true,
	}
	client := &http.Client{
		Transport:&transport,
		Timeout: 5 * time.Second,
	}

	res, err := client.Post("https://oapi.dingtalk.com/robot/send?access_token=2eb8253aae5237588004af68512f5fa6205fe2f6b4f08fc15d603287e0376d40", "application/json;charset=utf-8", body)
	if err != nil {
		gLog.Error("testsend error1:%s", err)
		return
	}
	result, err := ioutil.ReadAll(res.Body)
	res.Body.Close()
	if err != nil {
		gLog.Error("testsend error2:%s", err)
		return
	}
	gLog.Info("testsend :%s", result)
}
