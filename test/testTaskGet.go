package main

import (
	//"../taskManager"
	// "encoding/json"
	"fmt"
	"taskService/taskManager"
)

func main() {

	// for i := 0; i < 20; i++ {
	task := taskManager.GetTask(1)
	fmt.Println(task)
	// 	b, err := json.Marshal(&task)
	// 	if err != nil {
	// 		fmt.Println(err)
	// 	} else {
	// 		fmt.Println(string(b))
	// 	}

	// }

}
