package wanba

import (
	"bytes"
	"crypto/hmac"
	"crypto/md5"
	"crypto/sha1"
	"encoding/base64"
	"encoding/hex"
	"fmt"
	"github.com/yhhaiua/engine/gjson"
	"github.com/yhhaiua/engine/log"
	"github.com/yhhaiua/recharge/logic/config"
	"github.com/yhhaiua/recharge/logic/control/backstage"
	"github.com/yhhaiua/recharge/logic/model"
	"io/ioutil"
	"net/http"
	"strconv"
	"time"
)
var gLog = log.GetLogger()
type WanBa struct {
	config *config.Config
}

func (wanba *WanBa)Init(config *config.Config){
	wanba.config = config
}
//生成sig
func (wanba *WanBa) sigCreate(urlinit,urlstr string) string {
	str0 :="POST"
	str1:= model.EncodeURIComponent(urlinit)
	str2 := model.EncodeURIComponent(urlstr)
	str3 := str0+"&"+str1+"&"+str2
	key :=  []byte("DWB13t84CoEL8eax&")
	mac := hmac.New(sha1.New, key)
	mac.Write([]byte(str3))
	ucnc := base64.StdEncoding.EncodeToString(mac.Sum(nil))
	fmt.Println(ucnc)
	return model.EncodeURIComponent(ucnc)
}

//玩吧处理
func (wanba *WanBa)RechargeDeal(w http.ResponseWriter, r *http.Request) (int,string)  {

	ret,value := wanba.inquiryRet(w,r)
	if ret == 0{
		playerid := r.FormValue("playerid")

		if(model.CheckMap(playerid)){
			return -100,"有订单正在处理"
		}

		model.AddMap(playerid)

		billno :=model.CreateorderId(playerid,"wb")

		money := value
		ret,value = wanba.deductionRet(w,r,money,billno)
		if ret == 0{
			wanba.sendGm(w,r,money,billno)
			return 0,"sucess"
		}else{
			model.DelMap(playerid)
		}
	}
	return ret,value
}

//判断玩家金钱
func (wanba *WanBa) inquiryRet(w http.ResponseWriter,r *http.Request) (int,string) {
	openid := r.FormValue("openid")
	sign := r.FormValue("sign")
	stime := r.FormValue("time")
	playerid := r.FormValue("playerid")
	itemid := r.FormValue("itemid")
	//openid+itemid+time+playerid+key
	md5str := openid + itemid + stime + playerid + wanba.config.Clientkey
	md5Ctx := md5.New()
	md5Ctx.Write([]byte(md5str))
	cipherStr := md5Ctx.Sum(nil)
	mysigon := hex.EncodeToString(cipherStr)
	if(mysigon != sign){
		gLog.Error("md5 error : me:%s,client:%s",mysigon,sign)
		return -100,"md5 error"
	}
	id,_:= strconv.Atoi(itemid)
	value,ok := wanba.config.GetMoney(id)
	if !ok{
		gLog.Error("没有对应的商品:%s",itemid)
		return -100,"no item"
	}
	return 0,strconv.Itoa(value)
}

func (wanba *WanBa)runSendGm(billno,money,itemid,playerid string)  {
	for  {
		time.Sleep(10*time.Second)
		ret,errStr := backstage.Sendgm("/tm_charge/pay/by/wanba",billno,money,itemid,playerid,wanba.config.Rechargekey,wanba.config.Gmhost)
		if ret == 0 || ret == 1{
			chargelog:= new(model.RechargeLog)
			err := chargelog.Update(billno,money,itemid,playerid,errStr,wanba.config.Operatorid)
			if(err != nil){
				gLog.Error("chargelog.Insert: %s",err)
			}
			model.DelMap(playerid)
			go backstage.Testsend("订单补发:"+errStr)
			break
		}else{
			go backstage.Testsend("订单补发失败:"+errStr)
		}
	}
}
//发送gm
func (wanba *WanBa) sendGm(w http.ResponseWriter,r *http.Request,money,billno string) (bool) {
	playerid := r.FormValue("playerid")
	itemid := r.FormValue("itemid")
	ret,errStr := backstage.Sendgm("/tm_charge/pay/by/wanba",billno,money,itemid,playerid,wanba.config.Rechargekey,wanba.config.Gmhost)
	if ret != 0{
		chargelog:= new(model.RechargeLog)
		err:=chargelog.Insert(billno,money,itemid,playerid,errStr,wanba.config.Operatorid)
		if(err != nil){
			gLog.Error("chargelog.Insert: %s",err)
		}
		if ret == -1{
			go wanba.runSendGm(billno,money,itemid,playerid)
		}else{
			model.DelMap(playerid)
		}
		go backstage.Testsend("充值错误:"+errStr)
	}else{
		model.DelMap(playerid)
	}
	return ret == 0
}


//玩吧扣款
func (wanba *WanBa)deductionRet(w http.ResponseWriter, r *http.Request,money,billno string) (int,string) {

	openid := r.FormValue("openid")
	openkey := r.FormValue("openkey")
	appid := r.FormValue("appid")
	userip := r.FormValue("userip")
	//count := r.FormValue("count")
	zoneid := r.FormValue("zoneid")
	pf := r.FormValue("pf")
	info := wanba.deduction(openid,openkey,appid,userip,zoneid,money,billno,pf)
	if(info != nil){
		jsondata, err := gjson.NewJSONByte(info)
		if err != nil {
			gLog.Error("deductionRet NewJsonByte: %s",err)
			return -100,"deductionRet json error"
		}
		code := jsondata.Getint("code")
		if code == 0{
			return 0,""
		}else {
			message:= jsondata.Getstring("message")
			if code == 1004{
				return -101,strconv.Itoa(code)+","+message
			}else if code == 1002{
				if message == "白名单用户额度不够"{
					return -101,strconv.Itoa(code)+","+message
				}else{
					return -100,strconv.Itoa(code)+","+message
				}
			}else{
				if message == ""{
					message = jsondata.Getstring("msg")
				}
				return -100,strconv.Itoa(code)+","+message
			}

		}
	}
	return -100,"deduction error"
}
//玩吧扣除玩家星币
func (wanba *WanBa)deduction(openid,openkey,appid,userip,zoneid,money,billno,pf string) []byte {

	//https://api.urlshare.cn/v3/user/buy_playzone_item?
	//	billno=xxxxx&
	//		openid=B624064BA065E01CB73F835017FE96FA&
	//		zoneid=1&
	//		openkey=5F154D7D2751AEDC8527269006F290F70297B7E54667536C&
	//		appid=2&
	//		itemid=10&
	//		count=1&
	//		sig=VrN%2BTn5J%2Fg4IIo0egUdxq6%2B0otk%3D&
	//		pf=wanba_ts&
	//		format=json&
	//		userip=112.90.139.30

	var buffer bytes.Buffer
	buffer.WriteString("appid=")
	buffer.WriteString(appid)
	buffer.WriteString("&billno=")
	buffer.WriteString(billno)
	buffer.WriteString("&count=")
	buffer.WriteString(money)
	buffer.WriteString("&format=json")
	buffer.WriteString("&itemid=")
	if(zoneid == "1"){
		buffer.WriteString("38008")
	}else if(zoneid == "2"){
		buffer.WriteString("38011")
	}
	buffer.WriteString("&openid=")
	buffer.WriteString(openid)
	buffer.WriteString("&openkey=")
	buffer.WriteString(openkey)
	buffer.WriteString("&pf=")
	buffer.WriteString(pf)
	buffer.WriteString("&userip=")
	buffer.WriteString(userip)
	buffer.WriteString("&zoneid=")
	buffer.WriteString(zoneid)

	sig:= wanba.sigCreate("/v3/user/buy_playzone_item",buffer.String())
	buffer.WriteString("&sig=")
	buffer.WriteString(sig)

	sendStr := "https://api.urlshare.cn/v3/user/buy_playzone_item?" + buffer.String()
	gLog.Info(sendStr)

	transport := http.Transport{
		DisableKeepAlives: true,
	}
	client := &http.Client{
		Transport:&transport,
		Timeout: 10 * time.Second,
	}
	//req, err := client.Get(sendStr)
	req, err := client.Post("https://api.urlshare.cn/v3/user/buy_playzone_item","application/x-www-form-urlencoded",&buffer)
	if err == nil {
		defer req.Body.Close()
		body, err := ioutil.ReadAll(req.Body)
		if err == nil {
			gLog.Info(string(body))
			return body
		}else{
			gLog.Error("deduction error2:%s",err)
		}
	}else{
		gLog.Error("deduction error1:%s",err)
	}
	return nil
}

//玩吧补单
func (wanba *WanBa)MakeUpOrder(w http.ResponseWriter, r *http.Request)(int,string)  {
	playerid := r.FormValue("playerid")
	itemid := r.FormValue("itemid")
	id,_:= strconv.Atoi(itemid)
	value,ok := wanba.config.GetMoney(id)
	if !ok{
		gLog.Error("没有对应的商品:%s",itemid)
		return -100,"no item"
	}
	if(model.CheckMap(playerid)){
		return -100,"有订单正在处理"
	}

	model.AddMap(playerid)

	billno := r.FormValue("billno")
	if billno == ""{
		billno = model.CreateorderId(playerid,"wb")
	}
	ret,errStr := backstage.Sendgm("/tm_charge/pay/by/wanba",billno,strconv.Itoa(value),itemid,playerid,wanba.config.Rechargekey,wanba.config.Gmhost)
	go backstage.Testsend("gm补单充值:"+errStr)

	model.DelMap(playerid)

	if ret != 0{
		return -100,errStr
	}else{
		return 0,errStr
	}
}