# gogo

implement C by Go

imspired & reference https://github.com/rui314/8cc

## using example

add

```shell
$ echo '1+2' | go run . > gogo.s
$ cat gogo.s # check
.text
        .global mymain
mymain:
        mov $2, %eax
        push %rax
        mov $1, %eax
        pop %rbx
        add %ebx, %eax
        ret
$ gcc -o gogo c/driver.c gogo.s
$ ./gogo
3
```

declaration

```
$ echo 'int a = 1+1; a+2' | go run . > gogo.s
$ cat gogo.s # check
.text
        .global mymain
mymain:
        mov $1, %eax
        push %rax
        mov $1, %eax
        pop %rbx
        add %ebx, %eax
        mov %eax, -4(%rbp)
        mov $2, %eax
        push %rax
        mov %eax, -0(%rbp)
        pop %rbx
        add %ebx, %eax
        ret
$ gcc -o gogo c/driver.c gogo.s
$ ./gogo
```
