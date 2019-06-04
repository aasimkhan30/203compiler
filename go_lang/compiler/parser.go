package compiler

import (
	"fmt"
	"strings"
)

type keyMap struct {
	name  string
	value bool
}

type Parser struct {
	cg     CodeGenerator
	lex    Lex
	global []keyMap
	local  []keyMap
}

func (p Parser) initialize(fileString string, outputFile string) {
	p.lex = Lex{}
	p.cg = CodeGenerator{}
	p.lex.initialize(fileString)
	p.cg.initialize(outputFile)
}

func (p Parser) localCheck(name string) bool {
	for _, element := range p.local {
		if element.name == name {
			return true
		}
	}
	return false
}

func (p Parser) globalCheck(name string) bool {
	for _, element := range p.global {
		if element.name == name {
			return true
		}
	}
	return false
}

func (p Parser) localParams() int {
	var c = 0
	for _, l := range p.local {
		if l.value == false {
			c++
		}
	}
	return c
}

func (p Parser) varOffset(name string) int {
	var offset = 0
	var i = 0
	var li = 0
	for _, element := range p.local {
		if element.name == name {
			if element.value == true {
				return (i + 2) * p.cg.stack
			} else {
				return -1 * (li + 1) * p.cg.stack
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

func (p Parser) varOffsetString(offset int) string {
	if offset < 0 {
		return string(offset)
	} else if offset > 0 {
		return "+" + string(offset)
	} else {
		return ""
	}
	return ""
}

func (p Parser) checkToken(token string, expected []string) bool {
	for _, element := range expected {
		if token == element {
			return true
		}
	}
	fmt.Println("Expected " + strings.Join(expected, ",") + " Found" + token)
	return false
}

func (p Parser) parseSource() {
	p.stdFunctions()
	p.program()
}

func (p Parser) stdFunctions() {
	p.global = append(p.global, keyMap{"printf", true})
}

func (p Parser) program() {
	p.cg.write_code_to_file(".intel_syntax noprefix")
	p.lex.next()
	for p.lex.getToken() != "" {
		var name = p.lex.next()
		var current = p.lex.next()
		if current == "(" {
			p.decf(name)
		} else if current == "=" {
			return
		}
	}
}

func (p Parser) decf(name string) {
	fmt.Println("Declare Function called on" + name)
	p.local = nil
	p.cg.write_code_to_file(".section .text")
	p.global = append(p.global, keyMap{name, true})
	for p.lex.getToken() != ")" {
		var datatype = p.lex.next()
		if datatype == ")" {
			break
		}
		var lName = p.lex.next()
		p.local = append(p.local, keyMap{lName, true})
		_ = p.lex.next()
	}
	p.checkToken(p.lex.getToken(), []string{")"})
	var tag = p.cg.formLabel()
	p.cg.writeLabel(tag)
	p.checkToken(p.lex.next(), []string{"{"})
	p.statements(p.lex.getToken())
	p.cg.functionEnd()
	p.cg.functionBegin(name, p.localParams(), tag)
}

func (p Parser) statements(token string) {
	fmt.Println("Statements")
	if token == "{" {
		p.lex.next()
		for p.lex.getToken() != "}" {
			p.statements(p.lex.getToken())
		}
		p.checkToken(p.lex.getToken(), []string{"}"})
		p.lex.next()
		fmt.Println("End")
	} else if token == "int" || token == "char" {
		var name = p.lex.next()
		p.local = append(p.local, keyMap{name, false})
		var equal = p.lex.next()
		if equal == "=" {
			p.expr(p.lex.next())
			var offset = p.varOffset(name)
			p.cg.write_code_to_file("\tmov [ebp" + p.varOffsetString(offset) + "], eax")
			p.checkToken(p.lex.getToken(), []string{";"})
			p.lex.next()
		}
	} else if token == "if" {
		p.branch(token)
	} else if token == "while" {
		p.loop(token)
	} else if token == "return" {
		p.expr(p.lex.next())
		p.checkToken(p.lex.getToken(), []string{";"})
		p.lex.next()
	} else {
		p.expr(token)
		p.checkToken(p.lex.getToken(), []string{";"})
		p.lex.next()
	}
}

func (p Parser) expr(token string) {
	fmt.Println("Expr")
	p.expr(token)
}

func (p Parser) expr1(token string) {
	fmt.Println("Expr 1")
	p.expr2(token)
	if p.lex.getToken() == "=" {
		var name = token
		var expr = p.lex.next()
		p.expr2(expr)
		var offset = p.varOffset(name)
		p.cg.write_code_to_file("\tmov [ebp" + p.varOffsetString(offset) + "], eax")
	}
}

func (p Parser) expr2(token string) {
	fmt.Println("Expr 2")
	p.expr3(token)
	var operator = p.lex.getToken()
	if operator == "||" || operator == "&&" {
		var label = p.cg.formLabel()
		p.cg.write_code_to_file("\t cmp eax, 0")
		if operator == "||" {
			p.cg.write_code_to_file("\tjnz" + label)
		} else {
			p.cg.write_code_to_file("\tjz" + label)
		}
		p.expr3(p.lex.next())
		p.cg.writeLabel(label)
	}
}

func (p Parser) expr3(token string) {
	fmt.Println("Expr 3")
	p.expr4(token)
	var operator = p.lex.getToken()
	if operator == "<=" || operator == ">=" || operator == "<" || operator == ">" || operator == "!=" || operator == "==" {
		p.cg.write_code_to_file("\tpush eax")
		p.expr4(p.lex.next())
		p.cg.write_code_to_file("\tmov ebx, eax")
		p.cg.write_code_to_file("\tpop eax")
		var instruction = p.cg.instructions(operator)
		p.cg.write_code_to_file("\tcmp eax, ebx")
		p.cg.write_code_to_file("\t" + instruction + "al")
		p.cg.write_code_to_file("\tmovzx eax, al")
	}
}

func (p Parser) expr4(token string) {
	fmt.Println("Expr 4")
	p.expr5(token)
	var operator = p.lex.getToken()
	for operator == "+" || operator == "-" {
		p.cg.write_code_to_file("\tpush eax")
		p.expr5(p.lex.next())
		p.cg.write_code_to_file("\tmov ebx, eax")
		p.cg.write_code_to_file("\tpop eax")
		var instruction = p.cg.instructions(operator)
		p.cg.write_code_to_file("\t" + instruction + " eax, ebx")
		operator = p.lex.getToken()
	}
}

func (p Parser) expr5(token string) {
	fmt.Println("Expr 5")
	p.expr6(token)
	var operator = p.lex.getToken()
	for operator == "*" || operator == "/" {
		p.cg.write_code_to_file("\tpush eax")
		p.expr6(p.lex.next())
		p.cg.write_code_to_file("\tmov ebx, eax")
		p.cg.write_code_to_file("\tpop eax")
		var instructions = p.cg.instructions(operator)
		p.cg.write_code_to_file("\t" + instructions + " ebx")
		operator = p.lex.getToken()
	}
}

func (p Parser) expr6(token string) {
	fmt.Println("Expr 6")
	if token == "-" {
		p.expr6(p.lex.next())
		p.cg.write_code_to_file("\tneg eax")
	} else {
		p.unary(token)
		var operator = p.lex.getToken()
		if operator == "++" || operator == "--" {
			var instruction = p.cg.instructions(operator)
			p.cg.write_code_to_file("\t" + instruction + " eax")
			var offset = p.varOffset(token)
			p.cg.write_code_to_file("\tmov [ebp" + p.varOffsetString(offset) + "], eax")
			p.lex.next()
		}
	}

}

func (p Parser) unary(token string) {
	fmt.Println("Unary Called")
	var name = token
	var tokenType = p.lex.getType()
	switch tokenType {
	case 2:
		p.cg.write_code_to_file("\tmov eax, " + p.lex.getToken())
	case 4:
		var label = p.cg.formLabel()
		p.cg.write_code_to_file("\tmov eax, offset " + label)
		p.cg.write_code_to_file(".section .data")
		p.cg.write_code_to_file(label + ": .asciz " + p.lex.getToken())
		p.cg.write_code_to_file(".section .text")
	case 1:
		if p.localCheck(name) {
			var offset = p.varOffset(name)
			p.cg.write_code_to_file("\tmov eax, [ebp" + p.varOffsetString(offset) + "]")
		} else if p.globalCheck(name) {
			p.cg.write_code_to_file("\tlea eax, [" + name + "]")
		}
		var current = p.lex.next()
		if current == "(" {
			p.funCall(name, p.lex.getToken())
		}

	}

}

func (p Parser) funCall(name string, token string) {
	fmt.Println("Function Called")
	p.cg.writeLabel("\tpush eax")
	p.lex.next()
	var s = p.cg.formLabel()
	fmt.Println("Start Label", s)
	var e = p.cg.formLabel()
	fmt.Println("End label", e)
	var i = 0
	if p.lex.getToken() != ")" {
		p.cg.write_code_to_file("\tjmp" + s)
		var current = p.cg.formLabel()
		p.cg.writeLabel(current)
		p.expr(p.lex.getToken())
		p.cg.write_code_to_file("\tpush eax")
		i++
		p.cg.write_code_to_file("\tjmp " + e)
		var prev = current
		for p.lex.getToken() == "," {
			current = p.cg.formLabel()
			p.cg.writeLabel(current)
			p.expr(p.lex.next())
			p.cg.write_code_to_file("\tpush eax")
			i++
			p.cg.write_code_to_file("\tjmp " + prev)
			prev = current
		}
		fmt.Println("Start label" + s)
		p.cg.writeLabel(s)
		p.cg.write_code_to_file("\tjmp " + prev)
		p.cg.writeLabel(e)
	}
	p.cg.write_code_to_file("\tcall [esp+" + string(i*p.cg.stack) + "]")
	p.cg.write_code_to_file("\tadd esp, " + string((i+1)*p.cg.stack))
	p.checkToken(p.lex.getToken(), []string{")"})
	p.lex.next()
}

func (p Parser) branch(token string) {
	fmt.Println("If Called")
	p.checkToken(token, []string{"if"})
	p.checkToken(p.lex.next(), []string{"("})
	p.expr(p.lex.next())
	var e = p.cg.formLabel()
	p.cg.write_code_to_file("\tcmp eax, 0")
	p.cg.write_code_to_file("\tjz " + e)
	p.checkToken(p.lex.getToken(), []string{")"})
	var i = p.cg.formLabel()
	p.cg.writeLabel(i)
	p.statements(p.lex.next())
	var end = p.cg.formLabel()
	p.cg.write_code_to_file("jmp " + end)
	p.cg.writeLabel(e)
	if p.lex.getToken() == "else" {
		p.lex.next()
		p.statements(p.lex.getToken())
	}
	p.cg.writeLabel(end)
}

func (p Parser) loop(token string) {
	fmt.Println("While Called")
	p.checkToken(token, []string{"while"})
	var open = p.lex.next()
	p.checkToken(open, []string{"("})
	var e = p.cg.formLabel()
	var c = p.cg.formLabel()
	p.cg.writeLabel(c)
	p.expr(p.lex.next())
	p.cg.write_code_to_file("\tcmp eax, 0")
	p.cg.write_code_to_file("\tjz " + e)
	p.checkToken(p.lex.getToken(), []string{")"})
	var b = p.cg.formLabel()
	p.cg.writeLabel(b)
	p.statements(p.lex.next())
	p.cg.write_code_to_file("\tjmp " + c)
	p.cg.writeLabel(e)
}
