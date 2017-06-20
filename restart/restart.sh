kill -s 9 `ps -aux | grep taskService | awk '{print $2}'`
nohup go run ../taskService.go &
