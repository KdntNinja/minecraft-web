#!/bin/zsh
# Kill any process using port 3000, then start the Webcraft server

PORT=3000

# Find and kill process using the port
PID=$(lsof -ti tcp:$PORT)
if [ -n "$PID" ]; then
  echo "Killing process on port $PORT (PID: $PID)"
  kill -9 $PID
else
  echo "No process found on port $PORT"
fi
