package taskResult

import (
	"sync"
	"taskService/mysqlServer"
)

var mu sync.Mutex

type Result struct {
	Task_id, Success_status, Ip, Port, Elapsed, Area, Device,Account string

}

var resultSlice []Result

func HandlerResult(task_id string, success_status string, ip string, port string, elapsed string, area string, device string,account string) error {
	var re Result = Result{Task_id: task_id,
		Success_status: success_status,
		Ip:             ip,
		Port:           port,
		Elapsed:        elapsed,
		Area:           area,
		Device:         device,
		Account:	account}
	mu.Lock()
	defer mu.Unlock()
	resultSlice = append(resultSlice, re)
	if len(resultSlice) >= 1 {
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
	updateAccSql := "UPDATE eb_account SET exec_count = exec_count + 1, success_count=success_count + (?) WHERE account = ?"
	var insertVals [][]interface{}
	var updateVals [][]interface{}
	var updateAccountVals [][]interface{}
	for _, result := range resultSlice {
		var val = []interface{}{result.Task_id, result.Success_status, result.Ip, result.Port, result.Elapsed, result.Area, result.Device}
		insertVals = append(insertVals, val)
		updateVals = append(updateVals, []interface{}{result.Success_status, result.Task_id})
		if result.Account != ""{

			updateAccountVals = append(updateAccountVals, []interface{}{result.Success_status,result.Account})
		}

	}
	err := mysqlServer.MysqlServer.ExecBatch(sql, insertVals)
	if err != nil {
		return err
	}
	err = mysqlServer.MysqlServer.ExecBatch(updateSql, updateVals)
	if err != nil {
		return err
	}
	if len(updateAccountVals) >0 {
		err = mysqlServer.MysqlServer.ExecBatch(updateAccSql, updateAccountVals)
	}

	if err != nil {
		println(err.Error())
		return err
	}
	return nil
}
