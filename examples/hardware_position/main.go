package main

import (
	"fmt"
	"log"
	"time"

	"github.com/enneket/amap"
	position_v1 "github.com/enneket/amap/api/position/v1"
)

func main() {
	// 创建配置
	config := amap.NewConfig("your_amap_api_key")
	config.Timeout = 30 * time.Second

	// 初始化客户端
	client, err := amap.NewClient(config)
	if err != nil {
		log.Fatal(err)
	}

	// 硬件定位V1示例
	req := &position_v1.HardwarePositionRequest{
		WiFi:       "mac:ssid,rssi;mac2:ssid2,rssi2",
	}

	resp, err := client.HardwarePosition(req)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("=== 硬件定位V1结果 ===")
	fmt.Printf("精度: %d 米\n", resp.Accuracy)
}
