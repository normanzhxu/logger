
# DBOPS

## func ErrorToDB(traingID, businessID, logType string, logTakeTime int, msg string, args ...interface{}) { }

> The log format of the required output will be stored in the database

- What 它是什么？

将标准输出的 错误 / 详细 日志信息存储到数据库

- Why 为什么要写？

便于通过各类信息标签筛选选出想要的日志信息

- How 如何使用？

在想要存储业务的信息的业务代码中直接调用ErrorToDB / InfoToDB 函数，其会将日志信息存储到数据库

## func InfoToDB(traingID, businessID, logType string, logTakeTime int, msg string, args ...interface{}) { }

> 同上 日志级别为 INFO
