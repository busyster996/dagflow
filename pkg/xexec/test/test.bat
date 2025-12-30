@echo off
echo Hello World
echo 参数1: %1
echo 参数2: %2
echo 所有参数: %*
REM 从标准输入读取每一行
for /f "delims=" %%i in ('more') do (
    echo Line: %%i
)

REM 持续ping
ping -t 1.1.1.1
