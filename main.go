package main

import (
	"github.com/yhhaiua/engine/log"
	"github.com/yhhaiua/recharge/logic"
	"time"
)

var gLog = log.GetLogger()

func main() {

	gLog.Config("config/log4j.xml")
	if logic.Instance().LogicInit(){
		gLog.Info("recharge 启动成功")
	}else{
		gLog.Error("recharge 启动失败")
	}
	time.Sleep(3*time.Second)
}
