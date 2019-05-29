:loop
 curl -i http://localhost:9090/GetActiveJob
sleep 1
 curl -i http://localhost:9090/GetFinishedJobs
sleep 1
 curl -i http://localhost:9090/GetWaitingJobs
sleep 1

goto loop
