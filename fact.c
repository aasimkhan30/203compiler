int factorial(int arg){
    int i = 1;
    int sum = 1;
    while(i<=arg){
        printf("%d\n", sum);
        sum = sum * i;
        i++;
    }
    return sum;
}

int main () {
    int result = factorial(5);
    printf("%d\n", result);
}