package main

import (
	"fmt"
	"log"
	"time"

	"github.com/enneket/amap"
	"github.com/enneket/amap/api/bus/line_keyword"
)

func main() {
	// 创建配置
	config := &amap.Config{
		Key:     "your_amap_api_key",
		Timeout: 30 * time.Second,
	}

	// 初始化客户端
	client, err := amap.NewClient(config)
	if err != nil {
		log.Fatal(err)
	}

	// 公交路线关键字查询示例
	req := &line_keyword.LineKeywordRequest{
		Keywords: "1路",
		City:     "北京",
		Offset:   "10",
		Page:     "1",
	}

	resp, err := client.BusLineKeyword(req)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("=== 公交路线关键字查询结果 ===")
	fmt.Printf("总结果数: %d\n", resp.Count)
	for i, line := range resp.Lines {
		fmt.Printf("公交路线 %d:\n", i+1)
		fmt.Printf("名称: %s\n", line.Name)
		fmt.Printf("类型: %s\n", line.Type)
		fmt.Printf("首班车时间: %s\n", line.FirstTime)
		fmt.Printf("末班车时间: %s\n", line.LastTime)
		fmt.Println()
	}
}
