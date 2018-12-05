package model

import "github.com/astaxie/beego/orm"
import _ "github.com/go-sql-driver/mysql"

//SqlConfig 连接配置
type SqlConfig struct {
	Shost     string //ipport
	Sdbname   string //数据库名
	Suser     string //用户名
	Spassword string //密码
	Maxopen   int    //最大连接数
	Maxidle   int    //最大空闲数
}

//InitDB 初始配置
func (sqlconfig *SqlConfig) InitDB() error  {
	err := orm.RegisterDriver("mysql", orm.DRMySQL)
	if(err != nil){
		return err
	}
	orm.RegisterModel(new(RechargeLog))

	dataSource := sqlconfig.Suser + ":"+sqlconfig.Spassword+"@tcp("+sqlconfig.Shost+")/"+sqlconfig.Sdbname+"?charset=utf8&loc=Local"
	err = orm.RegisterDataBase("default", "mysql", dataSource)
	if(err != nil){
		return err
	}
	orm.SetMaxIdleConns("default",sqlconfig.Maxidle)
	orm.SetMaxOpenConns("default",sqlconfig.Maxopen)

	err = orm.RunSyncdb("default", false, true)
	return err
}
