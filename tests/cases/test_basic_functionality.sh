#!/bin/sh

echo "Testing basic functionality..."

EXECUTABLE_PATH="$1"

echo "Running: $EXECUTABLE_PATH -w 800 -h 600 -i 1000000 -t 4 -g -gamma 2.2 -o test_output.png -ap \"0.5,0.0,0.0,0.0,0.5,0.1/-0.3,0.0,0.0,0.0,-0.3,0.0\" -f \"linear:1.0,sinusoidal:0.1,spherical:1.0\""
"$EXECUTABLE_PATH" -w 800 -h 600 -i 1000000 -t 4 -g -gamma 2.2 -o test_output.png -ap "0.5,0.0,0.0,0.0,0.5,0.1/-0.3,0.0,0.0,0.0,-0.3,0.0" -f "linear:1.0,sinusoidal:0.1,spherical:1.0"

EXIT_CODE=$?
if [ $EXIT_CODE -eq 0 ]; then
    echo "✓ Application exited successfully (exit code: $EXIT_CODE)"
else
    echo "✗ Application failed with exit code: $EXIT_CODE"
    exit 1
fi

if [ -f "test_output.png" ]; then
    echo "✓ Image file 'test_output.png' was created"
else
    echo "✗ Image file 'test_output.png' was not created"
    exit 1
fi

echo "Basic functionality test passed!"
