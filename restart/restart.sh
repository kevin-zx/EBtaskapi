kill -s 9 `ps -aux | grep taskHttpService | awk '{print $2}'`
nohup go run ../taskHttpService.go &
