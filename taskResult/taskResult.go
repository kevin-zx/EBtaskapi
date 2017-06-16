package taskResult

import (
	"sync"
	"taskService/mysqlServer"
)

var mu sync.Mutex

type Result struct {
	Task_id, Success_status, Ip, Port, Elapsed, Area, Device,Account,Error_type string

}

var resultSlice []Result

func HandlerResult(task_id string, success_status string, ip string, port string, elapsed string, area string, device string,account string, error_type string) error {
	var re Result = Result{Task_id: task_id,
		Success_status: success_status,
		Ip:             ip,
		Port:           port,
		Elapsed:        elapsed,
		Area:           area,
		Device:         device,
		Account:	account,
		Error_type: error_type,
	}
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
	sql := "INSERT INTO eb_result (`task_id`,`success_status`,`ip`,`port`,`elapsed`,`insert_date`,`area`,`device`,`error_type`) VALUE (?,?,?,?,?,NOW(),?,?,?)"
	updateSql := "update eb_task set task_success_times = task_success_times + ?,task_execed_times = task_execed_times + 1 where task_id = ?"
	updateAccSql := "UPDATE eb_account SET exec_count = exec_count + 1, success_count=success_count + (?) WHERE account = ?"
	var insert_values [][]interface{}
	var update_values [][]interface{}
	var updateAccount_values [][]interface{}
	for _, result := range resultSlice {
		var val = []interface{}{result.Task_id, result.Success_status, result.Ip, result.Port, result.Elapsed, result.Area, result.Device,result.Error_type}
		insert_values = append(insert_values, val)
		update_values = append(update_values, []interface{}{result.Success_status, result.Task_id})
		if result.Account != ""{
			updateAccount_values = append(updateAccount_values, []interface{}{result.Success_status, result.Account})
		}
	}
	err := mysqlServer.MysqlServerInstance.ExecBatch(sql, insert_values)
	if err != nil {
		return err
	}
	err = mysqlServer.MysqlServerInstance.ExecBatch(updateSql, update_values)
	if err != nil {
		return err
	}
	if len(updateAccount_values) >0 {
		err = mysqlServer.MysqlServerInstance.ExecBatch(updateAccSql, updateAccount_values)
	}

	if err != nil {
		println(err.Error())
		return err
	}
	return nil
}
