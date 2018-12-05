package config

import (
	"encoding/xml"
	"github.com/yhhaiua/engine/gjson"
	"github.com/yhhaiua/log4go"
	"github.com/yhhaiua/recharge/logic/model"
	"io/ioutil"
	"strconv"
)

type Config struct {
	Sport        string             //http端口
	Gmhost		 string				//发送的后台地址
	Rechargekey  string
	Clientkey 	 string
	Operatorid 	 string				//对应平台
	Chargeconfig map[int]int
	sqlconfig model.SqlConfig
}

type UnitOne struct {
	Id   string `xml:"id,attr"`
	Money string `xml:"money,attr"`
}

type UnitConfig struct {
	Unit []UnitOne `xml:"unit"`
}

func (config *Config) ConfigInit() bool {
	path := "./config/config.json"
	key := "recharge"

	configdata, err := ioutil.ReadFile(path)
	if err != nil {
		log4go.Error("Failed to open config file '%s': %s\n", path, err)
		return false
	}

	jsondata, err := gjson.NewJSONByte(configdata)
	if err != nil {
		log4go.Error("Failed to NewJsonByte config file '%s': %s\n", path, err)
		return false
	}
	keydata := gjson.NewGet(jsondata, key)

	if !keydata.IsValid() {
		log4go.Error("Failed1 to config file '%s'", path)
		return false
	}

	data := gjson.NewGetindex(keydata, 0)

	if !data.IsValid(){
		log4go.Error("Failed2 to config file '%s'", path)
		return false
	}

	config.Sport = data.Getstring("port")
	config.Clientkey = data.Getstring("clientkey")
	config.Gmhost = data.Getstring("gmhost")
	config.Rechargekey = data.Getstring("rechargekey")
	config.Operatorid = data.Getstring("operatorid")

	mysqldata := gjson.NewGet(data, "mysql")
	if !mysqldata.IsValid() {
		log4go.Error("Failed to mysql config file '%s'", path)
		return false
	}
	config.sqlconfig.Shost = mysqldata.Getstring("host")
	config.sqlconfig.Sdbname = mysqldata.Getstring("dbname")
	config.sqlconfig.Suser = mysqldata.Getstring("user")
	config.sqlconfig.Spassword = mysqldata.Getstring("password")
	config.sqlconfig.Maxopen = mysqldata.Getint("open")
	config.sqlconfig.Maxidle = mysqldata.Getint("idle")

	err = config.sqlconfig.InitDB()
	if err != nil{
		log4go.Error("Failed to mysql InitDB file '%s',err:%s", path,err)
		return false
	}

	return config.configXml()
}


func (config *Config) configXml() bool {
	path := "./config/charge.xml"
	content, err := ioutil.ReadFile(path)
	if err != nil {
		log4go.Error("Failed to open config file '%s': %s\n", path, err)
		return false
	}
	var tempConfig	 UnitConfig
	err = xml.Unmarshal(content, &tempConfig)
	if err != nil {
		log4go.Error("Failed to Unmarshal config file '%s': %s\n", path, err)
		return false
	}
	config.Chargeconfig = make(map[int]int)
	for _,temp := range tempConfig.Unit {
		id,_:=strconv.Atoi(temp.Id)
		money,_:= strconv.Atoi(temp.Money)
		config.Chargeconfig[id] = money
	}
	log4go.Info("charge.xml count:%d", len(config.Chargeconfig))
	return true
}

//获取物品价格
func (config *Config) GetMoney(itemid int) (int,bool)  {
	value,ok := config.Chargeconfig[itemid]
	return value,ok
}
