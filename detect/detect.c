int main(int argc, char *argv[], char **environ)
{
    int i, j;
    char env[] = "LD_PRELOAD";
    for(i = 0; environ[i]; i++)
    {
        for(j = 0; env[j] != '\0' && environ[i][j] != '\0'; j++)
            if(env[j] != environ[i][j])
                break;

        if(env[j] == '\0')
        {
            return 1;
        }
    }
    return 0;
}
