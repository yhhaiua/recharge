package model

import (
	"github.com/astaxie/beego/orm"
	"strconv"
	"time"
)

type RechargeLog struct {
	OrderId string		`orm:"pk"`
	Time    time.Time 	`orm:"auto_now_add;type(datetime)"`
	Money int
	GoodsId int
	PlayerId int64
	ErrorInfo string	`orm:"type(text)"`
	Operatorid string
}

func (u *RechargeLog) TableName() string {
	return "recharge_log"
}

//Insert 插入一条日志
func (charge *RechargeLog)Insert(OrderId ,Money,GoodsId ,PlayerId ,ErrorInfo,Operatorid string) error  {

	money,_:= strconv.Atoi(Money)
	goodsid,_:= strconv.Atoi(GoodsId)
	playerid, _ := strconv.ParseInt(PlayerId, 10, 64)

	charge.OrderId = OrderId
	charge.Money = money
	charge.GoodsId = goodsid
	charge.PlayerId = playerid
	charge.ErrorInfo = ErrorInfo
	charge.Operatorid = Operatorid
	charge.Time = time.Now()
	//orm.Debug = true
	o := orm.NewOrm()
	_, err := o.Insert(charge)
	return err
}
//Update 更新一条日志
func (charge *RechargeLog)Update(OrderId ,Money,GoodsId ,PlayerId ,ErrorInfo,Operatorid string) error  {

	money,_:= strconv.Atoi(Money)
	goodsid,_:= strconv.Atoi(GoodsId)
	playerid, _ := strconv.ParseInt(PlayerId, 10, 64)

	charge.OrderId = OrderId
	charge.Money = money
	charge.GoodsId = goodsid
	charge.PlayerId = playerid
	charge.ErrorInfo = ErrorInfo
	charge.Operatorid = Operatorid
	charge.Time = time.Now()
	//orm.Debug = true
	o := orm.NewOrm()
	_, err := o.Update(charge)
	return err
}