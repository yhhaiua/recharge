package main

import (
	"github.com/yhhaiua/log4go"
	"github.com/yhhaiua/recharge/logic"
	"time"
)

func main() {

	log4go.LoadConfiguration("config/log4j.xml")
	if logic.Instance().LogicInit(){
		log4go.Info("recharge 启动成功")
	}else{
		log4go.Error("recharge 启动失败")
	}
	time.Sleep(3*time.Second)
}
