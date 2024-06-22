package server

import (
	"encoding/json"
	"fmt"
	"strings"
	"sync"
)

const (
	pushConn int = iota
	dbConn
	apiConn
	sdkConn
)

type Config struct {
	sync.RWMutex
	gatherJobScriptTemplate string   //采集任务转换脚本
	callJobScriptTemplate   string   //采集任务转换脚本
	db                      DBConfig //数据库连接信息
	callDb                  DBConfig //数据库连接信息
}

var config = &Config{}

func FetchConfig() {

	// 读取响应
	var (
		gatherJob Job
		callJob   Job
	)

	config.Set("gatherJobScriptTemplate", gatherJob.getValueByName("convert_script", "config_data_convert").Content) //采集任务转换脚本
	config.Set("callJobScriptTemplate", callJob.getValueByName("convert_script", "config_data_convert").Content)     //回调任务转换脚本
	if config.Get("gatherJobConnType") == dbConn {
		config.Set("db", DBConfig{
			DriverName: gatherJob.getValueByName("driver_name", "config_source_db").Content,
			User:       gatherJob.getValueByName("username", "config_source_db").Content,
			Passwd:     gatherJob.getValueByName("password", "config_source_db").Content,
			Addr:       gatherJob.getValueByName("server_addr", "config_source_db").Content,
			DBName:     gatherJob.getValueByName("database", "config_source_db").Content,
			Port:       gatherJob.getValueByName("port", "config_source_db").Content,
			QuerySql:   getSQLFromRule([]byte(gatherJob.getValueByName("query_rule_json", "config_source_db").Content), "querySql"),
			UpdateSql:  getSQLFromRule([]byte(gatherJob.getValueByName("update_rule_json", "config_source_db").Content), "updateSql"),
		})
	}
	if config.Get("callbackJobConnType") == dbConn {
		config.Set("callDb", DBConfig{
			DriverName: callJob.getValueByName("driver_name", "config_source_db").Content,
			User:       callJob.getValueByName("username", "config_source_db").Content,
			Passwd:     callJob.getValueByName("password", "config_source_db").Content,
			Addr:       callJob.getValueByName("server_addr", "config_source_db").Content,
			DBName:     callJob.getValueByName("database", "config_source_db").Content,
			Port:       callJob.getValueByName("port", "config_source_db").Content,
			QuerySql:   getSQLFromRule([]byte(callJob.getValueByName("query_rule_json", "config_source_db").Content), "querySql"),
			UpdateSql:  getSQLFromRule([]byte(callJob.getValueByName("update_rule_json", "config_source_db").Content), "updateSql"),
		})
	}
	//调用定时任务，并且重载InvoCrontab的配置
	InvoCrontab("* * * * * *")
}

func getSQLFromRule(ruleData []byte, fieldName string) string {
	var err error
	var rules map[string]interface{}
	funcName := "getSQLFromRule"
	if err = json.Unmarshal(ruleData, &rules); err != nil {
		fmt.Printf("%s err(%+v)\n", funcName, err)
		return ""
	}

	return getValueFromNestedMap(rules, strings.Split(fieldName, "."))
}

func getValueFromNestedMap(rules map[string]interface{}, fieldNames []string) string {
	funcName := "getValueFromNestedMap"
	if len(fieldNames) == 0 {
		return ""
	}

	key := fieldNames[0]
	if value, exists := rules[key]; exists {
		if len(fieldNames) == 1 {
			if str, ok := value.(string); ok {
				return str
			}
			fmt.Printf("%s error: field %s is not a string\n", funcName, key)
			return ""
		}

		// Recurse for nested maps
		if nestedMap, ok := value.(map[string]interface{}); ok {
			return getValueFromNestedMap(nestedMap, fieldNames[1:])
		}
		fmt.Printf("%s error: field %s is not a map\n", funcName, key)
		return ""
	}
	fmt.Printf("%s error: field %s does not exist\n", funcName, key)
	return ""
}

func (c *Config) Set(key string, value interface{}) {
	c.Lock()
	defer c.Unlock()

	switch key {
	case "gatherJobScriptTemplate":
		c.gatherJobScriptTemplate = value.(string)
	case "callJobScriptTemplate":
		c.callJobScriptTemplate = value.(string)
	case "db":
		c.db = value.(DBConfig)
	case "callDb":
		c.callDb = value.(DBConfig)
	default:
		fmt.Println("Unsupported key:", key)
	}
}

func (c *Config) Get(key string) interface{} {
	c.RLock()
	defer c.RUnlock()

	switch key {
	case "gatherJobScriptTemplate":
		return c.gatherJobScriptTemplate
	case "callJobScriptTemplate":
		return c.callJobScriptTemplate
	case "db":
		return c.db
	case "callDb":
		return c.callDb
	default:
		fmt.Println("Unsupported key:", key)
		return nil
	}
}
