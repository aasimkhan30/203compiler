## PyCC - Python Based C Compiler
## GoCC - Go Based C Compiler


Given an input C language file, PyCC and GoCC will generate an assemble code (3 address code) which the gcc can assemble into an executable file.

# Stage 1 - Running the Compilers:
  Example C File : fib.c
  
  A) To run GoCC:
   ./203cc.sh -l go fib.c
   
  B) To run PyCC:
  ./203cc.sh -l python fib.c
  
# Stage 2 - Linking
 gcc -m32 fib.s -o fib

# Stage 3 - Executing
 ./fib 
  
  
# Features Supported:

Functions: if, while.

Integer, character, true and false literals. 

Operators:= ,  ||, &&, ==, !=, <, >=, +, -, *, ++, --, !, - (unary), [], ()

In Built functions :
strlen, strcpy...
fprintf, fputsf, puts....

Local and global variable parameters

String literals
