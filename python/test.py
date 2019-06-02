from parser import parser
from code_generator import code_generator
from lex import lexer
from code_generator import code_generator




#Test Cases for lexer:

def t1():
    file = open("sample.c", "r")
    lines = ' '.join(file.readlines())
    l = lexer(lines)
    while l.next()is not None:
        print(l.get_token())


def t2():
    cg = code_generator("t2")
    cg.write_code_to_file("Hello World")
    assert (cg.instructions("+") == "add")
    assert (cg.form_label() == "L0")
    assert (cg.form_label() == "L1")
    L3 = cg.form_label()
    cg.write_label(L3)
    cg.write_function_begin("test_fun", 3, L3)
    cg.write_function_end()

def t3():
    file = open("sample.c", "r")
    lines = ' '.join(file.readlines())
    p = parser (lines, "t3")
    p.parse_source()
    print(p.scope_local)

# t1()
# t2()
t3()