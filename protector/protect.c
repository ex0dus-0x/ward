#define _GNU_SOURCE
#include <stdio.h>
#include <stdlib.h>
#include <unistd.h>
#include <fcntl.h>

#include <sys/types.h>
#include <sys/wait.h>
#include <sys/mman.h>

#include <libelf.h>

#include "runtime.h"

#define TEMPFILE "tmp"

/* helper method to exit with message */
static void die(int res, const char *msg)
{
    printf("error: %s\n", msg);
    exit(res);
}

/* executes code in-memory using memfd_create */
void exec_safe(void *data)
{
    int fd;
    ssize_t br;
    pid_t child;

    /* stores output */
    char buf[BUFSIZ] = "";

    /* create memfd */
    fd = memfd_create(TEMPFILE, 0);
    if (fd == -1)
        die(fd, "cannot create in-memory fd for code");

    /* fork a child of process */
    child = fork();
    if (child == -1) {
        close(fd);
        die(child, "cannot fork child process");
    }

    /* if fork successful, duplicate memory and execlp */
    if (!child) {
        dup2(fd, 1);
        close(fd);
        execlp("ls", "ls", NULL);
    }

    /* wait for child process to complete executing */
    waitpid(child, NULL, 0);

    /* read output */
    lseek(fd, 0, SEEK_SET);
    br = read(fd, buf, BUFSIZ);
    if (br < 0)
        exit(EXIT_FAILURE);
    buf[br] = '\0';

    /* print and close */
    printf("%s\n", buf);
    close(fd);
}


int main(int argc, char *argv[])
{
    int fd;

    Elf *e;

    /* open ourselves for reading */
    if ((fd = open(argv[0], O_RDONLY, 0)) < 0)
        die(-1, "cannot read ourselves as file");

    if (elf_version(EV_CURRENT) == EV_NONE)
        die(-1, elf_errmsg(-1));

    /* check if magic number is present */
    if ((e = elf_begin(fd, ELF_C_READ, NULL)) == NULL)
        die(-1, elf_errmsg(-1));

    /* close file and reuse for memfd */
    close(fd);
    return 0;
}
