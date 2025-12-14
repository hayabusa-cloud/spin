// ©Hayabusa Cloud Co., Ltd. 2025. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

//go:build arm64

package pause

// Assembly implementations in pause_arm64.s
func pause1()
func pauseN(cycles int)
