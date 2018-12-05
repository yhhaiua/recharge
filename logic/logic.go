package logic

import (
	"github.com/yhhaiua/engine/grouter"
	"github.com/yhhaiua/log4go"
	"github.com/yhhaiua/recharge/logic/config"
	"github.com/yhhaiua/recharge/logic/control"
	"net/http"
	"sync"
	"time"
)

var (
	instance *LogicSvr
	mu       sync.Mutex
)

//LogicSvr 服务器数据
type LogicSvr struct {
	config 	config.Config
	charge  control.ChargeControl
}

//Instance 实例化LogicSvr
func Instance() *LogicSvr {
	if instance == nil {
		mu.Lock()
		defer mu.Unlock()
		if instance == nil {
			instance = new(LogicSvr)
		}
	}
	return instance
}

//LogicInit 初始化
func (logic *LogicSvr) LogicInit() bool {
	if logic.config.ConfigInit(){
		logic.charge.Init(&logic.config)
		return  logic.routerInit()
	}
	return false
}

func (logic *LogicSvr) routerInit() bool{

	router := grouter.New()

	router.GET("/recharge", logic.charge.RechargeDeal)
	router.GET("/stopcharge", logic.charge.StopCharge)
	router.GET("/makeuporder", logic.charge.MakeUpOrder)

	log4go.Info("http监听开启%s", logic.config.Sport)
	log4go.Info("当前版本:v1.0.0")

	srv := &http.Server{
		ReadTimeout: 30 * time.Second,
		WriteTimeout: 30 * time.Second,
		Addr:logic.config.Sport,
		Handler : router,
	}

	err := srv.ListenAndServe()
	if err != nil {
		log4go.Error("http监听失败 %s", err)
		return false
	}
	return true
}