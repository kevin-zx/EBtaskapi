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
	println(r.RemoteAddr)
	r.Close = true
	// taskType := r.URL.Query().Get("task_type")
	task := taskManager.GetTask(1)
	// fmt.Println(task)
	b, err := json.Marshal(&task)
	// fmt.Println(r.RemoteAddr)
	r.Header.Set("Accept-Encoding", "")
	if err != nil {
		fmt.Println(err)
	} else {
		fmt.Println(string(b))
		io.WriteString(w, string(b))
	}

}

func getTaskByPlatform(w http.ResponseWriter, r *http.Request) {
	r.Close = true
	platformid := r.URL.Query().Get("platformid")
	device := r.URL.Query().Get("device")
	task := taskManager.GetTaskByType(platformid, device,"","")
	b, err := json.Marshal(&task)
	r.Header.Set("Accept-Encoding", "")
	if err != nil {
		fmt.Println(err)
	} else {
		fmt.Println(string(b))
		io.WriteString(w, string(b))
	}
}

func getTaskByArgs(w http.ResponseWriter, r *http.Request) {
	r.Close = true
	platformid := r.URL.Query().Get("platformid")
	device := r.URL.Query().Get("device")
	province := r.URL.Query().Get("province")
	city := r.URL.Query().Get("city")
	task := taskManager.GetTaskByType(platformid, device, province, city)
	b, err := json.Marshal(&task)
	r.Header.Set("Accept-Encoding", "")
	if err != nil {
		fmt.Println(err)
	} else {
		fmt.Println(string(b))
		io.WriteString(w, string(b))
	}
}

func handResult(w http.ResponseWriter, r *http.Request) {
	r.Close = true

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
	fmt.Printf("task_id:%s,ip:%s,port:%s,elapsed:%s,area:%s,device:%s,success_status:%s\r\n", task_id, ip, port, elapsed, area, device, success_status)
	err := taskResult.HandlerResult(task_id, success_status, ip, port, elapsed, area, device, account)
	if err == nil {
		io.WriteString(w, "{\"status\":\"ok\"}")
	} else {
		io.WriteString(w, "{\"status\":\"err\",\"message\":\""+err.Error()+"\"}")
	}
	// log.Info(fmt.Sprintf("remote:ip,", ...))
}

func main() {
	http.HandleFunc("/getTask", getTask)
	http.HandleFunc("/getTaskByPlatform", getTaskByPlatform)
	http.HandleFunc("/getTaskByArgs", getTaskByArgs)
	http.HandleFunc("/handResult", handResult)
	http.ListenAndServe(":19922", nil)
}
