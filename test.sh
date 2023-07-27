#!/bin/bash

function compile {
    echo "$1" | go run . > tmp.s
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

function testfail {
  expr="$1"
  echo "$expr" | go run . > /dev/null 2>&1
  if [ $? -eq 0 ]; then
    echo "Should fail to compile, but succeded: $expr"
    exit
  fi
}

# test expect expr

test 0 0
test 42 42
test hello '"hello"'
test "hello world" '"hello world"'

testfail '"abc' # 成功してしまう

rm -f tmp.out tmp.s
echo "All tests passed"
