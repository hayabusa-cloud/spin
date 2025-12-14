// ©Hayabusa Cloud Co., Ltd. 2025. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

// Package pause contains arch-specific pause/yield primitives.
// It is an internal package; external users should call spin.Pause.
package pause

// Pause delegates to the arch-specific implementations of pause1/pauseN.
// cycles <= 0 is treated as a single minimal pause.
func Pause(cycles int) {
	if cycles <= 1 {
		pause1()
		return
	}
	pauseN(cycles)
}
