'''
    This is the main file for the PyCC compiler.
    There are 3 steps for the entire compilation process.
    1. Lexer: Identifying token in the files.
    2. Parser: Parsing those tokens.
    3. Code Generation: Generate the actual assembly code for the compiler.

    Given an input file, PYCC will generate an assemble code (3 address code)
    which the gcc can assemble into an executable file.
    Hopefully, This file should run on bash.

    NO CODE OPTIMIZATIONS IS CONSIDERED AT THIS MOMENT.
'''
import os, sys

from code_generator import code_generator
from lex import lexer
from parser import parser

with open(sys.argv[1], "r") as input_file:
    print("Called on "+ sys.argv[1])
    file_str = input_file.read()
lex = lexer(file_str)
cg = code_generator(sys.argv[1].split(".")[0]+'.s')
parser = parser(lex, cg)
parser.parse_source()


# while(lex.next()):
#     pass


# cg = code_generator()
#
# cg.set_code_file(file_str.split(".")[0]+'.s')
