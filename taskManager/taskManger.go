package taskManager

import (
	"taskService/matchinfomanager"
	"taskService/mysqlServer"
	"fmt"
	"strconv"
	"sync"
	"time"
	"container/list"
	"github.com/kevin-zx/go-util/errorUtil"
	"github.com/kevin-zx/go-util/dateUtil"
	"strings"
)

type Task struct {
	Task_id      int
	Platform_id  int
	Task_type_id int
	Task_keyword string
	Browser_id    int
	Task_max_page int
	Task_status   int
	TaskMatchInfo []matchinfomanager.TaskMatchInfo
	Acc Account
}

type Account struct{
	Account string
	Password string
	Device string
	Platform string
}

var account_max_exec_time = 50
var taskList []Task
var index int
var mu sync.Mutex
var taskListLen int
var taskMapList map[int]*list.List = make(map[int]*list.List)

func ForbiddenAccount(account string){
	forbiddenAccountSql := "UPDATE eb_account SET status = 0 WHERE account = ?"
	mysqlServer.MysqlServerInstance.Exec(forbiddenAccountSql,account)
}

func GetTaskByType(platform_id string, device string, province string, city string) Task  {
  	var task Task
	var taskList *list.List
	var taskEle *list.Element
	task_type_int, err := strconv.Atoi(platform_id)
	if err != nil {
		task_type_int = -1
	}
	mu.Lock()
	defer mu.Unlock()
	if task_type_int == -1 && len(taskMapList) > 0 {
		for _,d := range taskMapList{
			taskList = d
			taskEle = d.Front()
			if taskEle == nil{
				continue
			}else {
				break
			}
		}
	}else if taskMapList[task_type_int] != nil && (taskMapList[task_type_int]).Len() > 0 {
		taskList = taskMapList[task_type_int]
		taskEle = taskList.Front()
	}else {
		getTaskQueue(platform_id)
	}
	if taskEle != nil{
		taskList.Remove(taskEle)
		task,_ = taskEle.Value.(Task)
		task.TaskMatchInfo = matchinfomanager.GetMatchInfo(task.Task_id)
		task.Acc = getAccount(device, platform_id, province, city)
	}
	return task
}

//根据device，platform,province，city获取account
func getAccount(device string, platform_id string, province string, city string) Account{

	platform := "taobao"
	var acc Account = Account{}
	if platform_id == "1" || platform_id == "2" || platform_id == "3"{
		platform = "taobao"
	}else {
		return Account{}
	}

	condition := ""
	if city != "" {
		if condition != "" {
			condition +=" AND"
		}
		condition += " city='" + city + "'"
	}else {
		return Account{}
	}
	//print(city)
	//如果条件里面有内容才会去获取condition
	if condition != "" {
		sql_format :="SELECT * FROM eb_account WHERE %s AND (exec_count/%d)< TIMESTAMPDIFF(MINUTE, DATE_FORMAT(NOW(),'%%Y-%%m-%%d 00:00:00'),NOW())/1440 AND status = 1 ORDER BY RAND() LIMIT 1"
		sql := fmt.Sprintf(sql_format, condition, account_max_exec_time)
		//println(sql)
		results, err :=  mysqlServer.MysqlServerInstance.SelectAll(sql)
		if err == nil{
			for _,result := range *results{
				acc = Account{Account:result["account"],Password:result["passwd"],Platform:platform,Device:device}
				return acc
			}
		}
	}

	return Account{}
}



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

//验证ip是否可用
func ValidateIp(remote_ip string) bool {
	two_hour_before_date := dateUtil.GetDeltaDate(-2*time.Hour)
	half_day_before_date := dateUtil.GetDeltaDate(-12*time.Hour)
	remote_ip = strings.Split(remote_ip,":")[0]
	//data,err := mysqlServer.MysqlServerInstance.SelectAll("SELECT * FROM eb_result WHERE ip = ? AND error_type = 1 and insert_date > ? LIMIT 1",remote_ip,two_hour_before_date)
	data,err := mysqlServer.MysqlServerInstance.SelectAll("SELECT * FROM eb_result WHERE ip = ? AND insert_date > ?  LIMIT 1",remote_ip,two_hour_before_date)

	data2,err := mysqlServer.MysqlServerInstance.SelectAll("SELECT * FROM eb_result WHERE ip = ? AND insert_date > ? AND error_type = 1  LIMIT 1",remote_ip,half_day_before_date)
	errorUtil.CheckErrorExit(err)
	return len(*data) == 0 && len(*data2) == 0
}

func reset_task() {
	var sql string = "update eb_task set task_execed_times = 0,task_success_times = 0"
	mysqlServer.MysqlServerInstance.Insert(sql)
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
	//println(sql)
	resultData, err := mysqlServer.MysqlServerInstance.SelectAll(sql)
	if err != nil {
		return nil, err
	}
	var taskResults []Task
	for _, data := range *resultData {
		var task Task
		task.Task_id, _ = strconv.Atoi(data["task_id"])
		task.Platform_id, _ = strconv.Atoi(data["platform_id"])
		task.Task_type_id, _ = strconv.Atoi(data["task_type_id"])
		task.Task_keyword = data["task_keyword"]
		task.Browser_id, _ = strconv.Atoi(data["browser_id"])
		task.Task_max_page, _ = strconv.Atoi(data["task_max_page"])
		task.Task_status, _ = strconv.Atoi(data["task_status"])
		taskResults = append(taskResults, task)

	}
	return taskResults, nil
}

func getTaskQueue(platform string){
	var sql string
	sqlFormatStr := `SELECT eb_task.task_id,
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
		%s
		AND (((eb_task.task_success_times +1)/(eb_task.task_exec_times + 1)) - (TIMESTAMPDIFF(MINUTE,date_format(now(),'%%Y-%%m-%%d 00:00:00'),NOW())+1)/1440) < 0
		ORDER BY (eb_task.task_success_times + 50)/(eb_task.task_execed_times + 50)
		DESC
		LIMIT 200`
	if platform != "" && platform != "-1" {
		platform = "AND eb_task.platform_id = "+ platform
		sql = fmt.Sprintf(sqlFormatStr, platform)
	}else {
		sql = fmt.Sprintf(sqlFormatStr, "")
	}

	//println(sql)
	resultData, err := mysqlServer.MysqlServerInstance.SelectAll(sql)

	if err != nil {
		println(err)
	}

	for _, data := range *resultData {
		var task Task
		task.Task_id, _ = strconv.Atoi(data["task_id"])
		task.Platform_id, _ = strconv.Atoi(data["platform_id"])
		task.Task_type_id, _ = strconv.Atoi(data["task_type_id"])
		task.Task_keyword = data["task_keyword"]
		task.Browser_id, _ = strconv.Atoi(data["browser_id"])
		task.Task_max_page, _ = strconv.Atoi(data["task_max_page"])
		task.Task_status, _ = strconv.Atoi(data["task_status"])
		//d,_ := json.Marshal(task)
		//println(string(d))
		if taskMapList[task.Platform_id] != nil{
			taskMapList[task.Platform_id].PushBack(task)
		}else {
			taskMapList[task.Platform_id] = list.New()
			taskMapList[task.Platform_id].PushBack(task)
		}
	}
}
