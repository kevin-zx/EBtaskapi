package main

import (
	"taskService/taskManager"
	"taskService/taskResult"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

func getTask(w http.ResponseWriter, r *http.Request) {
	//println(r.RemoteAddr)
	r.Close = true
	defer r.Body.Close()
	task := taskManager.Task{}
	if taskManager.ValidateIp(r.RemoteAddr){
		task = taskManager.GetTask(1)
		//println("获取到任务，ip 是 ", r.RemoteAddr,"---------------------------------------")
	}else{
		//println("没有获取到任务，ip 是 ", r.RemoteAddr,"---------------------------------------")
	}
	b, err := json.Marshal(&task)
	r.Header.Set("Accept-Encoding", "")
	if err != nil {
		fmt.Println(err)
	} else {
		//if task.Task_id != 0{
		//	fmt.Println(string(b))
		//}
		io.WriteString(w, string(b))
	}

	//r.
}

func getTaskByPlatform(w http.ResponseWriter, r *http.Request) {
	r.Close = true
	defer r.Body.Close()
	platformId := r.URL.Query().Get("platformId")
	device := r.URL.Query().Get("device")
	task := taskManager.GetTaskByType(platformId, device,"","")
	b, err := json.Marshal(&task)
	r.Header.Set("Accept-Encoding", "")
	if err != nil {
		fmt.Println(err)
	} else {
		//fmt.Println(string(b))
		io.WriteString(w, string(b))
	}

}

func getTaskByArgs(w http.ResponseWriter, r *http.Request) {
	r.Close = true
	defer r.Body.Close()
	//println(r.RemoteAddr)
	platformId := r.URL.Query().Get("platformId")
	device := r.URL.Query().Get("device")
	province := r.URL.Query().Get("province")
	city := r.URL.Query().Get("city")
	task := taskManager.GetTaskByType(platformId, device, province, city)
	b, err := json.Marshal(&task)
	r.Header.Set("Accept-Encoding", "")
	if err != nil {
		fmt.Println(err)
	} else {
		//fmt.Println(string(b))
		io.WriteString(w, string(b))
	}

}

func handResult(w http.ResponseWriter, r *http.Request) {
	r.Close = true
	defer r.Body.Close()
	task_id := r.URL.Query().Get("task_id")
	if task_id == "" {
		io.WriteString(w, "{\"status\":\"ok\",\"message\":\"no task_id\"}")
		return
	}
	success_status := r.URL.Query().Get("success_status")
	if success_status == "" {
		io.WriteString(w, "{\"status\":\"ok\",\"message\":\"no success_status\"}")
		return
	}
	ip := r.URL.Query().Get("ip")
	port := r.URL.Query().Get("port")
	elapsed := r.URL.Query().Get("elapsed")
	if elapsed == "" {
		io.WriteString(w, "{\"status\":\"ok\",\"message\":\"no elapsed\"}")
		return
	}

	area := r.URL.Query().Get("area")
	device := r.URL.Query().Get("device")
	account := r.URL.Query().Get("account")
	//error_type
	error_type := r.URL.Query().Get("error_type")

	fmt.Printf("task_id:%s,ip:%s,port:%s,elapsed:%s,area:%s,device:%s,success_status:%s, error_type:%s\r\n", task_id, ip, port, elapsed, area, device, success_status, error_type)
	err := taskResult.HandlerResult(task_id, success_status, ip, port, elapsed, area, device, account, error_type)
	if err == nil {
		io.WriteString(w, "{\"status\":\"ok\"}")
	} else {
		io.WriteString(w, "{\"status\":\"err\",\"message\":\""+err.Error()+"\"}")
	}

	// log.Info(fmt.Sprintf("remote:ip,", ...))
}
func forbiddenAccount(w http.ResponseWriter, r *http.Request) {
	r.Close = true
	defer r.Body.Close()

	account := r.URL.Query().Get("account")
	taskManager.ForbiddenAccount(account)
	io.WriteString(w,"{\"status\":\"ok\"}")

}
func main() {
	port := 19922
	http.HandleFunc("/getTask", getTask)
	http.HandleFunc("/getTaskByPlatform", getTaskByPlatform)
	http.HandleFunc("/getTaskByArgs", getTaskByArgs)
	http.HandleFunc("/handResult", handResult)
	http.HandleFunc("/forbiddenAccount", forbiddenAccount)
	fmt.Printf("listen at port %d", port)
	//http.DefaultClient.

	http.ListenAndServe(":19922", nil)
}
