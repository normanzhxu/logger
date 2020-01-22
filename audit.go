package main

import (
	"database/sql"
	"fmt"
	"os"
	"path/filepath"
	"pub/tokens"
	"strconv"
	"time"

	"github.com/emicklei/go-restful"
)

const (
	DebugLog = "Debug"
	InfoLog  = "Info"
	ErrorLog = "Error"
	FatalLog = "Fatal"
)
const (
	ActionAdd = "新增"
	ActionUpd = "更新"
	ActionDel = "删除"
)

// info日志
// AuditAddInfo 新增功能
func AuditAddInfo(req *restful.Request, businessID, content interface{}) {
	layerLast(perBase(req), makeCustom(InfoLog, ActionAdd, businessID, content))
}

// AuditUpdInfo 更新功能
func AuditUpdInfo(req *restful.Request, businessID, content interface{}) {
	layerLast(perBase(req), makeCustom(InfoLog, ActionUpd, businessID, content))
}

// AuditDelInfo 删除功能
func AuditDelInfo(req *restful.Request, businessID, content interface{}) {
	layerLast(perBase(req), makeCustom(InfoLog, ActionDel, businessID, content))
}

// error日志
// AuditAddError 新增功能
func AuditAddError(req *restful.Request, businessID, content interface{}) {
	layerLast(perBase(req), makeCustom(ErrorLog, ActionAdd, businessID, content))
}

// AuditUpdError 更新功能
func AuditUpdError(req *restful.Request, businessID, content interface{}) {
	layerLast(perBase(req), makeCustom(ErrorLog, ActionUpd, businessID, content))
}

// AuditDelError 删除功能
func AuditDelError(req *restful.Request, businessID, content interface{}) {
	layerLast(perBase(req), makeCustom(ErrorLog, ActionDel, businessID, content))
}

// fatal日志
// AuditAddFatal 新增功能
func AuditAddFatal(req *restful.Request, businessID, content interface{}) {
	layerLast(perBase(req), makeCustom(FatalLog, ActionAdd, businessID, content))
}

// AuditUpdFatal 更新功能
func AuditUpdFatal(req *restful.Request, businessID, content interface{}) {
	layerLast(perBase(req), makeCustom(FatalLog, ActionUpd, businessID, content))
}

// AuditDelFatal 删除功能
func AuditDelFatal(req *restful.Request, businessID, content interface{}) {
	layerLast(perBase(req), makeCustom(FatalLog, ActionDel, businessID, content))
}

// Debug日志
// AuditAddDebug 新增功能
func AuditAddDebug(req *restful.Request, businessID, content interface{}) {
	layerLast(perBase(req), makeCustom(DebugLog, ActionAdd, businessID, content))
}

// AuditUpdDebug 更新功能
func AuditUpdDebug(req *restful.Request, businessID, content interface{}) {
	layerLast(perBase(req), makeCustom(DebugLog, ActionUpd, businessID, content))
}

// AuditDelDebug 删除功能
func AuditDelDebug(req *restful.Request, businessID, content interface{}) {
	layerLast(perBase(req), makeCustom(DebugLog, ActionDel, businessID, content))
}

func phraseBID(businessID interface{}) string {
	var bidStr string
	switch businessID.(type) {
	case int64:
		bidStr = strconv.FormatInt(businessID.(int64), 10)
	case int:
		bidStr = strconv.Itoa(businessID.(int))
	case string:
		bidStr = businessID.(string)
	case []int64:
		for k, v := range businessID.([]int64) {
			if k > 0 {
				bidStr += ","
			}
			bidStr += strconv.FormatInt(v, 10)
		}
	case []int:
		for k, v := range businessID.([]int) {
			if k > 0 {
				bidStr += ","
			}
			bidStr += strconv.Itoa(v)
		}
	case []string:
		for k, v := range businessID.([]string) {
			if k > 0 {
				bidStr += ","
			}
			bidStr += v
		}
	default:
		bidStr = "undefined type"
	}
	return bidStr
}

func phraseCtx(content interface{}) string {
	// db := GetDB()
	// defer db.Close()
	var ctx string
	switch content.(type) {
	case map[string][]interface{}:
		for tName, pID := range content.(map[string][]interface{}) {
			ctx += "\n"
			ctx += singleTable(tName, pID)
		}
	default:
		ctx = "undefined type"
	}

	return ctx
}

func singleTable(tName string, pID []interface{}) string {
	var pIDStr string
	for k, v := range pID {
		if k > 0 {
			pIDStr += ","
		}
		switch v.(type) {
		case int64:
			pIDStr += strconv.FormatInt(v.(int64), 10)
		case int:
			pIDStr += strconv.Itoa(v.(int))
		case string:
			pIDStr += v.(string)
		}
	}
	// var singleCtx string

	// for k, singleID := range table {
	// 	// singleCtx += single(singleID)
	// }

	// tableID := ""
	// ID := ""
	// rows, err := db.Query(`SELECT * from ? WHERE id = ? and deleted=0 and enabled=1`, tableID, ID)
	// Assert(err)
	// dbContent := FetchRows(rows)
	// byteStr, err := json.MarshalIndent(dbContent, "", "	")
	// Assert(err)
	// ctx = string(byteStr)

	// return singleCtx
	return pIDStr
}

func makeCustom(level, actionType string, businessID, content interface{}) Custom {
	demo2 := Custom{
		BusinessID: phraseBID(businessID),
		Level:      level,
		ActionType: actionType,
		Content:    phraseCtx(content),
	}
	return demo2
}

func getUserID(req *restful.Request) (userID string) {
	tokenValue, err := tokens.ReadToken(req.Request.Header.Get("Authorization"))
	Assert(err)
	userID = strconv.FormatInt(tokenValue.UserID, 10)
	Log("ParseToken", "userId : %s", userID)
	return
}

func getClientIP(req *restful.Request) (clientIP string) {
	clientIP = req.Request.RemoteAddr
	return
}

func getResID(req *restful.Request) (resID int) {
	resIDStr := req.Request.Header.Get("ResourceID")
	resID, err := strconv.Atoi(resIDStr)
	Assert(err)
	return
}

// 头部解析信息
func perBase(req *restful.Request) Base {
	return makeBase(getResID(req), getClientIP(req), getUserID(req))
}

func makeBase(resID int, clientIP, userID string) Base {
	demo1 := Base{
		ResID:    resID,
		ResName:  "",
		ClienIP:  clientIP,
		UserID:   userID,
		UserName: "",
	}
	demo1.GetResName()
	demo1.GetUserName()
	return demo1
}

func layerLast(demo1 Base, demo2 Custom) {
	defer func() {
		if e := recover(); e != nil {
			Error("Audit", "%s", e.(error).Error())
		}
	}()
	recorder(demo1, UUID(8), demo2)
}

func (e *Base) GetResName() {
	db := GetDB()
	// defer db.Close()
	rows, err := db.Query(`SELECT name from auth_resource WHERE id = ? and deleted=0 and enabled=1`, e.ResID)
	Assert(err)
	data := FetchRows(rows)
	if len(data) != 1 {
		err := fmt.Errorf("请检查数据，当前菜单 %d 查询不到菜单名称或存在多个菜单名称", e.ResID)
		Assert(err)
	}
	e.ResName = data[0]["name"].(string)
}

func (e *Base) GetUserName() {
	db := GetDB()
	// defer db.Close()
	rows, err := db.Query(`SELECT user_name from auth_users WHERE id = ? and deleted=0 and enabled=1`, e.UserID)
	Assert(err)
	data := FetchRows(rows)
	if len(data) != 1 {
		err := fmt.Errorf("请检查数据，当前用户 %s 查询不到用户名称或存在多个用户名称", e.UserID)
		Assert(err)
	}
	e.UserName = data[0]["user_name"].(string)
}

type Base struct {
	ResID    int
	ResName  string
	ClienIP  string
	UserID   string
	UserName string
}

type Custom struct {
	BusinessID string
	Level      string
	ActionType string
	Content    string
}

func recorder(base Base, traceID string, custom Custom) {
	orm := cloudProjectEngine()
	_, err := orm.InsertOne(&SLog{
		LogResId:   base.ResID,
		LogResName: base.ResName,
		LogIp:      base.ClienIP,
		CreatedBy:  base.UserName,

		LogTraingId:   traceID,
		LogBusinessId: custom.BusinessID,
		LogLevel:      custom.Level,
		LogType:       custom.ActionType,
		LogMessage:    custom.Content,
		// LogTakeTime:   logTakeTime,
		// AppId         :,
		// OrgId         :,
		LogMicroUri: filepath.Base(os.Args[0]),
		CreatedTime: time.Now(),
	})
	Assert(err)
}

func FetchRows(rows *sql.Rows) []map[string]interface{} {
	defer func() {
		if e := recover(); e != nil {
			rows.Close()
			panic(e)
		}
	}()
	cols, err := rows.Columns()
	Assert(err)
	raw := make([][]byte, len(cols))
	ptr := make([]interface{}, len(cols))
	for i := range raw {
		ptr[i] = &raw[i]
	}
	var recs []map[string]interface{}
	for rows.Next() {
		Assert(rows.Scan(ptr...))
		rec := make(map[string]interface{})
		for i, r := range raw {
			if r == nil {
				rec[cols[i]] = nil
			} else {
				rec[cols[i]] = string(r)
			}
		}
		recs = append(recs, rec)
	}
	Assert(rows.Err())
	return recs
}
