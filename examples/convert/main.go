package main

import (
	"fmt"
	"log"
	"time"

	"github.com/enneket/amap"
	"github.com/enneket/amap/api/convert"
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

	// 坐标转换示例
	req := &convert.ConvertRequest{
		Locations: "116.481028,39.989643|116.514203,39.905409",
		CoordSys:  "gps",
		Output:    "json",
	}

	resp, err := client.Convert(req)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("=== 坐标转换结果 ===")
	fmt.Printf("转换状态: %s\n", resp.Status)
	fmt.Printf("转换结果: %s\n", resp.Locations)
}
