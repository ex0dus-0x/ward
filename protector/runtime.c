/* runtime.c
 * ==============
 * no-stdlib implementation of the runtime protection functionality
 */

extern char **environ;

/* stops any attempts to utilize the `LD_PRELOAD` envvar */
int check_preloading(void)
{
    long i, j;
    char env[] = "LD_PRELOAD";

    for (i = 0; environ[i]; i++) {

        /* check each char until we reach the null pointer */
        for (j = 0; environ[j] != '\0' && environ[i][j] != '\0'; j++) {
            if (env[j] != environ[i][j])
                break;
        }

        /* detected if every char matched */
        if (env[j] == '\0')
            return 1;
    }
    return 0;
}

/* introspect the running PID's memory mappings for unwarranted libraries */
int check_mmaps(int pfd)
{
    char buf[BUF_SIZE];

    /* attempt opening proc mappings for current pid
    mmap = fopen("/proc/self/maps", "r");
    if (mmap == NULL) {
        // TODO set perror
        return -1;
    }

    while (fgets(buffer, BUFFER_SIZE, memory_map) != NULL)
    {
    }
    */

    return 0;
}
