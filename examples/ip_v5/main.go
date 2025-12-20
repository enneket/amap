package main

import (
	"fmt"
	"log"
	"time"

	"github.com/enneket/amap"
	ip_v5 "github.com/enneket/amap/api/ip/v5"
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

	// IP定位V5示例
	req := &ip_v5.IPConfigRequest{
		IP: "114.247.50.2",
	}

	resp, err := client.IPV5Config(req)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("=== IP定位V5结果 ===")
	fmt.Printf("省份: %s\n", resp.Province)
	fmt.Printf("城市: %s\n", resp.City)
	fmt.Printf("区: %s\n", resp.District)
	fmt.Printf("经纬度: %s\n", resp.Location)
	fmt.Printf("ISP: %s\n", resp.ISP)
}
