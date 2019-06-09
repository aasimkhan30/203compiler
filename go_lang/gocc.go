package main

import (
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"strconv"
	"strings"
	"unicode"
)

var inputFile string
var fileString string
var stack int
var outputFile string
var label int
var fileVar *os.File
var ptr int
var tokenType int
var token string
var EOF, IDN, INT, CHAR, STR, OP, OTHERS int
var local []keyMap
var global []keyMap
var err error

func main() {
	//MAIN VARS
	var fileName = os.Args[1]
	fmt.Println("Compiling " + fileName)
	var b, _ = ioutil.ReadFile(fileName)
	fileString = string(b)
	//CODEGENERATOR VARS
	stack = 4
	outputFile = strings.Split(fileName, ".")[0]
	label = -1
	fileVar, err = os.Create(outputFile + ".s")
	if err != nil {
		//////////////fmt.println(fileVar, err)
	}
	ptr = 0
	tokenType = 6
	token = ""
	EOF = 0
	IDN = 1
	INT = 2
	CHAR = 3
	STR = 4
	OP = 5
	OTHERS = 6
	local = []keyMap{}
	global = []keyMap{}
	parseSource()

}

func write_code_to_file(str string) {
	n, err := io.WriteString(fileVar, str+"\n")
	if err != nil {
		fmt.Println(n, err)
	}
}

func instructions(operation string) string {
	switch operation {
	case "+":
		return "add"
	case "-":
		return "sub"
	case "*":
		return "imul"
	case "/":
		return "idiv"
	case ">":
		return "setg"
	case "<":
		return "setl"
	case ">=":
		return "setge"
	case "<=":
		return "setle"
	case "!=":
		return "setne"
	case "==":
		return "sete"
	case "++":
		return "inc"
	case "--":
		return "dec"
	}
	return ""
}
func formLabel() string {
	label++
	return "L" + strconv.Itoa(label)
}

func writeLabel(label string) {
	write_code_to_file(label + ":")
}

func functionBegin(name string, params int, label string) {
	write_code_to_file(".global " + name)
	write_code_to_file(name + ":")
	write_code_to_file("\t" + "push ebp")
	write_code_to_file("\t" + "mov ebp, esp")
	write_code_to_file("\t" + "sub esp, " + strconv.Itoa(params*4))
	write_code_to_file("\t" + "jmp " + label)
}

func functionEnd() {
	write_code_to_file("\t" + "mov esp, ebp")
	write_code_to_file("\t" + "pop ebp")
	write_code_to_file("\t" + "ret")
}

func getFile() {
}
func getToken() string {
	return token
}

func getType() int {
	return tokenType
}

func isAl(char byte) bool {
	if char >= 'a' && char <= 'z' || char >= 'A' && char <= 'Z' {
		return true
	}
	return false
}

func isNum(char byte) bool {
	if char >= '0' && char <= '9' {
		return true
	}
	return false
}

func isAlNum(char byte) bool {
	return isAl(char) || isNum(char)
}

func next() string {
	token = ""
	for ptr < len(fileString) && unicode.IsSpace(rune(fileString[ptr])) {
		ptr++
	}
	if ptr == len(fileString) {
		token = ""
		tokenType = EOF
		return ""
	}

	if isAl(fileString[ptr]) {
		token += string(fileString[ptr])
		ptr++
		for ptr < len(fileString) && isAlNum(fileString[ptr]) || fileString[ptr] == '_' {
			token += string(fileString[ptr])
			ptr++
		}
		tokenType = IDN
		return token
	}

	if isNum(fileString[ptr]) {
		token += string(fileString[ptr])
		ptr++
		for ptr < len(fileString) && isNum(fileString[ptr]) {
			token += string(fileString[ptr])
			ptr++
		}
		tokenType = INT
		return token
	}

	if fileString[ptr] == '"' || fileString[ptr] == '\'' {
		var terminator = fileString[ptr]
		token += string(fileString[ptr])
		ptr++
		for ptr < len(fileString) && fileString[ptr] != terminator {
			token += string(fileString[ptr])
			ptr++
		}
		token += string(fileString[ptr])
		ptr++
		if terminator == '"' {
			tokenType = STR
		}
		if terminator == '\'' {
			tokenType = CHAR
		}

		return token
	}
	var temp = fileString[ptr]
	if temp == '+' || temp == '-' || temp == '*' || temp == '/' || temp == '<' || temp == '>' || temp == '=' || temp == '&' || temp == '!' {
		token += string(temp)
		var terminator = temp
		ptr++
		if ptr < len(fileString) && fileString[ptr] == terminator || fileString[ptr] == '=' {
			token += string(fileString[ptr])
			ptr++
		}
		tokenType = OP

		return token
	}
	ptr++
	token += string(temp)
	tokenType = OTHERS
	//////fmt.println(tokenType)
	//////fmt.println(token)
	return token
}

type keyMap struct {
	name  string
	value bool
}

func localCheck(name string) bool {
	for _, element := range local {
		if element.name == name {
			return true
		}
	}
	return false
}

func globalCheck(name string) bool {
	for _, element := range global {
		if element.name == name {
			return true
		}
	}
	return false
}

func localParams() int {
	var c = 0
	for _, l := range local {
		if l.value == false {
			c++
		}
	}
	return c
}

func varOffset(name string) int {
	var offset = 0
	var i = 0
	var li = 0
	for _, element := range local {
		if element.name == name {
			if element.value == true {
				return (i + 2) * stack
			} else {
				return -1 * (li + 1) * stack
			}
		}
		if element.value == true {
			i++
		} else {
			li++
		}
	}
	return offset
}

func varOffsetString(offset int) string {
	if offset < 0 {
		return strconv.Itoa(offset)
	} else if offset > 0 {
		return "+" + strconv.Itoa(offset)
	} else {
		return ""
	}
	return ""
}

func checkToken(token string, expected []string) bool {
	for _, element := range expected {
		if token == element {
			return true
		}
	}
	return false
}

func parseSource() {
	stdFunctions()
	program()
	fileVar.Close()
}

func stdFunctions() {
	global = append(global, keyMap{"printf", true})
}

func program() {
	write_code_to_file(".intel_syntax noprefix")
	next()
	for getToken() != "" {
		var name = next()
		var current = next()
		if current == "(" {
			decf(name)
		} else if current == "=" {
			return
		}
	}
}

func decf(name string) {
	local = nil
	write_code_to_file(".section .text")
	global = append(global, keyMap{name, true})
	for getToken() != ")" {
		var datatype = next()
		if datatype == ")" {
			break
		}
		var lName = next()
		local = append(local, keyMap{lName, true})
		next()
	}
	checkToken(getToken(), []string{")"})
	var tag = formLabel()
	writeLabel(tag)
	checkToken(next(), []string{"{"})
	statements(getToken())
	functionEnd()
	functionBegin(name, localParams(), tag)
}

func statements(token string) {
	if token == "{" {
		next()
		for getToken() != "}" {
			statements(getToken())
		}
		checkToken(getToken(), []string{"}"})
		next()
	} else if token == "int" || token == "char" {
		var name = next()
		local = append(local, keyMap{name, false})
		var equal = next()
		if equal == "=" {
			expr(next())
			var offset = varOffset(name)
			write_code_to_file("\tmov [ebp" + varOffsetString(offset) + "], eax")
		}
		checkToken(getToken(), []string{";"})
		next()
	} else if token == "if" {
		branch(token)
	} else if token == "while" {
		loop(token)
	} else if token == "return" {
		expr(next())
		checkToken(getToken(), []string{";"})
		next()
	} else {
		expr(token)
		checkToken(getToken(), []string{";"})
		next()
	}
}

func expr(token string) {
	expr1(token)
}

func expr1(token string) {
	expr2(token)
	if getToken() == "=" {
		var name = token
		var expr = next()
		expr2(expr)
		var offset = varOffset(name)
		write_code_to_file("\tmov [ebp" + varOffsetString(offset) + "], eax")
	}
}

func expr2(token string) {
	expr3(token)
	var operator = getToken()
	if operator == "||" || operator == "&&" {
		var label = formLabel()
		write_code_to_file("\tcmp eax, 0")
		if operator == "||" {
			write_code_to_file("\tjnz" + label)
		} else {
			write_code_to_file("\tjz" + label)
		}
		expr3(next())
		writeLabel(label)
	}
}

func expr3(token string) {
	expr4(token)
	var operator = getToken()
	if operator == "<=" || operator == ">=" || operator == "<" || operator == ">" || operator == "!=" || operator == "==" {
		write_code_to_file("\tpush eax")
		expr4(next())
		write_code_to_file("\tmov ebx, eax")
		write_code_to_file("\tpop eax")
		var instruction = instructions(operator)
		write_code_to_file("\tcmp eax, ebx")
		write_code_to_file("\t" + instruction + " al")
		write_code_to_file("\tmovzx eax, al")
	}
}

func expr4(token string) {
	expr5(token)
	var operator = getToken()
	for operator == "+" || operator == "-" {
		write_code_to_file("\tpush eax")
		expr5(next())
		write_code_to_file("\tmov ebx, eax")
		write_code_to_file("\tpop eax")
		var instruction = instructions(operator)
		write_code_to_file("\t" + instruction + " eax, ebx")
		operator = getToken()
	}
}

func expr5(token string) {
	expr6(token)
	var operator = getToken()
	for operator == "*" || operator == "/" {
		write_code_to_file("\tpush eax")
		expr6(next())
		write_code_to_file("\tmov ebx, eax")
		write_code_to_file("\tpop eax")
		var instructions = instructions(operator)
		write_code_to_file("\t" + instructions + " ebx")
		operator = getToken()
	}
}

func expr6(token string) {
	if token == "-" {
		expr6(next())
		write_code_to_file("\tneg eax")
	} else {
		unary(token)
		var operator = getToken()
		if operator == "++" || operator == "--" {
			var instruction = instructions(operator)
			write_code_to_file("\t" + instruction + " eax")
			var offset = varOffset(token)
			write_code_to_file("\tmov [ebp" + varOffsetString(offset) + "], eax")
			next()
		}
	}

}

func unary(token string) {
	var name = token
	var unaryTokenType = getType()
	switch unaryTokenType {
	case 2:
		write_code_to_file("\tmov eax, " + getToken())
	case 4:
		var label = formLabel()
		write_code_to_file("\tmov eax, offset " + label)
		write_code_to_file(".section .data")
		write_code_to_file(label + ": .asciz " + getToken())
		write_code_to_file(".section .text")
	case 1:
		if localCheck(name) {
			var offset = varOffset(name)
			write_code_to_file("\tmov eax, [ebp" + varOffsetString(offset) + "]")
		} else {
			if globalCheck(name) {
				write_code_to_file("\tlea eax, [" + name + "]")
			}
		}
	}
	var current = next()
	if current == "(" {
		funCall(name, getToken())
	}

}

func funCall(name string, token string) {
	write_code_to_file("\tpush eax")
	next()
	var s = formLabel()
	var e = formLabel()
	var i = 0
	if getToken() != ")" {
		write_code_to_file("\tjmp " + s)
		var current = formLabel()
		writeLabel(current)
		expr(getToken())
		write_code_to_file("\tpush eax")
		i++
		write_code_to_file("\tjmp " + e)
		var prev = current
		for getToken() == "," {
			current = formLabel()
			writeLabel(current)
			expr(next())
			write_code_to_file("\tpush eax")
			i++
			write_code_to_file("\tjmp " + prev)
			prev = current
		}
		writeLabel(s)
		write_code_to_file("\tjmp " + prev)
		writeLabel(e)
	}
	write_code_to_file("\tcall [esp+" + strconv.Itoa(i*stack) + "]")
	write_code_to_file("\tadd esp, " + strconv.Itoa((i+1)*stack))
	checkToken(getToken(), []string{")"})
	next()
}

func branch(token string) {
	checkToken(token, []string{"if"})
	checkToken(next(), []string{"("})
	expr(next())
	var e = formLabel()
	write_code_to_file("\tcmp eax, 0")
	write_code_to_file("\tjz " + e)
	checkToken(getToken(), []string{")"})
	var i = formLabel()
	writeLabel(i)
	statements(next())
	var end = formLabel()
	write_code_to_file("\tjmp " + end)
	writeLabel(e)
	if getToken() == "else" {
		next()
		statements(getToken())
	}
	writeLabel(end)
}

func loop(token string) {
	checkToken(token, []string{"while"})
	var open = next()
	checkToken(open, []string{"("})
	var e = formLabel()
	var c = formLabel()
	writeLabel(c)
	expr(next())
	write_code_to_file("\tcmp eax, 0")
	write_code_to_file("\tjz " + e)
	checkToken(getToken(), []string{")"})
	var b = formLabel()
	writeLabel(b)
	statements(next())
	write_code_to_file("\tjmp " + c)
	writeLabel(e)
}
