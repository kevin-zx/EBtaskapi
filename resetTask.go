package main

import (

	"fmt"
	"time"
	"taskService/mysqlServer"
)

func main() {
	for {
		if time.Now().Hour() == 23 {
			if time.Now().Minute() == 59 {
				var sql string = "update eb_task set task_execed_times = 0,task_success_times = 0"
				mysqlServer.MysqlServer.Insert(sql)
				fmt.Println("restart")
				fmt.Println(time.Now())
				time.Sleep(60 * time.Second)

			}
		}
		time.Sleep(30 * time.Second)
	}

}
