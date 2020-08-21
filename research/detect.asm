BITS 64

SECTION .text
global _start

_start:
        push    rbp                                     ; 0000 _ 55
        mov     rbp, rsp                                ; 0001 _ 48: 89. E5
        mov     dword [rbp-24H], edi                    ; 0004 _ 89. 7D, DC
        mov     qword [rbp-30H], rsi                    ; 0007 _ 48: 89. 75, D0
        mov     qword [rbp-38H], rdx                    ; 000B _ 48: 89. 55, C8
        mov     rax, qword 4F4C4552505F444CH            ; 000F _ 48: B8, 4F4C4552505F444C
        mov     qword [rbp-13H], rax                    ; 0019 _ 48: 89. 45, ED
; Note: Length-changing prefix causes delay on Intel processors
        mov     word [rbp-0BH], 17473                   ; 001D _ 66: C7. 45, F5, 4441
        mov     byte [rbp-9H], 0                        ; 0023 _ C6. 45, F7, 00
        mov     dword [rbp-4H], 0                       ; 0027 _ C7. 45, FC, 00000000
        jmp     ?_007                                   ; 002E _ E9, 0000008D

?_001:  mov     dword [rbp-8H], 0                       ; 0033 _ C7. 45, F8, 00000000
        jmp     ?_003                                   ; 003A _ EB, 34

?_002:  mov     eax, dword [rbp-8H]                     ; 003C _ 8B. 45, F8
        cdqe                                            ; 003F _ 48: 98
        movzx   edx, byte [rbp+rax-13H]                 ; 0041 _ 0F B6. 54 05, ED
        mov     eax, dword [rbp-4H]                     ; 0046 _ 8B. 45, FC
        cdqe                                            ; 0049 _ 48: 98
        lea     rcx, [rax*8]                            ; 004B _ 48: 8D. 0C C5, 00000000
        mov     rax, qword [rbp-38H]                    ; 0053 _ 48: 8B. 45, C8
        add     rax, rcx                                ; 0057 _ 48: 01. C8
        mov     rcx, qword [rax]                        ; 005A _ 48: 8B. 08
        mov     eax, dword [rbp-8H]                     ; 005D _ 8B. 45, F8
        cdqe                                            ; 0060 _ 48: 98
        add     rax, rcx                                ; 0062 _ 48: 01. C8
        movzx   eax, byte [rax]                         ; 0065 _ 0F B6. 00
        cmp     dl, al                                  ; 0068 _ 38. C2
        jnz     ?_004                                   ; 006A _ 75, 3A
        add     dword [rbp-8H], 1                       ; 006C _ 83. 45, F8, 01
?_003:  mov     eax, dword [rbp-8H]                     ; 0070 _ 8B. 45, F8
        cdqe                                            ; 0073 _ 48: 98
        movzx   eax, byte [rbp+rax-13H]                 ; 0075 _ 0F B6. 44 05, ED
        test    al, al                                  ; 007A _ 84. C0
        jz      ?_005                                   ; 007C _ 74, 29
        mov     eax, dword [rbp-4H]                     ; 007E _ 8B. 45, FC
        cdqe                                            ; 0081 _ 48: 98
        lea     rdx, [rax*8]                            ; 0083 _ 48: 8D. 14 C5, 00000000
        mov     rax, qword [rbp-38H]                    ; 008B _ 48: 8B. 45, C8
        add     rax, rdx                                ; 008F _ 48: 01. D0
        mov     rdx, qword [rax]                        ; 0092 _ 48: 8B. 10
        mov     eax, dword [rbp-8H]                     ; 0095 _ 8B. 45, F8
        cdqe                                            ; 0098 _ 48: 98
        add     rax, rdx                                ; 009A _ 48: 01. D0
        movzx   eax, byte [rax]                         ; 009D _ 0F B6. 00
        test    al, al                                  ; 00A0 _ 84. C0
        jnz     ?_002                                   ; 00A2 _ 75, 98
        jmp     ?_005                                   ; 00A4 _ EB, 01

?_004:  nop                                             ; 00A6 _ 90
?_005:  mov     eax, dword [rbp-8H]                     ; 00A7 _ 8B. 45, F8
        cdqe                                            ; 00AA _ 48: 98
        movzx   eax, byte [rbp+rax-13H]                 ; 00AC _ 0F B6. 44 05, ED
        test    al, al                                  ; 00B1 _ 84. C0
        jnz     ?_006                                   ; 00B3 _ 75, 07
        mov     eax, 1                                  ; 00B5 _ B8, 00000001
        jmp     ?_008                                   ; 00BA _ EB, 29

?_006:  add     dword [rbp-4H], 1                       ; 00BC _ 83. 45, FC, 01
?_007:  mov     eax, dword [rbp-4H]                     ; 00C0 _ 8B. 45, FC
        cdqe                                            ; 00C3 _ 48: 98
        lea     rdx, [rax*8]                            ; 00C5 _ 48: 8D. 14 C5, 00000000
        mov     rax, qword [rbp-38H]                    ; 00CD _ 48: 8B. 45, C8
        add     rax, rdx                                ; 00D1 _ 48: 01. D0
        mov     rax, qword [rax]                        ; 00D4 _ 48: 8B. 00
        test    rax, rax                                ; 00D7 _ 48: 85. C0
        jne     ?_001                                   ; 00DA _ 0F 85, FFFFFF53
        mov     eax, 0                                  ; 00E0 _ B8, 00000000
?_008:  pop     rbp                                     ; 00E5 _ 5D
        ret                                             ; 00E6 _ C3
