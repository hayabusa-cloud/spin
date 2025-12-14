// ©Hayabusa Cloud Co., Ltd. 2025. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

//go:build loong64

#include "textflag.h"

// func pause1()
TEXT ·pause1(SB), NOSPLIT|NOFRAME, $0-0
    OR  R0, R0, R0
    RET

// func pauseN(cycles int)
TEXT ·pauseN(SB), NOSPLIT|NOFRAME, $0-8
    MOVV cycles+0(FP), R0
loop:
    OR    R0, R0, R0
    ADDV  $-1, R0, R0
    BNE   R0, loop
    RET
