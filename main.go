package main

import (
	"fmt"
	"net/http"

	"github.com/eric0571/tino/util/config"
	"github.com/eric0571/tino/util/system"

	"github.com/labstack/echo"
)

func UpdatePlayTiming() {
	host := "https://api.shunliandongli.com/v2/admin/"
	url := host + "live/liveList?key=&type=&value=&is_rec=&role_id=&page=1&page_size=10&total=&sort_value="
	headers := map[string]string{}
	cookie := map[string]string{}
	cookie["hashToken"] = "d97bd856617a057ee075377ba258bc14"
	ret, _ := system.HTTPPostJSON(url, string(""), headers, cookie)
	//fmt.Println("Query Result:", ret)
	//system.PrintInterface(ret)
	data1 := ret.(map[string]interface{})
	data := data1["data"].(map[string]interface{})
	list := data["list"].([]interface{})
	for i := 0; i < len(list); i++ {
		one := list[i].(map[string]interface{})
		fmt.Println("-------------------------------------------")
		fmt.Println(one["room_id"], one["id"])
		roomId := one["room_id"].(string)
		url = host + "live/updateRoomData?room_id=" + roomId
		retSubData, _ := system.HTTPPostJSON(url, string(""), headers, cookie)
		data2 := retSubData.(map[string]interface{})
		//system.PrintInterface(data2)
		code := data2["code"].(float64)
		fmt.Println("code:", code)
	}
}

func main() {
	UpdatePlayTiming()
	e := echo.New()
	e.GET("/", func(c echo.Context) error {
		return c.String(http.StatusOK, "Hello, World!")
	})
	system.JSONPrint("Ts", "22")
	config.HelloWorld()
	//e.Logger.Fatal(e.Start(":1323"))
}
