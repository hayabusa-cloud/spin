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
type Lock struct {
	_ noCopy
	n atomic.Uintptr
}

// Lock acquires the lock, spinning with an adaptive backoff.
func (sl *Lock) Lock() {
	for {
		n := sl.n.Add(1)
		if n == 1 {
			return
		} else if n < 4 {
			Pause(defaultPauseCycles)
			continue
		}
		runtime.Gosched()
	}
}

// Unlock releases the lock.
func (sl *Lock) Unlock() {
	sl.n.Store(0)
}

// Try attempts to acquire the lock without blocking.
// It returns true if the lock was acquired.
func (sl *Lock) Try() bool {
	return sl.n.Add(1) == 1
}
