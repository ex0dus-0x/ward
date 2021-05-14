/* runtime.c
 * ==============
 * Implements no-stdlib code injection checks to prevent executable from being traced.
 */
#include <fcntl.h>
#include <stdio.h>

#include "runtime.h"

/* static helper clone for str substitution */
static char* afterSubstr(char *str, const char *sub)
{
    int i, found;
    char *ptr;
    found = 0;
    for (ptr = str; *ptr != '\0'; ptr++)
    {
        found = 1;
        for(i = 0; found == 1 && sub[i] != '\0'; i++)
            if(sub[i] != ptr[i])
                found = 0;
        if(found == 1)
            break;
    }
    if (found == 0)
        return NULL;
    return ptr + i;
}


/* static helper for checking libc */
static int isLib(char *str, const char *lib)
{
    int i, found;
    static const char *end = ".so\n";
    char *ptr;

    // Trying to find lib in str
    ptr = afterSubstr(str, lib);
    if (ptr == NULL)
        return 0;

    // Should be followed by a '-'
    if (*ptr != '-')
        return 0;

    // Checking the first [0-9]+\.
    found = 0;
    for (ptr += 1; *ptr >= '0' && *ptr <= '9'; ptr++)
        found = 1;
    if (found == 0 || *ptr != '.')
        return 0;

    // Checking the second [0-9]+
    found = 0;
    for (ptr += 1; *ptr >= '0' && *ptr <= '9'; ptr++)
        found = 1;
    if (found == 0)
        return 0;

    // Checking if it ends with ".so\n"
    for (i = 0; end[i] != '\0'; i++)
        if (end[i] != ptr[i])
            return 0;

    return 1;
}


/* stops any attempts to utilize the `LD_PRELOAD` envvar */
int check_preloading(void)
{
    long i, j;
    char env[] = "LD_PRELOAD";

    for (i = 0; environ[i]; i++) {

        // check each char until we reach the null pointer
        for (j = 0; environ[j] != '\0' && environ[i][j] != '\0'; j++) {
            if (env[j] != environ[i][j])
                break;
        }

        // detected if every char matched
        if (env[j] == '\0')
            return 1;
    }
    return 0;
}

/* introspect the running PID's memory mappings for unwarranted libraries */
int check_mmaps(void)
{
    // stores result of execution
    char buffer[BUF_SIZE];

    // set after encounter libc shared object
    int after_libc = 0;

    // attempt opening proc mappings for current pid
    FILE *mmap;
    mmap = fopen("/proc/self/maps", "r");
    if (mmap == NULL) {
        return -1;
    }

    // iteratively read BUF_SIZE from mmap file
    while (fgets(buffer, BUF_SIZE, mmap) != NULL)
    {
        // check if libc so entry is found and set flag
        if (isLib(buffer, "libc"))
            after_libc = 1;

        // break once we reach ld entry
        if (isLib(buffer, "ld"))
            break;

        // check if not anonymous mapping
        if (after_libc && (afterSubstr(buffer, "00000000 00:00 0")  == NULL)) {
            return -1;
        }
    }
    return 0;
}
