# ward

Simple ELF runtime packer for creating stealthy droppers

## Introduction

This is a simple implementation of an ELF packer that creates stealthy droppers for loading
malicious ELFs in-memory. Useful for red teamers trying to proliferate a payload while evading
detection.

## How It Works

__ward__ compresses a target ELF executable and injects it into a stub program,
which uses a modified `PT_NOTE` infection technique to execute it in-memory with `memfd_create`
and `fexec`.

For instance, run __ward__ on a copy of `ls`:

```
$ ward pack ./ls
2021/04/14 20:26:07 Starting up ward
2021/04/14 20:26:07 Checking if valid ELF binary
2021/04/14 20:26:07 Provisioning stub program for packing
2021/04/14 20:26:07 Packing original executable into stub ./ls
2021/04/14 20:26:07 Finding PT_NOTE segment for injecting metadata
2021/04/14 20:26:07 Offset: 828304 Size: 141936
2021/04/14 20:26:07 Writing (not yet encoded) ELF to stub
2021/04/14 20:26:07 Done! Find the packed application at /home/alan/Code/ward/ls.packed
```

When you execute it now, the stub program will read the compressed executable from itself,
and create an anonymous file descriptor for execution. Once executed, the file will disappear
from the disk:

```
$ ./ls.packed
example  go.mod  go.sum  injector.go  ls  ls.packed  main.go  Makefile  README.md  stub  ward
```

## License

[MIT License](https://codemuch.tech/docs/license.txt)
