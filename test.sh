#!/bin/bash

function test {
    expected="$1"
    expr="$2"

    echo "$expr" | go run main.go > tmp.s
    if [ ! $? ]; then
        echo "Failed to compile $expr"
        exit
    fi

    gcc -o tmp.out driver.c tmp.s || exit
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

rm -f tmp.out tmp.s
echo "All tests passed"
