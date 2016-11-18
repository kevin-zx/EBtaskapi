package taskResult

import (
	// "fmt"
	// "strings"
	"../mysqlUtil"
	"sync"
)

var mu sync.Mutex

type Result struct {
	Task_id, Success_status, Ip, Port, Elapsed, Area, Device string
}

var resultSlice []Result

func HandlerResult(task_id string, success_status string, ip string, port string, elapsed string, area string, device string) error {
	var re Result = Result{Task_id: task_id,
		Success_status: success_status,
		Ip:             ip,
		Port:           port,
		Elapsed:        elapsed,
		Area:           area,
		Device:         device}
	mu.Lock()
	defer mu.Unlock()
	resultSlice = append(resultSlice, re)
	if len(resultSlice) >= 5 {
		err := insertResult()
		resultSlice = *new([]Result)
		if err != nil {
			return err
		}
	}
	return nil
}

func insertResult() error {
	sql := "INSERT INTO eb_result (`task_id`,`success_status`,`ip`,`port`,`elapsed`,`insert_date`,`area`,`device`) VALUE (?,?,?,?,?,NOW(),?,?)"
	updateSql := "update eb_task set task_success_times = task_success_times + ?,task_execed_times = task_execed_times + 1 where task_id = ?"

	var insertVals [][]interface{}
	var updateVals [][]interface{}
	for _, result := range resultSlice {
		var val = []interface{}{result.Task_id, result.Success_status, result.Ip, result.Port, result.Elapsed, result.Area, result.Device}
		insertVals = append(insertVals, val)
		updateVals = append(updateVals, []interface{}{result.Success_status, result.Task_id})
	}
	err := mysqlUtil.ExecBatch(sql, insertVals)
	if err != nil {
		return err
	}
	err = mysqlUtil.ExecBatch(updateSql, updateVals)
	if err != nil {
		return err
	}
	return nil
}
