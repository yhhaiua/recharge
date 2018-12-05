package model

import (
	"bytes"
	"net/http"
	"net/url"
	"strconv"
	"strings"
	"sync"
	"sync/atomic"
	"time"
)

var (
	increaseId    int64
	StopRecharge  bool
	players *sync.Map
)

func init()  {
	players = new(sync.Map)
}

func CheckMap(id string) bool{
	_,ok := players.Load(id)
	return ok
}

func AddMap(id string)  {
	players.Store(id,id)
}

func DelMap(id string)  {
	players.Delete(id)
}

func GetIncreaseID()int64  {
	return  atomic.AddInt64(&increaseId,1)
}

//生成订单号
func CreateorderId(playerId,str string) string {

	timestamp := time.Now().Unix()
	stime := strconv.FormatInt(timestamp, 10)
	var buffer bytes.Buffer
	buffer.WriteString(str)
	buffer.WriteString(playerId)
	buffer.WriteString(stime)
	buffer.WriteString( strconv.FormatInt(GetIncreaseID(),10))
	return buffer.String()
}

//urlencode
func EncodeURIComponent(str string) string {
	r := url.QueryEscape(str)
	r = strings.Replace(r, "+", "%20", -1)
	return r
}

//获取文件ip
func GetUserIp(r *http.Request) string  {
	userip := r.Header.Get("X-Real-IP")
	if userip == ""{
		userip = r.RemoteAddr
	}
	return userip
}