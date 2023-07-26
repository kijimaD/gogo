#!/bin/bash

function compile {
    echo "$1" | go run main.go > tmp.s
    if [ $? -ne 0 ]; then
        echo "Failed to compile $1"
        exit
    fi

    gcc -o tmp.out c/driver.c tmp.s
    if [ $? -ne 0 ]; then
        echo "GCC failed"
        exit
    fi
}

function test {
    expected="$1"
    expr="$2"
    compile "$expr"
    result="`./tmp.out`"
    if [ "$result" != "$expected" ]; then
        echo "Test failed: $expected expected but got $result"
        exit
    fi

    echo "✓"
}

# test expect expr

test 0 0
test 42 42
test hello "hello"
test "hello world" "hello world"

rm -f tmp.out tmp.s
echo "All tests passed"
