package server

import (
	"GoSimProject/lib"
	"encoding/json"
	"sync"
)

// Run 开发票函数
func Run(result map[string]interface{}, wg *sync.WaitGroup) {
	defer wg.Done()
	resultByte, _ := json.Marshal(result)
	//使用查询出的数据丢给js引擎处理，查询出的数据要存一下唯一标识，用来回写数据。
	lib.JsHandler(config.Get("gatherJobScriptTemplate").(string), string(resultByte))
	//通过quickjs引擎处理完的数据进行写表
}

func BatchActive() {
	invoData := make([]map[string]interface{}, 10)
	for _, data := range invoData {
		go Active(data["id"].(string))
	}
}
