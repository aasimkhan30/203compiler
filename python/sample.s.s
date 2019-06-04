.intel_syntax noprefix
.section .text
L0:
	mov eax, [ebp+8]
	push eax
	mov eax, [ebp+12]
	mov ebx, eax
	pop eax
	cmp eax, ebx
	setl al
	movzx eax, al
	cmp eax, 0
	jz L1
L2:
	mov eax, [ebp+12]
jmp L3
L1:
	mov eax, [ebp+8]
L3:
	mov esp, ebp
	pop ebp
	ret
.global maxInt
maxInt:
	push ebp
	mov ebp, esp
	sub esp, 0
	jmp L0
.section .text
L4:
	lea eax, [maxInt]
	push eax
	jmp L5
L7:
	mov eax, 3
	push eax
	jmp L6
L8:
	mov eax, 5
	push eax
	jmp L7
L5:
	jmp L8
L6:
	call [esp+8]
	add esp, 12
	mov [ebp-4], eax
	push eax
	jmp L9
L11:
	mov eax, offset L12
.section .data
L12: .asciz "%d"
.section .text
	push eax
	jmp L10
L13:
	mov eax, [ebp-4]
	push eax
	jmp L11
L9:
	jmp L13
L10:
	call [esp+8]
	add esp, 12
	mov esp, ebp
	pop ebp
	ret
.global main
main:
	push ebp
	mov ebp, esp
	sub esp, 4
	jmp L4
