
@echo off
set loopcount=%1
:loop
set /a loopcount=loopcount-1
echo %loopcount%
START /b curl http://localhost:9090/GetJobStatus?id=10
if %loopcount%==0 goto exitloop
goto loop
:exitloop
pause
