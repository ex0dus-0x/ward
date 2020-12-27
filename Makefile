CC = gcc
CFLAGS = -Wall -O2

all: cli

protect:
	$(CC) $(CFLAGS) -o protector/protector protector/protect.c protector/runtime.c -lelf

clean:
	rm -r ward

cli:
	go build .
