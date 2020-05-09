package system

import (
	"crypto/sha256"
	"encoding/json"
	"fmt"
)

//const APPID = "BuP97o68mC8Ty7ZGhlvBX3"
//const APPKEY = "WttOxmmTvfAUEAzK4sZPQ3"
//const MASTERSECRET = "zziz9gcP7D9GnKAqpHJMs5"
const GETTUIKEY = "GetTuiThirdToken"
const APPID = "npn3aAuSSs6L9SXBjPmej8"
const APPKEY = "4b1ufRe3o4ABGgirSmuyT4"
const MASTERSECRET = "wEir6XHGE97gK0HHDOfYE1"

var AuthToken string
var Title, Context string
var rdata []byte

// RspBody 个推Rsp body
// 个推请求返回的结构
// status : successed_offline 离线下发
//          successed_online 在线下发
//          successed_ignore 非活跃用户不下发
type RspBody struct {
	Result    string `json:"result"`
	TaskID    string `json:"taskid"`
	Desc      string `json:"desc"`
	Status    string `json:"status"`
	RequestID string `json:"requestID,omitempty"`
}
type MessageInput struct {
	Type    string `json:"type"`
	Scope   string `json:"scope"`
	Link    string `json:"link"`
	Logo    string `json:"logo"`
	Title   string `json:"title"`
	Context string `json:"context"`
	Cid     string `json:"cid"` //49e9cc0df7b75dcfa57e32dca40acc69
	Key     string `json:"key"`
}

// GeTui message
func GeTui(inputMsg MessageInput) error {
	AuthToken = GetAuthToken()
	Title = inputMsg.Title
	Context = inputMsg.Context
	if inputMsg.Type == "ALL" {
		PushApp()
	} else {
		PushSingle(inputMsg.Cid)
	}
	return nil
}

//Send to single user infomation
func NoticeSingleUserApp(cid, title, context string) error {
	//	return nil //close send getui

	Title = title
	Context = context
	AuthToken = GetAuthToken()
	PushSingle(cid)
	CloseAuth()
	return nil
}

func Gui() error {
	rdata = make([]byte, 1000)

	Title = "My Test title"
	Context = "从gocode发行包或者gocode的源代码库中拷贝emacs/go-autocomplete.ele文件到~/.emacs.d目录"
	AuthToken = GetAuthToken()
	PushApp()
	PushSingle("49e9cc0df7b75dcfa57e32dca40acc69")
	CloseAuth()
	return nil
}

func ParseJson(data []byte) (interface{}, error) {
	var jsonData interface{}
	err := json.Unmarshal(data, &jsonData)
	return jsonData, err
}

func CloseAuth() string {
	url := fmt.Sprintf("https://restapi.getui.com/v1/%s/auth_close", APPID)
	rspBody, err := HTTPSPost(url, "", AuthToken)
	if err != nil {
		return ""
	}
	fmt.Println(string(rspBody))
	ret := &RspBody{}
	err = json.Unmarshal(rspBody, ret)
	if err != nil {
		return fmt.Sprintf("[CloseAuth] 清空auth 请求返回的JSON无法解析, err: %s", err)
	}

	if ret.Result != "ok" {
		return fmt.Sprintf("[CloseAuth] 清空auth 失败, desc: %s", ret.Desc)
	}

	return string("ok")
}

func GetAuthToken() string {
	url := fmt.Sprintf("https://restapi.getui.com/v1/%s/auth_sign", APPID)
	sendDataJSON := `{
			"sign"           : "%s",
			"timestamp": "%d",
			"appkey"      : "%s" }`
	time := GetTimeInt()
	time = time * 1000

	wdata := fmt.Sprintf("%s%d%s", APPKEY, time, MASTERSECRET)
	sign := sha256.Sum256([]byte(wdata))
	ss := fmt.Sprintf("%x", sign)
	js := fmt.Sprintf(sendDataJSON, ss, time, APPKEY)

	ret, err := HTTPSPost(url, js, AuthToken)
	if err != nil {
		fmt.Println(rdata)
		return ""
	} else {
		if len(ret) > 0 {
			rd, _ := ParseJson(ret)
			data := rd.(map[string]interface{})
			isOk := data["result"].(string)
			if isOk == "ok" {
				return data["auth_token"].(string)
			}
		} else {
			fmt.Println("-----------------------no data------------------")
			return ""
		}
	}
	return ""
}

func PushSingle(cid string) error {
	json := `{
	                    "message":{
	                        "appkey":"%s",
	                        "is_offline":false,
	                        "msgtype":"transmission"
	                    },
	                    "transmission":{
	                        "transmission_type":false,
	                        "transmission_content":"{\"msgType\":\"emailConfirm\",\"actived\":true}",
	                        "duration_begin":"%s",
	                        "duration_end":"%s"
	                    },
	                    "push_info": {
	                            "aps": {
	                                "alert": {
	                                    "title": "%s",
	                                    "body": "%s"
	                                },
	                                "autoBadge": "+1",
	                                "content-available": 1
	                            }
	                         
	                        },
	                    "cid":"%s",
	                    "requestid":"%s"
	                 }`

	url := fmt.Sprintf("https://restapi.getui.com/v1/%s/push_single", APPID)
	time := GetTime()
	js := fmt.Sprintf(json, APPKEY, time, time, Title, Context, cid, GetRandomString(30))

	rdata, err := HTTPSPost(url, js, AuthToken)
	if err == nil {
		fmt.Println("Push Single return:", string(rdata))
	}
	return err
}

func SaveList() error {
	json := `{
                "message": {
                   "appkey": "%s",
                   "is_offline": true,
                   "offline_expire_time":10000000,
                   "msgtype": "notification"
                },
                "notification": {
                    "style": {
                        "type": 0,
                        "text": "text",
                        "title": "%s"
                    },
                    "transmission_type": true,
                    "transmission_content": "%s"
                }
           }`

	url := fmt.Sprintf("https://restapi.getui.com/v1/%s/save_list_body", APPID)
	js := fmt.Sprintf(json, APPKEY, Title, Context)

	rdata, err := HTTPSPost(url, js, AuthToken)
	if err == nil {
		fmt.Println("SaveList Ret:", string(rdata))
	}
	return err
}

func PushList(taskId string) error {
	json := `{"cid":["%s"],"taskid":"%s", "need_detail":true }`
	url := fmt.Sprintf("https://restapi.getui.com/v1/%s/push_list", APPID)
	js := fmt.Sprintf(json, "", taskId)
	rdata, err := HTTPSPost(url, js, AuthToken)
	if err == nil {
		fmt.Println("PushList return:", string(rdata))
	}
	return err
}

func PushApp() error {
	json := `{
                "message": {
                   "appkey": "%s",
                   "is_offline": false,
                   "msgtype": "notification"
                },
                "notification": {
                    "style": {
                        "type": 0,
                        "text": "text",
                        "title": "%s"
                    },
                    "transmission_type": true,
                    "transmission_content": "%s"
                },
              "condition":[{"key":"phonetype", "values":["ANDROID"], "opt_type":0}],
                "requestid":"%s"
                }`

	url := fmt.Sprintf("https://restapi.getui.com/v1/%s/push_app", APPID)
	js := fmt.Sprintf(json, APPKEY, Title, Context, GetRandomString(30))

	ret, err := HTTPSPost(url, js, AuthToken)
	if err == nil {
		fmt.Println("Push App Return:", string(ret))
	}
	return err
}
