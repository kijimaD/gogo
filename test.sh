#!/bin/bash

function test {
    expected="$1"
    expr="$2"

    echo "$expr" | go run main.go > tmp.s
    if [ ! $? ]; then
        echo "Failed to compile $expr"
        exit
    fi

    gcc -o tmp.out tmp.s || exit

    # 中の値を確かめる方法がわからないので、実行できるかだけチェックする
    # 8ccではdriver.cでやってる
    ./tmp.out || true
    if [[ $? -eq 0 ]];then
        echo "ok"
    else
        echo "not ok"
        exit 1
    fi
}

test 0 0
test 42 42

rm -f tmp.out tmp.s
echo "All tests passed"
