// copy from https://github.com/rui314/8cc
// アセンブラで定義された関数を実行して、標準出力することで動作確認できるようにする

#include <stdio.h>

extern int mymain(void);

int main(int argc, char **argv) {
  int val = mymain();
  printf("%d\n", val);
  return 0;
}
