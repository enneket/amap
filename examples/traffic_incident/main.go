package main

import (
	"fmt"
	"log"
	"time"

	"github.com/enneket/amap"
	"github.com/enneket/amap/api/traffic_incident"
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

	// 交通事件查询示例
	req := &traffic_incident.TrafficIncidentRequest{
		Level:  "0",
		Type:   "0",
		Rectangle: "116.351147,39.966309,116.357136,39.968722",
	}

	resp, err := client.TrafficIncident(req)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("=== 交通事件查询结果 ===")
	fmt.Printf("状态: %s\n", resp.Status)
}
