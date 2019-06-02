from code_generator import code_generator
from lex import lexer

class parser:
    def __init__(self, lex_str, cg_str):
        self.lex = lexer(lex_str)
        self.cg = code_generator(cg_str)
        self.scope_global = []
        self.scope_local = []
        self.errors = []

    def local_check(self, name):
        for e in self.scope_local:
            if e[0] == name:
                return True
        return False

    def global_check(self, name):
        for e in self.scope_global:
            if e[0] == name:
                return True
        return False

    def local_params(self):
        c = 0
        for l in self.scope_local:
            if l[1] == False:
                c += 1
        return c

    def var_offset(self, name):
        i = 0
        li = 0
        for e in self.scope_local:
            if e[0] == name:
                if e[1] == True:
                    return (i+2)*self.cg.stack
                else:
                    return -1 * (li + 1)*self.cg.stack
            if e[1]:
                i += 1
            else:
                li += 1

    def var_offset_str(self, offset):
        if offset < 0:
            return  str(offset)
        elif offset > 0:
            return "+"+str(offset)
        else:
            return ""

    def check_token(self, token, expected):
        if token not in expected:
            print("Expected : ", expected, " Found: ", token)
            return False
        else:
            return True

    # Main function where we are getting the source file.
    def parse_source(self):
        self.std_fun()
        print(self.scope_global)
        self.program()

    # Just adding the global functions to the scope
    def std_fun(self):
        self.scope_global.append(("printf", True))
        self.scope_global.append(("malloc", True))
        self.scope_global.append(("calloc", True))
        self.scope_global.append(("free", True))
        self.scope_global.append(("atoi", True))
        self.scope_global.append(("isaplha", True))
        self.scope_global.append(("isdigit", True))
        self.scope_global.append(("isalnum", True))
        self.scope_global.append(("strcpy", True))
        self.scope_global.append(("strlen", True))
        self.scope_global.append(("strcpy", True))
        self.scope_global.append(("strchr", True))
        self.scope_global.append(("strdup", True))
        self.scope_global.append(("fopen", True))
        self.scope_global.append(("fclose", True))
        self.scope_global.append(("fgetc", True))
        self.scope_global.append(("ungetc", True))
        self.scope_global.append(("feof", True))
        self.scope_global.append(("fputs", True))
        self.scope_global.append(("fprintf", True))
        self.scope_global.append(("puts", True))


    def program(self):
        self.cg.write_code_to_file(".intel_syntax noprefix")
        while self.lex.next() != None:
            self.declare(self.lex.get_token())

    def declare(self, token):
        # Now the token received must be of a dataype
        self.check_token(token, {"void", "int", "char"})
        name = self.lex.next()
        # User can either declare a fun or var
        temp = self.lex.next()
        self.check_token(temp, {"(", "="})
        if temp == "(":
            self.decf(name)
        else:
            pass

    def decf(self, name):
        # Clearing local variables from previous fun
        self.scope_local.clear()
        # Adding syntax
        self.cg.write_code_to_file(".section .text")
        # Adding current function
        self.scope_global.append((name, True))
        #Since we have already passed ( now passing till )
        while self.lex.next() != ")":
            datatype = self.lex.get_token()
            name = self.lex.next()
            self.scope_local.append((name, True))
            comma = self.lex.next()
        self.check_token(self.lex.get_token(), {")"})

        tag = self.cg.form_label()
        self.cg.write_code_to_file(tag)
        # Cheking for open brackets
        self.check_token(self.lex.next(), {"{"})
        # Parsing inner statements
        self.statements(self.lex.get_token())
        # Adding epilogue
        self.cg.write_function_end()
        # Adding function prologue
        self.cg.write_function_begin(name, self.local_params(), tag)
        pass

    def statements(self, token):
        if token == '{':
            self.lex.next()
            while self.lex.get_token() != '}':
                self.statements(self.lex.get_token())
            self.check_token(self.lex.get_token(), {"}"})
            self.lex.next()
        elif token == "int" or token == "char":
            name = self.lex.next()
            self.scope_local.append((name, False))
            equal = self.lex.next()
            if equal == "=":
                self.expr(self.lex.next())
                offset = self.var_offset(name)
                self.cg.write_code_to_file("\t mov [ebp"+ self.var_offset_str(offset)+"], eax")
            self.check_token(self.lex.get_token(), ";")
            self.lex.next()
        elif token == "if":
            self.branch(token)
        elif token == "while":
            self.loop(token)
        elif token  == "return":
            self.expr(self.lex.next())
            self.check_token(self.lex.get_token(), ";")
            self.lex.next()
        else:
            self.expr(token)
            self.check_token(self.lex.get_token(), ";")
            self.lex.next()
        pass

    def expr(self, token):
        print("EXPR Called")
        self.expr1(token)

    def expr1(self, token):
        self.expr2(token)
        if self.lex.get_token() == '=':
            name = token
            expr = self.lex.next()
            self.expr2(expr)
            offset = self.var_offset(name)
            self.cg.write_code_to_file("\t mov [ebp" + self.var_offset_str(offset) + "], eax")

    def expr2(self, token):
        self.expr3(token)
        operator = self.lex.get_token()
        if operator == "||" or  operator == "&&":
            label = self.cg.form_label()
            self.cg.write_code_to_file("\tcmp eax, 0")
            if operator == "||":
                self.cg.write_code_to_file("\tjnz "+label)
            if operator == "||":
                self.cg.write_code_to_file("\tjz "+label)
            self.expr3(self.lex.next())
            self.cg.write_label(label)

    def expr3(self, token):
        self.expr4(token)
        operator = self.lex.get_token()
        if operator == '<=' or operator == '>=' or operator == '<' or operator == '>' or operator == '!=' or operator == '==':
            self.cg.write_code_to_file("\tpush eax")
            self.expr4(self.lex.next())
            self.cg.write_code_to_file("\tmov ebx, eax")
            self.cg.write_code_to_file("\tpop eax")
            instruction = self.cg.instructions(operator)
            self.cg.write_code_to_file("\tcmp eax, ebx")
            self.cg.write_code_to_file("\t"+instruction+" al")
            self.cg.write_code_to_file("\tmovzx eax, al")


    def expr4(self, token):
        self.expr5(token)
        operator = self.lex.get_token()
        while operator == "+" or operator == "-":
            self.cg.write_code_to_file("\tpush eax")
            self.expr5(self.lex.next())
            self.cg.write_code_to_file("\tmov ebx, eax")
            self.cg.write_code_to_file("\tpop eax")
            instruction = self.cg.instructions(operator)
            self.cg.write_code_to_file("\t"+instruction+" eax, ebx")
            operator = self.lex.get_token()

    def expr5(self, token):
        self.expr6(token)
        operator = self.lex.get_token()
        while operator == "*" or operator == "/":
            self.cg.write_code_to_file("\tpush eax")
            self.expr6(self.lex.next())
            self.cg.write_code_to_file("\tmov ebx, eax")
            self.cg.write_code_to_file("\tpop eax")
            instruction = self.cg.instructions(operator)
            self.cg.write_code_to_file("\t" + instruction + " ebx")
            operator = self.lex.get_token()

    def expr6(self, token):
        if token == "-":
            self.expr6(self.lex.next())
            self.cg.write_code_to_file("\t neg eax")
        else:
            self.unary(token)
            operator = self.lex.get_token()
            if operator == "++" or operator == "--":
                instruction = self.cg.instructions(operator)
                self.cg.write_code_to_file("\t"+instruction+" eax")
                offset = self.var_offset(token)
                self.cg.write_code_to_file("\t mov [ebp" + self.var_offset_str(offset) + "], eax")

    def unary(self, token):
        name = token
        type = self.lex.token_type
        if type == 2:
            self.cg.write_code_to_file("\tmov eax, "+self.lex.get_token())
        elif type == 4:
            label = self.cg.form_label()
            self.cg.write_code_to_file("\tmov eax, offset "+label)
            self.cg.write_code_to_file(".section .data")
            self.cg.write_code_to_file(label+": .asciz "+ self.lex.get_token())
            self.cg.write_code_to_file(".section .text")
        elif type == 1:
            if self.local_check(name):
                offset = self.var_offset(name)
                self.cg.write_code_to_file("\t mov eax, [ebp" + self.var_offset_str(offset) + "], eax")
            elif self.global_check(name):
                self.cg.write_code_to_file("\tlea eax, ["+name+"]")
        curr = self.lex.next()
        if curr == "(":
            self.fun_call(name, self.lex.get_token())

    def fun_call(self, name, token):
        self.cg.write_code_to_file("\tpush eax")
        self.lex.next()
        s = self.cg.form_label()
        e = self.cg.form_label()
        i = 0
        if self.lex.get_token() != ')':
            self.cg.write_code_to_file("\tjmp "+s)
            current = self.cg.form_label()
            self.cg.write_label(current)
            self.expr(self.lex.get_token())
            self.cg.write_code_to_file("\tpush eax")
            i += 1
            self.cg.write_code_to_file("\tjmp "+e)
            prev = current
            while self.lex.get_token() != ",":
                current = self.cg.form_label()
                self.cg.write_label(current)
                self.expr(self.lex.next())
                self.cg.write_code_to_file("\t jmp "+prev)
                prev = current
            self.cg.write_label(s)
            self.cg.write_code_to_file("\tjmp "+prev)
            self.cg.write_label(e)
        self.cg.write_code_to_file("\t call [esp+"+str(i*self.cg.stack)+"]")
        self.cg.write_code_to_file("\t add esp, "+ str((i+1)*self.cg.stack))
        self.check_token(self.lex.get_token(), {")"})
        self.lex.next()
        pass

    def branch(self, token):
        self.check_token(token, {"if"})
        self.check_token(self.lex.next(), {"("})
        self.expr(self.lex.next())
        e = self.cg.form_label()
        self.cg.write_code_to_file("\t cmp eax, 0")
        self.cg.write_code_to_file("\tjz "+e)
        self.check_token(self.lex.get_token(), {")"})
        i = self.cg.form_label()
        self.cg.write_label(i)
        self.statements(self.lex.next())
        self.cg.write_code_to_file("jmp "+e)
        self.cg.write_label(e)
        if self.lex.get_token() == "else":
            self.lex.next()
            self.statements(self.lex.get_token())
        end = self.cg.form_label()
        self.cg.write_label(end)


    def loop(self, token):
        self.check_token(token, {"while"})
        open = self.lex.next()
        self.check_token(open, {"("})
        e = self.cg.form_label()
        c = self.cg.form_label()
        self.cg.write_label(c)
        self.expr(self.lex.next())
        self.cg.write_code_to_file("\tcmp eax, 0")
        self.cg.write_code_to_file("\tjz "+e)
        self.check_token(self.lex.get_token(), ")")
        b = self.cg.form_label()
        e = self.cg.form_label()
        self.statements(self.lex.next())
        self.cg.write_code_to_file("\t jmp "+c)
        self.cg.write_label("\t jmp " + c)





