package compiler

import "unicode"

type Lex struct {
	ptr                                  int
	fileString                           string
	tokenType                            int
	token                                string
	EOF, IDN, INT, CHAR, STR, OP, OTHERS int
}

func (l Lex) initialize(fileString string) {
	l.ptr = 0
	l.token = ""
	l.tokenType = 6
	l.fileString = fileString
	l.EOF = 0
	l.IDN = 2
	l.CHAR = 3
	l.STR = 4
	l.OP = 5
	l.OTHERS = 6
}

func (l Lex) getToken() string {
	return l.token
}

func (l Lex) getType() int {
	return l.tokenType
}

func (l Lex) isAl(char byte) bool {
	if char >= 'a' && char <= 'z' || char >= 'A' && char <= 'Z' {
		return true
	}
	return false
}

func (l Lex) isNum(char byte) bool {
	if char >= '0' && char <= '9' {
		return true
	}
	return false
}

func (l Lex) isAlNum(char byte) bool {
	return l.isAl(char) || l.isNum(char)
}

func (l Lex) next() string {
	l.token = ""
	for l.ptr < len(l.fileString) && unicode.IsSpace(rune(l.fileString[l.ptr])) {
		l.ptr++
	}

	if l.ptr == len(l.fileString) {
		l.token = ""
		l.tokenType = l.EOF
		return ""
	}

	if l.isAl(l.fileString[l.ptr]) {
		l.token += string(l.fileString[l.ptr])
		l.ptr++
		for l.ptr < len(l.fileString) && l.isAlNum(l.fileString[l.ptr]) || l.fileString[l.ptr] == '_' {
			l.token += string(l.fileString[l.ptr])
			l.ptr++
		}
		l.tokenType = l.IDN
		return l.token
	}

	if l.isNum(l.fileString[l.ptr]) {
		l.token += string(l.fileString[l.ptr])
		l.ptr++
		for l.ptr < len(l.fileString) && l.isNum(l.fileString[l.ptr]) {
			l.token += string(l.fileString[l.ptr])
			l.ptr++
		}
		l.tokenType = l.INT
		return l.token
	}

	if l.fileString[l.ptr] == '"' || l.fileString[l.ptr] == '\'' {
		var terminator = l.fileString[l.ptr]
		l.token += string(l.fileString[l.ptr])
		l.ptr++
		for l.ptr < len(l.fileString) && l.fileString[l.ptr] != terminator {
			l.token += string(l.fileString[l.ptr])
			l.ptr++
		}
		l.token += string(l.fileString[l.ptr])
		l.ptr++
		if terminator == '"' {
			l.tokenType = l.STR
		}
		if terminator == '\'' {
			l.tokenType = l.CHAR
		}
		return l.token
	}
	var temp = l.fileString[l.ptr]
	if temp == '+' || temp == '-' || temp == '*' || temp == '/' || temp == '<' || temp == '>' || temp == '=' || temp == '&' || temp == '!' {
		l.token += string(temp)
		var terminator = temp
		l.ptr++
		if l.ptr < len(l.fileString) && l.fileString[l.ptr] == terminator || l.fileString[l.ptr] == '=' {
			l.token += string(l.fileString[l.ptr])
			l.ptr++
		}
		l.tokenType = l.OP
		return l.token
	}
	l.ptr++
	l.token += string(temp)
	l.tokenType = l.OTHERS
	return l.token
}
