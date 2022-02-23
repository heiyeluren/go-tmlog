## go-tmlog 内置配置解析库 - tmconf



```go
package main

import (
    "flag"
    "fmt"
    "utility/configs"
    tmconfig "github.com/heiyeluren/go-tmlog/configs"
)

// 配置文件对象初始化
func newConfig() *tmconfig.Config {
    var F string
    //从命令参数读取 -f 后面跟随的配置文件路径，类似于：~/etc/tmlog.conf
    flag.StringVar(&F, "f", "", "config file")
    flag.Parse()
    if F == "" {
        panic("usage: ./app-bin -f etc/tmlog.conf")
    }
    //调用 tmconfig 读取配置文件
    Config := configs.NewConfig()
    if err := Config.Load(F); err != nil {
        panic(err.Error())
    }
    
    //读取所有配置选项，是一个 map
    allconf := Config.GetAll()
    
    //读取单个配置选项
    pid_file := Config.Get("pid_file")
    
    //读取单个配置选项int型，另外也可以读取 Int64
    port := Config.GetInt("port")
    
    //读取一批数据为Slice，通过某个切割符
    list := Config.GetSlice("ip_list", ",")
           
    //返回对象
    return Config
}

//主函数
func main() {
  // 配置初始化对象
  var Config = newConfig()
  
  //输出对象
  fmt.Println(Config)  
}


```

