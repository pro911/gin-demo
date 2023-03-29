package main

import (
	"fmt"
	"github.com/robfig/cron/v3"
	"time"
)

func main() {
	// 创建一个新的Cron对象
	c := cron.New()

	// 添加一个每秒执行一次的任务
	c.AddFunc("@every 1s", func() {
		fmt.Println("Task executed at", time.Now())
	})

	// 启动Cron
	c.Start()

	// 等待一段时间，以便查看输出结果
	time.Sleep(time.Second * 1000)

	// 停止Cron
	c.Stop()
}
