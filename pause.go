// ©Hayabusa Cloud Co., Ltd. 2025. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

package spin

import (
	pausepkg "code.hybscloud.com/spin/internal/pause"
)

const defaultPauseCycles = 30

// Pause executes CPU pause instructions to reduce energy consumption in spin-wait loops.
// It must not block or yield the scheduler.
//
// Defaults to 30 cycles if not specified. Uses optimized assembly on amd64/arm64.
//
// Usage:
//
//	Pause()     // 30 cycles (default)
//	Pause(1)    // 1 cycle
//	Pause(50)   // 50 cycles
func Pause(cycles ...int) {
	n := defaultPauseCycles
	if len(cycles) > 0 && cycles[0] > 0 {
		n = cycles[0]
	}
	pausepkg.Pause(n)
}
