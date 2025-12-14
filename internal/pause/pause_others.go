// ©Hayabusa Cloud Co., Ltd. 2025. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

//go:build !amd64 && !arm64 && !riscv64 && !ppc64le && !s390x && !loong64 && !386 && !wasm

package pause

import "runtime"

func pause1() {
	runtime.Gosched()
}

func pauseN(int) {
	runtime.Gosched()
}
