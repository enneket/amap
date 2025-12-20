package main

import (
	"fmt"
	"log"
	"time"

	"github.com/enneket/amap"
	"github.com/enneket/amap/api/weatherinfo"
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

	// 天气信息示例
	req := &weatherinfo.WeatherinfoRequest{
		City:       "北京",
		Extensions: "base",
	}

	resp, err := client.Weatherinfo(req)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("=== 天气信息结果 ===")
	fmt.Printf("城市: %s\n", resp.Weatherinfo.City)
	fmt.Printf("天气: %s\n", resp.Weatherinfo.Weather)
	fmt.Printf("温度: %s°C\n", resp.Weatherinfo.Temp)
	fmt.Printf("风向: %s\n", resp.Weatherinfo.Winddirection)
	fmt.Printf("风力: %s\n", resp.Weatherinfo.Windpower)
	fmt.Printf("湿度: %s\n", resp.Weatherinfo.Humidity)
	fmt.Printf("发布时间: %s\n", resp.Weatherinfo.Time)
}
