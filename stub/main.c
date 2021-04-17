/*
 * stub.c
 * ===============
 *
 * Implementation of an application that implements self-protection techniques,
 * while unpacking the original executable and executing it filelessly.
 */

#define _GNU_SOURCE
#include <stdio.h>
#include <stdlib.h>
#include <unistd.h>
#include <fcntl.h>

#include <sys/types.h>
#include <sys/mman.h>

#include <libelf.h>
#include <gelf.h>
#include <zlib.h>

#include "runtime.h"

#define TEMPFILE "[kworker/1:1]"
#define MIN(x, y) x > y ? y : x

/* helper method to exit with message */
static void die(int res, const char *msg)
{
    printf("Error: %s\n", msg);
    exit(res);
}


/* decompress executable blob for memfd */
static void decompress()
{


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

/* handles anti-analysis checks */
void __attribute__ ((constructor)) premain()
{
}

int main(int argc, char *argv[], char *envp[])
{
    int fd;
    Elf *e;

    // open ourselves for reading
    if ((fd = open(argv[0], O_RDONLY, 0)) < 0)
        die(-1, "cannot read ourselves as file");

    // check if valid binary version
    if (elf_version(EV_CURRENT) == EV_NONE)
        die(-1, elf_errmsg(-1));

    // check if ELF magic number is present
    if ((e = elf_begin(fd, ELF_C_READ, NULL)) == NULL)
        die(-1, elf_errmsg(-1));

    // parse out number of program headers
    size_t n;
    int ret = elf_getphdrnum(e, &n);
    if (ret != 0)
        die(-1, "cannot parse any program headers");
    
    // get the first PT_NOTE segment we find
    GElf_Phdr* phdr = NULL;
    for (size_t i = 0; i < n; i++) {
        GElf_Phdr tmp;
        if (!gelf_getphdr(e, i, &tmp))
            die(-1, "cannot get program header");

        if (tmp.p_type == PT_NOTE) {
            phdr = &tmp;
            break;
        }
    }

    if (!phdr)
        die(-1, "cannot find PT_NOTE segment to further parse");

    // get attributes for PT_NOTE segment
    Elf64_Off offset = phdr->p_offset;
    Elf64_Xword size = phdr->p_filesz;

    // read packed executable from file offset
    char data[size];
    lseek(fd, 0, SEEK_SET);
    lseek(fd, offset, SEEK_SET);
    pread(fd, (void*) data, size, offset);

    close(fd);

    // decompress the data parsed

    // create anonymous file
    fd = memfd_create(TEMPFILE, 0);
    if (fd == -1)
        die(fd, "cannot create in-memory fd for code");

    // delete ourselves
    if (remove(argv[0]) != 0)
        die(fd, "cannot remove self");

    // write ELF data to in memory fd and execute
    write_fd(fd, data, size - 1);
    {
        //const char *args[] = {TEMPFILE, NULL};
        argv[0] = TEMPFILE;
        fexecve(fd, (char * const *) argv, (char * const *) envp);
    }
    close(fd);
    return 0;
}
