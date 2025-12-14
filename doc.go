// ©Hayabusa Cloud Co., Ltd. 2025. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

// Package spin provides minimal spin-based primitives for performance-critical
// code paths:
//   - Lock  — non-fair spinlock for extremely short critical sections
//   - Wait  — adaptive spin-wait helper for tight polling loops
//   - Pause — architecture-specific CPU hint for busy-wait loops
//   - Yield — cooperative yield/sleep for non-hot paths
//
// Design notes
//   - Hot paths should prefer Pause (light CPU hint) and Wait (adaptive backoff)
//     instead of ad-hoc loops with runtime.Gosched().
//   - Lock is intentionally non-fair and intended only for very short critical
//     sections where OS mutex overhead is prohibitive.
//   - No allocations in hot paths; functions return immediately and never block
//     unless explicitly documented (Yield with a positive sleep duration).
//
// Architectures: amd64, arm64, 386, arm, riscv64, ppc64le, s390x, loong64, wasm.
package spin
