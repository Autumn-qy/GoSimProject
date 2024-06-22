package server

import (
	"fmt"
	"sync"
)

func invoiceAdd(jsonData string, wg *sync.WaitGroup, semaphore chan struct{}, id string) {
	defer wg.Done()
	// 从信号量通道获取一个令牌
	semaphore <- struct{}{}
	defer func() {
		// 将令牌归还到信号量通道
		<-semaphore
		fmt.Println("任务结束", id)
	}()
	fmt.Println("任务开始", id)
	//下面是开发票的实现
	return
}

// Active 主动查询发票
func Active(serialNo string) {
	//主动查询发票的实现
}
