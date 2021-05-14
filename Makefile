CC = gcc
CFLAGS = -O2

all: cli

test:
	$(CC) $(CFLAGS) -o stub/stub stub/main.c stub/runtime.c -lelf
	rm -f stub/stub

clean:
	rm -f ward *_out

cli:
	go build .
