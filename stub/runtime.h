#ifndef RUNTIME_H
#define RUNTIME_H

#define BUF_SIZE 2048

extern char **environ;

/* anti-tampering: check if LD_PRELOAD envvar is set */
int check_preloading(void);

/* anti-tampering: check if other executable code is mmapped */
int check_mmaps(void);

#endif
