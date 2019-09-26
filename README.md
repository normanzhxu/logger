# Connext-Audit 

    Code or Architecture level needs to be processed to collect log information
    Error Handling, Logging and Instrumentation for Golang

> 内容主要有：
> 
> - 研发过程中，何种代码或架构层次需要进行收集日志信息的处理
> 
>       - 区分不同的业务处理的情况
>       - 区分不同日志级别处理的情况
>       - 开设debug模式以及输出信息至log文件或db的使用
> 
> - 研发过程中，针对上述不同的处理方式介绍connext_log包的使用方式
>       
>       - logger functions
>       - error functions
>       - dbop functions
> 

## 不同业务日志处理的情况

> 何时何地需要设置并处理日志信息

1. 接受前端request 请求中params 参数
    - 使用：

            curOrgType := req.Request.URL.Query().Get("orgType")
            zhlog.Log("curOrgType","这是案例中的组织类型 %s",curOrgType)
    - 结果：
        
            2019/09/15 16:46:17 [Info] 这是案例中的组织类型 company
            TRACE_ID: curOrgType

2. 返回前端response 中重要params 参数
    
    > 重要参数泛指的是 此条数据所在数据库中唯一主键 或 数据变更项的返回 或 操作类型rsp 返回 code
    - 使用1：

            zhlog.Log("SucData","操作成功后的数据ID: %s",data.ID)
            response.RspSucRestData(rsp, "", data))
    - 结果1：
        
            2019/09/15 16:46:17 [Info] 操作成功后的数据ID: 926
            TRACE_ID: SucData
    - 使用2：

			zhlog.Error("FailList", "错误信息：%s", e.(error).Error())
			response.RspFailRestData(rsp, e.(error).Error())
    - 结果2：

            2019/09/26 06:11:44 [ERROR] 错误信息：CheckEvth TableLists 传入错误数据，无法识别
            (/rdb/apprdb/appcheck.go:141)    api-service-platform/rdb/apprdb.CheckEvth
            (/models/app/appsetmodel.go:34)  api-service-platform/models/app.SetApp
            (/models/app/appsetmodel.go:25)  api-service-platform/models/app.UpdApp
            (/handle/app/app_handle.go:156)  api-service-platform/handle/app.(*App).SetApp

3. 获取用户认证或验证信息等
    > 用户token的收集或用户信息的参数
    - 使用1：

            func GetURLList(){
                enrollNum, authKey:= getformDB()
                zhlog.Dbg("traceID", "enrollment:%s, key:%s", enrollNum, authKey)
            }
    - 结果1：
    
            2019/09/15 16:46:17 [Debug]enrollment: V570, key: secret 
            TRACE_ID:traceID (model/cloudBillModel/init_data.go:27) GetURLList
    - 使用2：

                func Test1(){
                    req.Header.Set("Authorization", authKey)
                    zhlog.Dbg("Authorization", "key:%s", authKey)
                }
    - 结果2：
    
            2019/09/15 16:46:17 [Debug]key: secret 
            TRACE_ID:Authorization (model/cloudBillModel/init_data.go:65) Test1

4. 业务上逻辑主动返回的操作原因
    - 使用1：

            initStartTime := billAttr[0].InitStartDate
	        if initStartTime == "" {
		        zhlog.Assert(errors.New("无法获取初始化时间，请在账户管理界面设置初始化时间"))
	        }
    - 结果1：
    
            2019/09/26 06:11:44 [ERROR] 无法获取初始化时间，请在账户管理界面设置初始化时间
            (/rdb/apprdb/appcheck.go:141)    api-service-platform/rdb/apprdb.CheckEvth
            (/models/app/appsetmodel.go:34)  api-service-platform/models/app.SetApp
            (/handle/app/app_handle.go:156)  api-service-platform/handle/app.(*App).SetApp

5. 建立db 链接
    
    - 使用1：

            func mysqlEngine(){
                ....
                dburl := User + ":" + Password + "@tcp(" + Host + ":" + strconv.Itoa(int(Port)) + ")/" + Name + "?charset=utf8"
                zhlog.Dbg("db path","db connect info : %s",dburl)
                ....
            }
    - 结果1：
        
            2019/09/15 16:46:17 [Debug] db connect info : connext:qwe@tcp(localhost:3306/test?charset=utf8)
            TRACE_ID: db path (rdb/dbEngine.go:53) mysqlEngine
	
6. 命令行参数设置以及全局变量或常量的显示
    - 使用1：

            func echo(){
                ....
                reqURL := fmt.Sprintf("%s/rest/%s/usage-reports", common.AzureDomain, enrollNum)
	            zhlog.Dbg("reqUrl", "domain: %s", common.Domain)
                ....
            }
    - 结果1：
        
            2019/09/15 16:46:17 [Debug] domain: http://ea.connext.com
            TRACE_ID: db path (handle/testhandle/connect.go:53) echo	


## 各类输出日志的级别分配
> 不同的日志信息需要如何处理

- 属于Info 日志级别
>
>1. 接受前端request 请求中params 参数
>2. 返回前端response 中重要params 参数

- 属于Dbg 日志级别
>1. 获取用户认证或验证信息等
>2. 建立db 链接
>3. 命令行参数设置以及全局变量或常量的显示

- 属于Error 日志级别
>1. 发生任意的逻辑错误 如：数组越界或空指针异常或断言错等
>2. 业务上逻辑主动返回的操作原因


## 开设debug模式以及输出信息至log文件或db的使用

- 开设debug模式日志
    
    - 使用1：

            func main(){
                SetDebugPro(true)
            }
    - 使用2：

            ./api-servive-platform -debug true

- 选择日志信息是否进入或db

    - 使用1：

            func main(){
                SetDBPro(true)
            }
    - 使用2：

            ./api-servive-platform -LogDB true