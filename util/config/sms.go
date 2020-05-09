package config

import (
	"encoding/json"
	"fmt"

	"crypto/hmac"
	"crypto/md5"
	"crypto/sha256"
	"encoding/hex"
	"io"
	"io/ioutil"
	"math/rand"
	"net/http"
	"strings"
	"time"

	"github.com/goroom/aliyun_sms"
	"github.com/goroom/logger"

	"crypto/tls"
)

// SMSInfo is the details for the SMS server
type SMSInfo struct {
	SignName        string
	TempletCode     string
	AccessKey       string
	AccessKeySecret string
}

// SendMobileMessage sends an email
func (e SMSInfo) SendMobileMessage(nationcode, to, code string) error {
	nationcode = strings.Trim(nationcode, " ")
	nationcode = strings.Trim(nationcode, "+")
	UsePhp := false
	if UsePhp {
		URL := `http://iot.roabay.com/lib/aliyun-php-sdk-sms/sms.php`
		text := fmt.Sprintf("%s%s%s", to, code, "viewmobile")
		key := GetMD5([]byte(text))
		Param := fmt.Sprintf("?to=%s&code=%s&key=%s&v=%s", to, code, key, GetRandomString(5))
		URL = fmt.Sprintf("%s%s", URL, Param)
		_, err := SmsRequestServer(URL, "", "GET")
		return err
	} else {
		if nationcode == "86" {
			TempletCode := e.TempletCode
			AccessKey := e.AccessKey
			AccessKeySecret := e.AccessKeySecret
			SignName := e.SignName
			aliyun_sms, err := aliyun_sms.NewAliyunSms(SignName, TempletCode, AccessKey, AccessKeySecret)
			if err != nil {
				logger.Error(err)
				return err
			}
			cjson := `{"code":"%s"}`
			data := fmt.Sprintf(cjson, code)
			err = aliyun_sms.Send(to, data)
			if err != nil {
				logger.Error(err)
				return err
			}
			logger.Error("Success")
		} else {
			//腾讯云发送国际短信national mobile sms
			appid := "1400046321"
			strMobile := to                                 //tel的mobile字段的内容
			strAppKey := "472b2bac4c6e79cd9340fabee51283b6" //sdkappid对应的appkey，需要业务方高度保密
			strRand := GetRandomNumberString(10)            //url中的random字段的值
			strTime := time.Now().Unix()                    //unix时间戳
			strSrc := fmt.Sprintf("appkey=%s&random=%s&time=%d&mobile=%s", strAppKey, strRand, strTime, strMobile)
			sig := getSha256Code(strSrc)
			//			fmt.Println("SignSrc:", strSrc)
			//			fmt.Println("Str Sig:", sig)
			url := fmt.Sprintf("https://yun.tim.qq.com/v5/tlssmssvr/sendsms?sdkappid=%s&random=%s", appid, strRand)
			params := make([]string, 1)
			params[0] = code
			tel := make(map[string]interface{})
			tel["nationcode"] = nationcode
			tel["mobile"] = strMobile
			data := make(map[string]interface{})
			data["tel"] = tel
			data["type"] = 0
			if nationcode == "86" {
				data["tpl_id"] = 53238
				data["sign"] = "络奇"
			} else {
				data["tpl_id"] = 53237
				data["sign"] = "Looki"
			}

			data["params"] = params
			// data["msg"] = fmt.Sprintf("[Looki] You vcode is:%s", code)
			data["sig"] = string(sig)
			data["time"] = strTime
			data["extend"] = ""
			data["ext"] = ""
			jdata, _ := json.Marshal(data)
			//			fmt.Println(url, strAppKey)
			strRet, err := HTTPSPost(url, string(jdata))
			fmt.Println(string(jdata))
			fmt.Println("Request Ret:", string(strRet))
			if err == nil {
				var jsonData interface{}
				json.Unmarshal(strRet, &jsonData)
				retData := jsonData.(map[string]interface{})
				if retData["errmsg"].(string) != "OK" {
					err1 := fmt.Errorf("Send sms error")
					//					fmt.Println("Error:", retData["errmsg"], " result:", retData["result"])
					fmt.Println("Error:", retData)
					return err1
				}
			} else {
				err1 := fmt.Errorf("HTTPS Send sms error:%s", err)
				logger.Error("Https send error:", err)
				return err1
			}
		}

	}
	return nil
}

func SmsRequestServer(url, data, method string) (string, error) {
	client := &http.Client{}
	req, err := http.NewRequest(method, url, strings.NewReader(data))
	req.Header.Add("User-Agent", "Aliyun sms")
	req.Header.Add("Content-Type", "application/x-www-form-urlencoded")

	resp, err := client.Do(req)
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		// handle error
		fmt.Println("Send sms http request error:", err)
	}

	return string(body), err
}
func getHmacCode(s string) string {
	h := hmac.New(sha256.New, []byte("ourkey"))
	io.WriteString(h, s)
	return fmt.Sprintf("%x", h.Sum(nil))
}

func getSha256Code(s string) string {
	h := sha256.New()
	h.Write([]byte(s))
	return fmt.Sprintf("%x", h.Sum(nil))
}
func HTTPSPost(url, data string) ([]byte, error) {
	tr := &http.Transport{
		TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
	}
	client := &http.Client{Transport: tr}
	req, err := http.NewRequest("POST", url, strings.NewReader(data))
	//	req.Header["authtoken"] = []string{token}
	req.Header.Add("Content-Type", "application/json")
	req.Header.Add("Accept", "*/*")
	req.Header.Add("Connection", "Keep-Alive")
	req.Header.Add("User-Agent", "Golang https client")
	req.Header.Add("Cache-control", "no-cache")
	req.Header.Set("Accept-Charset", "utf-8")
	req.Header.Set("Accept-Encoding", "gzip,deflate,br")

	resp, err := client.Do(req)
	if err != nil {
		fmt.Println(err)
		return nil, err
	}
	defer resp.Body.Close()
	return ioutil.ReadAll(resp.Body)
}

// GetRandomString 生成随机字符串
func GetRandomString(length int) string {
	str := "0123456789abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"
	bytes := []byte(str)
	result := []byte{}
	r := rand.New(rand.NewSource(time.Now().UnixNano()))
	for i := 0; i < length; i++ {
		result = append(result, bytes[r.Intn(len(bytes))])
	}
	return string(result)
}

func GetRandomNumberString(length int) string {
	str := "0123456789"
	bytes := []byte(str)
	result := []byte{}
	r := rand.New(rand.NewSource(time.Now().UnixNano()))
	for i := 0; i < length; i++ {
		result = append(result, bytes[r.Intn(len(bytes))])
	}
	return string(result)
}

func GetMD5(text []byte) string {
	sum := md5.Sum(text)

	return hex.EncodeToString(sum[:])
}
