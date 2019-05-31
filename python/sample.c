void main(){
    int a,b,c;
    printf("Hello world\n");
    b = 2;
    c = 1;
    if (c == 1)
        a = b + c;
    else
        a = b - c;
    printf("a = %d\n", a);
    while (a < 10){
        printf("203:  %d\n", a);
        a = a + 1;
    }
    return;
}