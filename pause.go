// ©Hayabusa Cloud Co., Ltd. 2025. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

package spin

import (
	pausepkg "code.hybscloud.com/spin/internal/pause"
)

const defaultPauseCycles = 30

// Pause executes CPU pause instructions to reduce energy consumption in spin-wait loops.
// On hardware-pause targets, it does not block or yield the scheduler.
// On wasm and unsupported fallback targets, it may yield via runtime.Gosched.
//
// Defaults to 30 pause hints if not specified. The cycles parameter is a
// historical repeat-count name, not a calibrated CPU-cycle delay. Uses optimized
// assembly on supported hardware-pause targets.
//
// Usage:
//
//	Pause()     // 30 pause hints (default)
//	Pause(1)    // 1 pause hint
//	Pause(50)   // 50 pause hints
func Pause(cycles ...int) {
	n := defaultPauseCycles
	if len(cycles) > 0 && cycles[0] > 0 {
		n = cycles[0]
	}
	pausepkg.Pause(n)
}
