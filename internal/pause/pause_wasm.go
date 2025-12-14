// ©Hayabusa Cloud Co., Ltd. 2025. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

//go:build wasm

package pause

import "runtime"

// WebAssembly does not support inline assembly and has no CPU pause hint.
// To avoid burning cycles in tight spins, yield to the Go scheduler.
// Advanced users should combine Pause with higher-level backoff (e.g., Wait/Yield)
// when targeting wasm to prevent busy-waiting on the single-threaded runtime.
func pause1() { runtime.Gosched() }

// For wasm we intentionally collapse multi-cycle pauses into a single
// scheduler yield to minimize overhead and power usage.
func pauseN(int) { runtime.Gosched() }
