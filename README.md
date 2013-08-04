go-tmlog
========

Golang implement logging library (the require golang >= 1.0)

Golang使用的日志库


最后更新：2013/8/3 v1.0.0 <br />
下载地址：<a href="http://heiyeluren.googlecode.com/files/go-tmlog-v1.0.0.zip">Download</a>


<pre>

<b>[ Golang日志库: go-tmlog 使用说明文件 ]</b>

@author: heiyeluren
@date: 2013/8/3


<b>1. 基本架构:</b>
   a. 采用 日志输入(客户端) -> 日志输出(服务器端) 的架构，日志输入可以任意调用，日志输出是一个单独协程在工作,
   b. 能够保证日志保证时序性，并且保证客户端可以无限量的写入日志，不用担心阻塞而影响性能.
   c. 应用场景: 适合高性能日志打印的场合，按照测试，能够在每秒1万次请求的后端服务上进行日志打印，不会对性能有太多影响


<b>2. 主要功能:</b>
   a. 日志类型: 可以打印5种类型(notice/trace/debug/warning/fatal)的日志，并且代码很容易新增类型，不过基本够用了, 同时可以配置里定制那些日志类型需要记录, 一般推荐最少最少记录 notice/warning/fatal 三种日志
   b. 日志文件: 不同日志类型可以单独输出到指定日志文件中, 一般建议 notice/trace/debug 放一个日志文件, warning/fatal 放一个日志文件
   c. 日志切割: 支持按照 天/小时/10分钟 三种粒度自动进行日志文件切割，方便控制日志文件大小
   d. 日志刷盘: 可以指定日志刷盘的时间, 缺省1秒，建议不超过3秒; 如果当前日志达到缓存90%占用，会自动刷盘，保证不会阻塞写日志操作
   e. 调试模式: 支持调试模式，可以在运行中在终端输出一些信息，方便监测


<b>3. 配置说明</b>
  #--------------
	#日志操作配置
	#--------------

	#日志文件位置, 可以配置不同的日志消息类型到不同的日志文件 (例：/var/log/heiyeluren.log)
	log_notice_file_path	= log/heiyeluren.log
	log_debug_file_path	    = log/heiyeluren.log
	log_trace_file_path	    = log/heiyeluren.log
	log_fatal_file_path	    = log/heiyeluren.log.wf
	log_warning_file_path	= log/heiyeluren.log.wf

	#日志文件切割周期（1天:day; 1小时:hour; 10分钟:ten）
	log_cron_time = day

	#日志chan队列的buffer长度，在高并发服务器下，建议不要少于10240, 越大越好
	#不建议多于1024000，测试最长: 67021478 (超过这个值会无法启动)
	log_chan_buff_size = 1024000

	#日志刷盘的间隔时间，单位:毫秒，建议500~5000毫秒(0.5s-5s)，建议不超过30秒
	log_flush_timer = 1000

	#是否开启日志库调试模式(会在终端打印一些日志, 1:开启, 0:关闭)
	log_debug_open = 0

	#输出日志的级别 (fatal:1,warngin:2,notice:4,trace:8,debug:16)
	#级别描述主要是确定需要打印什么级别的日志，数字配置是一个需要打印日志级别数字的"或"操作总数(简单理解为加)
	#如果只打印错误警告和notice日志则是7，如果需要打印所有日志则是31，如果只需打印除trace以外的日志，则是23
	#如果不想输出任何日志，可以设置为0，特别在性能测试的时候，屏蔽刷日志带来的影响
	log_level = 31
 

<b>4. 调用示例 </b> (具体可以参考源码中的main.go的调用)
	================= 启动日志工作协程代码 =================
    import "heiyeluren/tmlog"

    // 传递给日志类的配置项
    // 说明: 这些配置实际可以放置到配置文件里，不过我们这里只是演示，直接就生成map了
    logConf := map[string]string{
        "log_notice_file_path":  "log/heiyeluren.log",
        "log_debug_file_path":   "log/heiyeluren.log",
        "log_trace_file_path":   "log/heiyeluren.log",
        "log_fatal_file_path":   "log/heiyeluren.log.wf",
        "log_warning_file_path": "log/heiyeluren.log.wf",
        "log_cron_time":         "day",
        "log_chan_buff_size":    "1024",
        "log_flush_timer":       "1000",
        "log_debug_open":        "1",
        "log_level":             "31",
    }
    // 启动 tmlog 工作协程, 可以理解为tmlog的服务器端
    go tmlog.Log_Run(logConf)


	================= 打印日志代码调用示例 =================
    //在main主协程里打印日志, 可以理解为tmlog的客户端往服务器端输入日志
    //说明: 这些代码可以在单独的任何非 tmlog 协程之外工作, 包括main协程, 或者是某些业务处理协程    

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

	//go调用函数
    func LogTest() {
        logHandle3 := tmlog.NewLogger("987654321")

        logHandle3.Notice("[logger=logHandle3 msg='The notice message is test']")
        logHandle3.Warning("[logger=logHandle3 msg='The warning message is test']")
    }    

<b>5. 运行测试代码</b>
	更多测试使用案例可以参考源码中的 main.go 文件，同时可以直接make运行该代码
	下载代码，然后执行如下脚本运行测试代码：
	Windows: cd 代码目录 && make && bin\test.exe
    Unix: cd 代码目录 && make && bin/test
   
	然后查看代码中 log 目录下是否生成了 .log 和 .log.wf 的日志文件，可以反复运行 test程序，然后使用 tail 或记事本持续观察日志文件变化。
	
	注意：以上操作必须保证安装了go 1.0以上的编译器的基础之上，同时能够正常访问go命令的情况下才能正常运行。如果不了解如何使用golang编译工具，请预先学习一下。

<b>6. 其他说明</b>
	作者: heiyeluren
	博客: http://blog.csdn.net/heiyeshuwu
	微博: http://weibo.com/heiyeluren


</pre>
