package main

import (
	"fmt"
	"os"
	"path/filepath"
	"time"
)

// SLog 日志表
type SLog struct {
	LogId         int       `json:"log_id"`          //日志ID
	LogResId      int       `json:"log_res_id"`      //资源ID
	LogResName    string    `json:"log_res_name"`    //资源名称
	LogTraingId   string    `json:"log_traing_id"`   //业务跟踪ID
	LogBusinessId string    `json:"log_business_id"` //业务主键
	LogMicroUri   string    `json:"log_micro_uri"`   //微服务URI
	LogLevel      string    `json:"log_level"`       //日志级别
	LogType       string    `json:"log_type"`        //日志类型
	LogMessage    string    `json:"log_message"`     //日志内容
	LogIp         string    `json:"log_ip"`          //请求IP
	LogTakeTime   int       `json:"log_take_time"`   //耗时
	AppId         int       `json:"app_id"`          //应用ID
	OrgId         int       `json:"org_id"`          //所属机构
	Deleted       int       `json:"deleted"`         //删除标志
	CreatedBy     string    `json:"created_by"`      //创建人
	CreatedTime   time.Time `json:"created_time"`    //创建时间
}

// ErrorToDB Error log inTo DB
func ErrorToDB(traingID, businessID, logType string, logTakeTime int, msg string, args ...interface{}) {
	message := Trace("[ERROR]"+msg, args...).Error()
	if len(traingID) > 0 {
		message += "\nTRACE_ID:" + traingID
	}
	hostName, _ := os.Hostname()
	session := getSession()
	_, err := session.InsertOne(&SLog{
		// LogResId      :,
		LogResName:    filepath.Base(os.Args[0]),
		LogTraingId:   traingID,
		LogBusinessId: businessID,
		LogMicroUri:   filepath.Base(os.Args[0]),
		LogLevel:      "ERROR",
		LogType:       logType,
		LogMessage:    message,
		LogIp:         hostName,
		LogTakeTime:   logTakeTime,
		// AppId         :,
		// OrgId         :,
		CreatedTime: time.Now(),
	})
	if err != nil {
		session.Rollback()
		panic(err)
	}
	session.Commit()
}

// InfoToDB info into DB
func InfoToDB(traingID, businessID, logType string, LogTakeTime int, msg string, args ...interface{}) {
	if len(args) > 0 {
		msg = fmt.Sprintf(msg, args...)
	}
	hostName, _ := os.Hostname()
	session := getSession()
	_, err := session.InsertOne(&SLog{
		// LogResId      :,
		LogResName:    filepath.Base(os.Args[0]),
		LogTraingId:   traingID,
		LogBusinessId: businessID,
		LogMicroUri:   filepath.Base(os.Args[0]),
		LogLevel:      "INFO",
		LogType:       logType,
		LogMessage:    msg,
		LogIp:         hostName,
		LogTakeTime:   LogTakeTime,
		// AppId         :,
		// OrgId         :,
		CreatedTime: time.Now(),
	})
	if err != nil {
		session.Rollback()
		panic(err)
	}
	session.Commit()
}
