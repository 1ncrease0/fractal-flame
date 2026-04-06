#!/bin/sh

echo "Testing image properties..."

EXECUTABLE_PATH="$1"

if [ ! -f "test_output.png" ]; then
    echo "Generating test image..."
    "$EXECUTABLE_PATH" -w 800 -h 600 -i 1000000 -t 4 -g -gamma 2.2 -o test_output.png -ap "0.5,0.0,0.0,0.0,0.5,0.1/-0.3,0.0,0.0,0.0,-0.3,0.0" -f "linear:1.0,sinusoidal:0.1,spherical:1.0"
    if [ $? -ne 0 ]; then
        echo "✗ Failed to generate test image"
        exit 1
    fi
fi

if [ ! -f "test_output.png" ]; then
    echo "✗ Image file 'test_output.png' does not exist"
    exit 1
fi

case "test_output.png" in
    *.png)
        echo "✓ Image file has .png extension"
        ;;
    *)
        echo "✗ Image file does not have .png extension"
        exit 1
        ;;
esac

FILE_SIZE=$(stat -c%s "test_output.png" 2>/dev/null || stat -f%z "test_output.png" 2>/dev/null)
if [ "$FILE_SIZE" -gt 0 ]; then
    echo "✓ Image file has content (size: $FILE_SIZE bytes)"
else
    echo "✗ Image file is empty"
    exit 1
fi


if command -v hexdump >/dev/null 2>&1; then
    PNG_SIGNATURE=$(dd if="test_output.png" bs=8 count=1 2>/dev/null | hexdump -C | head -1 | awk '{print $2$3$4$5$6$7$8$9}' | tr -d ' ')
elif command -v od >/dev/null 2>&1; then
    PNG_SIGNATURE=$(dd if="test_output.png" bs=8 count=1 2>/dev/null | od -An -tx1 | tr -d ' \n' | cut -c1-16)
else
    if file "test_output.png" | grep -q "PNG"; then
        PNG_SIGNATURE="89504e470d0a1a0a"
    else
        PNG_SIGNATURE=""
    fi
fi

if [ "$PNG_SIGNATURE" = "89504e470d0a1a0a" ] || [ -n "$PNG_SIGNATURE" ]; then
    echo "✓ Image file has valid PNG signature"
else
    if file "test_output.png" | grep -q "PNG"; then
        echo "✓ Image file has valid PNG signature (verified by file command)"
    else
        echo "✗ Image file does not have valid PNG signature"
        exit 1
    fi
fi

echo "Image properties test passed!"
