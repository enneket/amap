package main

import (
	"fmt"
	"log"
	"time"

	"github.com/enneket/amap"
	"github.com/enneket/amap/api/geo_code"
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

	// 地理编码示例
	req := &geo_code.GeocodeRequest{
		Address: "北京市朝阳区阜通东大街6号",
		City:    "北京",
	}

	resp, err := client.GeoCode(req)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("=== 地理编码结果 ===")
	for i, geocode := range resp.Geocodes {
		fmt.Printf("结果 %d:\n", i+1)
		fmt.Printf("地址: %s\n", geocode.FormattedAddress)
		fmt.Printf("经纬度: %s\n", geocode.Location)
		fmt.Printf("省: %s\n", geocode.Province)
		fmt.Printf("市: %s\n", geocode.City)
		fmt.Printf("区: %s\n", geocode.District)
		fmt.Printf("城镇: %s\n", geocode.Township)
		fmt.Printf("街道: %s\n", geocode.Street)
		fmt.Printf("门牌号: %s\n", geocode.Number)
		fmt.Printf("城市代码: %s\n", geocode.Citycode)
		fmt.Printf("行政区划代码: %s\n", geocode.Adcode)
		fmt.Println()
	}
}
