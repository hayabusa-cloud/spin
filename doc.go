// ©Hayabusa Cloud Co., Ltd. 2025. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

// Package spin provides synchronization primitives for short local waits.
//
// [Lock] is a non-fair lock for very short critical sections. [Lock.Lock]
// waits for ownership and [Lock.Try] makes a single acquisition attempt.
//
// [Wait] pauses briefly and yields to the Go scheduler more often as contention
// persists. [Pause] issues processor hints on supported targets. [Yield] sleeps
// for a positive duration or yields to the scheduler for a non-positive duration.
//
// Architectures: amd64, arm64, 386, arm, riscv64, ppc64le, s390x, loong64, wasm.
package spin
