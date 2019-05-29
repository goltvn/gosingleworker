# gosingleworker
single worker queue 

My first golang project, aim is to create a rest API that will take incoming requests, and store these in a persistant storage (database), and have a single worker process each incoming task one by one.

performing one task can take a long time, and there can only be one worker and one job processed at a time.

Test of starting jobs (windows batch) - this starts job no 100 - there are scripts that automate alot of requests ( loop-set.cmd )

START /i /b curl http://localhost:9090/StartJob?id=100


to get some status values can call : 

 curl -i http://localhost:9090/GetActiveJob

 curl -i http://localhost:9090/GetFinishedJobs
 
 curl -i http://localhost:9090/GetWaitingJobs
