#define _GNU_SOURCE
#include <stdio.h>
#include <stdlib.h>
#include <unistd.h>
#include <sys/types.h>
#include <sys/wait.h>
#include <sys/mman.h>

#define TEMPFILE "tmp"

int main(int argc, char *argv[])
{
    int fd;
    pid_t child;
    char buf[BUFSIZ] = "";
    ssize_t br;

    fd = memfd_create(TEMPFILE, 0);
    if (fd == -1)
        exit(EXIT_FAILURE);

    child = fork();
    if (!child) {
        dup2(fd, 1);
        close(fd);
        execlp(argv[1], argv[1], NULL);
        exit(EXIT_FAILURE);
    } else if (child == -1)
        exit(EXIT_FAILURE);

    waitpid(child, NULL, 0);

    lseek(fd, 0, SEEK_SET);
    br = read(fd, buf, BUFSIZ);
    if (br < 0)
        exit(EXIT_FAILURE);

    buf[br] = '\0';
    printf("child said: '%s'", buf);
    return 0;
}
