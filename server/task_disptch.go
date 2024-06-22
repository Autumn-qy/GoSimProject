package server

import (
	"sync"
	"time"
)

// InvoCrontab 作用是接收一个cron配置，判断是否与原先配置相同，不同进行停止原先定时任务，重启一个新的定时任务
func InvoCrontab(cron string) {

	err := store()
	if err != nil {
		return
	}

}

// Worker 函数从任务通道中读取任务，并使用信号量控制任务执行
func Worker() {
	var wg sync.WaitGroup
	maxConcurrentWorkers := 10                             // 最大并发工作协程数量
	semaphore := make(chan struct{}, maxConcurrentWorkers) // 创建一个有缓冲的信号量通道
	invoiceApi := config.Get("invoiceApi").(APIClient)
	//通过判断某需要的消息是否为空，如果为空则等待
	for {
		invoiceApi = config.Get("invoiceApi").(APIClient)
		if invoiceApi.InterfaceAddr != "" {
			break
		}
		//等待一秒避免大量消耗CPU资源
		time.Sleep(1 * time.Second)
	}
	for {
		invoData := make([]map[string]interface{}, 10)
		if invoData != nil && len(invoData) != 0 {
			for _, data := range invoData {
				wg.Add(1)
				go invoiceAdd(data["invo_invoicingData"].(string), &wg, semaphore, data["invo_id"].(string))
			}
		}
		wg.Wait()
		//和上方同样效果
		time.Sleep(1 * time.Second)
	}
}

type SafeWriter struct {
	mu sync.Mutex
}

var writer = &SafeWriter{}

func store() error {
	var err error
	var wg sync.WaitGroup
	//进行加解锁，可以想下为什么要加锁
	writer.mu.Lock()
	defer writer.mu.Unlock()

	var service InvoiceService
	if config.Get("gatherJobConnType") == dbConn {
		service = config.Get("db").(DBConfig)
	} else if config.Get("gatherJobConnType") == apiConn {
		service = config.Get("api").(APIClient)
	} else if config.Get("gatherJobConnType") == sdkConn {
		service = config.Get("sdk").(SDKClient)
	}
	if service == nil {
		return err
	}
	var results []map[string]interface{}
	results = service.FetchData()
	//如果数据集为空则不继续处理
	if results == nil {
		return nil
	}
	for _, result := range results {
		wg.Add(1)
		Run(result, &wg)
	}
	wg.Wait()
	return err
}
