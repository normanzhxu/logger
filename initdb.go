package main

import (
	"database/sql"
	"fmt"
	"time"

	protoConfig "github.com/connext-cs/protocol/config"
	"github.com/connext-cs/pub/config"
	_ "github.com/go-sql-driver/mysql"
)

var (
	// DB is global db pool
	db       *sql.DB
	database *config.Database
)

func DbConfigSet(outdb *config.Database) {
	database = new(config.Database)
	database.Host = outdb.Host
	database.Port = outdb.Port
	database.User = outdb.User
	database.Password = outdb.Password
	database.Name = outdb.Name
	// fmt.Println("paasdb DbConfigSet, database:", database)
}

func init() {
	database = new(config.Database)
	database.Host = protoConfig.CMysqlHost()
	database.Port = uint16(protoConfig.CMysqlPort())
	database.Name = protoConfig.CCloudMysqlDatabase()
	database.User = protoConfig.CMysqlUserName()
	database.Password = protoConfig.CMysqlPasswd()
	// fmt.Println("paasdb init database:", database)
}

// InitDB init global db loc=Local&
func InitDB() (err error) {
	if db == nil {
		dbURI := fmt.Sprintf(
			"%s:%s@tcp(%s:%d)/%s?parseTime=true&charset=utf8",
			database.User,
			database.Password,
			database.Host,
			database.Port,
			database.Name,
		)
		// dbURI = "connextpaas:connext@0101@tcp(10.128.0.180:3306)/cloudproject?charset=utf8"
		db, err = sql.Open("mysql", dbURI)
		//fmt.Println("db:", db)
		if err != nil {
			//fmt.Println("Create DB error: ", err.Error())
			return err
		}
		if database.MaxIdleTime > 0 {
			db.SetConnMaxLifetime(time.Millisecond * time.Duration(database.MaxIdleTime))
		}
		if database.MaxIdle > 0 {
			db.SetMaxIdleConns(database.MaxIdle)
			db.SetMaxOpenConns(database.MaxOverflow + database.MaxIdle)
		} else {
			db.SetMaxIdleConns(10)
			db.SetMaxOpenConns(database.MaxOverflow + 10)
		}

		err = db.Ping()
		if err != nil {
			return
		}
	}
	return
}

func GetDB() *sql.DB {
	if db == nil {
		InitDB()
	}
	return db
}
