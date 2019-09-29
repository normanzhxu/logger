# Connext_Response
# Restful Http Response Handling

> The following requset/response are abbreviated as req/rsp

## Req

> 面向前端传入的Pagination 请使用以下包

    import (
            "github.com/connext-cs/pub/request"
    )
    type RequestPagination struct {
    	PageNo   uint16 `json:"pageno"`   //当前的页数
    	PageRows uint16 `json:"pagerows"` //每页的行数
    }

## Rsp Code

| CODE | 说明       | 备注     |
| -------- | ---------- | --------------------------- |
| 0 | 正确 correct | 系统状态正常；业务操作完成              |
| 1 | 错误 error |  发生意料之外的错误             |
| 2 | 错误 error | 系统崩溃              |
| 3 | 错误 error | 数据库执行有误              |
| 4 | 警告 warning | 连接第三方或访问失败              |
| 5 | 提示 info |  不允许的业务操作             |


- Code: 0
    >更新成功或插入成功，返回数据主键的id

	    id := orgmodel.UpdOrg(request)
	    response.RspSucRestData(rsp, "", id)
- Code: 1
    > recover()捕获到的任意的未识别错误

        response.RspFailRestData(rsp, e.(error).Error())   
- Code: 2
    > 判断系统状态异常的情况下

        level := runtime.CPUtrace()
        response.RspFailSys(rsp, level)  
- Code: 3
    > 数据库中返回了错误

        int,err := orm.InsertOne(&Test)
        if err != nil{
            return err
        }
        response.RspFailSQL(rsp, err.Error())  
- Code: 4
    > 访问第三方服务http请求出现问题

	    err := client.Do(request)
	    response.RspFailHttp(rsp,err)
- Code: 5
    > 产品设计逻辑不允许

	    str := errors.New("此操作不被允许")
	    response.RspReject(rsp, str)




## Rsp Data

- 约定如下：
    > 
    > 1. 业务操作正常，程序状态正常情况下，传入正确的对象值
    > 
    > 2. 增、更新、删除等操作需要返回数据的主键信息
    > 
    > 3. 不能传空值或空数组或空对象，不得为nil返回
    > 
    > 4. 发生非正常情况，也需要在data中返回数据的主键信息

- 返回结构体 

    | CODE | Data       | Error    | Pagination |
    | -------- | ---------- | --------------------------- | --------------------------- |
    | 0 | id:xx | "" |             |
    | 1 | id:xx | "[Error]:" |              |
    | 2 | id:xx | "[Error]:" |               |
    | 3 | id:xx | "[Error]:" |               |
    | 4 | id:xx | "[Warning]:" |               |
    | 5 | id:xx | "[Info]:" |              |



## Rsp Error

    无发生error 的情况下返回空字符串

    当发生error 的情况下返回string的字符串

## Rsp Pagination

    import (
        "github.com/connext-cs/pub/response"
    )
    type Pagination struct {
	    PageSize    int `json:"pagesize"`     //每页显示的数据条数, 50
	    Page        int `json:"page"`         //总页数 = 总行数/每页显示的数据条数
	    CurrentPage int `json:"current_page"` //当前页号
	    Total       int `json:"total"`        //总行数
    }

