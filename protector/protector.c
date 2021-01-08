/*
 * protector.c
 * ===============
 *
 * Implementation of a protector application that parses itself
 * for the original ELF binary, while enforcing a protection runtime to
 * mitigate code injection attacks.
 */

#define _GNU_SOURCE
#include <stdio.h>
#include <stdlib.h>
#include <unistd.h>
#include <fcntl.h>

#include <sys/types.h>
#include <sys/mman.h>

#include <libelf.h>

#include "runtime.h"

#define TEMPFILE "tmp"
#define MIN(x, y) x > y ? y : x

/* helper method to exit with message */
static void die(int res, const char *msg)
{
    printf("error: %s\n", msg);
    exit(res);
}


/* safely write buffer to a given input file descriptor */
static void write_fd(int fd, const char *str, size_t len)
{
    size_t cnt = 0;
    do {
        ssize_t result = write(fd, str + cnt, MIN(len - cnt, 0x7ffff000));
        if (result == -1)
            die(-1, "writing to memfd failed\n");
        cnt += result;
    } while (cnt != len);
}


/* executes code in-memory using memfd_create */
void exec_safe(const char *data)
{
    int fd;

    /* create memfd */
    fd = memfd_create(TEMPFILE, 0);
    if (fd == -1)
        die(fd, "cannot create in-memory fd for code");

    /* write ELF blob to in memory fd and execute */
    write_fd(fd, data, sizeof(data) - 1);
    {
        const char *argv[] = {TEMPFILE, NULL};
        const char *envp[] = {NULL};
        fexecve(fd, (char * const *) argv, (char * const *) envp);
    }
    close(fd);
}


int main(int argc, char *argv[])
{
    int fd;
    Elf *e;

    /* open ourselves for reading */
    if ((fd = open(argv[0], O_RDONLY, 0)) < 0)
        die(-1, "cannot read ourselves as file");

    /* check if valid ELF version */
    if (elf_version(EV_CURRENT) == EV_NONE)
        die(-1, elf_errmsg(-1));

    /* check if magic number is present */
    if ((e = elf_begin(fd, ELF_C_READ, NULL)) == NULL)
        die(-1, elf_errmsg(-1));

    /* iterate over program headers and find rewritten PT_NOTE */

    /* get virtual address offset and file size to read */

    /* read ELF file pointed to in-memory */

    /* close file and reuse for memfd */
    close(fd);
    return 0;
}
