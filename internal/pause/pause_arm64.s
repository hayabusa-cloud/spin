// ©Hayabusa Cloud Co., Ltd. 2025. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

//go:build arm64

#include "textflag.h"

// func pause1()
TEXT ·pause1(SB), NOSPLIT, $0-0
    YIELD
    RET

// func pauseN(cycles int)
TEXT ·pauseN(SB), NOSPLIT|NOFRAME, $0-8
    MOVD cycles+0(FP), R0
loop:
    YIELD
    SUBS $1, R0, R0
    BNE  loop
    RET
