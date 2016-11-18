package main

// import "database/sql"
// import _ "github.com/go-sql-driver/mysql"

import "fmt"
import (
	"../taskManager"
	// "log"
)

type test struct {
	name string
	dd   string
}

func main() {
	_, err := taskManager.GetTaskList()
	if err != nil {
		fmt.Println(err)
	} else {
		// fmt.Println(data)
	}
	// err := mysqlUtil.Insert("INSERT INTO test (`name`,`address`) VALUE (?,?)", "zx", "wd")
	// if err != nil {
	// 	fmt.Println(err)
	// }
	// value, err := mysqlUtil.Select("SELECT * FROM test")
	// if err != nil {
	// 	fmt.Println(err)
	// } else {
	// 	fmt.Println(value)
	// }

	// test1 := test{"name": "1", "dd": "2"}
	// fmt.Println(test)
}
