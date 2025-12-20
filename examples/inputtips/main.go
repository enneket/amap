package main

import (
	"fmt"
	"log"
	"time"

	"github.com/enneket/amap"
	inputtips "github.com/enneket/amap/api/input_tips"
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

	// 输入提示示例
	req := &inputtips.InputtipsRequest{
		Keywords: "天安门",
		City:     "北京",
		Datatype: "poi",
	}

	resp, err := client.Inputtips(req)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("=== 输入提示结果 ===")
	fmt.Printf("结果数量: %d\n", resp.Count)
	for i, tip := range resp.Tips {
		fmt.Printf("提示 %d:\n", i+1)
		fmt.Printf("名称: %s\n", tip.Name)
		fmt.Printf("地址: %s\n", tip.Address)
		fmt.Printf("类型: %s\n", tip.Type)
		fmt.Printf("经纬度: %s\n", tip.Location)
		fmt.Printf("城市: %s\n", tip.City)
		fmt.Printf("区域: %s\n", tip.District)
		fmt.Println()
	}
}
