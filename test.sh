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
test 42a '"42a"'
test hello '"hello"'
test "hello world" '"hello world"'
test 3 '1+2'
test 3 '1 + 2'
test 3 '1+ 2'
test 10 '1+2+3+4'
test 4 '1+2-3+4'
test 2 '5 - 3'
test 4 '5-1-1+1'

testfail 42a   # 引数として渡されるのは文字列としてのダブルクォートを含まない 42a
testfail "42a" # 引数として渡されるのは文字列としてのダブルクォートを含まない 42a
testfail '42a' # 引数として渡されるのは文字列としてのダブルクォートを含まない 42a
testfail '"abc'
testfail '1+'
testfail '1+"abc"'

rm -f tmp.out tmp.s
echo "All tests passed"
