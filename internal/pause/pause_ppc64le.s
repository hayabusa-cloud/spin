// ©Hayabusa Cloud Co., Ltd. 2025. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

//go:build ppc64le

#include "textflag.h"

// func pause1()
TEXT ·pause1(SB), NOSPLIT|NOFRAME, $0-0
    OR  R27, R27, R27
    RET

// func pauseN(cycles int)
TEXT ·pauseN(SB), NOSPLIT|NOFRAME, $0-8
    MOVD cycles+0(FP), R3
loop:
    OR   R27, R27, R27
    ADD  $-1, R3, R3
    CMP  R3, $0
    BNE  loop
    RET
