package matchinfomanager

import (
	"taskService/mysqlServer"
	"strconv"
	"time"
)

type TaskMatchInfo struct {
	Task_id       int
	Match_type_id int
	Match_info    string
	Match_method  int
}

// var mu sync.Mutex
var matchInfoMap map[int][]TaskMatchInfo = make(map[int][]TaskMatchInfo)
var times = time.Now().Unix()

func GetMatchInfo(taskId int) []TaskMatchInfo {
	if time.Now().Unix()-times > 30 || len(matchInfoMap) == 0 {
		times = time.Now().Unix()
		initMap()
	}
	return matchInfoMap[taskId]
}

func initMap() {
	matchInfoMap = getAllMatchInfoFromDB()
}

func getAllMatchInfoFromDB() map[int][]TaskMatchInfo {
	tmpMatchInfoMap := make(map[int][]TaskMatchInfo)
	resultData, _ := mysqlServer.MysqlServerInstance.SelectAll(`SELECT
		eb_task_match_info.task_id,
		match_type_id,
		match_info,
		match_method 
		FROM eb_task_match_info 
		LEFT JOIN eb_task 
		ON eb_task_match_info.task_id = eb_task.task_id
		 WHERE eb_task.task_status != 0`)
	for _, data := range *resultData {
		var matchInfo TaskMatchInfo
		matchInfo.Task_id, _ = strconv.Atoi(data["task_id"])
		matchInfo.Match_type_id, _ = strconv.Atoi(data["match_type_id"])
		matchInfo.Match_info = data["match_info"]
		matchInfo.Match_method, _ = strconv.Atoi(data["match_method"])
		tmpMatchInfoMap[matchInfo.Task_id] = append(tmpMatchInfoMap[matchInfo.Task_id], matchInfo)
	}
	return tmpMatchInfoMap
}
