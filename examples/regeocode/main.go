package main

import (
	"fmt"
	"log"
	"time"

	"github.com/enneket/amap"
	"github.com/enneket/amap/api/re_geo_code"
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

	// 逆地理编码示例
	req := &re_geo_code.ReGeocodeRequest{
		Location:   "116.481028,39.989643",
		Radius:     1000,
		Extensions: "all",
	}

	resp, err := client.ReGeocode(req)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("=== 逆地理编码结果 ===")
	regeo := resp.ReGeocode
	fmt.Printf("地址: %s\n", regeo.FormattedAddress)
	fmt.Printf("省: %s\n", regeo.AddressComponent.Province)
	fmt.Printf("市: %s\n", regeo.AddressComponent.City)
	fmt.Printf("区: %s\n", regeo.AddressComponent.District)
	fmt.Printf("城镇: %s\n", regeo.AddressComponent.Township)
	fmt.Printf("城市代码: %s\n", regeo.AddressComponent.Citycode)
	fmt.Printf("行政区划代码: %s\n", regeo.AddressComponent.Adcode)
}
