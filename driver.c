// copy from https://github.com/rui314/8cc
// アセンブラで定義された関数を実行して、標準出力することで動作確認できるようにする

#include <stdio.h>

#define WEAK __attribute__((weak))
extern int intfn(void) WEAK;

int main(int argc, char **argv) {
  if (intfn) {
    printf("%d\n", intfn());
  } else {
    printf("Should not happen");
  }
  return 0;
}
