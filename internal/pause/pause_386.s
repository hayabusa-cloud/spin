// ©Hayabusa Cloud Co., Ltd. 2025. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

//go:build 386

#include "textflag.h"

// func pause1()
TEXT ·pause1(SB), NOSPLIT, $0-0
    BYTE $0xF3 // REP prefix
    BYTE $0x90 // NOP
    RET

// func pauseN(cycles int)
TEXT ·pauseN(SB), NOSPLIT|NOFRAME, $0-4
    MOVL cycles+0(FP), AX
loop:
    BYTE $0xF3 // REP prefix
    BYTE $0x90 // NOP
    DECL AX
    JNZ  loop
    RET
