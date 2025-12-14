// ©Hayabusa Cloud Co., Ltd. 2025. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

//go:build s390x

package pause

// Assembly implementations in pause_s390x.s
func pause1()
func pauseN(cycles int)
