/**
 * @file: tmconfig.go
 * @package: tmconfig
 * @author: heiyeluren
 * @desc: config parse
 * @date: 2013/6/24
 * @history:
 *     2013/6/24 created file
 *     2022/2/24
 */
 
 
package tmconfig

// 1、所有配置加载到内存，
// 2、并推送到所有模块
// 3、配置内容实体

import (
    "errors"
    "io/ioutil"
    "strings"
)

var G_TMCONFIG map[string]map[string]string

// 内指定配置文件路径
func ParseConfPath(path string) (err error) {
    return parseFile(path)
}

// 解析配置文件成map[section]map[key]value的结构
func parseFile(path string) (err error) {
    if len(path) == 0 {
        path = "conf/conf.ini"
    }
    stream, err := ioutil.ReadFile(path)
    if err != nil {
        return errors.New("cannot load config file")
    }
    content := string(stream)
    lines := strings.Split(content, "\n")
    emptyRunes := "\n\r\v\t "
    section := ""
    G_TMCONFIG = make(map[string]map[string]string)
    for _, line := range lines {
        // 去掉注释
        if sinx := strings.Index(line, "#"); sinx > 0 {
            line = line[:sinx]
        }
        line = strings.Trim(line, emptyRunes)
        // 去除空行和注释
        if line == "" || line[0] == '#' || line[0] == ';' || line[0:2] == "//" {
            continue
        }
        // 识别section
        if line[0] == '[' && line[len(line)-1] == ']'{
            sect := line[1:len(line)-1]
            if _, ok := G_TMCONFIG[sect]; !ok {
                G_TMCONFIG[sect] = make(map[string]string)
            }
            // 更新目前解析的section
            section = sect
            continue
        //识别KV
        } else if cinx := strings.Index(line, "="); cinx > 0 &&
            cinx < len(line) - 1 && len(section) > 0 { // 等号不能出现在行首和行尾且当前的section非空
            key := strings.Trim(line[:cinx], emptyRunes)
            val := strings.Trim(line[cinx + 1:], emptyRunes)
            G_TMCONFIG[section][key] = val
        }
    }
    return nil
}

func ParseConf() (err error) {
    if len(G_TMCONFIG) > 0 {
        return
    }

    //默认的配置文件路径
    path := "../conf/conf.ini"
    return parseFile(path)
}

// 获取整个section的配置
func GetSection(section string) map[string]string {
    if cfm, ok := G_TMCONFIG[section]; ok {
        return cfm
    }

    return nil
}

// 获取某个section下key的值
func GetSectKey(section string, key string) string {
    if val, ok := G_TMCONFIG[section][key]; ok {
        return val

    }
    return ""
}
