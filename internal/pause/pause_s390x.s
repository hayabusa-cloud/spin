// ©Hayabusa Cloud Co., Ltd. 2025. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

//go:build s390x

#include "textflag.h"

// func pause1()
TEXT ·pause1(SB), NOSPLIT|NOFRAME, $0-0
    NOP
    RET

// func pauseN(cycles int)
TEXT ·pauseN(SB), NOSPLIT|NOFRAME, $0-8
    MOVD cycles+0(FP), R1
loop:
    NOP
    ADD  $-1, R1
    CMP  R1, $0
    BNE  loop
    RET
