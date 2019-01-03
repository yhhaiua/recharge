package control

import (
	"crypto/md5"
	"encoding/hex"
	"github.com/yhhaiua/engine/grouter"
	"github.com/yhhaiua/recharge/logic/model"
	"net/http"
)

//补单
func (contol *ChargeControl)MakeUpOrder(w http.ResponseWriter, r *http.Request, _ grouter.Params)  {
	operatorid := r.FormValue("operatorid")
	sign := r.FormValue("sign")
	stime := r.FormValue("time")
	itemid := r.FormValue("itemid")
	playerid := r.FormValue("playerid")

	if operatorid != contol.config.Operatorid{
		gLog.Error("makeUpOrder operatorid error : me:%s,client:%s,操作人:%s",contol.config.Operatorid,operatorid,model.GetUserIp(r))
		contol.send(w,-100,"operatorid error")
		return
	}
	//operatorid+stime+itemid+playerid+key
	md5str := operatorid + stime +itemid+playerid+contol.config.Rechargekey
	md5Ctx := md5.New()
	md5Ctx.Write([]byte(md5str))
	cipherStr := md5Ctx.Sum(nil)
	mysigon := hex.EncodeToString(cipherStr)
	if(mysigon != sign){
		gLog.Error("makeUpOrder md5 error : me:%s,client:%s,操作人:%s",mysigon,sign,model.GetUserIp(r))
		contol.send(w,-100,"md5 error")
		return
	}
	gLog.Info("makeUpOrder success : 操作人:%s",model.GetUserIp(r))
	switch operatorid {
	case "1":
		//玩吧渠道
		ret,value := contol.wanba.MakeUpOrder(w,r)
		contol.send(w,ret,value)
	default:
		gLog.Error("错误渠道请求:%s",operatorid)
		contol.send(w,-100," operatorid no have error")
	}
}
