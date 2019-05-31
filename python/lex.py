#Defining some tokens

EOF = 0
IDN = 1
INT = 2
CHAR = 3
STR = 4
OP = 5
OTHERS = 6

class lexer:
    def __init__(self, file_str):
        self.ptr = 0
        self.token = []
        self.token.type = None
        self.file_str = file_str
    def get_token(self):
        return self.token
    def get_type(self):
        return self.token.type
    def next(self):
        self.token.clear()
        # Detecting the end of file
        if self.ptr == len(self.file_str):
            self.token = None
            self.token_type = EOF
            return self.token
        # Detecting whitespaces and just moving ahead
        while self.file_str[self.ptr].isspace():
            self.ptr += 1
        # Detecting variables
        if self.file_str[self.ptr].isalpha():
            self.token.append(self.file_str[self.ptr])
            self.ptr += 1
            while self.ptr < len(self.file_str) and self.file_str[self.ptr].isalnum() or self.file_str[self.ptr] == '_':
                self.token.append(self.file_str[self.ptr])
                self.ptr += 1
            print(''.join(self.token))
            self.token_type = IDN
            return ''.join(self.token)
        # Detecting integers
        if self.file_str[self.ptr].isnumeric():
            self.token.append(self.file_str[self.ptr])
            self.ptr += 1
            while self.ptr < len(self.file_str) and self.file_str[self.ptr].isnumeric():
                self.token.append(self.file_str[self.ptr])
                self.ptr += 1
            print(''.join(self.token))
            self.token_type = INT
            return ''.join(self.token)
        # Detecting Strings
        if self.file_str[self.ptr] == '\"' or  self.file_str[self.ptr] == '\'':
            # storing the start quotes to find the end quotes too
            terminator = self.file_str[self.ptr]
            self.token.append(self.file_str[self.ptr])
            self.ptr += 1
            while self.ptr < len(self.file_str) and self.file_str[self.ptr] != terminator:
                self.token.append(self.file_str[self.ptr])
                self.ptr += 1
            self.token.append(self.file_str[self.ptr])
            self.ptr += 1
            print(''.join(self.token))
            if terminator == "\"":
                self.token_type = STR
            else:
                self.token_type = CHAR
            return ''.join(self.token)
        # Detecting operators
        temp = self.file_str[self.ptr]
        if temp == '+' or temp == '-' or temp == '*' or temp == '/' or temp == '<' or temp == '>' or temp == '=' or temp == '&' or temp == '!':
            self.token.append(temp)
            terminator = temp
            self.ptr += 1
            if self.ptr < len(self.file_str) and self.file_str[self.ptr] == terminator or self.file_str[self.ptr] == "=":
                self.token.append(self.file_str[self.ptr])
                self.ptr += 1
            print(''.join(self.token))
            self.token_type = OP
            return ''.join(self.token)
        # Detecting Not Operator
        if temp == '!':
            self.token.append(temp)
            self.ptr += 1
            print(''.join(self.token))
            self.token_type = OP
            return ''.join(self.token)
        self.ptr += 1
        print(temp)
        self.token_type = OTHERS
        return temp



