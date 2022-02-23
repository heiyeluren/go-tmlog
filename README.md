# 高性能Go日志库 - go-tmlog


### go-tmlog是什么
tmlog是一个使用Go语言开发，基于Channel高性能异步日志库，能够并发每秒单机写入10000+日志到磁盘不卡顿，使用简单，能够支持各种日志类型，并且支持日志轮转功能等等。是一个可靠小巧的日志库。

<br />

### Go版本要求
Golang implement logging library (the require golang >= 1.0)

<br />

### go-tmlog 使用

```go
package main

import (
    "github.com/heiyeluren/go-tmlog"
    "runtime"
    "time"
)

/**
 * 单独协程函数测试打印日志
 *
 */
func LogTest() {
    logHandle3 := tmlog.NewLogger("987654321")

    logHandle3.Notice("[logger=logHandle3 msg='The notice message is test']")
    logHandle3.Warning("[logger=logHandle3 msg='The warning message is test']")
}


/**
 * 测试服务进程的main函数
 */
func main() {

    /**
     * 以下为tmlog后端协程工作代码
     *
     * 说明: tmlog主工作协程必须在main协程里面启动工作, 否则无法完成后台日志打印的工作
     */

    // 传递给日志类的配置项 (配置项的含义参考 tmlog.go 代码里面的注释说明)
    // 说明: 这些配置实际可以放置到配置文件里，不过我们这里只是演示，直接就生成map了
    logConf := map[string]string{
        "log_notice_file_path":  "log/heiyeluren.log",
        "log_debug_file_path":   "log/heiyeluren.log",
        "log_trace_file_path":   "log/heiyeluren.log",
        "log_fatal_file_path":   "log/heiyeluren.log.wf",
        "log_warning_file_path": "log/heiyeluren.log.wf",
        "log_cron_time":         "day",
        "log_chan_buff_size":    "16",
        "log_flush_timer":       "1000",
        "log_debug_open":        "1",
        "log_level":             "31",
    }
    // 启动 tmlog 工作协程, 可以理解为tmlog的服务器端
    tmlog.Log_Run(logConf)

    /**
     * 以下为tmlog前端实际打印日志工作代码
     *
     * 说明: 这些代码可以在单独的任何非 tmlog 协程之外工作, 包括main协程, 或者是某些业务处理协程
     */

    //在main主协程里打印日志, 可以理解为tmlog的客户端往服务器端输入日志

    //打印日志1: 由tmlog生成logid, 生成log操作句柄, 打印notice和warning两条日志
    logHandle1 := tmlog.NewLogger("")
    logHandle1.Notice("[logger=logHandle1 msg='The notice message is test']")
    logHandle1.Warning("[logger=logHandle1 msg='The warning message is test']")

    //打印日志2: 调用方指定logid, 生成另外log操作句柄, 打印notice和warning两条日志
    //注意: 这里只是演示，为了追查问题方便, 一般情况不建议一个请求使用多个logid
    logHandle2 := tmlog.NewLogger("123456789")
    logHandle2.Notice("[logger=logHandle2 msg='The notice message is test']")
    logHandle2.Warning("[logger=logHandle2 msg='The warning message is test']")

    //打印日志3: 在一个单独协程里打印日志
    go LogTest()

    // hold 住 main协程不退出
    // 说明: 如果想一直运行程序, 可以把下面这行开启
    //select {}

    //sleep
    time.Sleep(time.Second * 5)

    //run done
    println("Log test programe run done.")
}


```

<br />

## go-tmlog 详细说明

### - [go-tmlog使用说明](https://github.com/heiyeluren/go-tmlog/blob/master/docs/tmlog-usage.md)


<br />


## 问题反馈和联系

- go-tmlog 技术交流

go-tmlog 目前在360公司部分线上项目在使用，当然多多稍稍也少不了一些问题和 bug，欢迎大家一起共创，或者直接提交 PR 等等。
有问题或者想要交流可以添加如下微信：
（如无法看到图片，请手工添加微信： xiehualiang）

<img src=https://raw.githubusercontent.com/heiyeluren/XMM/main/docs/img/xmm-wx.png width=30% />




