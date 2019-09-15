<!-- TOC -->

- [error](#error)
    - [func Trace(msg string, args ...interface{}) (logs Exception) { }](#func-tracemsg-string-args-interface-logs-exception--)
    - [func (e Exception) Error() string { }](#func-e-exception-error-string--)
    - [func Assert(err error) { }](#func-asserterr-error--)

<!-- /TOC -->

# error

## func Trace(msg string, args ...interface{}) (logs Exception) { }
> Trace returns the current stack trace

- What 它是什么？

Trace( )方法 是属于诊断日志的一种。

目的是将超出预期发生且无法被正常处理的异常，获取出其堆栈信息

其中将对触发函数所在的
**文件地址**
**函数名称**
**文件行数** 都进行有效的切分以及拼接，提供良好的分析

- Why 为什么要写？

确保开发运维人员快速准确定位线上问题

- How 如何使用？

需要查看变量或函数在堆栈中的变化以及值的时候，可以选择调用该函数

调用Trace() 函数后，会将返回一组函数触发层以及其上层直至main函数的堆栈信息

例如：

    msg:args
        (logger/test1.go:4) main.func1
        (logger/main.go:9) main.main

> 注意：这些堆栈信息将会以[]string 的形式返回

> 建议结合 func (e Exception) Error( ) { } 方法显示文本信息

## func (e Exception) Error() string { }

例如：

    Trace("[ERROR]"+msg, args...).Error()

## func Assert(err error) { }

> 判断err是否存在 ，存在则panic 触发recover

