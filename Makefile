CC = gcc
CFLAGS = -Wall -O2

all: cli protect

protect:
	$(CC) $(CFLAGS) -o protector/protector protector/protect.c protector/runtime.c -lelf

cli:
	go build .
