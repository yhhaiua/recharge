package backstage

import (
	"bytes"
	"crypto/md5"
	"encoding/hex"
	"github.com/yhhaiua/engine/log"
	"io/ioutil"
	"net/http"
	"strconv"
	"time"
)
var gLog = log.GetLogger()

func parseServer(playerId string) string  {
	pid,_:= strconv.ParseInt(playerId, 10, 64)
	serverid:= (pid % 1000000) / 100
	return  strconv.FormatInt(serverid, 10)
}
//像gm后台发送
func Sendgm(routers,orderId,money,goodsId,playerId,rechargekey,gmhost string) (int,string) {
	success:= -1
	errorStr := orderId+":error"
	//playerId+商品id+时间+订单号+key
	timestamp := time.Now().Unix()
	severid := parseServer(playerId)
	stime := strconv.FormatInt(timestamp, 10)
	md5str := playerId + goodsId + stime + orderId + rechargekey
	md5Ctx := md5.New()
	md5Ctx.Write([]byte(md5str))
	cipherStr := md5Ctx.Sum(nil)
	mysigon := hex.EncodeToString(cipherStr)

	var buffer bytes.Buffer
	buffer.WriteString("http://")
	buffer.WriteString(gmhost)
	buffer.WriteString(routers)
	buffer.WriteString("?orderId=")
	buffer.WriteString(orderId)
	buffer.WriteString("&serverId=")
	buffer.WriteString(severid)
	buffer.WriteString("&money=")
	buffer.WriteString(money)
	buffer.WriteString("&goodsId=")
	buffer.WriteString(goodsId)
	buffer.WriteString("&playerId=")
	buffer.WriteString(playerId)
	buffer.WriteString("&time=")
	buffer.WriteString(stime)
	buffer.WriteString("&sign=")
	buffer.WriteString(mysigon)
	gLog.Info(buffer.String())

	transport := http.Transport{
		DisableKeepAlives: true,
	}
	client := &http.Client{
		Transport:&transport,
		Timeout: 5 * time.Second,
	}
	req, err := client.Get(buffer.String())
	if err == nil {
		defer req.Body.Close()
		body, err := ioutil.ReadAll(req.Body)
		if err == nil {
			valueRet := string(body)
			if valueRet == "ok"{
				success = 0
				errorStr = orderId+":ok"
				gLog.Info(errorStr)
			}else{
				gLog.Error("订单号:orderId:%s,充值错误返回:%s",orderId,valueRet)
				errorStr = orderId+":"+valueRet
				if valueRet == "-2"{
					success = -1
				}else{
					success = 1
				}
			}
		}else{
			gLog.Error("sendgm error2 订单号:orderId:%s,error:%s",orderId,err)
			success = 1
			errorStr = orderId+":-10002"
		}
	}else{
		gLog.Error("sendgm error1 订单号:orderId:%s,error:%s",orderId,err)
		success = -1
		errorStr = orderId+":-10001"
	}

	return success,errorStr
}
