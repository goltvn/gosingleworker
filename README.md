# gosingleworker
single worker queue 

My first golang project, aim is to create a rest API that will take incoming requests, and store these in a persistant storage (database), and have a single worker process each incoming task one by one.

performing one task can take a long time, and there can only be one worker and one job processed at a time.

Test of starting jobs (windows batch) - loop-set.cmd

@echo off
set loopcount=%1
:loop
set /a loopcount=loopcount-1
echo %loopcount%
START /i /b curl http://localhost:9090/StartJob?id=%loopcount
if %loopcount%==0 goto exitloop
goto loop
:exitloop
pause

run ex with loop-set 100

it creates job 99 down to job 0.

Note if the database already have them created it will not create them again.

script to test getting status values :

get-act.cmd
:loop
 curl -i http://localhost:9090/GetActiveJob
sleep 1
 curl -i http://localhost:9090/GetFinishedJobs
sleep 1
 curl -i http://localhost:9090/GetWaitingJobs
sleep 1

it performs an api lookup every 1 second pr one of the 3 above to show active jobno, no of finished jobs, no of waiting jobs
