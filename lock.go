// ©Hayabusa Cloud Co., Ltd. 2025. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

package spin

import (
	"runtime"
	"sync/atomic"
)

// Lock is a minimal, non-fair spin lock intended for very short
// critical sections on hot paths. It avoids allocations and OS mutex
// overhead but should not be used as a general-purpose lock.
//
// Acquisition uses FAA (Fetch-And-Add) with a TTAS (test-and-test-and-set)
// slow path. FAA completes in a single atomic instruction (LOCK XADD on
// x86, LDADDAL on arm64) — under contention, CAS produces O(n²) cache-line
// invalidations while FAA keeps traffic at O(n). The TTAS slow path spins
// on read-only Load (MESI Shared) to avoid bouncing the cache line, and
// only attempts FAA when the lock appears free.
type Lock struct {
	_ noCopy
	n atomic.Uintptr
}

// Lock acquires the lock.
// Fast path: single FAA. Slow path: TTAS (test-and-test-and-set) —
// spin on read-only Load to keep the cache line in MESI Shared state,
// then attempt FAA only when the lock appears free.
func (sl *Lock) Lock() {
	if sl.n.Add(1) == 1 {
		return
	}
	sl.lockSlow()
}

func (sl *Lock) lockSlow() {
	for i := 0; ; i++ {
		// Spin on Load (read-only): multiple cores hold the line in
		// Shared state with zero cross-core invalidations. Only
		// transition to Add (write) when the lock appears released.
		if sl.n.Load() == 0 {
			if sl.n.Add(1) == 1 {
				return
			}
		}
		if i >= 4 {
			runtime.Gosched()
		} else {
			Pause(defaultPauseCycles)
		}
	}
}

// Unlock releases the lock.
func (sl *Lock) Unlock() {
	sl.n.Store(0)
}

// Try attempts to acquire the lock without blocking.
// It returns true if the lock was acquired.
//
// Uses FAA rather than CAS: on x86 LOCK CMPXCHG takes exclusive
// cache-line ownership regardless of comparison outcome, so a failed
// CAS costs the same coherence traffic as FAA. FAA keeps the atomic
// pattern consistent across the Lock API.
//
// A failed Try increments n without a corresponding decrement.
// This is intentional: Unlock resets n to zero via Store(0),
// clearing all accumulated increments. The uintptr counter space
// is large enough that overflow is not a practical concern.
func (sl *Lock) Try() bool {
	return sl.n.Add(1) == 1
}
