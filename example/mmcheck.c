#include <stdlib.h>
#include <stdio.h>
#include <fcntl.h>
#include <sys/stat.h>

#define BUFFER_SIZE 256

// Avoid to use libc strstr
// Return a pointer after the first location of sub in str
char* afterSubstr(char *str, const char *sub)
{
    int i, found;
    char *ptr;
    found = 0;
    for(ptr = str; *ptr != '\0'; ptr++)
    {
        found = 1;
        for(i = 0; found == 1 && sub[i] != '\0'; i++)
            if(sub[i] != ptr[i])
                found = 0;
        if(found == 1)
            break;
    }
    if(found == 0)
        return NULL;
    return ptr + i;
}

// Try to match the following regexp: libname-[0-9]+\.[0-9]+\.so$
// Not using any libc function makes that code awful, I know
int isLib(char *str, const char *lib)
{
    int i, found;
    static const char *end = ".so\n";
    char *ptr;
    // Trying to find lib in str
    ptr = afterSubstr(str, lib);
    if(ptr == NULL)
        return 0;
    // Should be followed by a '-'
    if(*ptr != '-')
        return 0;
    // Checking the first [0-9]+\.
    found = 0;
    for(ptr += 1; *ptr >= '0' && *ptr <= '9'; ptr++)
        found = 1;
    if(found == 0 || *ptr != '.')
        return 0;
    // Checking the second [0-9]+
    found = 0;
    for(ptr += 1; *ptr >= '0' && *ptr <= '9'; ptr++)
        found = 1;
    if(found == 0)
        return 0;
    // Checking if it ends with ".so\n"
    for(i = 0; end[i] != '\0'; i++)
        if(end[i] != ptr[i])
            return 0;
    return 1;
}

int main()
{
    FILE *memory_map;
    char buffer[BUFFER_SIZE];
    int after_libc = 0;
    memory_map = fopen("/proc/self/maps", "r");
    if(memory_map == NULL)
    {
        printf("/proc/self/maps is unaccessible, probably a LD_PRELOAD attempt\n");
        return 1;
    }
    // Read the memory map line by line
    // Try to look for a library loaded in between the libc and ld
    while(fgets(buffer, BUFFER_SIZE, memory_map) != NULL)
    {
        // Look for a libc entry
        if(isLib(buffer, "libc"))
            after_libc = 1;
        else if(after_libc)
        {
            // Look for a ld entry
            if(isLib(buffer, "ld"))
            {
                // If we got this far then everythin is fine
                printf("Memory maps are clean\n");
                break;
            }
            // If it's not an anonymous memory map
            else if(afterSubstr(buffer, "00000000 00:00 0") == NULL)
            {
                // Something has been preloaded by ld.so
                printf("LD_PRELOAD detected through memory maps\n");
                break;
            }
        }
    }
}
