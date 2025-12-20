package main

import (
	"fmt"
	"log"
	"time"

	"github.com/enneket/amap"
	ip_v3 "github.com/enneket/amap/api/ip/v3"
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

	// IP定位V3示例
	req := &ip_v3.IPConfigRequest{
		IP: "114.247.50.2",
	}

	resp, err := client.IPConfig(req)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("=== IP定位V3结果 ===")
	fmt.Printf("状态: %s\n", resp.Status)
	fmt.Printf("省份: %s\n", resp.Province)
	fmt.Printf("城市: %s\n", resp.City)
	fmt.Printf("区域: %s\n", resp.District)
	fmt.Printf("ISP: %s\n", resp.ISP)
	fmt.Printf("经纬度: %s\n", resp.Location)
}
