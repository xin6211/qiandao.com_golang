package data

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"sort"
	"strconv"
	"time"
)

type reData struct {
	Mon    []float64 `json:"Mon"`
	Tues   []float64 `json:"Tues"`
	Wed    []float64 `json:"Wed"`
	Thur   []float64 `json:"Thur"`
	Fri    []float64 `json:"Fri"`
	Sat    []float64 `json:"Sat"`
	Sun    []float64 `json:"Sun"`
	Member []string  `json:"member"`
}

type memberCheckIn struct {
	Id    string
	Name  string
	Times [8][]int
}

func GetWeek(unixTime int) (int, string) {
	var weeks = map[int]string{
		1: "Mon",
		2: "Tues",
		3: "Wed",
		4: "Thur",
		5: "Fri",
		6: "Sat",
		7: "Sun",
	}
	intWeek := (unixTime%(7*24*60*60))/(24*60*60) + 4
	if intWeek >= 8 {
		intWeek -= 7
	}
	return intWeek, weeks[intWeek]
}

func GetName(id string) string {
	fp, err := ioutil.ReadFile("./static/members.json")
	if err != nil {
		fmt.Println("read error")
		return ""
	}
	members := make(map[string]interface{})
	errJ := json.Unmarshal(fp, &members)
	if errJ != nil {
		fmt.Println("error")
	}
	name := members["20"].(map[string]interface{})[id]
	if name == nil {
		name = members["21"].(map[string]interface{})[id]
	}
	return name.(string)
}

func CalDurationDay(cTimesDay []int) float64 {
	var duration float64
	if len(cTimesDay)%2 == 0 {
		for i := range cTimesDay {
			if i%2 != 0 {
				duration = duration + float64(cTimesDay[i]-cTimesDay[i-1])/3600
			}
		}
		return duration
	} else {
		var durations []float64
		var timeSlice []int
		for i := range cTimesDay {
			if i%2 == 0 {
				timeSlice = append(timeSlice, cTimesDay[:i]...)
				timeSlice = append(timeSlice, cTimesDay[i+1:]...)
				durations = append(durations, CalDurationDay(timeSlice))
				timeSlice = []int{}
			}
		}
		sort.Float64s(durations)
		var rst float64
		if len(durations) >= 2 {
			if len(durations)%2 == 0 {
				rst = (durations[len(durations)/2] + durations[len(durations)/2-1]) / 2
			} else {
				rst = durations[len(durations)/2-1]
			}
		} else if len(durations) == 1 {
			rst = durations[0]
		} else {
			rst = 0
		}
		return rst
	}
}

func GetDuration(checkInTimes []int) float64 {
	if len(checkInTimes) <= 1 {
		return 0
	}
	var duration float64
	var day = time.Unix(int64(checkInTimes[0]), 0).Day()
	var cTimesDay []int
	cTimesDay = append(cTimesDay, checkInTimes[0])
	for cTime := 1; cTime < len(checkInTimes); cTime++ {
		if time.Unix(int64(checkInTimes[cTime]), 0).Day() == day {
			cTimesDay = append(cTimesDay, checkInTimes[cTime])
		}
		if time.Unix(int64(checkInTimes[cTime]), 0).Day() != day {
			day = time.Unix(int64(checkInTimes[cTime]), 0).Day()
			duration = duration + CalDurationDay(cTimesDay)
			cTimesDay = []int{checkInTimes[cTime]}
			continue
		}
		if cTime == len(checkInTimes)-1 {
			duration = duration + CalDurationDay(cTimesDay)
		}
	}
	duration, _ = strconv.ParseFloat(fmt.Sprintf("%0.2f", duration), 1)
	return duration
}

func GetCheckInData(respData map[string]interface{}) map[string]*memberCheckIn {
	var membersCheckTime = make(map[string]*memberCheckIn)
	for _, v := range respData["checkindata"].([]interface{}) {
		value := v.(map[string]interface{})
		id := value["userid"].(string)
		if membersCheckTime[id] == nil {
			membersCheckTime[id] = &memberCheckIn{}
		}
		membersCheckTime[id].Id = id
		membersCheckTime[id].Name = GetName(id)
		weekInt, _ := GetWeek(int(value["checkin_time"].(float64)))
		membersCheckTime[id].Times[weekInt] = append(membersCheckTime[id].Times[weekInt], int(value["checkin_time"].(float64)))
	}
	return membersCheckTime
}

func MakeReturnJson(respData map[string]interface{}) []byte {
	var retData = reData{}
	membersCheckTime := GetCheckInData(respData)
	for v := range membersCheckTime {
		retData.Mon = append(retData.Mon, GetDuration(membersCheckTime[v].Times[1]))
		retData.Tues = append(retData.Tues, GetDuration(membersCheckTime[v].Times[2]))
		retData.Wed = append(retData.Wed, GetDuration(membersCheckTime[v].Times[3]))
		retData.Thur = append(retData.Thur, GetDuration(membersCheckTime[v].Times[4]))
		retData.Fri = append(retData.Fri, GetDuration(membersCheckTime[v].Times[5]))
		retData.Sat = append(retData.Sat, GetDuration(membersCheckTime[v].Times[6]))
		retData.Sun = append(retData.Sun, GetDuration(membersCheckTime[v].Times[7]))
		retData.Member = append(retData.Member, membersCheckTime[v].Name)
	}
	jsonReData, err := json.Marshal(retData)
	if err != nil {
		fmt.Println("json marashal error!")
	}
	return jsonReData
}
