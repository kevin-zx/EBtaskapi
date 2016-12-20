while true 
do
	#statements
	a=`curl http://127.0.0.1:19922/getTask`
	if [[ $a =~ "Task_id" ]] && [[ $a =~ "Task_id" ]] ; then
		echo $a
		date
		echo "sleep"
		sleep 20m
	else
		sh restart.sh >> restart.log &
		date
		echo "restart"
		sleep 30s
	fi
done
