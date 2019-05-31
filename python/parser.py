class parser:
    def __init__(self, lex, cg):
        self.lex = lex
        self.cg = cg
    def parse_source(self):
        self.std_fns()
        self.cg.write_code_to_file(".intel_syntax noprefix\n")
        # Should hold an datatype
        token =  self.lex.next()
        while token != None:
            name = self.lex.next() # Should hold a function or global variable name
            operator = self.lex.next() # Should hold = or (. If = then global var else function
            if operator == "(":
                token = self.lex.next();
            if operator == "=":
                token = self.lex.next();
    def std_fun(self):
        pass

