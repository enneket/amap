package main

import (
	"fmt"
	"log"
	"time"

	"github.com/enneket/amap"
	"github.com/enneket/amap/api/bus/station_keyword"
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

	// 公交站关键字查询示例
	req := &station_keyword.StationKeywordRequest{
		Keywords: "天安门",
		City:     "北京",
		Offset:   "10",
		Page:     "1",
	}

	resp, err := client.BusStationKeyword(req)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("=== 公交站关键字查询结果 ===")
	fmt.Printf("总结果数: %d\n", resp.Count)
	for i, station := range resp.Stations {
		fmt.Printf("公交站 %d:\n", i+1)
		fmt.Printf("名称: %s\n", station.Name)
		fmt.Printf("经纬度: %s\n", station.Location)
		fmt.Printf("地址: %s\n", station.Address)
		fmt.Println()
	}
}
