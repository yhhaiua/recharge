package control

import (
	"crypto/md5"
	"encoding/hex"
	"github.com/yhhaiua/engine/grouter"
	"github.com/yhhaiua/recharge/logic/model"
	"net/http"
)

//停止充值
func (contol *ChargeControl) StopCharge(w http.ResponseWriter, r *http.Request, _ grouter.Params)  {

	operatorid := r.FormValue("operatorid")
	sign := r.FormValue("sign")
	stime := r.FormValue("time")
	stype := r.FormValue("type")

	if operatorid != contol.config.Operatorid{
		gLog.Error("stopCharge operatorid error : me:%s,client:%s,操作人:%s",contol.config.Operatorid,operatorid,model.GetUserIp(r))
		contol.send(w,-100,"operatorid error")
		return
	}
	//operatorid+stime+stype+key
	md5str := operatorid + stime +stype+contol.config.Rechargekey
	md5Ctx := md5.New()
	md5Ctx.Write([]byte(md5str))
	cipherStr := md5Ctx.Sum(nil)
	mysigon := hex.EncodeToString(cipherStr)
	if(mysigon != sign){
		gLog.Error("stopCharge md5 error : me:%s,client:%s,操作人:%s",mysigon,sign,model.GetUserIp(r))
		contol.send(w,-100,"md5 error")
		return
	}
	if stype == "1"{
		model.StopRecharge = true
		gLog.Info("stopCharge stop success : 操作人:%s",model.GetUserIp(r))
		contol.send(w,0,"stop success")
	}else if stype == "0"{
		model.StopRecharge = false
		gLog.Info("stopCharge open success : 操作人:%s",model.GetUserIp(r))
		contol.send(w,0,"open success")
	}else{
		gLog.Error("stopCharge stype error :client:%s,操作人:%s",stype,model.GetUserIp(r))
		contol.send(w,-100,"stype error")
	}
}