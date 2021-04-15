# ward

a dumb ELF packer

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
and create an anonymous file descriptor for execution. What the user sees however is the reg

```
$ ./ls.packed
example  go.mod  go.sum  injector.go  ls  ls.packed  main.go  Makefile  README.md  stub  ward
```

__ward__ implements a code injection check as part of its "anti-analysis". This is rudimentary
and can be substituted for other techniques, obfuscations, or none at all.

## TODO

* [ ] Actual compression of some sort

* [ ] Stealthy dropping - execute will remove binary from disk, maybe option to read executable
from socket to C2.
