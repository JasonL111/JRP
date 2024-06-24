@echo off
REM Get the directory where the batch file is located
cd /d %~dp0

REM Delete the go.mod file (if it exists)
if exist go.mod del go.mod

REM Initialize the Go module
go mod init main

REM Organize Go module dependencies
go mod tidy

REM Build the Go project
go build

REM Run the executable file and redirect the output to result.txt
main.exe > result.txt

REM Display the completion message
echo Commands executed successfully and output saved to result.txt
pause