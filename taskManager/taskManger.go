package taskManager

import (
	"../matchinfomanager"
	"../mysqlUtil"
	"fmt"
	"strconv"
	"sync"
	"time"
)

type Task struct {
	Task_id      int
	Platform_id  int
	Task_type_id int
	Task_keyword string
	// task_main_info   string
	// task_target_info string
	Browser_id    int
	Task_max_page int
	Task_status   int
	TaskMatchInfo []matchinfomanager.TaskMatchInfo
}

// func (t Task) String() string {
// 	return fmt.Sprintf("task_id:%d,platform_id:%d,task_type_id:%d,"+
// 		"task_keyword:%s,browser_id:%d,task_max_page:%d,"+
// 		"task_status:%d\n",
// 		t.Task_id,
// 		t.Platform_id,
// 		t.Task_type_id,
// 		t.Task_keyword,
// 		t.Browser_id,
// 		t.Task_max_page,
// 		t.Task_status)
// }
var taskList []Task
var index int
var mu sync.Mutex
var taskListLen int

func GetTask(task_status int) Task {
	mu.Lock()

	var task Task
	if taskListLen > 0 && taskListLen-1 != index {
		index++
		task = taskList[index]
		task.TaskMatchInfo = matchinfomanager.GetMatchInfo(task.Task_id)
	} else {
		taskList, _ = getTaskList()
		if taskList == nil || len(taskList) == 0 {
			fmt.Println("get task failed ")
			time.Sleep(2 * time.Second)
			return task
		}
		taskListLen = len(taskList)
		index = 0
		task = taskList[index]
		task.TaskMatchInfo = matchinfomanager.GetMatchInfo(task.Task_id)
	}

	mu.Unlock()
	return task
}

func reset_task() {
	var sql string = "update eb_task set task_execed_times = 0,task_success_times = 0"
	mysqlUtil.Insert(sql)
	time.Sleep(60 * time.Second)

}

func getTaskList() ([]Task, error) {
	var sql string = `SELECT eb_task.task_id,
	eb_task.platform_id,
	eb_task.task_type_id,
	eb_task.task_keyword,
	eb_task.browser_id,
	eb_task.task_max_page,
	eb_task.task_status
	FROM 
	eb_task 
	WHERE 
	eb_task.task_status = 1
	AND (((eb_task.task_success_times +1)/(eb_task.task_exec_times + 1)) - (TIMESTAMPDIFF(MINUTE,date_format(now(),'%Y-%m-%d 00:00:00'),NOW())+1)/1440) < 0 
	ORDER BY (eb_task.task_success_times + 50)/(eb_task.task_execed_times + 50) 
	DESC 
	LIMIT 200
	`
	resultData, err := mysqlUtil.SelectAll(sql)
	if err != nil {
		return nil, err
	}
	var taskResults []Task
	// for i := 0; i < len(data); i++ {
	// 	var task Task
	// 	taskResults = append(taskResults)
	// 	task.Task_id, _ = strconv.Atoi(string(data[i][0]))
	// 	task.Platform_id, _ = strconv.Atoi(string(data[i][1]))
	// 	task.Task_type_id, _ = strconv.Atoi(string(data[i][2]))
	// 	task.Task_keyword = string(data[i][3])
	// 	// fmt.Println(task.Task_keyword)
	// 	task.Browser_id, _ = strconv.Atoi(string(data[i][4]))
	// 	task.Task_max_page, _ = strconv.Atoi(string(data[i][5]))
	// 	task.Task_status, _ = strconv.Atoi(string(data[i][6]))
	// 	taskResults = append(taskResults, task)
	// }
	for _, data := range *resultData {
		var task Task
		task.Task_id, _ = strconv.Atoi(data["task_id"])
		task.Platform_id, _ = strconv.Atoi(data["platform_id"])
		task.Task_type_id, _ = strconv.Atoi(data["task_type_id"])
		task.Task_keyword = data["task_keyword"]
		task.Browser_id, _ = strconv.Atoi(data["browser_id"])
		// fmt.Println(task.Browser_id)
		task.Task_max_page, _ = strconv.Atoi(data["task_max_page"])
		task.Task_status, _ = strconv.Atoi(data["task_status"])
		taskResults = append(taskResults, task)

	}
	return taskResults, nil
}
