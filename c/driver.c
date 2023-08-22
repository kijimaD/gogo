// copy from https://github.com/rui314/8cc
// アセンブラで定義された関数を実行して、標準出力することで動作確認できるようにする

#include <stdio.h>

extern int mymain(void);

int sum2(int a, int b) {
  return a + b;
}

int sum5(int a, int b, int c, int d, int e) {
  return a + b + c + d + e;
}

int main(int argc, char **argv) {
  printf("%d\n", mymain());
  return 0;
}
