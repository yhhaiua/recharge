package control

import (
	"encoding/json"
	"fmt"
	"github.com/yhhaiua/engine/grouter"
	"github.com/yhhaiua/log4go"
	"github.com/yhhaiua/recharge/logic/config"
	"github.com/yhhaiua/recharge/logic/control/wanba"
	"github.com/yhhaiua/recharge/logic/model"
	"net/http"
)

type ChargeControl struct {
	wanba wanba.WanBa
	config *config.Config
}

func (contol *ChargeControl)Init(config *config.Config)  {
	contol.config = config
	contol.wanba.Init(config)
}
// RetInfo 错误返回
type RetInfo struct {
	Ret int `json:"ret"`
	Msg string `json:"msg"`
}

func (contol *ChargeControl) RechargeDeal(w http.ResponseWriter, r *http.Request, _ grouter.Params) {

	operatorid := r.FormValue("operatorid")
	if operatorid != contol.config.Operatorid{
		contol.send(w,-100,"operatorid error")
		return
	}
	if model.StopRecharge{
		contol.send(w,-100,"stop charge")
		return
	}
	switch operatorid {
	case "1":
		//玩吧渠道
		ret,value := contol.wanba.RechargeDeal(w,r)
		contol.send(w,ret,value)
	default:
		log4go.Error("错误渠道请求:%s",operatorid)
		contol.send(w,-100," operatorid no have error")
	}
}



func (contol *ChargeControl)send(w http.ResponseWriter,ret int,msg string)  {
	w.Header().Set("Access-Control-Allow-Origin", "*")             //允许访问所有域
	var info RetInfo
	info.Ret = ret
	info.Msg = msg
	Message, err := json.Marshal(info)

	if err == nil {
		fmt.Fprintf(w, "%s", Message)
	}
}