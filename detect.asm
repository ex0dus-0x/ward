BITS 64


extern __stack_chk_fail                                 ; near
extern puts                                             ; near
extern _GLOBAL_OFFSET_TABLE_                            ; byte
extern environ                                          ; qword


SECTION .text
global main

main:   ; Function begin
        push    rbp                                     ; 0000 _ 55
        mov     rbp, rsp                                ; 0001 _ 48: 89. E5
        sub     rsp, 48                                 ; 0004 _ 48: 83. EC, 30
; Note: Address is not rip-relative
; Note: Absolute memory address without relocation
        mov     rax, qword [fs:abs 28H]                 ; 0008 _ 64 48: 8B. 04 25, 00000028
        mov     qword [rbp-8H], rax                     ; 0011 _ 48: 89. 45, F8
        xor     eax, eax                                ; 0015 _ 31. C0
        mov     rax, qword 4F4C4552505F444CH            ; 0017 _ 48: B8, 4F4C4552505F444C
        mov     qword [rbp-13H], rax                    ; 0021 _ 48: 89. 45, ED
; Note: Length-changing prefix causes delay on Intel processors
        mov     word [rbp-0BH], 17473                   ; 0025 _ 66: C7. 45, F5, 4441
        mov     byte [rbp-9H], 0                        ; 002B _ C6. 45, F7, 00
        mov     qword [rbp-28H], 0                      ; 002F _ 48: C7. 45, D8, 00000000
        jmp     ?_007                                   ; 0037 _ E9, 000000A2

?_001:  mov     qword [rbp-20H], 0                      ; 003C _ 48: C7. 45, E0, 00000000
        jmp     ?_003                                   ; 0044 _ EB, 36

?_002:  lea     rdx, [rbp-13H]                          ; 0046 _ 48: 8D. 55, ED
        mov     rax, qword [rbp-20H]                    ; 004A _ 48: 8B. 45, E0
        add     rax, rdx                                ; 004E _ 48: 01. D0
        movzx   edx, byte [rax]                         ; 0051 _ 0F B6. 10
        mov     rcx, qword [rel environ]                ; 0054 _ 48: 8B. 0D, 00000000(rel)
        mov     rax, qword [rbp-28H]                    ; 005B _ 48: 8B. 45, D8
        shl     rax, 3                                  ; 005F _ 48: C1. E0, 03
        add     rax, rcx                                ; 0063 _ 48: 01. C8
        mov     rcx, qword [rax]                        ; 0066 _ 48: 8B. 08
        mov     rax, qword [rbp-20H]                    ; 0069 _ 48: 8B. 45, E0
        add     rax, rcx                                ; 006D _ 48: 01. C8
        movzx   eax, byte [rax]                         ; 0070 _ 0F B6. 00
        cmp     dl, al                                  ; 0073 _ 38. C2
        jnz     ?_004                                   ; 0075 _ 75, 3C
        add     qword [rbp-20H], 1                      ; 0077 _ 48: 83. 45, E0, 01
?_003:  lea     rdx, [rbp-13H]                          ; 007C _ 48: 8D. 55, ED
        mov     rax, qword [rbp-20H]                    ; 0080 _ 48: 8B. 45, E0
        add     rax, rdx                                ; 0084 _ 48: 01. D0
        movzx   eax, byte [rax]                         ; 0087 _ 0F B6. 00
        test    al, al                                  ; 008A _ 84. C0
        jz      ?_005                                   ; 008C _ 74, 26
        mov     rdx, qword [rel environ]                ; 008E _ 48: 8B. 15, 00000000(rel)
        mov     rax, qword [rbp-28H]                    ; 0095 _ 48: 8B. 45, D8
        shl     rax, 3                                  ; 0099 _ 48: C1. E0, 03
        add     rax, rdx                                ; 009D _ 48: 01. D0
        mov     rdx, qword [rax]                        ; 00A0 _ 48: 8B. 10
        mov     rax, qword [rbp-20H]                    ; 00A3 _ 48: 8B. 45, E0
        add     rax, rdx                                ; 00A7 _ 48: 01. D0
        movzx   eax, byte [rax]                         ; 00AA _ 0F B6. 00
        test    al, al                                  ; 00AD _ 84. C0
        jnz     ?_002                                   ; 00AF _ 75, 95
        jmp     ?_005                                   ; 00B1 _ EB, 01

?_004:  nop                                             ; 00B3 _ 90
?_005:  lea     rdx, [rbp-13H]                          ; 00B4 _ 48: 8D. 55, ED
        mov     rax, qword [rbp-20H]                    ; 00B8 _ 48: 8B. 45, E0
        add     rax, rdx                                ; 00BC _ 48: 01. D0
        movzx   eax, byte [rax]                         ; 00BF _ 0F B6. 00
        test    al, al                                  ; 00C2 _ 84. C0
        jnz     ?_006                                   ; 00C4 _ 75, 13
        lea     rdi, [rel ?_010]                        ; 00C6 _ 48: 8D. 3D, 00000000(rel)
        call    puts                                    ; 00CD _ E8, 00000000(PLT r)
        mov     eax, 1                                  ; 00D2 _ B8, 00000001
        jmp     ?_008                                   ; 00D7 _ EB, 34

?_006:  add     qword [rbp-28H], 1                      ; 00D9 _ 48: 83. 45, D8, 01
?_007:  mov     rdx, qword [rel environ]                ; 00DE _ 48: 8B. 15, 00000000(rel)
        mov     rax, qword [rbp-28H]                    ; 00E5 _ 48: 8B. 45, D8
        shl     rax, 3                                  ; 00E9 _ 48: C1. E0, 03
        add     rax, rdx                                ; 00ED _ 48: 01. D0
        mov     rax, qword [rax]                        ; 00F0 _ 48: 8B. 00
        test    rax, rax                                ; 00F3 _ 48: 85. C0
        jne     ?_001                                   ; 00F6 _ 0F 85, FFFFFF40
        lea     rdi, [rel ?_011]                        ; 00FC _ 48: 8D. 3D, 00000000(rel)
        call    puts                                    ; 0103 _ E8, 00000000(PLT r)
        mov     eax, 0                                  ; 0108 _ B8, 00000000
?_008:  mov     rsi, qword [rbp-8H]                     ; 010D _ 48: 8B. 75, F8
; Note: Address is not rip-relative
; Note: Absolute memory address without relocation
        sub     rsi, qword [fs:abs 28H]                 ; 0111 _ 64 48: 2B. 34 25, 00000028
        jz      ?_009                                   ; 011A _ 74, 05
        call    __stack_chk_fail                        ; 011C _ E8, 00000000(PLT r)
?_009:  leave
        ret                                             ; 0122 _ C3
; main End of function


SECTION .data
SECTION .bss
SECTION .rodata

?_010:                                                  ; byte
        db 4CH, 44H, 5FH, 50H, 52H, 45H, 4CH, 4FH       ; 0000 _ LD_PRELO
        db 41H, 44H, 20H, 64H, 65H, 74H, 65H, 63H       ; 0008 _ AD detec
        db 74H, 65H, 64H, 20H, 74H, 68H, 72H, 6FH       ; 0010 _ ted thro
        db 75H, 67H, 68H, 20H, 65H, 6EH, 76H, 69H       ; 0018 _ ugh envi
        db 72H, 6FH, 6EH, 00H                           ; 0020 _ ron.

?_011:                                                  ; byte
        db 45H, 6EH, 76H, 69H, 72H, 6FH, 6EH, 6DH       ; 0024 _ Environm
        db 65H, 6EH, 74H, 20H, 69H, 73H, 20H, 63H       ; 002C _ ent is c
        db 6CH, 65H, 61H, 6EH, 00H                      ; 0034 _ lean.


