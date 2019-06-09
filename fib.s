.intel_syntax noprefix
.section .text
L0:
	mov eax, 0
	mov [ebp-4], eax
	mov eax, 0
	mov [ebp-8], eax
	mov eax, 1
	mov [ebp-12], eax
	mov eax, 1
	mov [ebp-16], eax
L2:
	mov eax, [ebp-16]
	push eax
	mov eax, [ebp+8]
	mov ebx, eax
	pop eax
	cmp eax, ebx
	setl al
	movzx eax, al
	cmp eax, 0
	jz L1
L3:
	mov eax, [ebp-4]
	mov eax, [ebp-8]
	mov [ebp-4], eax
	mov eax, [ebp-8]
	mov eax, [ebp-12]
	mov [ebp-8], eax
	mov eax, [ebp-12]
	mov eax, [ebp-4]
	push eax
	mov eax, [ebp-8]
	mov ebx, eax
	pop eax
	add eax, ebx
	mov [ebp-12], eax
	mov eax, [ebp-16]
	inc eax
	mov [ebp-16], eax
	jmp L2
L1:
	mov eax, [ebp-12]
	mov esp, ebp
	pop ebp
	ret
.global fib
fib:
	push ebp
	mov ebp, esp
	sub esp, 16
	jmp L0
.section .text
L4:
	lea eax, [fib]
	push eax
	jmp L5
L7:
	mov eax, 10
	push eax
	jmp L6
L5:
	jmp L7
L6:
	call [esp+4]
	add esp, 8
	mov [ebp-4], eax
	mov esp, ebp
	pop ebp
	ret
.global main
main:
	push ebp
	mov ebp, esp
	sub esp, 4
	jmp L4
