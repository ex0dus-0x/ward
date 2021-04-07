CC = gcc
CFLAGS = -Wall -O2

all: cli

protect:
	$(CC) $(CFLAGS) -o protector/protector protector/protector.c protector/runtime.c -lelf

clean:
	rm -f ward *_out

cli:
	go build .
