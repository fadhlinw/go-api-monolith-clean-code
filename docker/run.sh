while true; do
    echo "[run.sh] Run migration..."
    sql-migrate up
    echo "[run.sh] Starting debugging..."
    go run . app:serve &

    PID=$!

    inotifywait -e modify -e move -e create -e delete -e attrib --exclude '(__debug_bin|\.git)' -r .

    echo "[run.sh] Stopping process id: $PID"
    
    kill -9 $PID
done