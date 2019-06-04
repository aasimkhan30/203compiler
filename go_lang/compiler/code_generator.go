package compiler

import (
	"io/ioutil"
	"log"
	"strconv"
)

type CodeGenerator struct {
	stack      int
	outputFile string
	label      int
}

func (c CodeGenerator) initialize(name string) {
	c.stack = 4
	c.outputFile = name
	c.label = -1
}

func (c CodeGenerator) write_code_to_file(str string) {
	err := ioutil.WriteFile(c.outputFile, []byte(str), 0666)
	if err != nil {
		log.Fatal(err)
	}
}

func (c CodeGenerator) instructions(operation string) string {
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
func (c CodeGenerator) formLabel() string {
	c.label++
	return "L" + strconv.Itoa(c.label)
}

func (c CodeGenerator) writeLabel(label string) {
	c.write_code_to_file(label + ":\n")
}

func (c CodeGenerator) functionBegin(name string, params int, label string) {
	c.write_code_to_file(".global " + name + "\n")
	c.write_code_to_file(name + ":\n")
	c.write_code_to_file("\t" + "push ebp\n")
	c.write_code_to_file("\t" + "mov ebp, esp\n")
	c.write_code_to_file("\t" + "sub esp, " + strconv.Itoa(params*4) + "\n")
	c.write_code_to_file("\t" + "jmp" + label + "\n")
}

func (c CodeGenerator) functionEnd() {
	c.write_code_to_file("\t" + "mov esp, ebp\n")
	c.write_code_to_file("\t" + "pop ebp\n")
	c.write_code_to_file("\t" + "ret\n")
}
