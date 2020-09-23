# ward-protector

This ELF binary encapsulates over the original target binary being protected. It implements self-parsing to retrieve the original blob of data that represents the target executable, runs its own protection routines, and runs the executable as an _anonymous file_ with the `memfd_create` system call.

## Usage

`ward` will automatically compile and inject the source for the protector application. However, if you choose to do so manually yourself, make sure you have `libelf` installed and:

```
$ gcc -Wall protector/protect.c protector/runtime.c -lelf
```
