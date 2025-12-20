package main

import (
	"fmt"
	"log"
	"time"

	"github.com/enneket/amap"
	"github.com/enneket/amap/api/grasproad"
)

func main() {
	// 创建配置
	config := amap.NewConfig("your_amap_api_key")
	config.SecurityKey = "your_amap_security_key"
	config.Timeout = 30 * time.Second

	// 初始化客户端
	client, err := amap.NewClient(config)
	if err != nil {
		log.Fatal(err)
	}

	// 轨迹纠偏示例
	req := &grasproad.GraspRoadRequest{
		SID: "test_sid",
		Points:          "116.481028,39.989643,1546831454;116.481068,39.989683,1546831455;116.481198,39.989673,1546831456",
		CoordTypeInput: "gps",
	}

	resp, err := client.GraspRoad(req)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("=== 轨迹纠偏结果 ===")
	fmt.Printf("状态: %s\n", resp.Status)
}
