<!-- TOC -->

- [logger](#logger)
    - [func Error / Log(traceID, msg string, args ...interface{}) { }](#func-error--logtraceid-msg-string-args-interface--)
    - [func LogTaketime(startTime int64) int { }](#func-logtaketimestarttime-int64-int--)
    - [func SetDebugPro(onoff bool) { }](#func-setdebugproonoff-bool--)
    - [func DbgPro(prefix, msg string, args ...interface{}) { }](#func-dbgproprefix-msg-string-args-interface--)
    - [func Perf(tag string, work func()) { }](#func-perftag-string-work-func--)

<!-- /TOC -->
# logger

## func Error / Log(traceID, msg string, args ...interface{}) { }
> Error / Log output message with stack trace. prefix is added to every single line of log output for tracing purpose.

- What 它是什么？

标准输出 错误 / 详细 日志信息的函数

- How 如何使用？

可在业务代码中直接调用Error / Log 函数，其会将标准输出打印显示

例如：

    Error("AzureBillTrans", "账单清洗发生错误：%s", err.Error())
    Log("AzureBillTrans", "账单清洗状态：%s", statusStr)

## func LogTaketime(startTime int64) int { }
> Log take time util the func be called

- What 它是什么？

记录调用函数至此所消耗的时间

- How 如何使用？

将开始时间传入，返回即使此过程消耗时间
	
例如：

    start := time.Now().UnixNano() / 1e6
    tokenTime := logTaketime(start)

## func SetDebugPro(onoff bool) { }
> set Debug mode be on or off

- What 它是什么？

设置debug模式的开和关

- How 如何使用？

包初始化默认为FALSE
    
    var debugging bool 

传入bool 型参数即可改变全局变量debugging 的属性值

	SetDebugPro(true)


## func DbgPro(prefix, msg string, args ...interface{}) { }
> Debug mode Log printing

- What 它是什么？

根据debugging 的bool值执行与否，显示输入msg以及args的组合信息

- How 如何使用？

在业务代码需要打出Debug的地方，调用函数并设置traceID 以及 msg的信息等


例如：

    2019/09/15 16:46:17 [Debug]msgmsg
    TRACE_ID:Dbg (logger/test2.go:6) main >


## func Perf(tag string, work func()) { }
 
> Perf calculte work func used time

- What 它是什么？

利用debug 模式，计算某个函数的消费时间。

- How 如何使用？

在需要调试优化的业务代码中，添加traceID 以及 函数名

> 适用于初期的调试工作或优化工作，注意其中的func函数不需要其函数的返回值。

例如：

    2019/09/15 20:21:13 [Debug][EXEC]taggg
    TRACE_ID:prefix (logger/logger.go:46) main.Perf>
    2019/09/15 20:21:13 [Debug][DONE]taggg (elapsed: 0.000000)
    TRACE_ID:prefix (logger/logger.go:49) main.Perf>
