package main

import (
	"fmt"
	"log"
	"time"

	"github.com/enneket/amap"
	"github.com/enneket/amap/api/bus/line_id"
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

	// 公交路线ID查询示例
	req := &line_id.LineIDRequest{
		ID: "B000A837X5",
		City: "北京",
	}

	resp, err := client.BusLineID(req)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("=== 公交路线ID查询结果 ===")
	fmt.Printf("状态: %s\n", resp.Status)
}
