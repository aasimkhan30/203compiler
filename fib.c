int fib(int arg){
    int previouspreviousNumber = 0;
    int previousNumber = 0;
    int currentNumber = 1;
        int i = 1;
        while (i < arg) {
            previouspreviousNumber = previousNumber;
            previousNumber = currentNumber;
            currentNumber = previouspreviousNumber + previousNumber;
            i++;
        }
        return currentNumber;
}

int main () {
    int result = fib(10);
    printf("%d\n", result);
}