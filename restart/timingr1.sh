while true 
do
	kill -s 9 `ps -aux | grep taskHttpService | awk '{print $2}'`
	sleep 1h
done
