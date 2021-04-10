CC = gcc
CFLAGS = -Wall -O2

all: cli

stub:
	$(CC) $(CFLAGS) -o stub/stub stub/main.c stub/runtime.c -lelf

clean:
	rm -f ward *_out

cli:
	go build .
