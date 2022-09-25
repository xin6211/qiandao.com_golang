package request

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"time"
)

type RqData struct {
	OpenCheckInDataType int      `json:"opencheckindatatype"`
	StartTime           int      `json:"starttime"`
	EndTime             int      `json:"endtime"`
	UserList            []string `json:"useridlist"`
}

func AccessToken() (string, error) {
	getTokenTime := time.Now()
	var accessTokenUrl string
	corpId := "ww15631b0d64da0525"
	corpSecret := "wgFg1rZKpewyzyOpmy4kenGlTYD_16T7ij2VqYaAdMk"
	accessTokenUrl = fmt.Sprintf("https://qyapi.weixin.qq.com/cgi-bin/gettoken?corpid=%s&corpsecret=%s", corpId, corpSecret)
	resp, err := http.Get(accessTokenUrl)
	var accessToken string
	if err != nil {
		return "", fmt.Errorf("get access_token eeror!%v", err)
	}
	defer func(Body io.ReadCloser) {
		err := Body.Close()
		if err != nil {
			fmt.Printf("Get close error! %v", err)
		}
	}(resp.Body)
	var respCode = resp.StatusCode
	if respCode == 200 {
		respBody, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			return "", fmt.Errorf("decode body error! %v", err)
		}
		var respData = make(map[string]interface{})
		jsonDecodeErr := json.Unmarshal(respBody, &respData)
		if jsonDecodeErr != nil {
			return "", fmt.Errorf("json decode error! %v", jsonDecodeErr)
		}
		if respData["errmsg"] == "ok" {
			accessToken = respData["access_token"].(string)
		}
	}
	fmt.Printf("get token:%v\n", time.Since(getTokenTime))
	return accessToken, nil
}

func GetSourData(token string, rqJson *bytes.Reader) (map[string]interface{}, error) {
	getDataTime := time.Now()
	var dataUrl = "https://qyapi.weixin.qq.com/cgi-bin/checkin/getcheckindata?access_token="
	dataUrl = dataUrl + token
	resp, requestErr := http.Post(dataUrl, "application/json", rqJson)
	if requestErr != nil {
		return nil, requestErr
	}
	var data = make(map[string]interface{})
	respBody, readErr := ioutil.ReadAll(resp.Body)
	if readErr != nil {
		return nil, readErr
	}
	jsonErr := json.Unmarshal(respBody, &data)
	if jsonErr != nil {
		fmt.Printf("%v", jsonErr)
	}
	if data["errmsg"].(string) != "ok" {
		fmt.Printf("%v", data["errmsg"].(string))
		return nil, fmt.Errorf("request error %v", data["errcode"])
	}
	fmt.Printf("request source data: %v \n", time.Since(getDataTime))
	return data, nil
}
