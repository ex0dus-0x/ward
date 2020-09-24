CC := gcc
CFLAGS := -Wall -O2

all: cli protector

protector:
	$(CC) $(CFLAGS) protector/protect.c protector/runtime.c -lelf

cli:
	cargo build
