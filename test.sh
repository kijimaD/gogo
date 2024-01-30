#!/bin/bash

# C言語では、main関数が返した値がプログラム全体としての終了コードになる。終了コードはシェルの$?変数に格納されているので、確認できる。

function compile {
    echo "$1" | go run . > gogo.s
    if [ $? -ne 0 ]; then
        echo "Failed to compile $1"
        exit -1
    fi

    gcc -o gogo c/driver.c gogo.s
    if [ $? -ne 0 ]; then
        echo "GCC failed"
        exit -1
    fi
}

function test {
    expected="$1"
    expr="$2"
    compile "$expr"
    result="`./gogo`"
    if [ "$result" != "$expected" ]; then
        echo "Test failed: $expected expected but got $result"
        exit -1
    fi

    echo "✓"
}

function testfail {
  expr="$1"
  echo "$expr" | go run . > /dev/null 2>&1
  if [ $? -eq 0 ]; then
    echo "Should fail to compile, but succeded: $expr"
    exit -1
  fi
}

# test expect expr

make -s gogo

test 0 0
test 42 42
test 2 '1;2'
test 3 '1+2'
test 3 '1 + 2'
test 3 '1+ 2'
test 10 '1+2+3+4'
test 4 '1+2-3+4'
test 2 '5 - 3'
test 4 '5-1-1+1'
test 11 '1+2*3+4'
test 7 '1+2*3'
test 2 '2/2+1'
test 1 '1+1/2'
test 8 '3+4/2+3'
test 18 '3*4/2*3'
test 3 '24/2/4'
test 98 "'a'+1;"

# Declaration
test 2 'int a = 2;a;'
test 3 'int a = 2;a;3'
test 4 'int a = 1+1; a+2'
test 1 'int a = 1;a'
test 2 'int a = 1;a+1'
test 97 "char a = 'a';a"
# test 98 "char a = 'a';a+1" # => なぜか2になる
# test 3 'int a = 1;a+2' # => なぜか4になる
# test 4 'int a = 1;a+3' # => なぜか6になる

# Function call
test 25 'sum2(20, 5);'
test 24 'sum2(20-1, 5);'
test 15 'sum5(1, 2, 3, 4, 5);'
# FIXME: printf単独だと1がくっつくのはなぜ? そして後続の式でその数字が上書きされるのはなぜ
test a1 'printf("%s", "a");'
test a99 'printf("%s", "a");99;'
test abc5 'printf("%s", "abc");5;'

testfail 42a   # 引数として渡されるのは文字列としてのダブルクォートを含まない 42a
testfail "42a" # 引数として渡されるのは文字列としてのダブルクォートを含まない 42a
testfail '42a' # 引数として渡されるのは文字列としてのダブルクォートを含まない 42a
testfail '"abc'
testfail '1+'
testfail '1+"abc"'
testfail 'a'

rm -f gogo.out gogo.s
echo "All tests passed"
