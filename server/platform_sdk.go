package server

import (
	"GoSimProject/lib"
	"encoding/json"
	"fmt"
	"time"
)

type SDKRequest struct {
	GroovyCode string
	GroovyName string
}

type SDKClient struct {
	SDKUrl             string
	GetTokenGroovy     SDKRequest
	RefreshTokenGroovy SDKRequest
	GetDataGroovy      SDKRequest
}

func (sdk SDKClient) WriteBack(resp string) {
	var err error
	if config.Get("callSdkToken") == "" {
		var (
			sdkToken     string
			refreshToken string
			refreshTime  time.Time
		)
		sdkToken, refreshToken, refreshTime, err = getSdkToken(sdk.GetTokenGroovy.GroovyCode, sdk.GetTokenGroovy.GroovyName, sdk.SDKUrl)
		if err != nil {
			fmt.Println("获取sdkToken失败")
			return
		}
		config.Set("callSdkToken", sdkToken)
		config.Set("callSdkRefreshExpireAt", refreshTime)
		config.Set("callSdkRefreshToken", refreshToken)
	}
	lib.JsHandler(config.Get("callJobScriptTemplate").(string), resp)
	//使用JsHandler处理的数据进行接下来的sdk请求
	return
}

func (sdk SDKClient) FetchData() []map[string]interface{} {
	var err error
	if config.Get("sdkToken") == "" {
		var (
			sdkToken     string
			refreshToken string
			refreshTime  time.Time
		)
		sdkToken, refreshToken, refreshTime, err = getSdkToken(sdk.GetTokenGroovy.GroovyCode, sdk.GetTokenGroovy.GroovyName, sdk.SDKUrl)
		if err != nil {
			fmt.Println("获取sdkToken失败")
		}
		config.Set("sdkToken", sdkToken)
		config.Set("sdkRefreshExpireAt", refreshTime)
		config.Set("sdkRefreshToken", refreshToken)
	}
	var response []map[string]interface{}
	//通过sdk获取数据然后返回
	return response
}

func getSdkToken(GroovyCode string, GroovyName string, SDKUrl string) (sdkToken string, refreshToken string, refreshTime time.Time, err error) {
	type Result struct {
		ExpireAt        int64  `json:"expireAt"`
		RefreshExpireAt int64  `json:"refreshExpireAt"`
		RefreshToken    string `json:"refreshToken"`
		Token           string `json:"token"`
	}

	type Response struct {
		Result Result `json:"result"`
	}
	var response Response
	//result是通过请求sdk获取的token
	var result string

	jsonData, err := json.Marshal(result)
	if err != nil {
		fmt.Println("Error marshalling map to JSON:", err)
		return
	}
	err = json.Unmarshal(jsonData, &response)
	if err != nil {
		return "", "", time.Time{}, err
	}

	// 将毫秒转换为秒
	seconds := response.Result.RefreshExpireAt / 1000

	// 纳秒部分（余数部分转换为纳秒）
	nanoseconds := (response.Result.RefreshExpireAt % 1000) * int64(time.Millisecond)
	sdkToken = response.Result.Token
	refreshToken = response.Result.RefreshToken
	refreshTime = time.Unix(seconds, nanoseconds)
	return
}
