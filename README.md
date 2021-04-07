# ward

My implementation of an ELF packer in Golang

## Build

Requirements:

* `make`
* `clang`

Use the `Makefile`:

```
$ make
```

## Packing Technique

1. Compile a packer runtime app, which employs several checks to ensure implants aren't attempting to inject themselves into the app
2. Compress the original executable, and use `PT_NOTE` injection technique to hide statically in a code cave in the packer runtime, writing the file offset to a segment
3. When executed, the packed executable will retrieve the blob of data from the file offset and use `memfd_create` to execute the original entry point

## License

[MIT]()
