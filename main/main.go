package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"qiandao.com/data"
	"qiandao.com/request"
	"strconv"
	"time"
)

var accessToken string

func keepAccessToken() {
	for {
		var token, err = request.GetAcessToken() // get accessToken
		if err != nil {
			fmt.Printf("get access token error! %v", err)
		}
		accessToken = token
		time.Sleep(1 * time.Hour)
	}
}

func handleGet(writer http.ResponseWriter, requ *http.Request) {
	writer.Header().Set("Access-Control-Allow-Origin", "*")
	writer.Header().Add("Access-Control-Allow-Headers", "Content-Type,AccessToken,X-CSRF-Token, Authorization, Token")
	writer.Header().Add("Access-Control-Allow-Credentials", "true")
	writer.Header().Add("Access-Control-Allow-Methods", "POST, GET, OPTIONS, PUT, DELETE")
	writer.Header().Set("content-type", "application/json;charset=UTF-8")
	query := requ.URL.Query()
	beginTime, strConvErr := strconv.Atoi(query["begin_time"][0])
	if strConvErr != nil {
		log.Fatal(strConvErr.Error())
	}
	endTime, strConvErr := strconv.Atoi(query["end_time"][0])
	if strConvErr != nil {
		log.Fatal(strConvErr.Error())
	}
	fmt.Printf("%v%v ", requ.Host, requ.RequestURI)
	var rqData = request.RqData{3, beginTime, endTime, GetMembers(query["member"][0])} // request para
	jsonRqData, jsonEncodeErr := json.Marshal(rqData)
	if jsonEncodeErr != nil {
		fmt.Printf("json encode error! %v", jsonEncodeErr)
	}
	var checkInData, _ = request.GetSourData(accessToken, bytes.NewReader(jsonRqData)) // request
	//fmt.Println(data.MakeReturnJson(checkInData))                                    // resp json data
	_, err := writer.Write(data.MakeReturnJson(checkInData))
	if err != nil {
		log.Fatal("ListenAndServe: ", err.Error())
	}
}

func GetMembers(grade string) []string {
	fp, err := ioutil.ReadFile("./static/members.json")
	if err != nil {
		fmt.Println("read error")
		log.Fatal(err.Error())
	}
	members := make(map[string]interface{})
	errJ := json.Unmarshal(fp, &members)
	if errJ != nil {
		fmt.Println("error")
	}
	membersMap := members["20"].(map[string]interface{})
	ids := make([]string, 0, len(membersMap))
	for k := range membersMap {
		ids = append(ids, k)
	}
	return ids
}

func main() {
	go keepAccessToken()
	http.HandleFunc("/", handleGet)
	fmt.Println("start....")
	err := http.ListenAndServe(":8888", nil)
	if err != nil {
		fmt.Printf("%v", err)
	}
}
