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
		fmt.Println("update access token")
		var token, err = request.AccessToken() // get accessToken
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
	fmt.Printf("%v%v ", requ.Host, requ.RequestURI)
	beginTime, strConvErr := strconv.Atoi(query["begin_time"][0])
	if strConvErr != nil {
		log.Fatal(strConvErr.Error())
	}
	endTime, strConvErr := strconv.Atoi(query["end_time"][0])
	if strConvErr != nil {
		log.Fatal(strConvErr.Error())
	}
	var rqData = request.RqData{OpenCheckInDataType: 1, StartTime: beginTime, EndTime: endTime, UserList: GetMembers(query["member"][0])} // request para
	jsonRqData, jsonEncodeErr := json.Marshal(rqData)
	if jsonEncodeErr != nil {
		fmt.Printf("json encode error! %v", jsonEncodeErr)
	}
	var checkInData, _ = request.GetSourData(accessToken, bytes.NewReader(jsonRqData)) // request
	//fmt.Println(string(data.MakeReturnJson(checkInData)))                              // resp json data
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
	membersMap := members[grade].(map[string]interface{})
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
