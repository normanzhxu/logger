package main

import (
	"strconv"

	protoConfig "github.com/connext-cs/protocol/config"

	"github.com/go-xorm/xorm"
)

func getSession() (session *xorm.Session) {
	orm := cloudProjectEngine()
	orm.ShowSQL(false)
	session = orm.NewSession()
	defer session.Close()
	err := session.Begin()
	Assert(err)
	return
}

func cloudProjectEngine() *xorm.Engine {
	orm, err := mysqlEngine()
	Assert(err)
	return orm
}

func mysqlEngine() (*xorm.Engine, error) {
	Host := protoConfig.CMysqlHost()
	Port := uint16(protoConfig.CMysqlPort())
	Name := protoConfig.CCloudMysqlDatabase()
	User := protoConfig.CMysqlUserName()
	Password := protoConfig.CMysqlPasswd()
	dburl := User + ":" + Password + "@tcp(" + Host + ":" + strconv.Itoa(int(Port)) + ")/" + Name + "?charset=utf8"
	// return xorm.NewEngine("mysql", "connextpaas:connext@0101@tcp(127.0.0.1:3306)/cloudproject?charset=utf8")
	return xorm.NewEngine("mysql", dburl)
}
