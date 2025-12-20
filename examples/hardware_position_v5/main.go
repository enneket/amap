package main

import (
	"fmt"
	"log"
	"time"

	"github.com/enneket/amap"
	position_v5 "github.com/enneket/amap/api/position/v5"
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

	// 硬件定位V5示例
	req := &position_v5.HardwarePositionRequest{
		WiFi:       "mac:ssid,rssi;mac2:ssid2,rssi2",
	}

	resp, err := client.HardwarePositionV5(req)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("=== 硬件定位V5结果 ===")
	fmt.Printf("精度: %d 米\n", resp.Accuracy)
}
