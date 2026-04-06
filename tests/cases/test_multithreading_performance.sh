#!/bin/sh

echo "Testing multithreading performance..."

measure_time() {
    threads=$1
    output_file="test_output_${threads}_threads.png"
    echo "Running with $threads threads..."

    if date +%s.%N >/dev/null 2>&1; then
        START_TIME=$(date +%s.%N)
        "$EXECUTABLE_PATH" -w 1920 -h 1080 -i 10000000 -t "$threads" -g -gamma 2.2 -o "$output_file" -ap "0.5,0.0,0.0,0.0,0.5,0.1/-0.3,0.0,0.0,0.0,-0.3,0.0" -f "linear:1.0,sinusoidal:0.1,spherical:1.0"
        EXIT_CODE=$?
        END_TIME=$(date +%s.%N)
        DURATION=$(awk "BEGIN {printf \"%.3f\", $END_TIME - $START_TIME}")
    else
        START_TIME=$(date +%s)
        "$EXECUTABLE_PATH" -w 1920 -h 1080 -i 10000000 -t "$threads" -g -gamma 2.2 -o "$output_file" -ap "0.5,0.0,0.0,0.0,0.5,0.1/-0.3,0.0,0.0,0.0,-0.3,0.0" -f "linear:1.0,sinusoidal:0.1,spherical:1.0"
        EXIT_CODE=$?
        END_TIME=$(date +%s)
        DURATION=$((END_TIME - START_TIME))
    fi
    if [ $EXIT_CODE -eq 0 ]; then
        echo "✓ Completed with $threads threads in ${DURATION} seconds"
        echo "$threads,$DURATION" >> performance_results.csv
        return 0
    else
        echo "✗ Failed with $threads threads (exit code: $EXIT_CODE)"
        return 1
    fi
}

if [ -z "$1" ]; then
    echo "✗ Executable file path not provided."
    echo "Usage: $0 <path_to_executable_file>"
    exit 1
fi

EXECUTABLE_PATH="$1"

if [ ! -f "$EXECUTABLE_PATH" ]; then
    echo "✗ Executable file '$EXECUTABLE_PATH' does not exist."
    exit 1
fi

echo "threads,duration_seconds" > performance_results.csv

for threads in 1 2 4; do
    if ! measure_time "$threads"; then
        echo "Performance test failed for $threads threads"
        exit 1
    fi
done

echo ""
echo "Performance test results:"
echo "------------------------"
cat performance_results.csv

if [ -f performance_results.csv ] && [ -s performance_results.csv ]; then
    echo ""
    echo "Speedup analysis (relative to 1 thread):"
    echo "----------------------------------------"
    baseline=$(awk -F',' 'NR==2 {print $2}' performance_results.csv)
    if [ -n "$baseline" ] && [ "$baseline" != "0" ]; then
        awk -F',' -v baseline="$baseline" '
        NR==1 {printf "%-10s %-15s %-15s\n", "Threads", "Time (s)", "Speedup"}
        NR>1 && $2 != "" {
            speedup = baseline / $2
            printf "%-10s %-15s %-15.2fx\n", $1, $2, speedup
        }' performance_results.csv
    fi
fi

echo ""
echo "Multithreading performance test completed!"
