#include <stdio.h>
#include <stdlib.h>
#include <time.h>

int main(void) {
    srand(time(NULL));
    int i = 10;
    while(i--) printf("%d ",rand()%100);
    printf("\n");
    return 0;
}
