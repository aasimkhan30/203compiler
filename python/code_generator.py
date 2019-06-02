'''
    Instruction Reference From:
    https://en.wikipedia.org/wiki/X86_instruction_listings
'''
class code_generator:
    def __init__(self, name):
        self.stack = 4
        print("Calling Code Generator on file " + name)
        self.output_file = open(name + '.s', 'w')
        self.label = -1

    def write_code_to_file(self, str):
        self.output_file.write(str+"\n")

    def instructions(self, operation):
        if operation == '+':
            return 'add'
        if operation == '-':
            return 'sub'
        if operation == '++':
            return 'inc'
        if operation == '--':
            return 'dec'
        if operation == '*':
            return 'imul'
        if operation == '/':
            return 'idiv'
        if operation == '>':
            return 'setg'
        if operation == '<':
            return 'setl'
        if operation == '<=':
            return 'setle'
        if operation == '>=':
            return 'setge'
        if operation == '!=':
            return 'setne'
        if operation == '==':
            return 'sete'
        pass

    def form_label(self):
        self.label += 1
        return 'L'+str(self.label)

    def write_label(self, label):
        print("label ::::::::::::",label)
        self.write_code_to_file(label + ":")

    '''
        Referred This for function Prologue:
        https://en.wikipedia.org/wiki/Function_prologue

    '''

    def write_function_begin(self, name, params, label):
        self.write_code_to_file(".global " + name)
        self.write_code_to_file(name +":")
        self.write_code_to_file("\t" + "push ebp")
        self.write_code_to_file("\t" + "mov ebp, esp")
        self.write_code_to_file("\t" + "sub esp, " + str(params * 4))
        self.write_code_to_file("\t" + "jmp " + label)
    
    def write_function_end(self):
        self.write_code_to_file("\t" + "mov esp, ebp")
        self.write_code_to_file("\t" + "pop ebp")
        self.write_code_to_file("\t" + "ret")