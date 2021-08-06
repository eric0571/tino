package main

import (
	"fmt"

	"github.com/eric0571/tino/util/system"
)

const Host string = "http://api-front.v2.shunliandev.com"
const sendSms string = "/member/common/sendSmsCode?_v=3.1.1"

func UpdatePlayTiming() {
	host := "https://api.shunliandongli.com/v2/admin/"
	url := host + "live/liveList?key=&type=&value=&is_rec=&role_id=&page=1&page_size=10&total=&sort_value="
	headers := map[string]string{}
	cookie := map[string]string{}
	cookie["hashToken"] = "d97bd856617a057ee075377ba258bc14"
	fmt.Println(url)
	ret, _ := system.HTTPPostJSON(url, string(""), headers, cookie)
	fmt.Println("Query Result:", ret)
	system.PrintInterface(ret)
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

func Login() {
	TestHeader := map[string]string{}
	TestCookie := map[string]string{}
	url := Host + sendSms
	TestHeader["UserNameAgent"] = "kkkkkkk"
	TestCookie["hashToken"] = "d97bd856617a057ee075377ba258bc14"
	sendDataJSON := `{"mobile":"13588151390","msg":"hello."}`
	system.JSONPrint("SendDataJson:", sendDataJSON)
	//JsonPrint("SendDataJson:", sendDataJSON)
	//data := "&mobile=13588151392"
	//fmt.Println("Url:", url)
	ret, _ := system.HTTPPostJSON(url, sendDataJSON, TestHeader, TestCookie)
	system.PrintInterface(ret)
	// data1 := ret.(map[string]interface{})
	// data := data1["data"].(map[string]interface{})
	// code := data["code"].(string)
	// fmt.Println("code:", code)
}

func PostData() {
	Header := map[string]string{}
	Cookie := map[string]string{}
	Data := map[string]string{}
	url := "http://api.zko.shunliandev.com/customs/setPlatOrder"
	Data["UserNameAgent"] = "TestAngent"
	Data["id"] = "2"
	Data["push_time"] = fmt.Sprintf("%d", system.GetTimeInt())
	Cookie["hashToken"] = "d97bd856617a057ee075377ba258bc14"

	fmt.Println("Url:", url)
	fmt.Println("Data:")
	system.MapPrint(Data)
	ret, _ := system.HTTPPostMulForm(url, Data, Header, Cookie)
	fmt.Println("Return ret is :")
	fmt.Println(string(ret))
}

func PostWwwFromData() {
	Header := map[string]string{}
	Cookie := map[string]string{}
	Data := map[string]string{}
	url := "http://api.zko.shunliandev.com/customs/setPlatOrder"
	Data["UserNameAgent"] = "TestAngent"
	Data["id"] = "2"
	Data["jasonData"] = `{"name":"wanglinhui","sex":"m","age":40}`
	Data["push_time"] = fmt.Sprintf("%d", system.GetTimeInt())
	Cookie["hashToken"] = "d97bd856617a057ee075377ba258bc14"

	fmt.Println("Url:", url)
	fmt.Println("Data:")
	system.MapPrint(Data)
	ret, _ := system.HTTPPostWForm(url, Data, Header, Cookie)
	fmt.Println("Return ret is :")
	fmt.Println(string(ret))
}

func main() {
	// PostData()
	PostWwwFromData()
	//Login()
	//UpdatePlayTiming()
	// e := echo.New()
	// e.GET("/", func(c echo.Context) error {
	// 	return c.String(http.StatusOK, "Hello, World!")
	// })
	//system.JSONPrint("Ts", "22")
	//config.HelloWorld()
	//e.Logger.Fatal(e.Start(":1323"))
}
