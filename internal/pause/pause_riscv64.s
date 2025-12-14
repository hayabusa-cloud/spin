// ©Hayabusa Cloud Co., Ltd. 2025. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

//go:build riscv64

#include "textflag.h"

// func pause1()
TEXT ·pause1(SB), NOSPLIT|NOFRAME, $0-0
    FENCE
    RET

// func pauseN(cycles int)
TEXT ·pauseN(SB), NOSPLIT|NOFRAME, $0-8
    MOV  cycles+0(FP), X10
loop:
    FENCE
    ADDI $-1, X10, X10     // X10 = X10 - 1
    BNEZ X10, loop
    RET
